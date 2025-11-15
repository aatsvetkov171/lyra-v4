package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8000/ab"
	count := 1000

	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalTime time.Duration
	var maxTime time.Duration
	var errors int

	start := time.Now()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			reqStart := time.Now()
			resp, err := http.Get(url)
			elapsed := time.Since(reqStart)

			mu.Lock()
			if err != nil {
				errors++
			} else {
				resp.Body.Close()
				totalTime += elapsed
				if elapsed > maxTime {
					maxTime = elapsed
				}
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	avgTime := totalTime / time.Duration(count-errors)
	rps := float64(count-errors) / duration.Seconds()

	fmt.Println("----- РЕЗУЛЬТАТЫ -----")
	fmt.Printf("Всего запросов: %d\n", count)
	fmt.Printf("Ошибок: %d\n", errors)
	fmt.Printf("Среднее время ответа: %v\n", avgTime)
	fmt.Printf("Максимальное время ответа: %v\n", maxTime)
	fmt.Printf("Полное время теста: %v\n", duration)
	fmt.Printf("Примерная скорость: %.2f RPS\n", rps)
}
