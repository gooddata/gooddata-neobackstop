package operations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func HoverSelectors(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	if scenario.HoverSelectors != nil {
		fmt.Println(logPrefix, "hoverSelectors: processing", len(scenario.HoverSelectors), "hovers")

		for i, hs := range scenario.HoverSelectors {
			idx := strconv.Itoa(i + 1)

			if hs.WaitBefore != nil {
				delayMs := hs.WaitBefore.Milliseconds()
				fmt.Println(logPrefix, "hoverSelectors["+idx+"]: waitBefore starting sleep for", delayMs, "ms")
				time.Sleep(*hs.WaitBefore)
				fmt.Println(logPrefix, "hoverSelectors["+idx+"]: waitBefore ending sleep for", delayMs, "ms")
			}

			t0 := time.Now()
			fmt.Println(logPrefix, "hoverSelectors["+idx+"]: waiting for", hs.Selector)

			_, err := page.WaitForSelector(hs.Selector, playwright.PageWaitForSelectorOptions{
				State: playwright.WaitForSelectorStateAttached,
				// use default 30s timeout
			})
			if err != nil {
				e := "HoverSelector " + hs.Selector + " didn't appear"
				return &e
			}
			waitTime := time.Since(t0).Milliseconds()

			tHover := time.Now()
			err = page.Hover(hs.Selector)
			if err != nil {
				e := "HoverSelector " + hs.Selector + " couldn't be hovered"
				return &e
			}
			hoverTime := time.Since(tHover).Milliseconds()
			fmt.Println(logPrefix, "hoverSelectors["+idx+"]:", hs.Selector, "completed in", waitTime+hoverTime, "ms (wait:", waitTime, "ms, hover:", hoverTime, "ms)")

			sErr := postInteractionWait(logPrefix, page, scenario.PostInteractionWait)
			if sErr != nil {
				return sErr
			}

			if hs.WaitAfter != nil {
				delayMs := hs.WaitAfter.Milliseconds()
				fmt.Println(logPrefix, "hoverSelectors["+idx+"]: waitAfter starting sleep for", delayMs, "ms")
				time.Sleep(*hs.WaitAfter)
				fmt.Println(logPrefix, "hoverSelectors["+idx+"]: waitAfter ending sleep for", delayMs, "ms")
			}
		}
	}

	return nil
}
