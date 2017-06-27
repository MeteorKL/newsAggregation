package fetchNews

type News struct {
	ID     string `bson:"_id"`
	Title  string
	Docurl string
	Imgurl string
	Time   int64
	Tag    string
	Source string
}

const (
	collection = "news"
)
