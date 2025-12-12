package screenshotter

import (
	"fmt"

	"github.com/gooddata/gooddata-neobackstop/viewport"
	"github.com/playwright-community/playwright-go"
)

// takeStableScreenshot - the original concept for this function was to take multiple screenshots and stitch them
// together, but that might not be necessary
// If filePath is provided, saves to disk and returns nil bytes. Otherwise returns bytes in memory.
func takeStableScreenshot(page playwright.Page, filePath *string, originalViewport viewport.Viewport) ([]byte, error) {
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
	if filePath == nil {
		screenshotBytes, err = page.Screenshot(playwright.PageScreenshotOptions{
			FullPage: playwright.Bool(false),
		})
	} else {
		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path:     filePath,
			FullPage: playwright.Bool(false),
		})
	}
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
