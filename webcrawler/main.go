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

type webcrawlerRunner struct {
	Client     *http.Client
	domain     string
	outputFile *os.File
}

type webcrawlerConfig struct {
	url       string
	isVisited *safeIsVisited
}

type safeIsVisited struct {
	sync.RWMutex
	m map[string]bool
}

// Set - Set given value to the given key
func (c *safeIsVisited) set(key string, value bool) {
	c.Lock()
	// Lock so only one goroutine at a time can access the map
	c.m[key] = value
	c.Unlock()
}

// Get - Returns mapped value for the given key.
func (c *safeIsVisited) get(key string) bool {
	c.RLock()
	// Lock so only one goroutine at a time can access the map
	defer c.RUnlock()
	val := c.m[key]
	return val
}

// Run - Initial function, needs to create all required items, call the parser and wait for finish
func Run(url, outputFileName string) error {

	// Create the file for the sitemap
	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed writing to file '%s': %s", outputFileName, err)
	}
	fmt.Fprint(output, fmt.Sprintf("Sitemap for URL: %s \n", url))
	defer output.Close()
	// Create initial channel and configuration
	ch := make(chan string)
	config := webcrawlerConfig{
		url: url,
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}
	// Create the runner
	r := webcrawlerRunner{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		domain:     url,
		outputFile: output,
	}
	go r.parseURL(config, ch)
	// Wait until all channels have been closed
	for range ch {
	}
	log.Info("Web crawled has finished.")
	return nil
}

// getUnvisited - Remove already visited URLs from a given slice
func getUnvisited(allUrls []string, isVisited *safeIsVisited) (unvisitedUrls []string) {

	for _, url := range allUrls {
		if !isVisited.get(url) {
			unvisitedUrls = append(unvisitedUrls, url)
		}
	}
	return
}

// parseURL - Parse recursively an URL and send the desired URLs to the channel
func (r *webcrawlerRunner) parseURL(config webcrawlerConfig, ch chan string) {
	defer close(ch)
	// If URL has already visited do nothing and return
	// TODO TO avoid a race condition while get isVisited and set, do it in the same method, lock -> get-> if not exist -> set -> unlock
	if config.isVisited.get(config.url) {
		return
	}
	// If hasn't been visited yet add it to the visite map
	config.isVisited.set(config.url, true)
	// TODO Is fmt thread safe?
	fmt.Fprint(r.outputFile, fmt.Sprintf("%s \n", config.url))
	// Make the request
	// TODO Handle limit of sockets (file descriptors)
	response, err := r.Client.Get(config.url)
	if err != nil {
		// Log if error with one request. Do not exit as it could be legitime
		log.WithFields(log.Fields{
			"Url":   config.url,
			"Error": err,
		}).Info("Error processing url")
		return
	}
	// Close the connection at the end of the function
	defer response.Body.Close()
	// Add response (url) to the channel
	ch <- config.url
	// Get all the links which belong to the domain
	allUrls := utils.GetDomainLinks(response.Body, r.domain)
	// From the entire list, exclude those already visited
	unvisitedUrls := getUnvisited(allUrls, config.isVisited)
	// Create an array of channels with the lengh of th URL to be processed
	result := make([]chan string, len(unvisitedUrls))
	// Process each one in a new go routine
	for i, url := range unvisitedUrls {
		config.url = url
		result[i] = make(chan string)
		go r.parseURL(config, result[i])

	}
	// Loop through the channels and set the response in the channel
	for i := range result {
		for response := range result[i] {
			ch <- response
		}
	}
	return
}
