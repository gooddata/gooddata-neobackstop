package screenshotter

import (
	"fmt"
	"os"

	"github.com/gooddata/gooddata-neobackstop/viewport"
	"github.com/playwright-community/playwright-go"
)

func takeStableScreenshotBytes(page playwright.Page, originalViewport viewport.Viewport) ([]byte, error) {
	scrollHeightValue, err := page.Evaluate("() => document.documentElement.scrollHeight")
	if err != nil {
		return nil, err
	}

	var totalHeight int

	switch v := scrollHeightValue.(type) {
	case float64:
		totalHeight = int(v)
	case int:
		totalHeight = v
	case int32:
		totalHeight = int(v)
	case int64:
		totalHeight = int(v)
	default:
		panic(fmt.Sprintf("unexpected scrollHeight type: %T", v))
	}

	// Set the viewport to the full page height
	err = page.SetViewportSize(originalViewport.Width, totalHeight)
	if err != nil {
		return nil, err
	}

	// Now take a *regular* screenshot — NOT FullPage:true
	var screenshotBytes []byte
	screenshotBytes, err = page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	// Restore original viewport
	err = page.SetViewportSize(originalViewport.Width, originalViewport.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to restore viewport: %w", err)
	}

	return screenshotBytes, nil
}

func saveScreenshotBytes(bytes []byte, filePath string) error {
	return os.WriteFile(filePath, bytes, 0644)
}
