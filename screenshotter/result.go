package screenshotter

import (
	"github.com/gooddata/gooddata-neobackstop/internals"
)

type Result struct {
	Scenario                      *internals.Scenario
	Success                       bool
	Error                         *string
	FileName                      *string
	PreComputedMatch              *bool    // If set, comparison was already done
	PreComputedMismatchPercentage *float64 // Mismatch percentage if comparison was done
	PreComputedHasReference       *bool    // Whether reference existed during pre-computation
}

func buildResultFromScenario(scenario internals.Scenario, fileName *string, error *string) Result {
	return Result{
		Scenario: &scenario,
		Success:  error == nil,
		Error:    error,
		FileName: fileName,
	}
}
