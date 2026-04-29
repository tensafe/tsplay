package tsplay_core

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const emailDefaultConnection = "default"

type emailTLSMode string

const (
	emailTLSModeNone     emailTLSMode = "none"
	emailTLSModeStartTLS emailTLSMode = "starttls"
	emailTLSModeTLS      emailTLSMode = "tls"
)

type flowEmailConfig struct {
	To          []mail.Address
	CC          []mail.Address
	BCC         []mail.Address
	Subject     string
	Body        string
	HTML        string
	Attachments []emailAttachmentConfig
	Headers     map[string]string
	Connection  string
	FromEmail   string
	ReplyTo     string
	SMTP        *emailConnectionConfig
	TimeoutMS   int
}

type emailAttachmentConfig struct {
	Path        string
	Name        string
	ContentType string
}

type emailConnectionConfig struct {
	Name               string
	Host               string
	Port               int
	Username           string
	Password           string
	From               string
	TLSMode            emailTLSMode
	HelloName          string
	InsecureSkipVerify bool
	TimeoutMS          int
}

func send_email(L *lua.LState) int {
	values, err := emailValuesFromLua(L)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	if err := luaEmailExecutionAllowed(L, "send_email"); err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	config, err := normalizeSendEmailConfig(values)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	config, err = applyLuaEmailRuntimePolicy(L, config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	flowCtx := flowContextFromState(L)
	runCtx := context.Background()
	if flowCtx != nil && flowCtx.Context != nil {
		runCtx = flowCtx.Context
	}

	result, err := executeSendEmail(runCtx, config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func runFlowSendEmailStep(ctx *FlowContext, step FlowStep) (any, error) {
	values, err := resolvedSendEmailValues(ctx, step)
	if err != nil {
		return nil, err
	}
	config, err := normalizeSendEmailConfig(values)
	if err != nil {
		return nil, err
	}
	config, err = applyFlowEmailRuntimePolicy(ctx, config)
	if err != nil {
		return nil, err
	}
	runCtx := context.Background()
	if ctx != nil && ctx.Context != nil {
		runCtx = ctx.Context
	}
	return executeSendEmail(runCtx, config)
}

func emailValuesFromLua(L *lua.LState) (map[string]any, error) {
	if L == nil || L.GetTop() == 0 {
		return nil, fmt.Errorf("send_email requires either a config table or to/subject/body arguments")
	}

	first := luaValueToGo(L.CheckAny(1))
	if values, ok := first.(map[string]any); ok {
		return values, nil
	}

	values := map[string]any{"to": first}
	if L.GetTop() >= 2 {
		values["subject"] = luaValueToGo(L.CheckAny(2))
	}
	if L.GetTop() >= 3 {
		values["body"] = luaValueToGo(L.CheckAny(3))
	}
	if L.GetTop() >= 4 {
		values["connection"] = luaValueToGo(L.CheckAny(4))
	}
	if L.GetTop() >= 5 {
		values["timeout"] = luaValueToGo(L.CheckAny(5))
	}
	return values, nil
}

func resolvedSendEmailValues(ctx *FlowContext, step FlowStep) (map[string]any, error) {
	names := []string{
		"to",
		"cc",
		"bcc",
		"subject",
		"body",
		"html",
		"headers",
		"attachments",
		"connection",
		"from_email",
		"reply_to",
		"smtp",
		"timeout",
	}
	values := map[string]any{}
	for _, name := range names {
		value, ok, err := flowStepResolvedParam(ctx, step, name)
		if err != nil {
			return nil, err
		}
		if ok {
			values[name] = value
		}
	}
	return values, nil
}

func luaEmailExecutionAllowed(L *lua.LState, action string) error {
	flowCtx := flowContextFromState(L)
	if flowCtx != nil && flowCtx.Security != nil && !flowCtx.Security.AllowEmail {
		return fmt.Errorf("%s is disabled by security policy; set allow_email=true only for trusted flows", action)
	}
	return nil
}

func normalizeSendEmailConfig(values map[string]any) (flowEmailConfig, error) {
	config := flowEmailConfig{
		Headers: map[string]string{},
	}

	toValue, ok := values["to"]
	if !ok {
		return flowEmailConfig{}, fmt.Errorf("send_email requires to")
	}
	to, err := emailAddressListValue(toValue, "to")
	if err != nil {
		return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
	}
	config.To = to

	if ccValue, ok := values["cc"]; ok {
		cc, err := emailAddressListValue(ccValue, "cc")
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
		}
		config.CC = cc
	}
	if bccValue, ok := values["bcc"]; ok {
		bcc, err := emailAddressListValue(bccValue, "bcc")
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
		}
		config.BCC = bcc
	}

	subjectValue, ok := values["subject"]
	if !ok || strings.TrimSpace(fmt.Sprint(subjectValue)) == "" {
		return flowEmailConfig{}, fmt.Errorf("send_email requires subject")
	}
	if err := validateEmailHeaderText("subject", fmt.Sprint(subjectValue)); err != nil {
		return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
	}
	config.Subject = fmt.Sprint(subjectValue)

	if bodyValue, ok := values["body"]; ok {
		bodyText, ok := bodyValue.(string)
		if !ok {
			return flowEmailConfig{}, fmt.Errorf("send_email body must be a string")
		}
		config.Body = bodyText
	}
	if htmlValue, ok := values["html"]; ok {
		htmlText, ok := htmlValue.(string)
		if !ok {
			return flowEmailConfig{}, fmt.Errorf("send_email html must be a string")
		}
		config.HTML = htmlText
	}
	if strings.TrimSpace(config.Body) == "" && strings.TrimSpace(config.HTML) == "" {
		return flowEmailConfig{}, fmt.Errorf("send_email requires body or html")
	}

	if headersValue, ok := values["headers"]; ok {
		headers, err := normalizeEmailHeaders(headersValue)
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
		}
		config.Headers = headers
	}
	if attachmentsValue, ok := values["attachments"]; ok {
		attachments, err := normalizeEmailAttachments(attachmentsValue)
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
		}
		config.Attachments = attachments
	}
	if connectionValue, ok := values["connection"]; ok {
		config.Connection = strings.TrimSpace(fmt.Sprint(connectionValue))
	}
	if fromValue, ok := values["from_email"]; ok {
		text := strings.TrimSpace(fmt.Sprint(fromValue))
		if text == "" {
			return flowEmailConfig{}, fmt.Errorf("send_email from_email cannot be blank")
		}
		if _, err := parseSingleEmailAddress(text, "from_email"); err != nil {
			return flowEmailConfig{}, err
		}
		config.FromEmail = text
	}
	if replyToValue, ok := values["reply_to"]; ok {
		text := strings.TrimSpace(fmt.Sprint(replyToValue))
		if text == "" {
			return flowEmailConfig{}, fmt.Errorf("send_email reply_to cannot be blank")
		}
		if _, err := parseSingleEmailAddress(text, "reply_to"); err != nil {
			return flowEmailConfig{}, err
		}
		config.ReplyTo = text
	}
	if smtpValue, ok := values["smtp"]; ok {
		smtpConfig, err := normalizeInlineEmailConnectionConfig(smtpValue)
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email %w", err)
		}
		config.SMTP = &smtpConfig
	}
	if timeoutValue, ok := values["timeout"]; ok {
		timeoutMS, err := intParam(timeoutValue)
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email timeout %w", err)
		}
		if timeoutMS < 1 {
			return flowEmailConfig{}, fmt.Errorf("send_email timeout must be at least 1")
		}
		config.TimeoutMS = timeoutMS
	}

	return config, nil
}

func executeSendEmail(runCtx context.Context, config flowEmailConfig) (map[string]any, error) {
	connectionConfig := emailConnectionConfig{}
	var err error
	switch {
	case config.SMTP != nil && strings.TrimSpace(config.Connection) != "":
		baseConfig, resolveErr := resolveEmailConnectionConfig(config.Connection)
		if resolveErr != nil {
			return nil, resolveErr
		}
		connectionConfig, err = finalizeEmailConnectionConfig(
			mergeEmailConnectionConfig(baseConfig, *config.SMTP),
			false,
		)
		if err != nil {
			return nil, err
		}
	case config.SMTP != nil:
		inlineConfig := mergeEmailConnectionConfig(defaultEmailConnectionConfig("inline"), *config.SMTP)
		connectionConfig, err = finalizeEmailConnectionConfig(inlineConfig, false)
		if err != nil {
			return nil, err
		}
	case strings.TrimSpace(config.Connection) != "":
		connectionConfig, err = resolveEmailConnectionConfig(config.Connection)
		if err != nil {
			return nil, err
		}
	default:
		connectionConfig, err = resolveEmailConnectionConfig(config.Connection)
		if err != nil {
			return nil, err
		}
	}
	if config.TimeoutMS > 0 {
		connectionConfig.TimeoutMS = config.TimeoutMS
	}

	message, envelopeFrom, envelopeRecipients, err := buildSendEmailMessage(config, connectionConfig)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(connectionConfig.TimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if runCtx == nil {
		runCtx = context.Background()
	}

	address := net.JoinHostPort(connectionConfig.Host, strconv.Itoa(connectionConfig.Port))
	dialer := net.Dialer{Timeout: timeout}
	tlsConfig := &tls.Config{
		ServerName:         connectionConfig.Host,
		InsecureSkipVerify: connectionConfig.InsecureSkipVerify,
	}

	var conn net.Conn
	if connectionConfig.TLSMode == emailTLSModeTLS {
		conn, err = tls.DialWithDialer(&dialer, "tcp", address, tlsConfig)
	} else {
		conn, err = dialer.DialContext(runCtx, "tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("send_email connect %q: %w", connectionConfig.Name, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	client, err := smtp.NewClient(conn, connectionConfig.Host)
	if err != nil {
		return nil, fmt.Errorf("send_email open smtp client %q: %w", connectionConfig.Name, err)
	}
	defer client.Close()

	if strings.TrimSpace(connectionConfig.HelloName) != "" {
		if err := client.Hello(connectionConfig.HelloName); err != nil {
			return nil, fmt.Errorf("send_email hello %q: %w", connectionConfig.Name, err)
		}
	}
	if connectionConfig.TLSMode == emailTLSModeStartTLS {
		ok, _ := client.Extension("STARTTLS")
		if !ok {
			return nil, fmt.Errorf("send_email connection %q does not advertise STARTTLS", connectionConfig.Name)
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("send_email starttls %q: %w", connectionConfig.Name, err)
		}
	}
	if connectionConfig.Username != "" || connectionConfig.Password != "" {
		if connectionConfig.Username == "" {
			return nil, fmt.Errorf("send_email connection %q requires USERNAME when PASSWORD is set", connectionConfig.Name)
		}
		auth := smtp.PlainAuth("", connectionConfig.Username, connectionConfig.Password, connectionConfig.Host)
		if err := client.Auth(auth); err != nil {
			return nil, fmt.Errorf("send_email auth %q: %w", connectionConfig.Name, err)
		}
	}
	if err := client.Mail(envelopeFrom.Address); err != nil {
		return nil, fmt.Errorf("send_email MAIL FROM %q: %w", envelopeFrom.Address, err)
	}
	for _, recipient := range envelopeRecipients {
		if err := client.Rcpt(recipient.Address); err != nil {
			return nil, fmt.Errorf("send_email RCPT TO %q: %w", recipient.Address, err)
		}
	}
	writer, err := client.Data()
	if err != nil {
		return nil, fmt.Errorf("send_email DATA: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("send_email write message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("send_email finalize message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return nil, fmt.Errorf("send_email QUIT: %w", err)
	}

	return map[string]any{
		"ok":               true,
		"connection":       connectionConfig.Name,
		"from":             envelopeFrom.String(),
		"to":               emailAddressStrings(config.To),
		"cc":               emailAddressStrings(config.CC),
		"bcc_count":        len(config.BCC),
		"recipient_count":  len(envelopeRecipients),
		"attachment_count": len(config.Attachments),
		"subject":          config.Subject,
		"message_bytes":    len(message),
	}, nil
}

func buildSendEmailMessage(config flowEmailConfig, connection emailConnectionConfig) ([]byte, mail.Address, []mail.Address, error) {
	fromText := strings.TrimSpace(config.FromEmail)
	if fromText == "" {
		fromText = strings.TrimSpace(connection.From)
	}
	if fromText == "" {
		return nil, mail.Address{}, nil, fmt.Errorf("send_email requires from_email or TSPLAY_EMAIL_*_FROM")
	}
	fromAddress, err := parseSingleEmailAddress(fromText, "from_email")
	if err != nil {
		return nil, mail.Address{}, nil, err
	}

	var replyTo *mail.Address
	if strings.TrimSpace(config.ReplyTo) != "" {
		parsed, err := parseSingleEmailAddress(config.ReplyTo, "reply_to")
		if err != nil {
			return nil, mail.Address{}, nil, err
		}
		replyTo = &parsed
	}

	recipients := make([]mail.Address, 0, len(config.To)+len(config.CC)+len(config.BCC))
	recipients = append(recipients, config.To...)
	recipients = append(recipients, config.CC...)
	recipients = append(recipients, config.BCC...)

	headers := []struct {
		Name  string
		Value string
	}{
		{Name: "From", Value: fromAddress.String()},
		{Name: "To", Value: strings.Join(emailAddressStrings(config.To), ", ")},
		{Name: "Subject", Value: encodeEmailHeaderText(config.Subject)},
		{Name: "Date", Value: time.Now().Format(time.RFC1123Z)},
		{Name: "Message-ID", Value: buildEmailMessageID(fromAddress.Address)},
		{Name: "MIME-Version", Value: "1.0"},
	}
	if len(config.CC) > 0 {
		headers = append(headers, struct {
			Name  string
			Value string
		}{Name: "Cc", Value: strings.Join(emailAddressStrings(config.CC), ", ")})
	}
	if replyTo != nil {
		headers = append(headers, struct {
			Name  string
			Value string
		}{Name: "Reply-To", Value: replyTo.String()})
	}
	for key, value := range config.Headers {
		headers = append(headers, struct {
			Name  string
			Value string
		}{Name: key, Value: value})
	}

	body, contentType, transferEncoding, err := buildEmailBody(config)
	if err != nil {
		return nil, mail.Address{}, nil, err
	}
	headers = append(headers, struct {
		Name  string
		Value string
	}{Name: "Content-Type", Value: contentType})
	if transferEncoding != "" {
		headers = append(headers, struct {
			Name  string
			Value string
		}{Name: "Content-Transfer-Encoding", Value: transferEncoding})
	}

	var buffer bytes.Buffer
	for _, header := range headers {
		if _, err := fmt.Fprintf(&buffer, "%s: %s\r\n", header.Name, header.Value); err != nil {
			return nil, mail.Address{}, nil, err
		}
	}
	if _, err := buffer.WriteString("\r\n"); err != nil {
		return nil, mail.Address{}, nil, err
	}
	if _, err := buffer.Write(body); err != nil {
		return nil, mail.Address{}, nil, err
	}
	return buffer.Bytes(), fromAddress, recipients, nil
}

func buildEmailBody(config flowEmailConfig) ([]byte, string, string, error) {
	if len(config.Attachments) > 0 {
		return buildEmailBodyWithAttachments(config)
	}
	return buildInlineEmailBody(config)
}

func buildInlineEmailBody(config flowEmailConfig) ([]byte, string, string, error) {
	if strings.TrimSpace(config.Body) != "" && strings.TrimSpace(config.HTML) != "" {
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		parts := []struct {
			ContentType string
			Value       string
		}{
			{ContentType: "text/plain; charset=UTF-8", Value: config.Body},
			{ContentType: "text/html; charset=UTF-8", Value: config.HTML},
		}
		for _, part := range parts {
			header := textproto.MIMEHeader{}
			header.Set("Content-Type", part.ContentType)
			header.Set("Content-Transfer-Encoding", "quoted-printable")
			partWriter, err := writer.CreatePart(header)
			if err != nil {
				return nil, "", "", err
			}
			if err := writeQuotedPrintable(partWriter, part.Value); err != nil {
				return nil, "", "", err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, "", "", err
		}
		return buffer.Bytes(), fmt.Sprintf("multipart/alternative; boundary=%q", writer.Boundary()), "", nil
	}

	contentType := "text/plain; charset=UTF-8"
	value := config.Body
	if strings.TrimSpace(config.HTML) != "" {
		contentType = "text/html; charset=UTF-8"
		value = config.HTML
	}

	var buffer bytes.Buffer
	if err := writeQuotedPrintable(&buffer, value); err != nil {
		return nil, "", "", err
	}
	return buffer.Bytes(), contentType, "quoted-printable", nil
}

func buildEmailBodyWithAttachments(config flowEmailConfig) ([]byte, string, string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	bodyHeader := textproto.MIMEHeader{}
	bodyContentType, bodyTransferEncoding := "text/plain; charset=UTF-8", "quoted-printable"
	if strings.TrimSpace(config.Body) != "" && strings.TrimSpace(config.HTML) != "" {
		var bodyBuffer bytes.Buffer
		innerWriter := multipart.NewWriter(&bodyBuffer)
		for _, part := range []struct {
			ContentType string
			Value       string
		}{
			{ContentType: "text/plain; charset=UTF-8", Value: config.Body},
			{ContentType: "text/html; charset=UTF-8", Value: config.HTML},
		} {
			partHeader := textproto.MIMEHeader{}
			partHeader.Set("Content-Type", part.ContentType)
			partHeader.Set("Content-Transfer-Encoding", "quoted-printable")
			partWriter, err := innerWriter.CreatePart(partHeader)
			if err != nil {
				return nil, "", "", err
			}
			if err := writeQuotedPrintable(partWriter, part.Value); err != nil {
				return nil, "", "", err
			}
		}
		if err := innerWriter.Close(); err != nil {
			return nil, "", "", err
		}
		bodyHeader.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%q", innerWriter.Boundary()))
		bodyPartWriter, err := writer.CreatePart(bodyHeader)
		if err != nil {
			return nil, "", "", err
		}
		if _, err := bodyPartWriter.Write(bodyBuffer.Bytes()); err != nil {
			return nil, "", "", err
		}
	} else {
		bodyValue := config.Body
		if strings.TrimSpace(config.HTML) != "" {
			bodyContentType = "text/html; charset=UTF-8"
			bodyValue = config.HTML
		}
		bodyHeader.Set("Content-Type", bodyContentType)
		bodyHeader.Set("Content-Transfer-Encoding", bodyTransferEncoding)
		bodyPartWriter, err := writer.CreatePart(bodyHeader)
		if err != nil {
			return nil, "", "", err
		}
		if err := writeQuotedPrintable(bodyPartWriter, bodyValue); err != nil {
			return nil, "", "", err
		}
	}

	for _, attachment := range config.Attachments {
		attachmentBody, err := os.ReadFile(attachment.Path)
		if err != nil {
			return nil, "", "", fmt.Errorf("send_email attachment %q: %w", attachment.Path, err)
		}
		header := textproto.MIMEHeader{}
		contentType := attachment.ContentType
		if contentType == "" {
			contentType = detectAttachmentContentType(attachment.Name)
		}
		header.Set("Content-Type", fmt.Sprintf("%s; name=%q", contentType, attachment.Name))
		header.Set("Content-Transfer-Encoding", "base64")
		header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", attachment.Name))
		partWriter, err := writer.CreatePart(header)
		if err != nil {
			return nil, "", "", err
		}
		if err := writeBase64MIME(partWriter, attachmentBody); err != nil {
			return nil, "", "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", "", err
	}
	return buffer.Bytes(), fmt.Sprintf("multipart/mixed; boundary=%q", writer.Boundary()), "", nil
}

func writeQuotedPrintable(writer io.Writer, value string) error {
	qp := quotedprintable.NewWriter(writer)
	if _, err := qp.Write([]byte(value)); err != nil {
		_ = qp.Close()
		return err
	}
	return qp.Close()
}

func writeBase64MIME(writer io.Writer, value []byte) error {
	encoded := base64.StdEncoding.EncodeToString(value)
	for len(encoded) > 76 {
		if _, err := io.WriteString(writer, encoded[:76]+"\r\n"); err != nil {
			return err
		}
		encoded = encoded[76:]
	}
	if encoded != "" {
		if _, err := io.WriteString(writer, encoded+"\r\n"); err != nil {
			return err
		}
	}
	return nil
}

func resolveEmailConnectionConfig(connection string) (emailConnectionConfig, error) {
	name := strings.TrimSpace(connection)
	if name == "" {
		name = emailDefaultConnection
	}

	config := defaultEmailConnectionConfig(name)
	config.Host = lookupEmailConfigValue(name, "HOST")
	if portText := lookupEmailConfigValue(name, "PORT"); portText != "" {
		port, err := strconv.Atoi(strings.TrimSpace(portText))
		if err != nil {
			return emailConnectionConfig{}, fmt.Errorf("email connection %q env PORT must be an integer: %w", name, err)
		}
		config.Port = port
	}
	config.Username = lookupEmailConfigValue(name, "USERNAME")
	config.Password = lookupEmailConfigValue(name, "PASSWORD")
	config.From = lookupEmailConfigValue(name, "FROM")
	config.HelloName = lookupEmailConfigValue(name, "HELO")
	if modeText := lookupEmailConfigValue(name, "TLS_MODE"); modeText != "" {
		mode, err := parseEmailTLSMode(modeText)
		if err != nil {
			return emailConnectionConfig{}, fmt.Errorf("email connection %q env TLS_MODE %w", name, err)
		}
		config.TLSMode = mode
	}
	if insecureText := lookupEmailConfigValue(name, "INSECURE_SKIP_VERIFY"); insecureText != "" {
		insecure, err := parseEmailEnvBool(insecureText)
		if err != nil {
			return emailConnectionConfig{}, fmt.Errorf("email connection %q env INSECURE_SKIP_VERIFY %w", name, err)
		}
		config.InsecureSkipVerify = insecure
	}
	if timeoutText := lookupEmailConfigValue(name, "TIMEOUT_MS"); timeoutText != "" {
		timeoutMS, err := strconv.Atoi(strings.TrimSpace(timeoutText))
		if err != nil {
			return emailConnectionConfig{}, fmt.Errorf("email connection %q env TIMEOUT_MS must be an integer: %w", name, err)
		}
		if timeoutMS > 0 {
			config.TimeoutMS = timeoutMS
		}
	}
	config, err := finalizeEmailConnectionConfig(config, true)
	if err != nil {
		return emailConnectionConfig{}, fmt.Errorf("email connection %q %w", name, err)
	}
	return config, nil
}

func applyLuaEmailRuntimePolicy(L *lua.LState, config flowEmailConfig) (flowEmailConfig, error) {
	ctx := flowContextFromState(L)
	if ctx == nil || ctx.Security == nil {
		return config, nil
	}
	return applyFlowEmailRuntimePolicy(ctx, config)
}

func applyFlowEmailRuntimePolicy(ctx *FlowContext, config flowEmailConfig) (flowEmailConfig, error) {
	if ctx == nil || ctx.Security == nil {
		return config, nil
	}
	if !ctx.Security.AllowEmail {
		return flowEmailConfig{}, fmt.Errorf("send_email is disabled by security policy; set allow_email=true only for trusted flows")
	}
	if len(config.Attachments) == 0 {
		return config, nil
	}
	if !ctx.Security.AllowFileAccess {
		return flowEmailConfig{}, fmt.Errorf("send_email attachments are disabled by security policy; set allow_file_access=true only for trusted flows")
	}
	return rewriteEmailRuntimePaths(config, *ctx.Security)
}

func defaultEmailConnectionConfig(name string) emailConnectionConfig {
	name = strings.TrimSpace(name)
	if name == "" {
		name = emailDefaultConnection
	}
	return emailConnectionConfig{
		Name:      name,
		TLSMode:   emailTLSModeStartTLS,
		TimeoutMS: 10000,
	}
}

func finalizeEmailConnectionConfig(config emailConnectionConfig, fromEnv bool) (emailConnectionConfig, error) {
	if strings.TrimSpace(config.Host) == "" {
		if fromEnv {
			return emailConnectionConfig{}, fmt.Errorf("is not configured; set %s", strings.Join(emailConfigEnvHints(config.Name), " or "))
		}
		return emailConnectionConfig{}, fmt.Errorf("requires smtp.host or a configured connection")
	}
	if config.Port < 0 {
		return emailConnectionConfig{}, fmt.Errorf("port must be at least 1")
	}
	if config.Port == 0 {
		switch config.TLSMode {
		case emailTLSModeTLS:
			config.Port = 465
		case emailTLSModeNone:
			config.Port = 25
		default:
			config.Port = 587
		}
	}
	return config, nil
}

func mergeEmailConnectionConfig(base emailConnectionConfig, override emailConnectionConfig) emailConnectionConfig {
	merged := base
	if strings.TrimSpace(override.Name) != "" {
		merged.Name = override.Name
	}
	if strings.TrimSpace(override.Host) != "" {
		merged.Host = override.Host
	}
	if override.Port != 0 {
		merged.Port = override.Port
	}
	if strings.TrimSpace(override.Username) != "" {
		merged.Username = override.Username
	}
	if strings.TrimSpace(override.Password) != "" {
		merged.Password = override.Password
	}
	if strings.TrimSpace(override.From) != "" {
		merged.From = override.From
	}
	if override.TLSMode != "" {
		merged.TLSMode = override.TLSMode
	}
	if strings.TrimSpace(override.HelloName) != "" {
		merged.HelloName = override.HelloName
	}
	if override.InsecureSkipVerify {
		merged.InsecureSkipVerify = true
	}
	if override.TimeoutMS != 0 {
		merged.TimeoutMS = override.TimeoutMS
	}
	return merged
}

func normalizeInlineEmailConnectionConfig(value any) (emailConnectionConfig, error) {
	objectValue, err := objectMapValue(value, "smtp")
	if err != nil {
		return emailConnectionConfig{}, err
	}

	config := emailConnectionConfig{}
	for key, rawValue := range objectValue {
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "host":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.host must be a string")
			}
			config.Host = strings.TrimSpace(text)
		case "port":
			port, err := intParam(rawValue)
			if err != nil {
				return emailConnectionConfig{}, fmt.Errorf("smtp.port %w", err)
			}
			if port < 1 {
				return emailConnectionConfig{}, fmt.Errorf("smtp.port must be at least 1")
			}
			config.Port = port
		case "username":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.username must be a string")
			}
			config.Username = strings.TrimSpace(text)
		case "password":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.password must be a string")
			}
			config.Password = text
		case "from":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.from must be a string")
			}
			if strings.TrimSpace(text) != "" {
				if _, err := parseSingleEmailAddress(text, "smtp.from"); err != nil {
					return emailConnectionConfig{}, err
				}
			}
			config.From = strings.TrimSpace(text)
		case "tls_mode":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.tls_mode must be a string")
			}
			mode, err := parseEmailTLSMode(text)
			if err != nil {
				return emailConnectionConfig{}, fmt.Errorf("smtp.tls_mode %w", err)
			}
			config.TLSMode = mode
		case "helo":
			text, ok := rawValue.(string)
			if !ok {
				return emailConnectionConfig{}, fmt.Errorf("smtp.helo must be a string")
			}
			config.HelloName = strings.TrimSpace(text)
		case "insecure_skip_verify":
			flag, err := boolParam(rawValue)
			if err != nil {
				return emailConnectionConfig{}, fmt.Errorf("smtp.insecure_skip_verify %w", err)
			}
			config.InsecureSkipVerify = flag
		case "timeout_ms":
			timeoutMS, err := intParam(rawValue)
			if err != nil {
				return emailConnectionConfig{}, fmt.Errorf("smtp.timeout_ms %w", err)
			}
			if timeoutMS < 1 {
				return emailConnectionConfig{}, fmt.Errorf("smtp.timeout_ms must be at least 1")
			}
			config.TimeoutMS = timeoutMS
		default:
			return emailConnectionConfig{}, fmt.Errorf("smtp does not accept field %q", key)
		}
	}
	return config, nil
}

func normalizeEmailAttachments(value any) ([]emailAttachmentConfig, error) {
	switch typed := value.(type) {
	case string:
		attachment, err := normalizeEmailAttachment(typed)
		if err != nil {
			return nil, err
		}
		return []emailAttachmentConfig{attachment}, nil
	case map[string]any:
		attachment, err := normalizeEmailAttachment(typed)
		if err != nil {
			return nil, err
		}
		return []emailAttachmentConfig{attachment}, nil
	case []string:
		attachments := make([]emailAttachmentConfig, 0, len(typed))
		for _, item := range typed {
			attachment, err := normalizeEmailAttachment(item)
			if err != nil {
				return nil, err
			}
			attachments = append(attachments, attachment)
		}
		return attachments, nil
	case []any:
		attachments := make([]emailAttachmentConfig, 0, len(typed))
		for _, item := range typed {
			attachment, err := normalizeEmailAttachment(item)
			if err != nil {
				return nil, err
			}
			attachments = append(attachments, attachment)
		}
		return attachments, nil
	default:
		return nil, fmt.Errorf("attachments must be a string, object, or list")
	}
}

func normalizeEmailAttachment(value any) (emailAttachmentConfig, error) {
	switch typed := value.(type) {
	case string:
		path := strings.TrimSpace(typed)
		if path == "" {
			return emailAttachmentConfig{}, fmt.Errorf("attachment path cannot be blank")
		}
		return emailAttachmentConfig{
			Path: path,
			Name: filepath.Base(path),
		}, nil
	case map[string]any:
		pathValue, ok := typed["path"]
		if !ok {
			return emailAttachmentConfig{}, fmt.Errorf("attachment object requires path")
		}
		pathText, ok := pathValue.(string)
		if !ok {
			return emailAttachmentConfig{}, fmt.Errorf("attachment path must be a string")
		}
		pathText = strings.TrimSpace(pathText)
		if pathText == "" {
			return emailAttachmentConfig{}, fmt.Errorf("attachment path cannot be blank")
		}
		attachment := emailAttachmentConfig{
			Path: pathText,
			Name: filepath.Base(pathText),
		}
		if nameValue, ok := typed["name"]; ok {
			nameText, ok := nameValue.(string)
			if !ok {
				return emailAttachmentConfig{}, fmt.Errorf("attachment name must be a string")
			}
			nameText = strings.TrimSpace(nameText)
			if nameText != "" {
				attachment.Name = nameText
			}
		}
		if contentTypeValue, ok := typed["content_type"]; ok {
			contentTypeText, ok := contentTypeValue.(string)
			if !ok {
				return emailAttachmentConfig{}, fmt.Errorf("attachment content_type must be a string")
			}
			attachment.ContentType = strings.TrimSpace(contentTypeText)
		}
		return attachment, nil
	default:
		return emailAttachmentConfig{}, fmt.Errorf("attachment must be a string or object")
	}
}

func detectAttachmentContentType(name string) string {
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(name)))
	if strings.TrimSpace(contentType) == "" {
		return "application/octet-stream"
	}
	return contentType
}

func rewriteEmailRuntimePaths(config flowEmailConfig, policy FlowSecurityPolicy) (flowEmailConfig, error) {
	if len(config.Attachments) == 0 {
		return config, nil
	}
	rewritten := make([]emailAttachmentConfig, 0, len(config.Attachments))
	for _, attachment := range config.Attachments {
		resolved, err := resolveRuntimeFilePath(attachment.Path, flowFileInputPath, policy)
		if err != nil {
			return flowEmailConfig{}, fmt.Errorf("send_email attachment %q %w", attachment.Path, err)
		}
		attachment.Path = resolved
		rewritten = append(rewritten, attachment)
	}
	config.Attachments = rewritten
	return config, nil
}

func lookupEmailConfigValue(connection string, suffix string) string {
	for _, key := range emailConfigEnvKeys(connection, suffix) {
		if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func emailConfigEnvHints(connection string) []string {
	return emailConfigEnvKeys(connection, "HOST")
}

func emailConfigEnvKeys(connection string, suffix string) []string {
	normalized := normalizeEmailConnectionName(connection)
	if normalized == normalizeEmailConnectionName(emailDefaultConnection) {
		return []string{
			"TSPLAY_EMAIL_DEFAULT_" + suffix,
			"TSPLAY_EMAIL_" + suffix,
		}
	}
	return []string{"TSPLAY_EMAIL_" + normalized + "_" + suffix}
}

func normalizeEmailConnectionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = emailDefaultConnection
	}
	builder := strings.Builder{}
	for _, r := range strings.ToUpper(name) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		} else {
			builder.WriteByte('_')
		}
	}
	return strings.Trim(builder.String(), "_")
}

func parseEmailTLSMode(value string) (emailTLSMode, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "starttls":
		return emailTLSModeStartTLS, nil
	case "tls", "ssl", "smtps":
		return emailTLSModeTLS, nil
	case "none", "plain":
		return emailTLSModeNone, nil
	default:
		return "", fmt.Errorf("must be one of none, starttls, or tls")
	}
}

func parseEmailEnvBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("must be a boolean-like string")
	}
}

func emailAddressListValue(value any, field string) ([]mail.Address, error) {
	var inputs []string
	switch typed := value.(type) {
	case string:
		inputs = []string{typed}
	case []string:
		inputs = append(inputs, typed...)
	case []any:
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be a string or list of strings", field)
			}
			inputs = append(inputs, text)
		}
	default:
		return nil, fmt.Errorf("%s must be a string or list of strings", field)
	}

	addresses := make([]mail.Address, 0, len(inputs))
	for _, input := range inputs {
		if strings.TrimSpace(input) == "" {
			return nil, fmt.Errorf("%s cannot be blank", field)
		}
		parsed, err := mail.ParseAddressList(input)
		if err != nil {
			return nil, fmt.Errorf("%s contains an invalid email address: %w", field, err)
		}
		for _, address := range parsed {
			addresses = append(addresses, *address)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("%s must contain at least one email address", field)
	}
	return addresses, nil
}

func parseSingleEmailAddress(value string, field string) (mail.Address, error) {
	if err := validateEmailHeaderText(field, value); err != nil {
		return mail.Address{}, err
	}
	address, err := mail.ParseAddress(value)
	if err != nil {
		return mail.Address{}, fmt.Errorf("send_email %s contains an invalid email address: %w", field, err)
	}
	return *address, nil
}

func normalizeEmailHeaders(value any) (map[string]string, error) {
	headers, err := stringMapValue(value, "headers")
	if err != nil {
		return nil, err
	}
	normalized := map[string]string{}
	for key, headerValue := range headers {
		canonical := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(key))
		if canonical == "" {
			return nil, fmt.Errorf("headers contains a blank header name")
		}
		if isReservedEmailHeader(canonical) {
			return nil, fmt.Errorf("headers cannot override reserved header %q", canonical)
		}
		if err := validateEmailHeaderText("headers."+canonical, headerValue); err != nil {
			return nil, err
		}
		normalized[canonical] = headerValue
	}
	return normalized, nil
}

func isReservedEmailHeader(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "bcc", "cc", "content-transfer-encoding", "content-type", "date", "from", "message-id", "mime-version", "reply-to", "subject", "to":
		return true
	default:
		return false
	}
}

func validateEmailHeaderText(field string, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("send_email %s cannot contain newlines", field)
	}
	return nil
}

func encodeEmailHeaderText(value string) string {
	if value == "" {
		return value
	}
	if isASCIIEmailHeader(value) {
		return value
	}
	return mime.QEncoding.Encode("utf-8", value)
}

func isASCIIEmailHeader(value string) bool {
	for _, r := range value {
		if r > 127 {
			return false
		}
	}
	return true
}

func buildEmailMessageID(from string) string {
	domain := "localhost"
	if at := strings.LastIndex(from, "@"); at >= 0 && at < len(from)-1 {
		domain = from[at+1:]
	}
	return fmt.Sprintf("<tsplay.%d@%s>", time.Now().UnixNano(), domain)
}

func emailAddressStrings(addresses []mail.Address) []string {
	values := make([]string, 0, len(addresses))
	for _, address := range addresses {
		values = append(values, address.String())
	}
	return values
}
