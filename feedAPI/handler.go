package feedAPI

import (
	"fmt"
	"net/http"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/MeteorKL/koala"
	"github.com/MeteorKL/newsAggregation/mgoHelper"
)

// Feed api has params topTime, bottomTime, show_num, tag, source
func FeedHandlers() {
	koala.Get("/api/feed", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		findM := make(bson.M)
		if arr, ok := p.ParamGet["tag"]; ok {
			if arr[0] == "推荐" {
				if session := koala.PeekSession(r, "sessionID"); session != nil {
					tags := session.Values["tags"]
					if tags, ok := tags.([]interface{}); ok && len(tags) != 0 {
						findM["tag"] = bson.M{
							"$in": tags,
						}
					}
				}
			} else if arr[0] != "" {
				findM["tag"] = arr[0]
			}
		}
		var timeUpdate string
		if arr, ok := p.ParamGet["topTime"]; ok {
			timeUpdate = "topTime"
			i, err := strconv.Atoi(arr[0])
			if err != nil {
				println("error: ParamGet topTime is not int", err.Error())
			}
			findM["time"] = bson.M{
				"$gt": int64(i),
			}
		} else if arr, ok := p.ParamGet["bottomTime"]; ok {
			timeUpdate = "bottomTime"
			i, err := strconv.Atoi(arr[0])
			if err != nil {
				println("error: ParamGet bottomTime is not int", err.Error())
			}
			findM["time"] = bson.M{
				"$lt": int64(i),
			}
		}
		if arr, ok := p.ParamGet["source"]; ok {
			findM["source"] = arr[0]
		}
		show_num := 10
		if arr, ok := p.ParamGet["show_num"]; ok {
			var err error
			show_num, err = strconv.Atoi(arr[0])
			if err != nil {
				println("error: ParamGet show_num is not int", err.Error())
			}
		}
		fmt.Print(findM)
		news, _ := mgoHelper.MgoSelectAll("news", func(c *mgo.Collection) ([]map[string]interface{}, error) {
			var r []map[string]interface{}
			err := c.Find(findM).Sort("-time").Limit(show_num).All(&r)
			println(r)
			return r, err
		})
		if news == nil {
			koala.WriteJSON(w, map[string]interface{}{
				"status":  1,
				"message": "没有更多新闻了",
				"data":    nil,
			})
		} else {
			var newTime interface{}
			if timeUpdate == "topTime" {
				newTime = news[0]["time"]
			} else {
				newTime = news[len(news)-1]["time"]
			}
			koala.WriteJSON(w, map[string]interface{}{
				"status":  0,
				"message": "拉取新闻成功",
				"data": map[string]interface{}{
					"newTime": newTime,
					"content": news,
				},
			})
		}
	})
}
