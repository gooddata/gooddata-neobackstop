package operations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/playwright-community/playwright-go"
)

func ReadySelectors(logPrefix string, page playwright.Page, selectors []scenario.ReadySelector) *string {
	if len(selectors) == 0 {
		return nil
	}

	maxAttempts := 6
	timeout := float64(5000) // 5 seconds
	t0 := time.Now()
	fmt.Println(logPrefix, "readySelectors: waiting for", len(selectors), "selectors")

	for attempt := 0; attempt < maxAttempts; attempt++ {
		allPresent := true
		var failedSelector string

		for _, sel := range selectors {
			state := playwright.WaitForSelectorState(sel.State)
			tSel := time.Now()
			_, err := page.WaitForSelector(sel.Selector, playwright.PageWaitForSelectorOptions{
				State:   &state,
				Timeout: playwright.Float(timeout),
			})
			if err != nil {
				fmt.Println(logPrefix, "readySelectors: selector", sel.Selector, "failed attempt", attempt+1, "after", time.Since(tSel).Milliseconds(), "ms")
				allPresent = false
				failedSelector = sel.Selector
				break // restart from first selector
			}
			fmt.Println(logPrefix, "readySelectors: selector", sel.Selector, "appeared in", time.Since(tSel).Milliseconds(), "ms (attempt", strconv.Itoa(attempt+1)+")")
		}

		if allPresent {
			fmt.Println(logPrefix, "readySelectors: all", len(selectors), "selectors appeared in", time.Since(t0).Milliseconds(), "ms")
			return nil // all selectors present in same pass
		}

		// last attempt failed
		if attempt == maxAttempts-1 {
			e := "ReadySelector " + failedSelector + " didn't appear after " + strconv.Itoa(maxAttempts) + " attempts"
			return &e
		}
	}

	return nil
}
