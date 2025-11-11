package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8000/"
	count := 10

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp.Body.Close()
		}(i)
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Printf("\n%d запросов за %v\n", count, end)
}
