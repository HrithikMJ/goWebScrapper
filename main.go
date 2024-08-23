package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// var detailedRottenLinks map[string][]string

func main() {
	var rottenLinks []string
	var mu sync.Mutex
	detailedRottenLinks := make(map[int][]string)
	urlFlag := flag.String("u", "", "The URL to scrape")
	flag.Parse()
	if *urlFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	parsedURL, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		os.Exit(1)
	}
	var reqCount int

	var st time.Time
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(parsedURL.Host),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"))
	c.OnHTML("a", func(e *colly.HTMLElement) {
		// fmt.Println(e.Attr("href"))

		e.Request.Visit(e.Attr("href"))

	})
	c.Limit(&colly.LimitRule{Parallelism: 1000, RandomDelay: 5 * time.Second})
	c.OnRequest(func(r *colly.Request) {
		reqCount++
		if reqCount == 1 {
			st = time.Now()
		}
		temtTime := int(time.Since(st).Seconds())
		if temtTime != 0 {
			fmt.Printf("\r%d req %d depth %vs passed\n", reqCount/temtTime, r.Depth, temtTime)
		} else {
			fmt.Printf("\r%d req %d depth %vs passed\n", reqCount, r.Depth, temtTime)
		}
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Url %v Returned Response status: %d\n", r.Request.URL, r.StatusCode)
	})
	c.OnError(func(r *colly.Response, err error) {
		// fmt.Println("Error:", err)
		a := fmt.Sprintf("%v", r.Request.URL)
		mu.Lock()
		defer mu.Unlock()
		if _, exists := detailedRottenLinks[r.StatusCode]; exists {
			detailedRottenLinks[r.StatusCode] = append(detailedRottenLinks[r.StatusCode], a)
		} else {
			detailedRottenLinks[r.StatusCode] = append(detailedRottenLinks[r.StatusCode], a)
		}
		rottenLinks = append(rottenLinks, a)
		fmt.Printf("Url %v Returned Response status: %d\n", r.Request.URL.RawPath, r.StatusCode)
	})
	c.Visit(*urlFlag)
	c.Wait()
	fmt.Println(detailedRottenLinks)
}
