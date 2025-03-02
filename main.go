package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type HackerNewsItem struct {
    Title string
    Site  string
    URL   string
}

func main() {
    // Initialize database
    database, err := InitDB("hackernews.db")
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer database.Close()

    // Create a new collector with rate limiting and configuration
    c := colly.NewCollector(
        colly.AllowURLRevisit(),
        colly.Async(true),
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
    )

    // Add rate limiting
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Parallelism: 2,
        Delay:       1 * time.Second,
    })

    // Set up scraping logic here
    c.OnHTML(".titleline", func(e *colly.HTMLElement) {
        item := HackerNewsItem{}

        // Extract title and URL from the first anchor tag
        titleLink := e.ChildAttr("a", "href")
        titleText := e.ChildText("a")
        item.Title = strings.TrimSpace(titleText)
        item.URL = strings.TrimSpace(titleLink)
        item.Site = strings.TrimSpace(e.ChildText(".sitestr"))

        // Only store items that have both title and URL
        if item.Title != "" && item.URL != "" {
            err := InsertItem(database, item.Title, item.Site, item.URL)
            if err != nil {
                if strings.Contains(err.Error(), "duplicate URL") {
                    log.Printf("Skipping duplicate: %s\n", item.Title)
                    return
                }
                log.Printf("Error storing item: %v\n", err)
                return
            }

            fmt.Printf("Stored: %s\n", item.Title)
            if item.Site != "" {
                fmt.Printf("Site: %s\n", item.Site)
            }
            fmt.Printf("URL: %s\n", item.URL)
            fmt.Println(strings.Repeat("-", 50))
            fmt.Println()
        }
    })

    // URLs to scrape
    urls := []string{
        "https://news.ycombinator.com/",
        "https://news.ycombinator.com/?p=2",
        "https://news.ycombinator.com/?p=3",
    }

    // Visit each URL
    for _, url := range urls {
        fmt.Printf("\nScraping page: %s\n", url)
        err := c.Visit(url)
        if err != nil {
            log.Printf("Error visiting %s: %v\n", url, err)
            continue
        }
    }

    // Wait for all scraping jobs to complete
    c.Wait()

    println("Job finished")
}