package tsplay_core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type BrowserVideoRecording struct {
	OutputPath string
	Dir        string
	Width      int
	Height     int
}

func PrepareBrowserVideoRecording(outputPath string, width int, height int) (*BrowserVideoRecording, error) {
	trimmed := strings.TrimSpace(outputPath)
	if trimmed == "" {
		return nil, nil
	}
	resolvedPath, err := filepath.Abs(trimmed)
	if err != nil {
		return nil, fmt.Errorf("resolve browser video output %q: %w", outputPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0755); err != nil {
		return nil, fmt.Errorf("create browser video output dir: %w", err)
	}
	return &BrowserVideoRecording{
		OutputPath: resolvedPath,
		Dir:        filepath.Dir(resolvedPath),
		Width:      width,
		Height:     height,
	}, nil
}

func (recording *BrowserVideoRecording) NewContextRecordVideo() *playwright.BrowserNewContextOptionsRecordVideo {
	if recording == nil {
		return nil
	}
	options := &playwright.BrowserNewContextOptionsRecordVideo{
		Dir: playwright.String(recording.Dir),
	}
	if recording.Width > 0 && recording.Height > 0 {
		options.Size = &playwright.BrowserNewContextOptionsRecordVideoSize{
			Width:  playwright.Int(recording.Width),
			Height: playwright.Int(recording.Height),
		}
	}
	return options
}

func (recording *BrowserVideoRecording) NewPageRecordVideo() *playwright.BrowserNewPageOptionsRecordVideo {
	if recording == nil {
		return nil
	}
	options := &playwright.BrowserNewPageOptionsRecordVideo{
		Dir: playwright.String(recording.Dir),
	}
	if recording.Width > 0 && recording.Height > 0 {
		options.Size = &playwright.BrowserNewPageOptionsRecordVideoSize{
			Width:  playwright.Int(recording.Width),
			Height: playwright.Int(recording.Height),
		}
	}
	return options
}

func (recording *BrowserVideoRecording) PersistentContextRecordVideo() *playwright.BrowserTypeLaunchPersistentContextOptionsRecordVideo {
	if recording == nil {
		return nil
	}
	options := &playwright.BrowserTypeLaunchPersistentContextOptionsRecordVideo{
		Dir: playwright.String(recording.Dir),
	}
	if recording.Width > 0 && recording.Height > 0 {
		options.Size = &playwright.BrowserTypeLaunchPersistentContextOptionsRecordVideoSize{
			Width:  playwright.Int(recording.Width),
			Height: playwright.Int(recording.Height),
		}
	}
	return options
}

func SaveBrowserVideo(page playwright.Page, outputPath string) (string, error) {
	if page == nil {
		return "", fmt.Errorf("page is nil")
	}
	video := page.Video()
	if video == nil {
		return "", fmt.Errorf("page video is unavailable")
	}
	if !page.IsClosed() {
		if err := page.Close(); err != nil {
			return "", fmt.Errorf("close page before saving video: %w", err)
		}
	}
	resolvedPath, err := filepath.Abs(strings.TrimSpace(outputPath))
	if err != nil {
		return "", fmt.Errorf("resolve browser video output %q: %w", outputPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0755); err != nil {
		return "", fmt.Errorf("create browser video output dir: %w", err)
	}
	if err := video.SaveAs(resolvedPath); err != nil {
		return "", fmt.Errorf("save page video to %q: %w", resolvedPath, err)
	}
	return resolvedPath, nil
}
