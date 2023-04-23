package cmd

import (
	"crawler/web"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var directory string

var wgetCommand = &cobra.Command{
	Use:   "wget [url]",
	Short: "Fast and simple URLs crawler",
	Args:  cobra.MinimumNArgs(1),
	Run:   wget,
}

func init() {
	wgetCommand.Flags().StringVarP(&directory, "directory", "d", "./", "Destination directory to download webpage to")

	rootCmd.AddCommand(wgetCommand)
}

// wget is the main function that will be called by the wget command
func wget(cmd *cobra.Command, args []string) {

	if len(args) > 1 {
		fmt.Printf("Error: too many arguments")
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{})

	directory := cmd.Flag("directory").Value.String()
	urlString := args[0]

	// Check on the directory input
	if _, err := os.Stat(directory); err != nil && os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", directory)
		return
	} else if err != nil {
		fmt.Printf("%s is not a valid path\n", directory)
		return
	}

	go initProcessLink(doneChan, urlString, urlString, directory)

	select {
	case <-sigChan: // Wait for Ctrl+C signal to stop the program
		fmt.Println("Received SIGINT, terminating...")
	case <-doneChan:
		fmt.Println("Processing completed")
	}

}

// worker is a worker function that will be called by the main goroutine
func worker(initialUrl, currentUrl string, directory string, wg *sync.WaitGroup) {
	defer wg.Done()

	processLink(initialUrl, currentUrl, directory)

}

// processLink is a function that will process the link
// if there is child link in the page, then it will create the workers to process the child links
func processLink(initialUrl, currentUrl, directory string) {
	// check on URL input
	parsedURL, err := url.Parse(currentUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		fmt.Printf("%s is not a valid URL\n", currentUrl)
		return
	}

	webpage, err := web.NewWebpage(parsedURL.String())
	if err != nil {
		fmt.Printf("Error initializing webpage %s: %+v\n", currentUrl, err)
		return
	}

	if err := webpage.Process(initialUrl, directory); err != nil {
		fmt.Printf("Error processing webpage %s: %+v\n", currentUrl, err)
		return
	}

	// create a wait group to synchronize the workers
	var wg sync.WaitGroup

	// start the workers
	for _, childUrl := range webpage.ChildURLs {
		wg.Add(1)
		go worker(initialUrl, childUrl, directory, &wg)
	}

	// wait for all the workers to finish
	wg.Wait()
}

// initProcessLink is a init wrapper function to call processLink,
// once it's done, then it will send a signal to doneChan to inform that the process has completed
func initProcessLink(doneChan chan struct{}, initialUrl, currentUrl, directory string) {
	processLink(initialUrl, currentUrl, directory)

	doneChan <- struct{}{}

}
