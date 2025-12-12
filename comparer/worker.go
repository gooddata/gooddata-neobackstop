package comparer

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gooddata/gooddata-neobackstop/config"
	"github.com/gooddata/gooddata-neobackstop/screenshotter"
)

func Run(c config.Config, jobs chan screenshotter.Result, wg *sync.WaitGroup, results chan Result, id int) {
	defer wg.Done()

	logPrefix := "comparer-" + strconv.Itoa(id) + " |"

	fmt.Println(logPrefix, "started")

	// iterate until channel is closed
	i := 0
	for job := range jobs {
		i++
		fmt.Println(logPrefix, "received job", i, "("+job.Scenario.Id+")")

		doJob(c, job, results)
	}

	fmt.Println(logPrefix, "finished")
}
