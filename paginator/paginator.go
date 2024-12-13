package s3port

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func restAPIPaginator(text string) {
	limiter := rate.NewLimiter(10, 10) // 10 requests per second, with a burst of 10

	for i := 0; i < 20; i++ {
		if err := limiter.Wait(nil); err != nil { // Block until a token is available
			fmt.Println("Rate limiter error:", err)
			continue
		}
		go func(i int) {
			resp, err := http.Get("https://example.com")
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", i, err)
				return
			}
			fmt.Printf("Request %d succeeded: %s\n", i, resp.Status)
			resp.Body.Close()
		}(i)
		time.Sleep(100 * time.Millisecond) // Just to demonstrate concurrency
	}

	time.Sleep(2 * time.Second) // To ensure goroutines finish
}
