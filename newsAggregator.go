package gositeaggregator

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

// NewsAggPage presents the page for the news aggregator.
type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

// Sitemapindex represents the main site sitemap which lists
// all other sitemaps.
type Sitemapindex struct {
	Locations []string `xml:"sitemap>loc"`
}

// News represents the Titles, Keywords, and Locations of a
// given sitemap.
type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

// NewsMap is a struct representing each individual News
// article.
type NewsMap struct {
	Keywords string
	Location string
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Whoa, Go is awesome!")
}

func About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Expert web design by Shaquizzle :)")
}

func newsRoutine(c chan News, Location string) {
	defer wg.Done()
	var n News
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()
	// Send News struct into channel
	c <- n
}

func NewsAgg(w http.ResponseWriter, r *http.Request) {
	var s Sitemapindex

	resp, err := http.Get("https://www.telegraph.co.uk/sitemap.xml")
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	xml.Unmarshal(bytes, &s)
	resp.Body.Close()

	newsMap := make(map[string]NewsMap)

	// Create a channel which takes in News structs, with a
	// buffer of size 512 (should be large enough to hold) all
	// news items.
	q := make(chan News, 128)

	// Create goroutine which will read all news items for
	// each (nested) sitemap.
	for _, Loc := range s.Locations {
		wg.Add(1)
		go newsRoutine(q, Loc)
	}

	// Wait for all goroutines to have been loaded onto channel
	wg.Wait()
	close(q)
	// Iterate over the channel and unpack each item into a
	// NewsMap key-value pair.
	for elem := range q {
		for idx := range elem.Keywords {
			newsMap[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
		}
	}

	p := NewsAggPage{Title: "Shaq's Concurrent News Aggregator", News: newsMap}
	t, _ := template.ParseFiles("newsaggtemplate.html")
	fmt.Println(t.Execute(w, p))
}