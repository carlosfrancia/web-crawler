package utils

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// GetDomainLinks - Get all links of a HTML Body
func GetDomainLinks(b io.Reader) (urls []string) {

	// var urls = []string{}
	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			} else {
				// TODO IF IS DOMAIN AND .....
				if strings.HasPrefix(url, "/") {
					// println(url)
					urls = append(urls, "https://www.monzo.com"+url)
				}
			}
		}
	}
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}
