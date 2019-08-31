module github.com/carlosfrancia/web-crawler

go 1.12

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/Sirupsen/logrus v0.0.0-20170215164324-7f4b1adc7917
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0 // indirect
)

replace github.com/carlosfrancia/web-crawler/webcrawler => ./web-crawler
