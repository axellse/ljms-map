package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"slices"
	"strconv"
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

var okExtensions = []string{
	"jpeg",
	"jpg",
	"png",
	"gif",
	"svg",
	"webp",
}

func DownloadImage(url string) (string, string) {
	fmt.Println(Ok, "Downloading high-res image", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(Error, err)
		return "", ""
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(Error, err)
		return "", ""
	}

	name := strconv.FormatInt(rand.Int64N(100_000_000), 16) 
	ext := ""
	if exti := slices.IndexFunc(okExtensions, func(ext string) bool {return resp.Header.Get("content-type") == "image/" + ext}); exti != -1 {
		name += "." + okExtensions[exti]
		ext = okExtensions[exti]
	}
	
	err = os.WriteFile("./output/images/" + name, out, 0600)
	if err != nil {
		fmt.Println(Error, err)
		return "", ""
	}
	
	return "./images/" + name, ext
}