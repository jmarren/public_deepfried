package util

import (
	"fmt"
	"net/url"
	"strings"
)

func UpdateQueryParam(currentUrl string, key string, value string) (string, error) {
	url, err := url.Parse(currentUrl)
	if err != nil {
		return "", err
	}
	q := url.Query()
	q.Set(key, value)
	url.RawQuery = q.Encode()
	return url.String(), nil
}

func ChangePathValue(url string, pathIndex int, newValue string) {
	pathArr := strings.Split(url, "/")
	fmt.Printf("pathArr: %v\n", pathArr)
}

// func ChangePathValue(currentUrl string, pathParam string, newValue string) (string, error) {
// 	fmt.Printf("currentUrl: %s\npathParam: %s\nnewValue: %s\n", currentUrl, pathParam, newValue)
// 	currUrl, err := url.Parse(currentUrl)
// 	Eprint(err)
// 	fmt.Printf("currUrl: %v\n", currUrl)
// 	newUrl, err := url.JoinPath(currentUrl, "1")
// 	Eprint(err)
//
// 	return newUrl, nil
// }
