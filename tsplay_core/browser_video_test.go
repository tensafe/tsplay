package tsplay_core

import (
	"path/filepath"
	"testing"
)

func TestPrepareBrowserVideoRecording(t *testing.T) {
	output := filepath.Join(t.TempDir(), "videos", "lesson-10.webm")
	recording, err := PrepareBrowserVideoRecording(output, 1280, 720)
	if err != nil {
		t.Fatalf("PrepareBrowserVideoRecording: %v", err)
	}
	if recording == nil {
		t.Fatalf("expected recording")
	}
	if recording.OutputPath != output {
		t.Fatalf("OutputPath = %q, want %q", recording.OutputPath, output)
	}
	if recording.Dir != filepath.Dir(output) {
		t.Fatalf("Dir = %q, want %q", recording.Dir, filepath.Dir(output))
	}
	recordVideo := recording.NewContextRecordVideo()
	if recordVideo == nil {
		t.Fatalf("expected RecordVideo")
	}
	if recordVideo.Dir == nil || *recordVideo.Dir != filepath.Dir(output) {
		t.Fatalf("Dir = %#v, want %q", recordVideo.Dir, filepath.Dir(output))
	}
	if recordVideo.Size == nil ||
		recordVideo.Size.Width == nil ||
		*recordVideo.Size.Width != 1280 ||
		recordVideo.Size.Height == nil ||
		*recordVideo.Size.Height != 720 {
		t.Fatalf("unexpected size: %#v", recordVideo.Size)
	}
}

func TestPrepareBrowserVideoRecordingBlank(t *testing.T) {
	recording, err := PrepareBrowserVideoRecording("", 0, 0)
	if err != nil {
		t.Fatalf("PrepareBrowserVideoRecording blank: %v", err)
	}
	if recording != nil {
		t.Fatalf("expected nil recording")
	}
}
