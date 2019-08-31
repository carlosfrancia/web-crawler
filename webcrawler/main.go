package webcrawler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

// TODO CAN I USE A STRUCT- DELETE OTHERWISE
type webcrawlerRunner struct {
	url        string
	outputFile *os.File
	visited    string
}

// // Runner2 - TODO
// type Runner2 struct {
// 	Url        string
// 	OutputFile *os.File
// 	Visited    string
// }

//TODO

// // NewWebCrawlerRunner - OPTION 2 Create a runner
// func NewWebCrawlerRunner(url string) (r *Runner2, err error) {

// 	// TODO SHOULD I TAKE THIS METHOD TO THE MAIN AND CALL THIS FIRST AS CREATE NEW RUNNER
// 	//
// 	// WHICH WILL RETURN A WEBCRAWLERRUNNER AND THEN DO R.PARSEPAGE?

// 	return &Runner2{
// 		Url:        url,
// 		OutputFile: nil,
// 		Visited:    "to be done",
// 	}, nil

// }

// Run - Main control function, has to decide whether to continue or exit at each step
func Run(url, outputFileName string) error {

	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed writing to file '%s': %s", outputFileName, err)
	}
	defer output.Close()
	r := webcrawlerRunner{
		url:        url,
		outputFile: output,
		visited:    "to be done",
	}
	r.parsePage()
	log.Info("Web crawled has finished.")
	return nil
}

func (r *webcrawlerRunner) processElement(index int, element *goquery.Selection) {
	// Check if it is a legitime link
	href, exists := element.Attr("href")
	//	if exists && noExternal(href) && isNotProcessed(href) {
	// add href to Map
	// print link
	if exists {
		// fmt.Println(href)
		fmt.Fprint(r.outputFile, fmt.Sprintf("%s \n", href))
	}
}

func (r *webcrawlerRunner) parsePage() {
	// ERROR MANAGEMENT!!!! This can return an error i.e if an empty URL is passed. I think the error is handled
	// Make HTTP request. What actually happens if one of the link is broken? Will it fail? You can do something instead of
	// log.Fatal and continue. Catch the error and continue

	log.WithField("Url", r.url).Info("Starting parsing URL")
	fmt.Fprint(r.outputFile, fmt.Sprintf("Sitemap for URL: %s \n", r.url))

	response, err := http.Get(r.url)
	if err != nil {
		log.WithFields(log.Fields{
			"Url":   r.url,
			"Error": err,
		}).Fatal("Error processing url")
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them
	var aTags = document.Find("a")

	// put in a map

	// link1 > false
	// link2 > false

	aTags.Each(r.processElement)
}

// func (r *Runner2) processElement2(index int, element *goquery.Selection) {
// 	// Check if it is a legitime link
// 	href, exists := element.Attr("href")
// 	//	if exists && noExternal(href) && isNotProcessed(href) {
// 	// add href to Map
// 	// print link
// 	if exists {
// 		fmt.Println(href)
// 	}
// }

// // ParsePage2 -
// func (r *Runner2) ParsePage2() {
// 	// ERROR MANAGEMENT!!!! This can return an error i.e if an empty URL is passed. I think the error is handled
// 	// Make HTTP request. What actually happens if one of the link is broken? Will it fail? You can do something instead of
// 	// log.Fatal and continue. Catch the error and continue

// 	log.WithField("Url", r.Url).Info("Starting parsing URL")
// 	fmt.Fprint(r.OutputFile, fmt.Sprintf("Sitemap for URL: %s \n", r.Url))

// 	response, err := http.Get(r.Url)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"Url":   r.Url,
// 			"Error": err,
// 		}).Fatal("Error processing url")
// 	}
// 	defer response.Body.Close()

// 	// Create a goquery document from the HTTP response
// 	document, err := goquery.NewDocumentFromReader(response.Body)
// 	if err != nil {
// 		log.Fatal("Error loading HTTP response body. ", err)
// 	}

// 	// Find all links and process them
// 	var aTags = document.Find("a")

// 	// put in a map

// 	// link1 > false
// 	// link2 > false

// 	aTags.Each(r.processElement2)
// }
