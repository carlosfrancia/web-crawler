package utils

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HTMLTestSuite - d
type HTMLTestSuite struct {
	suite.Suite
}

func (suite *HTMLTestSuite) TestWhenAllLinksAreValidReturnAll() {

	expected := []string{"https://www.foo.com/", "https://www.foo.com/bar", "https://www.foo.com/foo/bar", "https://www.foo.com/foo/bar"}
	tmpfile, err := os.Open("../examples/example_all_valid_links.html")
	if err != nil {
		log.Fatal(err)
	}
	actual := GetDomainLinks(tmpfile, "https://www.foo.com")

	assert.Len(suite.T(), actual, 4)
	assert.Equal(suite.T(), actual, expected)
}

func (suite *HTMLTestSuite) TestWhenExternalLinksReturnOnlyDomainLinks() {

	expected := []string{"https://www.foo.com/bar"}
	tmpfile, err := os.Open("../examples/example_with_external_links.html")
	if err != nil {
		log.Fatal(err)
	}
	actual := GetDomainLinks(tmpfile, "https://www.foo.com")

	assert.Len(suite.T(), actual, 1)
	assert.Equal(suite.T(), actual, expected)
}

func (suite *HTMLTestSuite) TestWhenNoDomainLinksReturnZero() {

	tmpfile, err := os.Open("../examples/example_with_no_domain_links.html")
	if err != nil {
		log.Fatal(err)
	}
	actual := GetDomainLinks(tmpfile, "https://www.foo.com")

	assert.Len(suite.T(), actual, 0)

}

func (suite *HTMLTestSuite) TestWhenNoLinksReturnZero() {

	tmpfile, err := os.Open("../examples/example_with_no_links.html")
	if err != nil {
		log.Fatal(err)
	}
	actual := GetDomainLinks(tmpfile, "https://www.foo.com")

	assert.Len(suite.T(), actual, 0)

}

// TestHTMLTestSuite - Suite to test HTML package
func TestHTMLSuite(t *testing.T) {
	suite.Run(t, new(HTMLTestSuite))
}
