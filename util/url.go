package util

import (
	"fmt"
	"net/url"
	"strings"
)

type IURL interface {
	EscapedPath() string
	EscapedFragment() string
	String() string
	Redacted() string
	IsAbs() bool
	Query() url.Values
	RequestURI() string
	Hostname() string
	UnmarshalBinary(text []byte) error
	MarshalBinary() (text []byte, err error)
	Port() string
}

type IValues interface {
	Set(key, value string)
	String() string
	Get(key string) string
	Add(key, value string)
	Del(key string)
	Has(key string) bool
	Encode() string
}

type ParseIURLType func(rawURL string) (IURL, error)

var urlParse ParseIURLType = func(rawURL string) (IURL, error) {
	return url.Parse(rawURL)
}

func UpdateQueryParam(currentUrl string, key string, value string) (string, error) {
	url, err := urlParse(currentUrl)
	if err != nil {
		return "", err
	}
	q := url.Query()
	q.Set(key, value)
	return fmt.Sprintf("%s?%s", currentUrl, q.Encode()), nil
}

func ChangePathValue(url string, pathIndex int, newValue string) {
	pathArr := strings.Split(url, "/")
	fmt.Printf("pathArr: %v\n", pathArr)
}
