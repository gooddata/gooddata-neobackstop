package screenshotter

import (
	"fmt"
	"log"
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

func Job(logPrefix string, saveDir string, viewportLabel string, page playwright.Page, job internals.Scenario, results chan Result, mode string, conf config.Config) {
	if _, err := page.Goto(job.Url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Panicf("could not goto: %v", err)
	}

	t0 := time.Now()

	// readySelector
	sErr := operations.ReadySelector(page, job.ReadySelector)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// readySelectors
	sErr = operations.ReadySelectors(page, job.ReadySelectors)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// reloadAfterReady
	sErr = operations.ReloadAfterReady(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// delay
	if job.Delay != nil {
		fmt.Println(logPrefix, "sleep start")
		time.Sleep(job.Delay.PostReady)
		fmt.Println(logPrefix, "sleep end")
	}

	// keyPressSelector
	sErr = operations.KeyPressSelector(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// hoverSelector
	sErr = operations.HoverSelector(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// hoverSelectors
	sErr = operations.HoverSelectors(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// clickSelector
	sErr = operations.ClickSelector(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// clickSelectors
	sErr = operations.ClickSelectors(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	// scroll to selector
	sErr = operations.ScrollToSelector(page, job)
	if sErr != nil {
		fmt.Println(logPrefix, *sErr+", exiting quietly without a screenshot")
		results <- buildResultFromScenario(job, nil, sErr)
		return
	}

	fmt.Println(logPrefix, "operations completed in", time.Since(t0))

	if job.Delay != nil {
		time.Sleep(job.Delay.PostOperation)
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
		panic(err.Error())
	}

	_, err = page.Evaluate(`() => {
    const fullHeight = document.documentElement.scrollHeight + 'px';
    document.documentElement.style.height = fullHeight;
    document.body.style.height = fullHeight;
}`)
	if err != nil {
		panic(err.Error())
	}

	_, err = page.Evaluate(`() => {
  return new Promise(resolve => {
    setTimeout(() => requestAnimationFrame(resolve), 0);
  });
}`)
	if err != nil {
		panic(err.Error())
	}

	safeCombinedName := job.Id + "_0_document_0_" + cleanText(viewportLabel)
	fileName := "storybook_" + string(job.Browser) + "_" + safeCombinedName + ".png"

	t1 := time.Now()

	if mode == "test" {
		// in test mode: capture to memory and compare immediately
		screenshotBytes, err := takeStableScreenshot(page, nil, job.Viewport)
		if err != nil {
			log.Panicf("could not take screenshot: %v", err)
		}
		fmt.Println(logPrefix, "capture took", time.Since(t1))

		var preComputedMatch *bool
		var preComputedMismatch *float64
		var preComputedHasRef *bool
		referencePath := conf.BitmapsReferencePath + "/" + fileName

		if _, err := os.Stat(referencePath); err == nil {
			// reference exists, compare immediately
			preComputedHasRef = new(bool)
			*preComputedHasRef = true

			referenceImg, err := utils.LoadImage(referencePath)
			if err != nil {
				log.Panicf("could not load reference image: %v", err)
			}

			testImg, err := utils.DecodeImageFromBytes(screenshotBytes)
			if err != nil {
				log.Panicf("could not decode screenshot: %v", err)
			}

			_, mismatch := utils.DiffImagesPink(referenceImg, testImg)
			preComputedMismatch = &mismatch

			matches := (job.MisMatchThreshold != nil && *job.MisMatchThreshold >= mismatch) ||
				(job.MisMatchThreshold == nil && mismatch == 0)
			preComputedMatch = &matches

			if matches {
				if conf.HtmlReport.ShowSuccessfulTests {
					testPath := conf.BitmapsTestPath + "/" + fileName
					err = os.WriteFile(testPath, screenshotBytes, 0644)
					if err != nil {
						log.Panicf("could not save test screenshot: %v", err)
					}
				}
			} else {
				// mismatch detected - save to disk immediately and free memory
				testPath := conf.BitmapsTestPath + "/" + fileName
				err = os.WriteFile(testPath, screenshotBytes, 0644)
				if err != nil {
					log.Panicf("could not save test screenshot: %v", err)
				}
			}
		} else if os.IsNotExist(err) {
			// reference doesn't exist - save to disk immediately and free memory
			preComputedHasRef = new(bool)
			*preComputedHasRef = false
			testPath := conf.BitmapsTestPath + "/" + fileName
			err = os.WriteFile(testPath, screenshotBytes, 0644)
			if err != nil {
				log.Panicf("could not save test screenshot: %v", err)
			}
		} else {
			log.Panicf("could not check reference image: %v", err)
		}

		result := buildResultFromScenario(job, &fileName, nil)
		result.PreComputedMatch = preComputedMatch
		result.PreComputedMismatchPercentage = preComputedMismatch
		result.PreComputedHasReference = preComputedHasRef
		results <- result
	} else {
		// In approve mode: save to disk (existing behavior)
		results <- buildResultFromScenario(job, &fileName, nil)
		filePath := saveDir + "/" + fileName
		fmt.Println(logPrefix, "saving", filePath)
		_, err = takeStableScreenshot(page, &filePath, job.Viewport)
		if err != nil {
			log.Panicf("could not take screenshot: %v", err)
		}
		fmt.Println(logPrefix, "saving took", time.Since(t1))
	}

	// Move mouse outside viewport to clear any hover states
	if job.HoverSelector != nil || job.HoverSelectors != nil || job.ClickSelector != nil || job.ClickSelectors != nil {
		err = page.Mouse().Move(-1, -1)
		if err != nil {
			log.Panicf("could not move mouse to neutral position: %v", err)
		}

		_, err = page.Evaluate(`() => {
  document.querySelectorAll(":hover").forEach(el =>
    el.dispatchEvent(new MouseEvent("mouseout"))
  )
}`)
		if err != nil {
			log.Panicf("could not cleanup hovers: %v", err)
		}
	}
}
