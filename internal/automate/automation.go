package automate

import (
	"log"
	"sync"
)

func RunAutomation() error {
	userCount := 1_000
	concurrency := 100 // how many run at once

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < userCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}        //acquire
			defer func() { <-sem }() //release

			if err := UserFlow(); err != nil {
				log.Println("user flow failed:", err)
			}

		}()
	}

	wg.Wait()
	return nil
}
