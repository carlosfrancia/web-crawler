package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

// TODO CAN I USE A STRUCT- DELETE OTHERWISE
type rootConfig struct {
	outputConfig *os.File
}

var output *os.File

// TODO IMPORTANT
// YOU DONT NEED CONFIG. YOU CAN USE VARIABLES, THEY ARE GLOBAL

func init() {

	// ?????????????????????
}

// Run - Main control function, has to decide whether to continue or exit at each step
func Run(url, outputFileName string) error {
	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed writing to file '%s': %s", outputFileName, err)
	}
	fmt.Fprint(output, fmt.Sprintf("Sitemap for URL: %s \n", url))
	defer output.Close()
	parsePage(url)
	log.Info("Web crawled has finished.")
	return nil
}

func processElement(index int, element *goquery.Selection) {
	// Check if it is a legitime link
	href, exists := element.Attr("href")
	//	if exists && noExternal(href) && isNotProcessed(href) {
	// add href to Map
	// print link
	if exists {
		fmt.Println(href)

	}
}

func parsePage(url string) {
	// ERROR MANAGEMENT!!!! This can return an error
	// Make HTTP request
	fmt.Fprint(output, "Starting parsePage")

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	var aTags = document.Find("a")

	// put in a map

	// link1 > false
	// link2 > false

	aTags.Each(processElement)

}
