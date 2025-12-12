package comparer

import (
	"os"

	"github.com/gooddata/gooddata-neobackstop/config"
	"github.com/gooddata/gooddata-neobackstop/screenshotter"
	"github.com/gooddata/gooddata-neobackstop/utils"
)

func doJob(c config.Config, job screenshotter.Result, results chan Result) {
	// If comparison was already done in screenshotter, use pre-computed results
	if job.PreComputedHasReference != nil {
		if !*job.PreComputedHasReference {
			// Reference doesn't exist - file was already saved to disk by screenshotter
			e := "Reference image does not exist"
			results <- Result{
				ScreenshotterResult: &job,
				Error:               &e,
			}
			return
		}

		// Reference exists and comparison was done
		if job.PreComputedMatch != nil && *job.PreComputedMatch {
			// Image matches - use pre-computed result
			results <- Result{
				ScreenshotterResult: &job,
				HasReference:        true,
				MatchesReference:    true,
				MismatchPercentage:  job.PreComputedMismatchPercentage,
			}
			return
		}

		// Mismatch was detected - file was already saved to disk by screenshotter
		// Load images from disk to create diff
		referenceImg, err := utils.LoadImage(c.BitmapsReferencePath + "/" + *job.FileName)
		if err != nil {
			panic(err)
		}

		testImg, err := utils.LoadImage(c.BitmapsTestPath + "/" + *job.FileName)
		if err != nil {
			panic(err)
		}

		// Create diff image
		diff, _ := utils.DiffImagesPink(referenceImg, testImg)
		err = utils.SaveImage(c.BitmapsTestPath+"/diff_"+*job.FileName, diff)
		if err != nil {
			panic(err.Error())
		}

		results <- Result{
			ScreenshotterResult: &job,
			HasReference:        true,
			DiffCreated:         true,
			MatchesReference:    false,
			MismatchPercentage:  job.PreComputedMismatchPercentage,
		}
		return
	}

	// Fallback: original comparison logic for cases where pre-computation wasn't done
	// This should not happen in normal operation since pre-computation is always done in test mode
	if _, err := os.Stat(c.BitmapsReferencePath + "/" + *job.FileName); os.IsNotExist(err) {
		// Reference doesn't exist - file should already be saved by screenshotter
		e := "Reference image does not exist"
		results <- Result{
			ScreenshotterResult: &job,
			Error:               &e,
		}
		return
	}

	referenceImg, err := utils.LoadImage(c.BitmapsReferencePath + "/" + *job.FileName)
	if err != nil {
		panic(err)
	}

	testImg, err := utils.LoadImage(c.BitmapsTestPath + "/" + *job.FileName)
	if err != nil {
		panic(err)
	}

	diff, mismatch := utils.DiffImagesPink(referenceImg, testImg)

	if (job.Scenario.MisMatchThreshold != nil && *job.Scenario.MisMatchThreshold >= mismatch) || (job.Scenario.MisMatchThreshold == nil && mismatch == 0) {
		// image matches - do not save test screenshot or diff
		results <- Result{
			ScreenshotterResult: &job,
			HasReference:        true,
			MatchesReference:    true,
			MismatchPercentage:  &mismatch,
		}
		return
	}

	err = utils.SaveImage(c.BitmapsTestPath+"/diff_"+*job.FileName, diff)
	if err != nil {
		panic(err.Error())
	}

	results <- Result{
		ScreenshotterResult: &job,
		HasReference:        true,
		DiffCreated:         true,
		MatchesReference:    false,
		MismatchPercentage:  &mismatch,
	}
}
