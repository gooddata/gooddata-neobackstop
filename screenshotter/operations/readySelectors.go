package operations

import (
	"strconv"

	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/playwright-community/playwright-go"
)

func ReadySelectors(page playwright.Page, selectors []scenario.ReadySelector) *string {
	if len(selectors) == 0 {
		return nil
	}

	maxAttempts := 6
	timeout := float64(5000) // 5 seconds

	for attempt := 0; attempt < maxAttempts; attempt++ {
		allPresent := true
		var failedSelector string

		for _, sel := range selectors {
			state := playwright.WaitForSelectorState(sel.State)
			_, err := page.WaitForSelector(sel.Selector, playwright.PageWaitForSelectorOptions{
				State:   &state,
				Timeout: playwright.Float(timeout),
			})
			if err != nil {
				allPresent = false
				failedSelector = sel.Selector
				break // restart from first selector
			}
		}

		if allPresent {
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
