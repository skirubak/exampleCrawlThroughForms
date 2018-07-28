package main

import (
	"bytes"
	"fmt"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
)

func main() {
	u,err := url.Parse("https://www.bon.texas.gov/forms/rnrslt.asp")
	if err != nil {
		return
	}
	f := url.Values{}
	f.Add("LicNumber","")
	f.Add("SSNumber","")
	f.Add("DOB","")
	f.Add("firstname","a")
	f.Add("lastname","c")
	
	req,e := http.NewRequest("POST",u.String(),strings.NewReader(f.Encode()))
	if e != nil {
		return
		}
	req.Header.Add("User-Agent","Mozilla/5.0 (X11; Linux armv7l) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.91 Safari/537.36")
	req.Header.Add("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Content-Type","application/x-www-form-urlencoded")
	req.Header.Add("Cookie","ASPSESSIONIDQWCDCDBS=OLOAOAOADEHANMHCJDOCFKAP")

	var HttpClient = & http.Client{}
	response,_ := HttpClient.Do(req)

	if response.StatusCode >=200 && response.StatusCode < 300 {
		if bd,_ := ioutil.ReadAll(response.Body); e != nil {
			return
		} else {
		if node, eP := html.Parse(bytes.NewBuffer(bd)); eP != nil {
		return
		} else {
			doc := goquery.NewDocumentFromNode(node)
			rows := doc.Find("table tbody tr td")
			rows.Each( func(i int, s *goquery.Selection) {
				if s.HasClass("content") == true {
					var f func(*html.Node)
					f = func(n *html.Node) {
						if n.Type== html.TextNode {
							fmt.Printf("%v,",n.Data)
						}
						if n.FirstChild != nil {
							for c := n.FirstChild; c != nil; c = c.NextSibling {
								f(c)
							}
						}
					}
					test := s.Children()
					for _, n := range test.Nodes {
						f(n)
						fmt.Printf("\n")
					}
				}
				})
			}
		}
	}
}
