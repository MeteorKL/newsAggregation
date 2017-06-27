package main

import (
	"time"

	"github.com/MeteorKL/koala"
	"github.com/MeteorKL/newsAggregation/feedAPI"
	"github.com/MeteorKL/newsAggregation/fetchNews"
	"github.com/MeteorKL/newsAggregation/userAPI"
)

func toutiao(hrefs []string) {
	for i := 0; i < len(hrefs); i++ {
		koala.GetRequest(hrefs[i])

	}
}

func main() {
	go func() {
		timer1 := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-timer1.C:
				fetchNews.FetchNewsFrom163()
				fetchNews.FetchNewsFromSina()
				fetchNews.FetchNewsFromToutiao()
			}
		}
	}()

	// r := koala.GetRequest(URL)
	// ioutil.WriteFile("test.html", r, 0666)
	// ExampleScrape(URL)
	// koala.GetRequest(URL)
	// url and the target path
	userAPI.UserHandlers()
	feedAPI.FeedHandlers()
	koala.RunWithLog("1123")
}
