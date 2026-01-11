package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetImageDataUri(iurl string) string {
	fmt.Println(Ok, "Downloading image", iurl)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(iurl)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}

	newLoc := resp.Header.Get("Location")
	if strings.Contains(newLoc, "/revision/latest") {
		newLoc = strings.ReplaceAll(newLoc, "/revision/latest", "/revision/latest/scale-to-width-down/150")
	}

	resp, err = http.Get(newLoc)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}

	return "data:" + resp.Header.Get("content-type") + ";base64," + base64.StdEncoding.EncodeToString(ba)
}