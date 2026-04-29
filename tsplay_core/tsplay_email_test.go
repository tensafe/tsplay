package tsplay_core

import (
	"bufio"
	"bytes"
	"io"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

type smtpTestMessage struct {
	From       string
	Recipients []string
	Data       []byte
}

type smtpTestServer struct {
	listener net.Listener
	mu       sync.Mutex
	messages []smtpTestMessage
}

func newSMTPTestServer(t *testing.T) *smtpTestServer {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen smtp test server: %v", err)
	}
	server := &smtpTestServer{listener: listener}
	go server.serve()
	return server
}

func (s *smtpTestServer) Addr() string {
	if s == nil || s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func (s *smtpTestServer) Close() {
	if s == nil || s.listener == nil {
		return
	}
	_ = s.listener.Close()
}

func (s *smtpTestServer) LastMessage(t *testing.T) smtpTestMessage {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.messages) == 0 {
		t.Fatalf("expected at least one smtp message")
	}
	return s.messages[len(s.messages)-1]
}

func (s *smtpTestServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}

func (s *smtpTestServer) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString("220 smtp.test ESMTP ready\r\n"); err != nil {
		return
	}
	if err := writer.Flush(); err != nil {
		return
	}

	current := smtpTestMessage{}
	inData := false
	var data bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")

		if inData {
			if line == "." {
				current.Data = append([]byte(nil), data.Bytes()...)
				s.mu.Lock()
				s.messages = append(s.messages, current)
				s.mu.Unlock()
				current = smtpTestMessage{}
				data.Reset()
				inData = false
				if _, err := writer.WriteString("250 OK\r\n"); err != nil {
					return
				}
				if err := writer.Flush(); err != nil {
					return
				}
				continue
			}
			if strings.HasPrefix(line, "..") {
				line = line[1:]
			}
			if _, err := data.WriteString(line + "\r\n"); err != nil {
				return
			}
			continue
		}

		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "EHLO "):
			if _, err := writer.WriteString("250-smtp.test\r\n250-AUTH PLAIN LOGIN\r\n250 OK\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "HELO "):
			if _, err := writer.WriteString("250 OK\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "AUTH "):
			if _, err := writer.WriteString("235 Authentication successful\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "MAIL FROM:"):
			current.From = smtpEnvelopeAddress(line[len("MAIL FROM:"):])
			if _, err := writer.WriteString("250 OK\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "RCPT TO:"):
			current.Recipients = append(current.Recipients, smtpEnvelopeAddress(line[len("RCPT TO:"):]))
			if _, err := writer.WriteString("250 OK\r\n"); err != nil {
				return
			}
		case upper == "DATA":
			inData = true
			if _, err := writer.WriteString("354 End data with <CR><LF>.<CR><LF>\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "RSET"):
			current = smtpTestMessage{}
			data.Reset()
			if _, err := writer.WriteString("250 OK\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "NOOP"):
			if _, err := writer.WriteString("250 OK\r\n"); err != nil {
				return
			}
		case strings.HasPrefix(upper, "QUIT"):
			_, _ = writer.WriteString("221 Bye\r\n")
			_ = writer.Flush()
			return
		default:
			if _, err := writer.WriteString("502 unsupported\r\n"); err != nil {
				return
			}
		}
		if err := writer.Flush(); err != nil {
			return
		}
	}
}

func smtpEnvelopeAddress(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "<")
	value = strings.TrimSuffix(value, ">")
	return strings.TrimSpace(value)
}

func mustEmailAddressList(t *testing.T, value any, field string) []mail.Address {
	t.Helper()
	addresses, err := emailAddressListValue(value, field)
	if err != nil {
		t.Fatalf("parse email addresses for %s: %v", field, err)
	}
	return addresses
}

func TestValidateFlowSecurityRejectsEmailByDefault(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "email_policy",
		Steps: []FlowStep{
			{
				Action: "send_email",
				With: map[string]any{
					"to":      "ops@example.com",
					"subject": "TSPlay test",
					"body":    "done",
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, DefaultFlowSecurityPolicy())
	if err == nil {
		t.Fatalf("expected email security policy error")
	}
	if !strings.Contains(err.Error(), "allow_email") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFlowSecurityRejectsEmailAttachmentsWithoutFileAccess(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "email_attachment_policy",
		Steps: []FlowStep{
			{
				Action: "send_email",
				With: map[string]any{
					"to":          "ops@example.com",
					"subject":     "TSPlay test",
					"body":        "done",
					"attachments": []any{"artifacts/report.txt"},
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
	err := ValidateFlowSecurity(flow, FlowSecurityPolicy{AllowEmail: true})
	if err == nil {
		t.Fatalf("expected file access security policy error")
	}
	if !strings.Contains(err.Error(), "allow_file_access") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFlowSendEmailAcceptsPlaceholderRecipientList(t *testing.T) {
	flow := &Flow{
		SchemaVersion: "1",
		Name:          "email_placeholder_recipient_list",
		Vars: map[string]any{
			"recipient_email": "ops@example.com",
		},
		Steps: []FlowStep{
			{
				Action: "send_email",
				With: map[string]any{
					"to":      []any{"{{recipient_email}}"},
					"subject": "TSPlay test",
					"body":    "done",
				},
			},
		},
	}

	if err := ValidateFlow(flow); err != nil {
		t.Fatalf("validate flow: %v", err)
	}
}

func TestRunFlowSendEmailStep(t *testing.T) {
	server := newSMTPTestServer(t)
	defer server.Close()

	host, portText, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("split smtp addr: %v", err)
	}
	t.Setenv("TSPLAY_EMAIL_ALERTS_HOST", host)
	t.Setenv("TSPLAY_EMAIL_ALERTS_PORT", portText)
	t.Setenv("TSPLAY_EMAIL_ALERTS_FROM", "TSPlay Bot <bot@example.com>")
	t.Setenv("TSPLAY_EMAIL_ALERTS_TLS_MODE", "none")

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "send_email_flow",
		Steps: []FlowStep{
			{
				Action:     "send_email",
				Connection: "alerts",
				With: map[string]any{
					"to":      []any{"alice@example.com", "Bob <bob@example.com>"},
					"cc":      "ops@example.com",
					"bcc":     []any{"audit@example.com"},
					"subject": "导入完成",
					"body":    "Rows: 2",
					"headers": map[string]any{"X-TSPlay-Run": "123"},
				},
				SaveAs: "email_result",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowEmail: true},
	})
	if err != nil {
		t.Fatalf("run flow: %v", err)
	}
	if result == nil {
		t.Fatalf("expected flow result")
	}

	message := server.LastMessage(t)
	if message.From != "bot@example.com" {
		t.Fatalf("smtp MAIL FROM = %q, want %q", message.From, "bot@example.com")
	}
	if len(message.Recipients) != 4 {
		t.Fatalf("smtp recipients = %#v", message.Recipients)
	}

	parsed, err := mail.ReadMessage(bytes.NewReader(message.Data))
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	if parsed.Header.Get("To") == "" || parsed.Header.Get("Cc") == "" {
		t.Fatalf("expected To and Cc headers, got %#v", parsed.Header)
	}
	if parsed.Header.Get("Bcc") != "" {
		t.Fatalf("unexpected Bcc header in message: %#v", parsed.Header.Get("Bcc"))
	}
	if parsed.Header.Get("X-Tsplay-Run") != "123" {
		t.Fatalf("unexpected custom header: %#v", parsed.Header.Get("X-Tsplay-Run"))
	}
	if parsed.Header.Get("Subject") == "" || !strings.Contains(string(message.Data), "=?utf-8?") {
		t.Fatalf("expected RFC 2047 encoded subject, got raw message %q", string(message.Data))
	}
	if parsed.Header.Get("Content-Transfer-Encoding") != "quoted-printable" {
		t.Fatalf("unexpected transfer encoding: %#v", parsed.Header.Get("Content-Transfer-Encoding"))
	}

	bodyReader := parsed.Body
	if strings.EqualFold(parsed.Header.Get("Content-Transfer-Encoding"), "quoted-printable") {
		bodyReader = quotedPrintableReader(parsed.Body)
	}
	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if strings.TrimSpace(string(bodyBytes)) != "Rows: 2" {
		t.Fatalf("unexpected body: %q", string(bodyBytes))
	}

	rawResult, ok := result.Vars["email_result"].(map[string]any)
	if !ok {
		t.Fatalf("expected saved email_result, got %#v", result.Vars["email_result"])
	}
	switch got := rawResult["recipient_count"].(type) {
	case int:
		if got != 4 {
			t.Fatalf("unexpected flow result: %#v", rawResult)
		}
	case float64:
		if got != 4 {
			t.Fatalf("unexpected flow result: %#v", rawResult)
		}
	default:
		t.Fatalf("unexpected flow result type: %#v", rawResult)
	}
}

func TestRunFlowSendEmailStepWithInlineSMTP(t *testing.T) {
	server := newSMTPTestServer(t)
	defer server.Close()

	host, portText, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("split smtp addr: %v", err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("parse smtp port: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "send_email_inline_smtp_flow",
		Vars: map[string]any{
			"sender_email":    "bot@example.com",
			"sender_password": "mailpwd",
			"smtp_host":       host,
		},
		Steps: []FlowStep{
			{
				Action: "send_email",
				With: map[string]any{
					"to":      "alice@example.com",
					"subject": "inline smtp config",
					"body":    "sent with smtp block",
					"smtp": map[string]any{
						"host":     "{{smtp_host}}",
						"port":     port,
						"username": "{{sender_email}}",
						"password": "{{sender_password}}",
						"from":     "{{sender_email}}",
						"tls_mode": "none",
					},
				},
				SaveAs: "email_result",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowEmail: true},
	})
	if err != nil {
		t.Fatalf("run flow with inline smtp: %v", err)
	}
	if result == nil {
		t.Fatalf("expected flow result")
	}

	message := server.LastMessage(t)
	if message.From != "bot@example.com" {
		t.Fatalf("smtp MAIL FROM = %q, want %q", message.From, "bot@example.com")
	}
	if len(message.Recipients) != 1 || message.Recipients[0] != "alice@example.com" {
		t.Fatalf("smtp recipients = %#v", message.Recipients)
	}
}

func TestRunFlowSendEmailStepWithAttachment(t *testing.T) {
	server := newSMTPTestServer(t)
	defer server.Close()

	root := t.TempDir()
	attachmentPath := filepath.Join(root, "reports", "hello.txt")
	if err := os.MkdirAll(filepath.Dir(attachmentPath), 0755); err != nil {
		t.Fatalf("mkdir attachment dir: %v", err)
	}
	if err := os.WriteFile(attachmentPath, []byte("Hello attachment"), 0600); err != nil {
		t.Fatalf("write attachment: %v", err)
	}

	host, portText, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("split smtp addr: %v", err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("parse smtp port: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "send_email_attachment_flow",
		Vars: map[string]any{
			"smtp_host": host,
		},
		Steps: []FlowStep{
			{
				Action: "send_email",
				With: map[string]any{
					"to":      "alice@example.com",
					"subject": "attachment smtp config",
					"body":    "sent with attachment",
					"attachments": []any{
						map[string]any{
							"path":         "reports/hello.txt",
							"name":         "custom-name.txt",
							"content_type": "text/plain",
						},
					},
					"smtp": map[string]any{
						"host":     "{{smtp_host}}",
						"port":     port,
						"username": "bot@example.com",
						"password": "mailpwd",
						"from":     "bot@example.com",
						"tls_mode": "none",
					},
				},
				SaveAs: "email_result",
			},
		},
	}

	result, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{
			AllowEmail:      true,
			AllowFileAccess: true,
			FileInputRoot:   root,
			FileOutputRoot:  root,
		},
	})
	if err != nil {
		t.Fatalf("run flow with attachment: %v", err)
	}
	if result == nil {
		t.Fatalf("expected flow result")
	}

	message := server.LastMessage(t)
	raw := string(message.Data)
	for _, want := range []string{
		"Content-Type: multipart/mixed;",
		"Content-Disposition: attachment; filename=\"custom-name.txt\"",
		"Content-Type: text/plain; name=\"custom-name.txt\"",
		"SGVsbG8gYXR0YWNobWVudA==",
	} {
		if !strings.Contains(raw, want) {
			t.Fatalf("expected %q in message:\n%s", want, raw)
		}
	}

	rawResult, ok := result.Vars["email_result"].(map[string]any)
	if !ok {
		t.Fatalf("expected saved email_result, got %#v", result.Vars["email_result"])
	}
	if got, ok := rawResult["attachment_count"].(int); ok {
		if got != 1 {
			t.Fatalf("unexpected attachment_count in flow result: %#v", rawResult)
		}
	} else if got, ok := rawResult["attachment_count"].(float64); ok {
		if got != 1 {
			t.Fatalf("unexpected attachment_count in flow result: %#v", rawResult)
		}
	} else {
		t.Fatalf("unexpected attachment_count type in flow result: %#v", rawResult)
	}
}

func TestRunFlowLuaSendEmailHonorsAllowEmail(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	flow := &Flow{
		SchemaVersion: "1",
		Name:          "lua_email_policy",
		Steps: []FlowStep{
			{
				Action: "lua",
				Code: `return send_email({
  to = "ops@example.com",
  subject = "TSPlay test",
  body = "done"
})`,
			},
		},
	}

	_, err := RunFlowInStateWithOptions(L, flow, FlowRunOptions{
		Security: &FlowSecurityPolicy{AllowLua: true},
	})
	if err == nil {
		t.Fatalf("expected allow_email runtime error")
	}
	if !strings.Contains(err.Error(), "allow_email") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildSendEmailMessageMultipartAlternative(t *testing.T) {
	config := flowEmailConfig{
		To:      mustEmailAddressList(t, "ops@example.com", "to"),
		Subject: "导出完成",
		Body:    "plain text",
		HTML:    "<p>html body</p>",
	}
	connection := emailConnectionConfig{
		From: "TSPlay Bot <bot@example.com>",
	}

	message, from, recipients, err := buildSendEmailMessage(config, connection)
	if err != nil {
		t.Fatalf("build message: %v", err)
	}
	if from.Address != "bot@example.com" {
		t.Fatalf("unexpected from: %#v", from)
	}
	if len(recipients) != 1 || recipients[0].Address != "ops@example.com" {
		t.Fatalf("unexpected recipients: %#v", recipients)
	}
	raw := string(message)
	for _, want := range []string{
		"Content-Type: multipart/alternative;",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Type: text/html; charset=UTF-8",
		"Subject: =?utf-8?",
	} {
		if !strings.Contains(raw, want) {
			t.Fatalf("expected %q in message:\n%s", want, raw)
		}
	}
}

func quotedPrintableReader(reader io.Reader) io.Reader {
	return quotedprintable.NewReader(bufio.NewReader(reader))
}

func TestResolveEmailConnectionConfigUsesStructuredEnv(t *testing.T) {
	t.Setenv("TSPLAY_EMAIL_REPORTS_HOST", "smtp.example.com")
	t.Setenv("TSPLAY_EMAIL_REPORTS_PORT", "2525")
	t.Setenv("TSPLAY_EMAIL_REPORTS_FROM", "reports@example.com")
	t.Setenv("TSPLAY_EMAIL_REPORTS_TLS_MODE", "none")
	t.Setenv("TSPLAY_EMAIL_REPORTS_TIMEOUT_MS", "9000")

	config, err := resolveEmailConnectionConfig("reports")
	if err != nil {
		t.Fatalf("resolve email config: %v", err)
	}
	if config.Host != "smtp.example.com" || config.Port != 2525 {
		t.Fatalf("unexpected email config: %#v", config)
	}
	if config.TimeoutMS != 9000 || config.TLSMode != emailTLSModeNone {
		t.Fatalf("unexpected email config details: %#v", config)
	}
}

func TestNormalizeSendEmailConfigRejectsReservedHeaders(t *testing.T) {
	_, err := normalizeSendEmailConfig(map[string]any{
		"to":      "ops@example.com",
		"subject": "done",
		"body":    "ok",
		"headers": map[string]any{"Subject": "override"},
	})
	if err == nil {
		t.Fatalf("expected reserved header error")
	}
	if !strings.Contains(err.Error(), "reserved header") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormalizeSendEmailConfigSupportsInlineSMTP(t *testing.T) {
	config, err := normalizeSendEmailConfig(map[string]any{
		"to":      "ops@example.com",
		"subject": "done",
		"body":    "ok",
		"smtp": map[string]any{
			"host":     "smtp.qq.com",
			"port":     465,
			"username": "sender@qq.com",
			"password": "mailpwd",
			"from":     "sender@qq.com",
			"tls_mode": "tls",
		},
	})
	if err != nil {
		t.Fatalf("normalize send_email config: %v", err)
	}
	if config.SMTP == nil {
		t.Fatalf("expected inline smtp config")
	}
	if config.SMTP.Host != "smtp.qq.com" || config.SMTP.Port != 465 {
		t.Fatalf("unexpected smtp config: %#v", config.SMTP)
	}
	if config.SMTP.Password != "mailpwd" || config.SMTP.TLSMode != emailTLSModeTLS {
		t.Fatalf("unexpected smtp config details: %#v", config.SMTP)
	}
}

func TestNormalizeSendEmailConfigSupportsAttachments(t *testing.T) {
	config, err := normalizeSendEmailConfig(map[string]any{
		"to":      "ops@example.com",
		"subject": "done",
		"body":    "ok",
		"attachments": []any{
			map[string]any{
				"path":         "artifacts/report.xlsx",
				"name":         "final-report.xlsx",
				"content_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			},
		},
	})
	if err != nil {
		t.Fatalf("normalize send_email config with attachments: %v", err)
	}
	if len(config.Attachments) != 1 {
		t.Fatalf("expected one attachment, got %#v", config.Attachments)
	}
	if config.Attachments[0].Name != "final-report.xlsx" {
		t.Fatalf("unexpected attachment config: %#v", config.Attachments[0])
	}
}

func TestNormalizeSendEmailConfigSupportsSingleAttachmentObject(t *testing.T) {
	config, err := normalizeSendEmailConfig(map[string]any{
		"to":      "ops@example.com",
		"subject": "done",
		"body":    "ok",
		"attachments": map[string]any{
			"path":         "artifacts/report.xlsx",
			"name":         "final-report.xlsx",
			"content_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
	})
	if err != nil {
		t.Fatalf("normalize send_email config with single attachment object: %v", err)
	}
	if len(config.Attachments) != 1 {
		t.Fatalf("expected one attachment, got %#v", config.Attachments)
	}
	if config.Attachments[0].Path != "artifacts/report.xlsx" {
		t.Fatalf("unexpected attachment config: %#v", config.Attachments[0])
	}
}

func TestSMTPTestServerAddrIsNumericPort(t *testing.T) {
	server := newSMTPTestServer(t)
	defer server.Close()

	_, portText, err := net.SplitHostPort(server.Addr())
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	if _, err := strconv.Atoi(portText); err != nil {
		t.Fatalf("expected numeric port, got %q", portText)
	}
}
