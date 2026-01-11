package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetImageDataUri(iurl string) string {
	resp1, err := http.Get(iurl)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}

	newLoc := resp1.Request.URL.String()
	if strings.Contains(newLoc, "/revision/latest") {
		newLoc = strings.ReplaceAll(newLoc, "/revision/latest", "/revision/latest/scale-to-width-down/100")
	}
	fmt.Println(Ok, "Downloading image", newLoc)

	resp, err := http.Get(newLoc)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}
	if !strings.HasPrefix(resp.Header.Get("content-type"), "image") {
		return ""
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(Error, err)
		return ""
	}

	return "data:" + resp.Header.Get("content-type") + ";base64," + base64.StdEncoding.EncodeToString(ba)
}