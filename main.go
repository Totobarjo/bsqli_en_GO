package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/projectdiscovery/goflags"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	banner = ` _____         _ _ 
| __  |___ ___| |_|
| __ -|_ -| . | | |
|_____|___|_  |_|_|
            |_|      made by Coffinxp :)` + "\n"
)

type Config struct {
	UrlString   string // Single URL to scan.
	UrlFile     string // Text file containing a list of URLs to scan.
	PayloadFile string // Text file containing the payloads to append to the URLs.
	CookieP     string // Cookie to include in the GET request.
	ThreadP     int    // Number of concurrent threads
}

func performRequest(url, payload, cookie string) {
	urlWithPayload := url + payload
	startTime := time.Now()

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
		},
	}
	req, err := http.NewRequest("GET", urlWithPayload, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: "cookie", Value: cookie})
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	errMessage := ""
	if !success {
		errMessage = resp.Status
	}

	respTime := time.Since(startTime)

	display(success, urlWithPayload, errMessage, respTime)
}

func display(success bool, urlWithPayload, errMessage string, respTime time.Duration) {
	if respTime >= 10*time.Second {
		fmt.Printf("%s✔️  SQLi Found! URL: %s - Response Time: %v seconds%s\n", Yellow, urlWithPayload, respTime.Seconds(), Reset)
	} else {
		fmt.Printf("%s❌ Not Vulnerable. URL: %s - Response Time: %v seconds%s\n", Red, urlWithPayload, respTime.Seconds(), Reset)
	}
}


func main() {
	var (
		cfg Config
	)

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("bsqli, Perform GET requests to multiple URLs with different payloads.")
	flagSet.CreateGroup("input", "INPUT OPTIONS",
		flagSet.StringVarP(&cfg.UrlFile, "l", "list", "", "Text file containing a list of URLs to scan."),
		flagSet.StringVarP(&cfg.UrlString, "u", "url", "", "Single URL to scan."),
		flagSet.StringVarP(&cfg.PayloadFile, "p", "payloads", "", "Text file containing the payloads to append to the URLs."),
		flagSet.StringVarP(&cfg.CookieP, "c", "cookie", "", "Cookie to include in the GET request."),
	)
	flagSet.CreateGroup("optimizations", "OPTIMIZATIONS OPTIONS",
		flagSet.IntVarP(&cfg.ThreadP, "t", "threads", 40, "Number of concurrent threads"),
	)

	_ = flagSet.Parse()

	if cfg.UrlFile == "" && cfg.UrlString == "" {
		fmt.Printf("[%sErr%s] -u/-l is Required.\n", Red, Reset)
		os.Exit(1)
	} else if cfg.PayloadFile == "" {
		fmt.Printf("[%sErr%s] -p is Required.\n", Red, Reset)
		os.Exit(1)
	}

	if cfg.ThreadP <= 0 {
		fmt.Printf("[%sErr%s] Thread count must be a positive integer.\n", Red, Reset)
		os.Exit(1)
	}

	var urls []string
	if cfg.UrlString != "" {
		urls = append(urls, cfg.UrlString)
	} else if cfg.UrlFile != "" {
		file, err := ioutil.ReadFile(cfg.UrlFile)
		if err != nil {
			fmt.Printf("[%sErr%s] Error reading file: %s\n", Red, Reset, err)
			os.Exit(1)
		}
		urls = append(urls, strings.Split(strings.TrimSpace(string(file)), "\n")...)
	}

	file, err := ioutil.ReadFile(cfg.PayloadFile)
	if err != nil {
		fmt.Printf("[%sErr%s] Error reading file: %s\n", Red, Reset, err)
		os.Exit(1)
	}
	payloads := strings.Split(strings.TrimSpace(string(file)), "\n")

	fmt.Print(banner)

	var wg sync.WaitGroup
	requestsChan := make(chan struct {
		url     string
		payload string
	}, len(urls)*len(payloads))

	// Goroutines to process URLs and payloads concurrently
	for i := 0; i < cfg.ThreadP; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range requestsChan {
				performRequest(req.url, req.payload, cfg.CookieP)
			}
		}()
	}

	// Feeding URLs and Payloads into the combined channel
	go func() {
		for _, url := range urls {
			for _, payload := range payloads {
				requestsChan <- struct {
					url     string
					payload string
				}{url, payload}
			}
		}
		close(requestsChan)
	}()

	wg.Wait()
}
