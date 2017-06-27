package userAPI

import (
	"github.com/MeteorKL/newsAggregation/mgoHelper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func updateTag(nickname string, password string, tags []string) bool {
	err := mgoHelper.MgoUpdateOne(collection, func(c *mgo.Collection) error {
		return c.Update(
			bson.M{"nickname": nickname, "password": password},
			bson.M{"$set": bson.M{"tags": tags}},
		)
	})
	if err != nil {
		println("error: loginCheck", err.Error())
		return false
	}
	return true
}
