package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseAVFoundationDevices(t *testing.T) {
	output := `[AVFoundation indev @ 0x123] AVFoundation video devices:
[AVFoundation indev @ 0x123] [0] Capture screen 0
[AVFoundation indev @ 0x123] [1] FaceTime HD Camera
[AVFoundation indev @ 0x123] AVFoundation audio devices:
[AVFoundation indev @ 0x123] [0] MacBook Pro Microphone
[in#0 @ 0x123] Error opening input: Input/output error`

	video, audio := parseAVFoundationDevices(output)

	wantVideo := []string{"Capture screen 0", "FaceTime HD Camera"}
	wantAudio := []string{"MacBook Pro Microphone"}
	if !reflect.DeepEqual(video, wantVideo) {
		t.Fatalf("video = %#v, want %#v", video, wantVideo)
	}
	if !reflect.DeepEqual(audio, wantAudio) {
		t.Fatalf("audio = %#v, want %#v", audio, wantAudio)
	}
}

func TestParseAVFoundationDevicesWithANSI(t *testing.T) {
	output := "\x1b[0;35m[AVFoundation indev @ 0x123] \x1b[0mAVFoundation video devices:\n" +
		"\x1b[0;35m[AVFoundation indev @ 0x123] \x1b[0m[0] Capture screen 0\n" +
		"\x1b[0;35m[AVFoundation indev @ 0x123] \x1b[0mAVFoundation audio devices:\n" +
		"\x1b[0;35m[AVFoundation indev @ 0x123] \x1b[0m[0] BlackHole 2ch"

	video, audio := parseAVFoundationDevices(output)
	if !reflect.DeepEqual(video, []string{"Capture screen 0"}) {
		t.Fatalf("video = %#v", video)
	}
	if !reflect.DeepEqual(audio, []string{"BlackHole 2ch"}) {
		t.Fatalf("audio = %#v", audio)
	}
}

func TestBuildScreenRecordFFmpegArgs(t *testing.T) {
	args := buildScreenRecordFFmpegArgs(screenRecordOptions{
		InputSpec:     "Capture screen 0:none",
		FrameRate:     30,
		VideoSize:     "1728x1117",
		CaptureCursor: true,
		MaxDuration:   5 * time.Second,
		CRF:           20,
		Preset:        "fast",
	}, "/tmp/tutorial.mp4")

	want := []string{
		"-y",
		"-hide_banner",
		"-f", "avfoundation",
		"-capture_cursor", "1",
		"-framerate", "30",
		"-video_size", "1728x1117",
		"-t", "5",
		"-i", "Capture screen 0:none",
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "20",
		"-movflags", "+faststart",
		"/tmp/tutorial.mp4",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestFormatFFmpegSeconds(t *testing.T) {
	if got := formatFFmpegSeconds(1500 * time.Millisecond); got != "1.5" {
		t.Fatalf("formatFFmpegSeconds(1.5s) = %q", got)
	}
	if got := formatFFmpegSeconds(5 * time.Second); got != "5" {
		t.Fatalf("formatFFmpegSeconds(5s) = %q", got)
	}
}
