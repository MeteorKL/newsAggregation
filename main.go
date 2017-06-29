package main

import (
	"net/http"
	"time"

	"github.com/MeteorKL/koala"
	"github.com/MeteorKL/newsAggregation/feedAPI"
	"github.com/MeteorKL/newsAggregation/fetchNews"
	"github.com/MeteorKL/newsAggregation/userAPI"
)

// cd react && npm run build && cd ..
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

	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("./react/dist"))))
	// r := koala.GetRequest(URL)
	// ioutil.WriteFile("test.html", r, 0666)
	// ExampleScrape(URL)
	// koala.GetRequest(URL)
	// url and the target path

	userAPI.UserHandlers()
	feedAPI.FeedHandlers()

	koala.RenderPath = "react/dist/"
	koala.Get("/", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "index.html", nil)
	})

	koala.RunWithLog("1123")
}
