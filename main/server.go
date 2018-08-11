package main

import (
	gsa "sentdextuts/gositeaggregator"
	"net/http"
)

func main() {
	http.HandleFunc("/", gsa.Home)
	http.HandleFunc("/about/", gsa.About)
	http.HandleFunc("/agg/", gsa.NewsAgg)
	http.ListenAndServe(":8000", nil)
}
