package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const ROCKET_HOST string = "http://rocketpun.ch"

type rocketJobHtmlParser struct {
	Depart    string
	StartPage int
	EndPage   int
	Keywords  []string
}

func NewRocketJobHtmlParser() *rocketJobHtmlParser {
	obj := &rocketJobHtmlParser{Depart: "개발자", StartPage: 1, EndPage: 10}
	obj.Keywords = []string{"재택", "원격근무"}
	return obj
}

func (r *rocketJobHtmlParser) TotalGet() {
	r.searchDetailKeyword(r.searchList())
}

func (r rocketJobHtmlParser) searchList() []*html.Node {
	var articles []*html.Node

	matcher := func(n *html.Node) bool {
		if n.Parent != nil && n.Parent.Parent != nil &&
			n.Parent.Parent.DataAtom == atom.Div && scrape.Attr(n.Parent.Parent, "class") == "hr_list" &&
			n.Parent.DataAtom == atom.Div && scrape.Attr(n.Parent, "class") == "hr_contents" &&
			n.DataAtom == atom.A &&
			n.FirstChild != nil && n.FirstChild.FirstChild != nil &&
			n.FirstChild.FirstChild.NextSibling != nil &&
			n.FirstChild.FirstChild.NextSibling.DataAtom == atom.Div &&
			n.FirstChild.FirstChild.NextSibling.FirstChild != nil &&
			n.FirstChild.FirstChild.NextSibling.FirstChild.DataAtom == atom.Div &&
			scrape.Attr(n.FirstChild.FirstChild.NextSibling.FirstChild, "class") == "hr_text_job" &&
			n.FirstChild.FirstChild.NextSibling.FirstChild.NextSibling != nil &&
			strings.Contains(scrape.Text(n.FirstChild.FirstChild.NextSibling.FirstChild), r.Depart) {
			return true
		} else {
			return false
		}
	}

	for i := r.StartPage; i <= r.EndPage; i++ {
		LIST_URL := ROCKET_HOST + "/recruit/list/" + strconv.Itoa(i) + "/"
		resp, err := http.Get(LIST_URL)
		if err != nil {
			continue
		}
		root, err := html.Parse(resp.Body)
		if err != nil {
			continue
		}
		articleArray := scrape.FindAll(root, matcher)
		for _, article := range articleArray {
			articles = append(articles, article)
		}
	}

	return articles
}

func (r rocketJobHtmlParser) searchDetailKeyword(articles []*html.Node) {
	remoteJobRepository := NewRemoteJobRepository()
	remoteJobRepository.Open()
	defer remoteJobRepository.Close()
	cnt := 0
	updateTimeStamp := time.Now()
	for _, article := range articles {
		company := scrape.Text(article.FirstChild.FirstChild.NextSibling.FirstChild.NextSibling)
		var updateStr string
		if article.NextSibling != nil &&
		article.NextSibling.NextSibling != nil &&
		article.NextSibling.NextSibling.FirstChild != nil &&
		article.NextSibling.NextSibling.FirstChild.FirstChild != nil {
			updateStr = scrape.Text(article.NextSibling.NextSibling.FirstChild.FirstChild)
		}
		
		dtailLink := scrape.Attr(article, "href")
		companyResp, err := http.Get(ROCKET_HOST + dtailLink)
		if err != nil {
			continue
		}
		companyRoot, err := html.Parse(companyResp.Body)
		if err != nil {
			continue
		}
		detailText := scrape.Text(companyRoot)
		isKeyword := false
		for _, keyword := range r.Keywords {
			isKeyword = strings.Contains(detailText, keyword)
		}
		if isKeyword {
			url := ROCKET_HOST + dtailLink
			updateSlice := strings.Fields(updateStr)
			var updateDate string
			if len(updateSlice) >= 2 {
				updateDate = updateSlice[1]
				updateDate = strings.TrimSpace(updateDate)
			}
			
			fmt.Println(strconv.Itoa(cnt) + ":" + updateDate + ":" + company + ":" + url)
			remoteJobRepository.Save(url, company, updateDate, updateTimeStamp)
			cnt++
		}
	}
	fmt.Println("=============================");
}
