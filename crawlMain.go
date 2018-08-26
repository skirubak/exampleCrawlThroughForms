package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"runtime"
	"time"

	"gocrawl"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"os"
	"os/exec"
	"github.com/ledongthuc/pdf"
	"bytes"
	"bufio"
)

var gCrawler *gocrawl.Crawler
var inputFilename string

type Ext struct {
	*gocrawl.DefaultExtender
}
func (e *Ext) Visited(ctx *gocrawl.URLContext,  harvested interface{}) {
	har := harvested.([]*url.URL)
	for _, element := range har {
		if strings.ToLower(element.Scheme) == "mailto" {
			gCrawler.EmailAddress[strings.ToLower(element.String())[7:len(element.String())]] +=1
			} else if ( len(element.Path) >= 4 && element.Path[len(element.Path)-4:] == ".pdf" ) {
			//fmt.Printf("Visited: pdf file %s\n",element.String())
			gCrawler.PdfFiles[strings.ToLower(element.String())] +=1
			} else if ( len(element.Path) >= 8 && strings.Contains(element.Path,"mailto" ) ) {
			//fmt.Printf("Visited : email id %s",strings.ToLower(element.Path)[8:len(element.Path)])
			gCrawler.EmailAddress[strings.ToLower(element.Path)[8:len(element.Path)]] +=1
			}
		}
	//fmt.Printf("Visisted : %s - %v\n", ctx.URL(),harvested)
     }

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	reg,_ := regexp.Compile("[^\n\r a-zA-Z0-9]+")
	//fmt.Printf("Visit: %s\n", ctx.URL())

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	//fmt.Printf("HeapAlloc=%02fMB; Sys=%02fMB\n", float64(stats.HeapAlloc)/1024.0/1024.0, float64(stats.Sys)/1024.0/1024.0)
	var content string
	plainText := strings.ToLower(reg.ReplaceAllString(doc.Find("body").Text(),""))
	if ( len(ctx.URL().Path) >= 4 && ctx.URL().Path[len(ctx.URL().Path)-4:] == ".pdf" ) {
		out, err := exec.Command("uuidgen").Output()
		if err == nil {
			filename:=string(out)
			filename=filename[1:len(filename)-1]+".pdf"
			body,_ := ioutil.ReadAll(res.Body)
			_ = ioutil.WriteFile(filename,body,0644)
			content, _ = readPdf(filename)
			content = strings.ToLower(reg.ReplaceAllString(content,""))
			//fmt.Printf("%s\n",content)
			_ = os.Remove(filename)
			}
		//fmt.Printf("Visit: pdf: %s\n",ctx.URL())
		}
		

	file , err1 := os.Open(inputFilename)
	if err1 != nil {
		return nil, true
		}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(),",")
		firstName, lastName := " "+strings.ToLower(s[0])+" ", " "+strings.ToLower(s[1])+" "
		matched, err := regexp.MatchString(firstName+".*"+lastName,plainText)
		if err == nil && matched == true {
			fmt.Printf("Visit: Name: %s %s Matched in %s\n",firstName,lastName,ctx.URL())
			}
		if ( len(ctx.URL().Path) >= 4 && ctx.URL().Path[len(ctx.URL().Path)-4:] == ".pdf" ) {
			matched, err = regexp.MatchString(firstName+".*"+lastName,content)
			if err == nil && matched == true {
				fmt.Printf("Visit: Name: %s %s Matched in %s\n",firstName,lastName,ctx.URL())
				}
			}
		}


	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	if ctx.URL().Host == "website" {
		return true
	}
	return false
}

func readPdf(path string) (string, error) {
    f,r, err := pdf.Open(path)
    defer f.Close()
    if err != nil {
        return "", err
    }
    var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
	return buf.String(), nil
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 5000
	inputFilename = os.Args[1]

	log.Print("starting crawl...")
	gCrawler = gocrawl.NewCrawlerWithOptions(opts)
	if err := gCrawler.Run("https://website"); err != nil {
		log.Print(err)
	}
}
