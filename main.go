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
	var result []string
	for i := 1; i < 10; i++ {
		fmt.Println(i)
		var url string
		if i == 1 {
			url = "https://scpfoundation.net/scp-series"
		} else {
			url = "https://scpfoundation.net/scp-series-" + strconv.Itoa(i)
		}
		response, err := http.Get(url)
		if err != nil {
			fmt.Errorf("some error request: %v", err)
		}
		defer response.Body.Close()
		bytesData, err := io.ReadAll(response.Body)
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
	}
	return result
}

type SCPUnit struct {
	Name        string
	Class       string
	Structure   string
	Filial      string
	Anomaly     string
	Subject     []string
	Discription string
	SpecialCOD  string
	Property    []string
	Link        string
}

func parseGetCurrentSCP(data []string) []SCPUnit {
	var result []SCPUnit
	for _, name := range data {
		fmt.Println("Парсинг обьекта %v", name[1:])
		url := "https://scpfoundation.net" + name
		unit := SCPUnit{
			Name:        "",
			Class:       "",
			Structure:   "",
			Filial:      "",
			Anomaly:     "",
			Discription: "",
			SpecialCOD:  "",
			Link:        url,
		}
		response, err := http.Get(url)
		if err != nil {
			fmt.Errorf("some error request: %v", err)
		}
		bytesData, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		htmlString := string(bytesData)
		doc, err := html.Parse(strings.NewReader(htmlString))
		if err != nil {
			fmt.Errorf("some error htlm parse, %v", err)
		}
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode {
				for _, attr := range n.Attr {
					if n.Data == "div" {
						if attr.Key == "id" {
							if attr.Val == "page-title" {
								if n.FirstChild != nil {
									unit.Name = strings.TrimSpace(n.FirstChild.Data)
								}
							}
							if attr.Val == "page-content" {
								if n.FirstChild != nil {
									fmt.Println(n.FirstChild.Data)
								}
							}
							// if attr.Val == "page-content" {
							// 	if n.FirstChild != nil {
							// 		fmt.Println("YA TUT")
							// 		fmt.Println(n.FirstChild.Data)
							// 	}
							// }
						}
					}
					if n.Data == "a" {
						if len(strings.Split(attr.Val, "/system:page-tags/tag/")) > 1 {
							value := strings.Split(strings.Split(attr.Val, "/system:page-tags/tag/")[1], "#pages")[0]
							switch strings.Split(value, ":")[0] {
							case "структура":
								unit.Structure = strings.Split(value, ":")[1]
							case "класс":
								unit.Class = strings.Split(value, ":")[1]
							case "филиал":
								unit.Filial = strings.Split(value, ":")[1]
							case "свойство":
								unit.Property = append(unit.Property, strings.Split(value, ":")[1])
							case "аномалия":
								unit.Anomaly = strings.Split(value, ":")[1]
							case "тематика":
								unit.Subject = append(unit.Subject, strings.Split(value, ":")[1])
							}
							fmt.Println(strings.Split(strings.Split(attr.Val, "/system:page-tags/tag/")[1], "#pages")[0])
						}
					}
					if n.Data == "p" {
					}
				}
			}
			if n.Type == html.TextNode && n.Data == "Описание:" {
				unit.Discription = strings.TrimSpace(n.Parent.Parent.LastChild.Data)
			}
			if n.Type == html.TextNode && n.Data == "Особые условия содержания:" {
				unit.SpecialCOD = strings.TrimSpace(n.Parent.Parent.LastChild.Data)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		// fmt.Println()
		result = append(result, unit)
		break
	}
	return result
}

func main() {
	data := []string{"/scp-1122"}
	res := parseGetCurrentSCP(data)
	fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
