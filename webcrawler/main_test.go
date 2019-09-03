package webcrawler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HTMLTestSuite - Type
type WebcrawlerTestSuite struct {
	suite.Suite
}

// TestRunCreatesOutputFile - Test function run (only creates the output file)
func (suite *WebcrawlerTestSuite) TestRunCreatesOutputFile() {

	_ = Run("foo", "outputFile.txt")
	content, err := ioutil.ReadFile("outputFile.txt")
	defer func() {
		if _err := os.Remove("outputFile.txt"); _err != nil {
			log.Fatal(_err)
		}
	}()
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), string(content), "Sitemap for URL: foo")

}

// TestAllValidAndRepeteadLinks - Test with an HTML that contains all valid links plus one repeated
func (suite *WebcrawlerTestSuite) TestAllValidAndRepeteadLinks() {

	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if _err := os.Remove(tmpfile.Name()); _err != nil {
			log.Fatal(_err)
		}
	}()

	client := NewTestClient(func(req *http.Request) *http.Response {
		contents, _ := ioutil.ReadFile("../examples/example_all_valid_links.html")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(string(contents))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	config := webcrawlerConfig{
		url: "https://www.foo.com/",
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}
	// Use Client
	runner := webcrawlerRunner{
		Client:     client,
		domain:     "https://www.foo.com",
		outputFile: tmpfile,
	}
	ch := make(chan string)
	go runner.parseURL(config, ch)
	for range ch {
	}
	content, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(suite.T(), err)
	numberOfLines := countChars(string(content), '\n')
	assert.Contains(suite.T(), string(content), "https://www.foo.com/")
	assert.Contains(suite.T(), string(content), "https://www.foo.com/foo/bar")
	assert.Contains(suite.T(), string(content), "https://www.foo.com/bar")
	assert.Equal(suite.T(), 3, numberOfLines)
}

// TestOneAndExternalLinks - Test with an HTML that contains all valid links plus one repeated
func (suite *WebcrawlerTestSuite) TestOneAndExternalLinks() {

	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if _err := os.Remove(tmpfile.Name()); _err != nil {
			log.Fatal(_err)
		}
	}()

	client := NewTestClient(func(req *http.Request) *http.Response {
		contents, _ := ioutil.ReadFile("../examples/example_with_external_links.html")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(string(contents))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	config := webcrawlerConfig{
		url: "https://www.foo.com/",
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}
	// Use Client
	runner := webcrawlerRunner{
		Client:     client,
		domain:     "https://www.foo.com",
		outputFile: tmpfile,
	}
	ch := make(chan string)
	go runner.parseURL(config, ch)
	for range ch {
	}
	content, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(suite.T(), err)
	numberOfLines := countChars(string(content), '\n')
	assert.Contains(suite.T(), string(content), "https://www.foo.com/bar")
	// Result is 2 because the original domain is always part of the result.
	assert.Equal(suite.T(), 2, numberOfLines)
}

// TestOnlyExternalLinks - Test with an HTML that contains only external links
func (suite *WebcrawlerTestSuite) TestOnlyExternalLinks() {

	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if _err := os.Remove(tmpfile.Name()); _err != nil {
			log.Fatal(_err)
		}
	}()

	client := NewTestClient(func(req *http.Request) *http.Response {
		contents, _ := ioutil.ReadFile("../examples/example_with_no_domain_links.html")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(string(contents))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	config := webcrawlerConfig{
		url: "https://www.foo.com/",
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}
	// Use Client
	runner := webcrawlerRunner{
		Client:     client,
		domain:     "https://www.foo.com",
		outputFile: tmpfile,
	}
	ch := make(chan string)
	go runner.parseURL(config, ch)
	for range ch {
	}
	content, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(suite.T(), err)
	numberOfLines := countChars(string(content), '\n')
	// The original domain is always part of the result
	assert.Equal(suite.T(), 1, numberOfLines)
}

// TestNoLinks - Test with an HTML that contains none links
func (suite *WebcrawlerTestSuite) TestNoLinks() {

	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if _err := os.Remove(tmpfile.Name()); _err != nil {
			log.Fatal(_err)
		}
	}()
	client := NewTestClient(func(req *http.Request) *http.Response {
		contents, _ := ioutil.ReadFile("../examples/example_with_no_links.html")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(string(contents))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	config := webcrawlerConfig{
		url: "https://www.foo.com/",
		isVisited: &safeIsVisited{
			m: make(map[string]bool),
		},
	}
	// Use Client
	runner := webcrawlerRunner{
		Client:     client,
		domain:     "https://www.foo.com",
		outputFile: tmpfile,
	}
	ch := make(chan string)
	go runner.parseURL(config, ch)
	for range ch {
	}
	content, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(suite.T(), err)
	numberOfLines := countChars(string(content), '\n')
	// The original domain is always part of the result
	assert.Equal(suite.T(), 1, numberOfLines)
}

func countChars(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// TestHTMLTestSuite - Suite to test the main Webcrawler package
func TestWebcrawlerSuite(t *testing.T) {
	suite.Run(t, new(WebcrawlerTestSuite))
}
