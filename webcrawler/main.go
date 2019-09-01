package webcrawler

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/carlosfrancia/web-crawler/utils"
)

type webcrawlerConfig struct {
	domain     string
	url        string
	outputFile *os.File
	isVisited  *safeIsVisited
}

type safeIsVisited struct {
	sync.RWMutex
	m map[string]bool
}

// Set - Set given value to the given key
func (c *safeIsVisited) set(key string, value bool) {
	c.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.m[key] = value
	c.Unlock()
}

// Value returns mapped value for the given key.
func (c *safeIsVisited) get(key string) bool {
	c.RLock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.RUnlock()
	val := c.m[key]
	return val
}

// Run - Main control function, has to decide whether to continue or exit at each step
func Run(url, outputFileName string) error {

	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed writing to file '%s': %s", outputFileName, err)
	}
	fmt.Fprint(output, fmt.Sprintf("Sitemap for URL: %s \n", url))
	defer output.Close()

	ch := make(chan string)
	config := webcrawlerConfig{
		domain:     url,
		url:        url,
		outputFile: output,
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}

	go parseURL(config, ch)
	for resp := range ch {
		// TODO This is needed, should I write file here?
		fmt.Println("RESULT" + resp)
	}
	log.Info("Web crawled has finished.")
	return nil
}

func getUnvisited(allUrls []string, isVisited *safeIsVisited) (unvisitedUrls []string) {

	for _, url := range allUrls {
		if !isVisited.get(url) {
			unvisitedUrls = append(unvisitedUrls, url)
		}
	}
	return
}

func parseURL(config webcrawlerConfig, ch chan string) {
	defer close(ch)
	if config.isVisited.get(config.url) {
		return
	}
	config.isVisited.set(config.url, true)
	fmt.Fprint(config.outputFile, fmt.Sprintf("%s \n", config.url))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Get(config.url)
	if err != nil {
		print("ERROR processing URL:" + config.url + "\n")
		log.WithFields(log.Fields{
			"Url":   config.url,
			"Error": err,
		}).Info("Error processing url")
		return
	}
	defer response.Body.Close()
	ch <- config.url
	allUrls := utils.GetDomainLinks(response.Body)
	// println(len(allUrls))
	// println(allUrls)

	unvisitedUrls := getUnvisited(allUrls, config.isVisited)
	result := make([]chan string, len(unvisitedUrls))
	for i, url := range unvisitedUrls {
		config.url = url
		result[i] = make(chan string)
		go parseURL(config, result[i])

	}
	for i := range result {
		for response := range result[i] {
			ch <- response
		}
	}
	return
}
