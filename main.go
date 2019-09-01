package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/carlosfrancia/web-crawler/webcrawler"
)

var flags = flag.NewFlagSet("Web-crawler", flag.ExitOnError)
var url string
var outputFile string
var logLevel string

type webCrawlerConfig struct {
	url        string
	outputFile string
}

// TODO IMPORTANT
// YOU DONT NEED CONFIG. YOU CAN USE VARIABLES, THEY ARE GLOBAL

func init() {
	flags.StringVar(
		&logLevel, "log-level", "info", fmt.Sprintf("Log level, valid "+
			"values are %+v", log.AllLevels),
	)
	flags.StringVar(
		&url, "url", "", "The URL to be crawled",
	)
	flags.StringVar(
		&outputFile, "output-file", "sitemap.txt", "The filename where the sitemap will be written",
	)

	flags.Usage = func() {
		fmt.Printf("Usage of Web-crawler:\n")
		flags.PrintDefaults()
		fmt.Println("\n web-crawler -url http://www.monzo.com [-output-file mono-sitemap.txt]")
		//		fmt.Println("\n Web-crawler -url www.monzo.com -output-file monzo-site.txt")
		fmt.Println()
	}
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Starting
	log.SetOutput(os.Stderr)
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(ll)

	config, err := newWebcrawlerConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Start crawling")
	err = webcrawler.Run(config.url, config.outputFile)
	// 	err = Run(config)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func newWebcrawlerConfig() (*webCrawlerConfig, error) {

	// TODO Check URL format. error if it doesn't start with http, https or www.
	// If starts with www add http
	// ACCEPT ONLY HTTP or HTTPS. ERROR IF NOT

	switch {
	case (url == ""):
		flags.Usage()
		return nil, fmt.Errorf("Url to be crawled must be specified")
	case (!strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "https://")):
		flags.Usage()
		return nil, fmt.Errorf("Url must start with HTTP or HTTPS")
	}
	return &webCrawlerConfig{
		url:        url,
		outputFile: outputFile,
	}, nil
}

// // TODO - DECIDE WETHER TO USE MAIN OR MAIN2
// func main2() {
// 	// Starting
// 	log.SetOutput(os.Stderr)
// 	ll, err := log.ParseLevel(logLevel)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.SetLevel(ll)

// 	config, err := newWebcrawlerConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Info("Start crawling")
// 	r, err := webcrawler.NewWebCrawlerRunner(config.url)

// 	output, err := os.Create(config.outputFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer output.Close()

// 	r.OutputFile = output
// 	r.ParsePage2()
// 	// 	err = Run(config)
// 	if err != nil {
// 		log.Fatalf("ERROR: %s", err)
// 	}
// }
