# Truvity CLI

## About

Truvity CLI is a tool that accepts a starting URL and a
destination directory. The crawler will then download the page at the URL, save it in
the destination directory, and then recursively proceed to any valid links in this page.

## How to run

Prerequisite: Go version 1.18^

1. In the root directory, run `go run main.go` and you will be greeted with the following CLI screen in your terminal

```
Truvity CLI is a tool to crawl urls

Usage:
  truvity [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  wget        Fast and simple URLs crawler

Flags:
  -h, --help   help for truvity

Use "truvity [command] --help" for more information about a command.
```

2. Run `go run main.go wget [urls]` to start crawling and download the webpage. Example command:

```
go run main.go wget https://www.brandeis.edu/student-affairs -d ../tmp
```

## Unit test

1. Run command in the root directory:

```
sh test.sh
```
