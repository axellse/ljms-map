package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var Info = "[\x1b[34mINFO\x1b[0m]"
var Ok = "[\x1b[32mOK\x1b[0m]"
var Error = "[\x1b[31mERROR\x1b[0m]"
var Warn = "\x1b[33m[WARN\x1b[0m]"

var WikiBase = flag.String("wiki", "https://latejulymidsummer.fandom.com/wiki/%s", "wiki url, with article name subsituted by %s")

var Dreams = []Dream{}
var Connections = []Connection{}

func ParseConnections(src string, c []string) {
	for _, connection := range c {
		Connections = append(Connections, Connection{
			FromID: src,
			ToID: connection,
		})

		if slices.ContainsFunc(Dreams, func(d Dream) bool {return strings.EqualFold(d.Id, connection)}) {
			continue
		}

		resultDream, newConnections, err := CrawlPage(connection, *WikiBase)
		if err != nil {
			fmt.Println(err, err)
			continue
		}

		if slices.ContainsFunc(Dreams, func(d Dream) bool {return strings.EqualFold(d.Id, connection)}) {
			continue
		}

		Dreams = append(Dreams, resultDream)

		ParseConnections(connection, newConnections)
	}
}

func main() {
	fmt.Println("welcome to ljms-map!")
	fmt.Println("written by axell - axell.me")
	fmt.Println()
	centerPage := flag.String("center", "The_nexus", "default page to start on")
	flag.Parse()
	fmt.Println("wiki base url is", *WikiBase + ", starting on", *centerPage)
	err := os.MkdirAll("./output", 0600)
	if err != nil {
		fmt.Println(err, err)
		os.Exit(1)
	}
	err = os.MkdirAll("./cache", 0600)
	if err != nil {
		fmt.Println(err, err)
		os.Exit(1)
	}

	startTime := time.Now()

	doCrawl := true
	if dreamsba, err := os.ReadFile("./cache/dreams.json"); err == nil {
		if connba, err := os.ReadFile("./cache/connections.json"); err == nil {
			doCrawl = false
			fmt.Println(Info, "Found cache, skipping crawling (delete cache directory to recrawl)")
			
			err := json.Unmarshal(dreamsba, &Dreams)
			if err != nil {
				fmt.Println(Error, "Failed reading cache: recrawling anyways.")
				doCrawl = true
			}

			if !doCrawl {
				err = json.Unmarshal(connba, &Connections)
				if err != nil {
					fmt.Println(Error, "Failed reading cache: recrawling anyways.")
					doCrawl = true
				}
			}
		}
	}

	if doCrawl {
		d, c, err := CrawlPage(*centerPage, *WikiBase)
		if err != nil {
			fmt.Println(err, err)
			os.Exit(1)
		}

		Dreams = append(Dreams, d)
		ParseConnections(d.Id, c)

		dreamsJson, err := json.MarshalIndent(Dreams, "", "  ")
		if err != nil {
			fmt.Println(err, err)
			os.Exit(1)
		}

		connectionsJson, err := json.MarshalIndent(Connections, "", "  ")
		if err != nil {
			fmt.Println(err, err)
			os.Exit(1)
		}

		fmt.Println(Ok, "Done Crawling!")
		fmt.Println("----------------------------------------")
		fmt.Println(Info, "Writing cache...")

		err = os.WriteFile("./cache/dreams.json", dreamsJson, 0600)
		if err != nil {
			fmt.Println(err, err)
			os.Exit(1)
		}
		err = os.WriteFile("./cache/connections.json", connectionsJson, 0600)
		if err != nil {
			fmt.Println(err, err)
			os.Exit(1)
		}
		fmt.Println(Ok, "Cache written.")
	}
	fmt.Println("----------------------------------------")
	fmt.Println(Info, "Now building output...")
	GenerateOutput(Dreams, Connections)
	fmt.Println(Ok, "Finished in", strconv.FormatFloat(time.Since(startTime).Seconds(), 'g', -1, 64) + "s")
}
        