package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func parseGetListSCP() []string {
	response, err := http.Get("https://scpfoundation.net/scp-series")
	if err != nil {
		fmt.Errorf("some error request: %v", err)
	}
	defer response.Body.Close()
	bytesData, err := io.ReadAll(response.Body)
	var result []string
	htmlString := string(bytesData)
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		fmt.Errorf("some error htlm parse, %v", err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if strings.Contains(attr.Val, "/scp-") {
						part := strings.Split(attr.Val, "-")
						if len(part) > 1 {
							_, err := strconv.Atoi(part[1])
							if err == nil {
								result = append(result, attr.Val)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return result
}

func main() {
	arraySCP := parseGetListSCP()
	fmt.Println(arraySCP)
}
