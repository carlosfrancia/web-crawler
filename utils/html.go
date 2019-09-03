package utils

import (
	"io"
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/html"
)

// GetDomainLinks - Get all links which belong to the domain of a HTML Body
func GetDomainLinks(b io.Reader, domain string) (urls []string) {
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			//Document ended. We are done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			// Check if it itis an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			ok, href := getHref(t)
			if !ok {
				continue
			} else {
				isDomain, link := isValidDomainLink(href, domain)
				if isDomain {
					urls = append(urls, link)
				}
			}
		}
	}
}

// getHref - Get href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// isDomainLink - Return wether an href is a valid URL and belongs to the domain or not
func isValidDomainLink(href string, domain string) (bool, string) {
	if strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "//") {
		href = domain + href
	}
	if govalidator.IsURL(href) && strings.HasPrefix(href, domain) {
		return true, href
	}
	return false, ""
}
