package main

import (
	"strings"
	"testing"
)

func TestValidateBrowserCDPFlagOptions(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		endpointSet    bool
		port           int
		portSet        bool
		launch         bool
		executable     string
		executableSet  bool
		userDataDir    string
		userDataDirSet bool
		videoOutput    string
		wantErr        string
		shouldPass     bool
	}{
		{
			name:       "absent_port_zero_is_default",
			shouldPass: true,
		},
		{
			name:    "explicit_zero_rejected",
			portSet: true,
			wantErr: "between 1 and 65535",
		},
		{
			name:    "negative_port_rejected",
			port:    -1,
			portSet: true,
			wantErr: "between 1 and 65535",
		},
		{
			name:    "overflow_port_rejected",
			port:    65536,
			portSet: true,
			wantErr: "between 1 and 65535",
		},
		{
			name:       "max_port_allowed",
			port:       65535,
			portSet:    true,
			shouldPass: true,
		},
		{
			name:     "endpoint_and_port_conflict",
			endpoint: "http://127.0.0.1:9222",
			port:     9222,
			portSet:  true,
			wantErr:  "cannot be used together",
		},
		{
			name:        "explicit_blank_endpoint_rejected",
			endpointSet: true,
			wantErr:     "browser-cdp-endpoint cannot be blank",
		},
		{
			name:        "explicit_whitespace_endpoint_rejected",
			endpoint:    "  \t",
			endpointSet: true,
			wantErr:     "browser-cdp-endpoint cannot be blank",
		},
		{
			name:     "invalid_endpoint_rejected",
			endpoint: "127.0.0.1:70000/json/version",
			wantErr:  "invalid port",
		},
		{
			name:     "remote_launch_endpoint_rejected",
			endpoint: "http://192.0.2.1:9222",
			launch:   true,
			wantErr:  "only start or reuse a local browser",
		},
		{
			name:     "local_launch_endpoint_without_port_rejected",
			endpoint: "http://127.0.0.1",
			launch:   true,
			wantErr:  "explicit port",
		},
		{
			name:       "remote_attach_endpoint_allowed",
			endpoint:   "http://192.0.2.1:9222",
			shouldPass: true,
		},
		{
			name:       "executable_implies_launch_remote_endpoint_rejected",
			endpoint:   "http://192.0.2.1:9222",
			executable: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			wantErr:    "only start or reuse a local browser",
		},
		{
			name:          "explicit_blank_executable_rejected",
			executableSet: true,
			wantErr:       "browser-cdp-executable cannot be blank",
		},
		{
			name:          "explicit_whitespace_executable_rejected",
			executable:    "  \n",
			executableSet: true,
			wantErr:       "browser-cdp-executable cannot be blank",
		},
		{
			name:        "user_data_dir_implies_launch_remote_endpoint_rejected",
			endpoint:    "http://192.0.2.1:9222",
			userDataDir: "profiles/cdp",
			wantErr:     "only start or reuse a local browser",
		},
		{
			name:           "explicit_blank_user_data_dir_rejected",
			userDataDirSet: true,
			wantErr:        "browser-cdp-user-data-dir cannot be blank",
		},
		{
			name:           "explicit_whitespace_user_data_dir_rejected",
			userDataDir:    " \t ",
			userDataDirSet: true,
			wantErr:        "browser-cdp-user-data-dir cannot be blank",
		},
		{
			name:        "video_without_cdp_allowed",
			videoOutput: "artifacts/browser.webm",
			shouldPass:  true,
		},
		{
			name:        "video_with_cdp_port_rejected",
			port:        9222,
			portSet:     true,
			videoOutput: "artifacts/browser.webm",
			wantErr:     "browser-video-output",
		},
		{
			name:        "video_with_cdp_launch_rejected",
			launch:      true,
			videoOutput: "artifacts/browser.webm",
			wantErr:     "browser-video-output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			launch := tt.launch || strings.TrimSpace(tt.executable) != "" || strings.TrimSpace(tt.userDataDir) != ""
			endpointSet := tt.endpointSet || strings.TrimSpace(tt.endpoint) != ""
			executableSet := tt.executableSet || strings.TrimSpace(tt.executable) != ""
			userDataDirSet := tt.userDataDirSet || strings.TrimSpace(tt.userDataDir) != ""
			err := validateBrowserCDPFlagOptions(tt.endpoint, endpointSet, tt.port, tt.portSet, launch, tt.executable, executableSet, tt.userDataDir, userDataDirSet, tt.videoOutput)
			if tt.shouldPass {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}
