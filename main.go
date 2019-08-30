package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/carlosfrancia/web-crawler/cmd"
)

var flags = flag.NewFlagSet("Web-crawler", flag.ExitOnError)
var url string
var outputFile string

type webCrawlerConfig struct {
	url        string
	outputFile string
}

// TODO IMPORTANT
// YOU DONT NEED CONFIG. YOU CAN USE VARIABLES, THEY ARE GLOBAL

func init() {
	flags.StringVar(
		&url, "url", "", "The URL to be crawled",
	)
	flags.StringVar(
		&outputFile, "output-file", "sitemap.txt", "The filename where the sitemap will be written",
	)

	flags.Usage = func() {
		fmt.Printf("Usage of Web-crawler:\n")
		flags.PrintDefaults()
		fmt.Println("\n web-crawler -url www.monzo.com")
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

	config, err := newWebcrawlerConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Start crawling")

	err = cmd.Run(config.url, config.outputFile)
	// 	err = Run(config)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

}

func newWebcrawlerConfig() (*webCrawlerConfig, error) {

	switch {
	case (url == ""):
		flags.Usage()
		return nil, fmt.Errorf("Url to be crawled must be specified")
	}
	return &webCrawlerConfig{
		url:        url,
		outputFile: outputFile,
	}, nil
}
