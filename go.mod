module github.com/carlosfrancia/web-crawler

go 1.12

require (
	github.com/Sirupsen/logrus v0.0.0-20170215164324-7f4b1adc7917
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a
	golang.org/x/sys v0.0.0-20191025090151-53bf42e6b339 // indirect
)

replace github.com/carlosfrancia/web-crawler/webcrawler => ./webcrawler

replace github.com/carlosfrancia/web-crawler/utils => ./utils
