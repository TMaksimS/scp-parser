package scpclient

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type ScpClient struct {
	URL     string
	Headers map[string]string
	Client  *http.Client
}

func (client *ScpClient) ParseGetListSCP() []string {
	var result []string
	for i := 1; i < 10; i++ {
		fmt.Println(i)
		var url string
		if i == 1 {
			url = client.URL + "/scp-series"
		} else {
			url = client.URL + "/scp-series-" + strconv.Itoa(i)
		}
		response, err := client.Client.Get(url)
		for key, value := range client.Headers {
			response.Header.Set(key, value)
		}
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

func (client *ScpClient) ParseGetCurrentSCP(data string) SCPUnit {
	fmt.Println("Парсинг обьекта %v", data[1:])
	url := client.URL + data
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("erroe created req: %v", err)
	}
	for key, value := range client.Headers {
		req.Header.Set(key, value)
	}
	if err != nil {
		fmt.Errorf("some error request: %v", err)
	}
	response, err := client.Client.Do(req)
	if err != nil {
		fmt.Println("error sending req: %v", err)
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
						default:
							fmt.Println(strings.Split(strings.Split(attr.Val, "/system:page-tags/tag/")[1], "#pages")[0])
						}
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
	return unit
}
