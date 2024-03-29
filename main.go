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
		fmt.Println("\n web-crawler -url http://www.url.com [-output-file my-sitemap.txt]")
		fmt.Println()
	}
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Starting
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
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func newWebcrawlerConfig() (*webCrawlerConfig, error) {
	switch {
	case (url == ""):
		flags.Usage()
		return nil, fmt.Errorf("Url to be crawled must be specified")
	case (!strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://")):
		flags.Usage()
		return nil, fmt.Errorf("Url must start with HTTP or HTTPS")
	}
	return &webCrawlerConfig{
		url:        url,
		outputFile: outputFile,
	}, nil
}
