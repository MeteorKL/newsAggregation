package feedAPI

import (
	"net/http"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/MeteorKL/koala"
	"github.com/MeteorKL/newsAggregation/mgoHelper"
)

// Feed api has params behot_time, show_num, tag, source
func FeedHandlers() {
	koala.Get("/api/feed", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		findM := make(bson.M)
		if arr, ok := p.ParamGet["tag"]; ok {
			findM["tag"] = arr[0]
		}
		if arr, ok := p.ParamGet["source"]; ok {
			findM["source"] = arr[0]
		}
		behot_time := time.Now().Unix()
		show_num := 10
		var err error
		if arr, ok := p.ParamGet["behot_time"]; ok {
			i, err := strconv.Atoi(arr[0])
			behot_time = int64(i)
			if err != nil {
				println("error: ParamGet behot_time is not int", err.Error())
			}
		}
		findM["time"] = bson.M{
			"$lt": behot_time,
		}
		if arr, ok := p.ParamGet["show_num"]; ok {
			show_num, err = strconv.Atoi(arr[0])
			if err != nil {
				println("error: ParamGet show_num is not int", err.Error())
			}
		}
		news, err := mgoHelper.MgoSelectAll("news", func(c *mgo.Collection) ([]map[string]interface{}, error) {
			var r []map[string]interface{}
			err := c.Find(findM).Sort("-time").Limit(show_num).All(&r)
			println(r)
			return r, err
		})
		if news == nil {
			koala.WriteJSON(w, map[string]interface{}{
				"message": "fail",
				"data":    news,
			})
		} else {
			koala.WriteJSON(w, map[string]interface{}{
				"message": "success",
				"data":    news,
			})
		}
	})
}
