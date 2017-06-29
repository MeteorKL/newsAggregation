package fetchNews

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"

	"github.com/MeteorKL/koala"
	"github.com/MeteorKL/newsAggregation/mgoHelper"
	iconv "github.com/djimenez/iconv-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 网易新闻 取每一天的新闻
// http://temp.163.com/special/00804KVA/cm_yaowen_03.js?callback=data_callback

type Keyword163 struct {
	Akey_link string
	Keyname   string
}
type News163 struct {
	Add1        string
	Add2        string
	Add3        string
	Channelname string
	Commenturl  string
	Digest      string
	Docurl      string
	Imgurl      string
	Keywords    []Keyword163
	Label       string
	Newstype    string
	Pics        []string
	Tienum      int
	Time        string
	Title       string
	Tlastid     string
	Tlink       string
}

type Data163 struct {
	Data []News163
}

const (
	source163 = "网易新闻"
)

var lastFecthTime int64

func fetchNewsFrom163Tag(tag string, tagurl string) {
	URL := "http://temp.163.com/special/00804KVA/cm_" + tagurl + ".js?callback=data_callback"
	input := koala.GetRequest(URL)
	output, err := iconv.ConvertString(string(input), "gb2312", "utf-8")
	if err != nil {
		println("error: iconv.ConvertString()", err.Error())
	}
	var data Data163
	if len(output) < 14 {
		println("error: fetchNewsFrom", tagurl, tag)
		return
	}
	output = "{\"data\":" + output[14:len(output)-1] + "}"
	json.Unmarshal([]byte(output), &data)

	n, err := mgoHelper.MgoQuery(collection, func(c *mgo.Collection) (interface{}, error) {
		var r interface{}
		err := c.Find(bson.M{"source": source163, "tag": tag}).Select(bson.M{"time": 1}).Sort("-time").One(&r)
		return r, err
	})
	if err != nil { // not found
		lastFecthTime = 0
	} else {
		lastFecthTime = n.(bson.M)["time"].(int64)
	}

	var news []interface{}
	var fetchNews int
	var newFetchTime int64
	for _, data := range data.Data {
		t, err := time.ParseInLocation("01/02/2006 15:04:05", data.Time, time.Local)
		if err != nil {
			continue
		}
		if t.Unix() <= lastFecthTime {
			break
		}
		fetchNews++
		if fetchNews == 1 {
			newFetchTime = t.Unix()
		}
		println(data.Time, data.Title)

		h := md5.New()
		io.WriteString(h, source163+tag+data.Time+data.Title)
		news = append(news, News{
			ID:     hex.EncodeToString(h.Sum(nil)),
			Source: source163,
			Tag:    tag,
			Title:  data.Title,
			Docurl: data.Docurl,
			Imgurl: data.Imgurl,
			Time:   t.Unix(),
		})
	}
	println(source163, tag, fetchNews)
	if fetchNews > 0 {
		lastFecthTime = newFetchTime

		err = mgoHelper.MgoInsert(collection, func(c *mgo.Collection) error {
			err := c.Insert(news...)
			return err
		})
		if err != nil {
			println("error: Insert", collection, err.Error())
		}
	}
}

func FetchNewsFrom163() {

	fetchNewsFrom163Tag("国内", "guonei")
	fetchNewsFrom163Tag("国际", "guoji")
	fetchNewsFrom163Tag("军事", "war")
	fetchNewsFrom163Tag("财经", "money")
	fetchNewsFrom163Tag("科技", "tech")
	fetchNewsFrom163Tag("体育", "sports")
	fetchNewsFrom163Tag("娱乐", "ent")
}
