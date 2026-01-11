package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func GetWikiText(id string, base string) (string, error) {
	pageUrl := strings.Replace(base + "?action=edit", "%s", id, 1)
	fmt.Println(Info, "Hitting", pageUrl)

	resp, err := http.Get(pageUrl)
	if err != nil {
		fmt.Println(Error, err)
		return "", err
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(Error, err)
		return "", err
	}

	page := string(ba)

	start := strings.Index(page, "<textarea readonly=\"\" accesskey=")
	if start == -1 {
		fmt.Println(Error, "source not found")
		return "", err
	}

	elementDeclarationEnd := strings.Index(page[start:], ">")
	if elementDeclarationEnd == -1 {
		fmt.Println(Error, "elementDeclarationEnd not found")
		return "", err
	}
	elementDeclarationEnd += start + 1

	end := strings.Index(page[elementDeclarationEnd:], "</textarea>")
	if end == -1 {
		fmt.Println(Error, "source end not found")
		return "", err
	}
	end += elementDeclarationEnd

	return strings.ReplaceAll(page[elementDeclarationEnd:end], "\r", ""), nil
}

var ExtractTitle = regexp.MustCompile(`{{DISPLAYTITLE:(.*?)}}`)
var ExtractConnections = regexp.MustCompile(`{{[\S\s]*\|connections=(.*?)(}|=)`) //{{Dream[\S\s]*\|connections=(.*?)(}|=)
var ConnectionSplitter = regexp.MustCompile(`\]\].*?\[\[`)  
var ExtractImage = regexp.MustCompile(`{{[\S\s]*\|image1=\[?\[?(File:)?([^\]\[\{\}]*)`) //{{Dream.*\|image1=\[?\[?F?i?l?e?:?(.*?)\|
var ExtractLinkTitle = regexp.MustCompile(`\[?\[?([^|\[\]]*)([^\[\]]*)\]?\]?`)

func ModifyTitle(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

func ExtractDream(wikiText string) (Dream, []string, error) {
	wikiText = strings.ReplaceAll(wikiText, "\n", "")
	title := ExtractTitle.FindStringSubmatch(wikiText)
	connectionsRawOutput := ExtractConnections.FindStringSubmatch(wikiText)
	connections := []string{}
	if len(connectionsRawOutput) > 1 {
		connectionsRaw := connectionsRawOutput[1]

		if strings.LastIndex(connectionsRaw, "|") != -1 {
			connectionsRaw = connectionsRaw[:strings.LastIndex(connectionsRaw, "|")]
		}

		for _, connectionRaw := range ConnectionSplitter.Split(connectionsRaw, -1) {
			id := ExtractLinkTitle.FindStringSubmatch(connectionRaw)
			connections = append(connections, ModifyTitle(id[1]))
		}
	} else {
		fmt.Println(Info, "Dead end found!")
	}

	dream := Dream{}
	if ExtractImage.Match([]byte(wikiText)) {
		image := ExtractImage.FindStringSubmatch(wikiText)[2]
		if strings.Contains(image, "|") {
			image = image[:strings.Index(image, "|")]
		}
		dream.ImageViewLink = strings.ReplaceAll(*WikiBase, "%s", "File:" + image)
		dream.ImageHqLink = strings.ReplaceAll(*WikiBase, "%s", "Special:FilePath/" + image)
		dream.Image = GetImageDataUri(dream.ImageHqLink)
	} else {
		fmt.Println(Warn, "No image found!")
	}

	if len(title) > 1 {
		dream.Name = title[1]
	}

	return dream, connections, nil
}

var ExtractRedirect = regexp.MustCompile(`#REDIRECT \[\[(.*)\]\]`)

func CrawlPage(id string, base string) (d Dream, c []string, e error) {
	wikiText, err := GetWikiText(id, base)
	if err != nil {
		return Dream{}, []string{}, err
	}

	redir := ExtractRedirect.FindStringSubmatch(wikiText)
	if len(redir) > 1 {
		fmt.Println("Redirecting", id, "to", ModifyTitle(redir[1]))
		page, c, e := CrawlPage(ModifyTitle(redir[1]), base)
		page.Id = id
		return page, c, e
	}

	d, c, e = ExtractDream(wikiText)
	d.Id = id
	fmt.Println(Ok, "Crawled", id)
	return
}