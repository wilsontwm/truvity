# Truvity CLI

## About

Truvity CLI is a tool that will crawl the URLs being inputted and return the response body size in ascending order

## How to run

Prerequisite: Go version 1.18^

1. In the root directory, run `go run main.go` and you will be greeted with the following CLI screen in your terminal

```
Truvity CLI is a tool to crawl urls

Usage:
  truvity [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  crawl       Fast and simple URLs crawler
  help        Help about any command

Flags:
  -h, --help   help for truvity

Use "truvity [command] --help" for more information about a command.
```

2. Run `go run main.go crawl [urls...]` to start crawling and get the response body size. Example command:

```
go run main.go crawl https://facebook.com https://google.com https://www.wikipedia.org/
```
