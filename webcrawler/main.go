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

type safeIsVisited struct {
	sync.RWMutex
	m map[string]bool
}

// type safeIsVisited struct {
// 	m   map[string]bool
// 	mux sync.Mutex
// }

var c chan string

// Run - Main control function, has to decide whether to continue or exit at each step
func Run(url, outputFileName string) error {

	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed writing to file '%s': %s", outputFileName, err)
	}
	fmt.Fprint(output, fmt.Sprintf("Sitemap for URL: %s \n", url))
	defer output.Close()

	ch := make(chan string)
	isVisited := safeIsVisited{
		m: make(map[string]bool),
	}

	go parseURL(url, output, isVisited, ch)
	for resp := range ch {
		fmt.Println("RESULT" + resp)

	}
	log.Info("Web crawled has finished.")
	return nil
}

func isHrefVisited(href string, isVisited safeIsVisited) bool {
	isVisited.RLock()
	result := isVisited.m[href]
	return result
}

func getUnvisited(allUrls []string, isVisited safeIsVisited) (unvisitedUrls []string) {

	for _, url := range allUrls {
		if !isHrefVisited(url, isVisited) {
			unvisitedUrls = append(unvisitedUrls, url)
		}
	}
	print("univisited")
	print(unvisitedUrls)
	return
}

func parseURL(url string, output *os.File, isVisited safeIsVisited, ch chan string) {

	defer close(ch)
	println(url)
	isVisited.RLock()
	if isVisited.m[url] {
		isVisited.RUnlock()
		return
	}
	isVisited.RUnlock()
	isVisited.Lock()
	isVisited.m[url] = true
	isVisited.Unlock()
	fmt.Fprint(output, fmt.Sprintf("%s \n", url))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		print("ERROR processing URL:" + url + "\n")
		log.WithFields(log.Fields{
			"Url":   url,
			"Error": err,
		}).Info("Error processing url")
		return
	}
	defer response.Body.Close()
	ch <- url
	allUrls := utils.GetDomainLinks(response.Body)
	println(len(allUrls))
	println(allUrls)
	if len(allUrls) == 0 {
		return
	}
	unvisitedUrls := getUnvisited(allUrls, isVisited)
	result := make([]chan string, len(unvisitedUrls))
	for i, url := range unvisitedUrls {

		result[i] = make(chan string)
		go parseURL(url, output, isVisited, result[i])

	}

	for i := range result {
		for response := range result[i] {
			ch <- response
		}
	}

	return

}
