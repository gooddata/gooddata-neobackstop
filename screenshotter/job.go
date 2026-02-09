package screenshotter

import (
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/gooddata/gooddata-neobackstop/config"
	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/gooddata/gooddata-neobackstop/screenshotter/operations"
	"github.com/gooddata/gooddata-neobackstop/utils"
	"github.com/playwright-community/playwright-go"
)

func cleanText(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, " ", "_"), "/", "_"), ",", "_")
}

func loadScenarioAndCaptureScreenshot(logPrefix string, page playwright.Page, job internals.Scenario) ([]byte, *string) {
	strPtr := func(s string) *string { return &s }

	tTotal := time.Now()
	tPageLoad := time.Now()
	t1 := time.Now()
	fmt.Println(logPrefix, "screenshot: capturing")

	fmt.Println(logPrefix, "pageLoad: navigating to", job.Url)

	if _, err := page.Goto(job.Url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return nil, strPtr(fmt.Sprintf("could not goto: %v", err))
	}
	fmt.Println(logPrefix, "pageLoad: completed in", time.Since(tPageLoad).Milliseconds(), "ms")

	// readySelector
	sErr := operations.ReadySelector(logPrefix, page, job.ReadySelector)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// readySelectors
	sErr = operations.ReadySelectors(logPrefix, page, job.ReadySelectors)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// reloadAfterReady
	sErr = operations.ReloadAfterReady(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	t0 := time.Now()

	// delay
	if job.Delay != nil && job.Delay.PostReady > 0 {
		delayMs := job.Delay.PostReady.Milliseconds()
		fmt.Println(logPrefix, "postReadyDelay: starting sleep for", delayMs, "ms")
		time.Sleep(job.Delay.PostReady)
		fmt.Println(logPrefix, "postReadyDelay: ending sleep for", delayMs, "ms")
	}

	// keyPressSelector
	sErr = operations.KeyPressSelector(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// hoverSelector
	sErr = operations.HoverSelector(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// hoverSelectors
	sErr = operations.HoverSelectors(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// clickSelector
	sErr = operations.ClickSelector(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// clickSelectors
	sErr = operations.ClickSelectors(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	// scroll to selector
	sErr = operations.ScrollToSelector(logPrefix, page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		return nil, sErr
	}

	fmt.Println(logPrefix, "allOperations: completed in", time.Since(t0).Milliseconds(), "ms (incl. load+ready:", time.Since(tTotal).Milliseconds(), "ms)")

	if job.Delay != nil && job.Delay.PostOperation > 0 {
		delayMs := job.Delay.PostOperation.Milliseconds()
		fmt.Println(logPrefix, "postOperationDelay: starting sleep for", delayMs, "ms")
		time.Sleep(job.Delay.PostOperation)
		fmt.Println(logPrefix, "postOperationDelay: ending sleep for", delayMs, "ms")
	}

	_, err := page.Evaluate(`() => {
    // Disable window.onresize
    window.onresize = null;

    // Override addEventListener for resize
    const originalAddEventListener = window.addEventListener;
    window.addEventListener = function(type, listener, options) {
        if (type === 'resize') return; // ignore resize listeners
        return originalAddEventListener.call(this, type, listener, options);
    };
}`)
	if err != nil {
		return nil, strPtr(err.Error())
	}

	_, err = page.Evaluate(`() => {
    const fullHeight = document.documentElement.scrollHeight + 'px';
    document.documentElement.style.height = fullHeight;
    document.body.style.height = fullHeight;
}`)
	if err != nil {
		return nil, strPtr(err.Error())
	}

	_, err = page.Evaluate(`() => {
  return new Promise(resolve => {
    setTimeout(() => requestAnimationFrame(resolve), 0);
  });
}`)
	if err != nil {
		return nil, strPtr(err.Error())
	}

	screenshotBytes, err := takeStableScreenshotBytes(page, job.Viewport)
	if err != nil {
		return nil, strPtr(fmt.Sprintf("could not take screenshot: %v", err))
	}
	fmt.Println(logPrefix, "screenshot: captured in", time.Since(t1).Milliseconds(), "ms")

	// Move mouse outside viewport to clear any hover states
	if job.HoverSelector != nil || job.HoverSelectors != nil || job.ClickSelector != nil || job.ClickSelectors != nil {
		err = page.Mouse().Move(-1, -1)
		if err != nil {
			return nil, strPtr(fmt.Sprintf("could not move mouse to neutral position: %v", err))
		}

		_, err = page.Evaluate(`() => {
  document.querySelectorAll(":hover").forEach(el =>
    el.dispatchEvent(new MouseEvent("mouseout"))
  )
}`)
		if err != nil {
			return nil, strPtr(fmt.Sprintf("could not cleanup hovers: %v", err))
		}
	}

	return screenshotBytes, nil
}

func matchWithReference(screenshotBytes []byte, referenceImg image.Image, threshold *float64) (*bool, *float64, error) {
	testImg, err := utils.DecodeImageFromBytes(screenshotBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode screenshot: %w", err)
	}
	_, mismatch := utils.DiffImagesPink(referenceImg, testImg)

	matches := (threshold != nil && *threshold >= mismatch) ||
		(threshold == nil && mismatch == 0)

	return &matches, &mismatch, nil
}

func referenceExists(referencePath string) bool {
	_, err := os.Stat(referencePath)
	return err == nil
}

func Job(logPrefix string, saveDir string, viewportLabel string, page playwright.Page, job internals.Scenario, results chan Result, mode string, conf config.Config) {
	screenshotBytes, scenarioErr := loadScenarioAndCaptureScreenshot(logPrefix, page, job)
	if scenarioErr != nil {
		results <- buildResultFromScenario(job, nil, scenarioErr)
		return
	}

	safeCombinedName := job.Id + "_0_document_0_" + cleanText(viewportLabel)
	fileName := "storybook_" + string(job.Browser) + "_" + safeCombinedName + ".png"
	referencePath := conf.BitmapsReferencePath + "/" + fileName

	if mode == "test" && referenceExists(referencePath) {
		referenceImg, err := utils.LoadImage(referencePath)
		if err != nil {
			errStr := fmt.Sprintf("could not load reference image: %v", err)
			results <- buildResultFromScenario(job, &fileName, &errStr)
			return
		}

		preComputedMatch, preComputedMismatch, err := matchWithReference(screenshotBytes, referenceImg, job.MisMatchThreshold)
		if err != nil {
			errStr := fmt.Sprintf("could not compare images: %v", err)
			results <- buildResultFromScenario(job, &fileName, &errStr)
			return
		}

		matches := preComputedMatch != nil && *preComputedMatch
		fmt.Println(logPrefix, "screenshot matches:", matches)

		realRetries := 0
		if !matches && job.RetryCount > 0 {
			for i := 0; i < job.RetryCount; i++ {
				realRetries++
				fmt.Println(logPrefix, "screenshot: retrying", i+1, "of", job.RetryCount)
				retryScreenshotBytes, scenarioErr := loadScenarioAndCaptureScreenshot(logPrefix, page, job)
				if scenarioErr != nil {
					continue
				}
				retryPreComputedMatch, retryPreComputedMismatch, retryErr := matchWithReference(retryScreenshotBytes, referenceImg, job.MisMatchThreshold)
				if retryErr != nil {
					continue
				}

				screenshotBytes = retryScreenshotBytes
				preComputedMatch = retryPreComputedMatch
				preComputedMismatch = retryPreComputedMismatch

				matches = preComputedMatch != nil && *preComputedMatch
				if matches {
					break
				}
			}
		}

		if !matches || conf.HtmlReport.ShowSuccessfulTests {
			testPath := conf.BitmapsTestPath + "/" + fileName
			if err := saveScreenshotBytes(screenshotBytes, testPath); err != nil {
				errStr := fmt.Sprintf("could not save test screenshot: %v", err)
				results <- buildResultFromScenario(job, &fileName, &errStr)
				return
			}
		}

		result := buildResultFromScenario(job, &fileName, nil)
		result.PreComputedMatch = preComputedMatch
		result.PreComputedMismatchPercentage = preComputedMismatch
		result.RetriesUsed = &realRetries
		hasRef := true
		result.PreComputedHasReference = &hasRef
		results <- result

	} else {
		testPath := conf.BitmapsTestPath + "/" + fileName
		if err := saveScreenshotBytes(screenshotBytes, testPath); err != nil {
			errStr := fmt.Sprintf("could not save test screenshot: %v", err)
			results <- buildResultFromScenario(job, &fileName, &errStr)
			return
		}
		results <- buildResultFromScenario(job, &fileName, nil)
	}
}
