package cmd

import (
	"cli/webpage"
	"fmt"
	"sort"

	"sync"

	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var crawlCommand = &cobra.Command{
	Use:   "crawl [urls...]",
	Short: "Fast and simple URLs crawler",
	Args:  cobra.MinimumNArgs(1),
	Run:   crawler,
}

func init() {
	rootCmd.AddCommand(crawlCommand)
}

func crawler(cmd *cobra.Command, urls []string) {
	var wg sync.WaitGroup

	webpages := make([]webpage.Webpage, len(urls))
	for i := range urls {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			wp := fetchFromUrl(urls[i])
			webpages[i] = wp
		}(i)
	}

	wg.Wait()

	webpagesWithoutError := funk.Filter(webpages, func(wp webpage.Webpage) bool {
		return wp.Error == nil
	}).([]webpage.Webpage)

	sort.Slice(webpagesWithoutError, func(i, j int) bool {
		return webpagesWithoutError[i].ResponseSize < webpagesWithoutError[j].ResponseSize

	})

	for _, webpage := range webpagesWithoutError {
		fmt.Printf("%s : %d \n", webpage.URL, webpage.ResponseSize)
	}
}

func fetchFromUrl(url string) webpage.Webpage {
	wp := webpage.NewWebpage(url)
	wp.Fetch()

	return *wp
}
