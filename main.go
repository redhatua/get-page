package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	prefixCss = "css"
	prefixJs  = "js"
	prefixImg = "img"
)

var DownloadDir string

func main() {
	var url = flag.String("u", "", "Target URL")
	var downloadDir = flag.String("d", "", "Download directory")

	flag.Parse()

	if "" == *downloadDir || "" == *url {
		flag.Usage()
		os.Exit(1)
	}

	DownloadDir = *downloadDir

	r, err := http.Get(*url)
	checkErr(err)

	defer r.Body.Close()
	if _, err := os.Stat(DownloadDir); os.IsNotExist(err) {
		os.Mkdir(DownloadDir, 0755)
	}
	res := parse(r.Body)
	checkErr(err)

	html, err := res.Html()
	checkErr(err)

	file, err := os.Create(fmt.Sprintf("%s/index.html", DownloadDir))
	checkErr(err)

	defer file.Close()
	fmt.Fprintf(file, html)

}

func parse(h io.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(h)
	if err != nil {
		panic(err)
	}
	doc.Find("link[rel=stylesheet]").Each(func(i int, s *goquery.Selection) {
		if t, e := s.Attr("href"); e {
			ss := strings.Split(t, "/")
			sc := strings.Split(ss[len(ss)-1], "?")
			n := sc[0]
			s.SetAttr("href", fmt.Sprintf("./%s/%s", prefixCss, n))
			download(t, prefixCss, n)
		}
	})

	doc.Find("script[type=text/javascript]").Each(func(i int, s *goquery.Selection) {
		if t, e := s.Attr("src"); e {
			ss := strings.Split(t, "/")
			sc := strings.Split(ss[len(ss)-1], "?")
			n := sc[0]
			s.SetAttr("src", fmt.Sprintf("./%s/%s", prefixJs, n))
			download(t, prefixJs, n)
		}
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if t, e := s.Attr("src"); e {
			ss := strings.Split(t, "/")
			sc := strings.Split(ss[len(ss)-1], "?")
			n := sc[0]
			s.SetAttr("src", fmt.Sprintf("./%s/%s", prefixImg, n))
			download(t, prefixImg, n)
		}
	})

	return doc
}

func download(target, prefix, name string) {
	var fullPath = fmt.Sprintf("%s/%s", DownloadDir, prefix)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.Mkdir(fullPath, 0755)
	}
	out, err := os.Create(fmt.Sprintf("%s/%s", fullPath, name))
	checkErr(err)
	defer out.Close()

	r, err := http.Get(target)
	checkErr(err)
	defer r.Body.Close()

	_, err = io.Copy(out, r.Body)
	checkErr(err)

	fmt.Printf("Saved %s into %s\n", name, prefix)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
