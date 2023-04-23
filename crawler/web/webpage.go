package web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Webpage struct {
	URL       string
	Body      []byte
	ChildURLs []string
}

// NewWebpage: instantiate a new Webpage object
// URL is parsed and cleaned by removing the query
func NewWebpage(inputUrl string) (*Webpage, error) {
	urlObj, err := url.Parse(inputUrl)
	if err != nil {
		return nil, err
	}

	// Clean the query
	urlObj.RawQuery = ""

	return &Webpage{
		URL: urlObj.String(),
	}, nil
}

// IsDownloaded: A flag to indiciate if the page has been downloaded before to the directory
func (wp *Webpage) IsDownloaded(directory string) bool {
	filename := url.QueryEscape(wp.URL)
	path := filepath.Join(directory, filename)

	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

// Process: Process the webpage
func (wp *Webpage) Process(initialURL string, directory string) error {
	if wp.IsDownloaded(directory) {
		fmt.Printf("%s is already downloaded\n", wp.URL)
		return nil
	}

	fmt.Printf("%s is downloading...\n", wp.URL)
	resp, err := http.Get(wp.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	wp.Body, err = ioutil.ReadAll(tee)
	if err != nil {
		return err
	}

	err = wp.download(directory)
	if err != nil {
		return err
	}

	err = wp.extractLinks(initialURL)
	if err != nil {
		return err
	}

	return nil
}

// download: Download the webpage to the directory
// the file name will be the escaped URL
func (wp *Webpage) download(directory string) error {
	// Create a new file to save the downloaded webpage
	filename := url.QueryEscape(wp.URL)
	path := filepath.Join(directory, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the HTTP response body to the file
	_, err = io.Copy(file, bytes.NewReader(wp.Body))
	if err != nil {
		return err
	}

	return nil
}

// extractLinks: Extract all the links from the current webpage that are from the initial Url
func (wp *Webpage) extractLinks(initialUrl string) error {

	links := make([]string, 0)
	doc, err := html.Parse(bytes.NewReader(wp.Body))
	if err != nil {
		return err
	}

	var parseLink func(*html.Node)

	parseLink = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				// Only append the child link if it has same prefix with the initial URL
				if a.Key == "href" && strings.HasPrefix(a.Val, initialUrl) {
					links = append(links, a.Val)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseLink(c)
		}
	}

	parseLink(doc)

	wp.ChildURLs = links
	return nil
}
