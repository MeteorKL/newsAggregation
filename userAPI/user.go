package userAPI

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/MeteorKL/newsAggregation/mgoHelper"
)

// 保证这两个字段是唯一的
// db.user.ensureIndex({ mail: 1,  },
//         {sparse: true, unique: true }
// )
// db.user.ensureIndex({ nickname: 1 },
//         {sparse: true, unique: true }
// )

type User struct {
	Mail     string
	NickName string
	Password string
	Tags     []string
}

const (
	collection = "user"
	sessionID  = "sessionID"
)

func loginCheck(nickname string, password string) map[string]interface{} {
	user, err := mgoHelper.MgoSelectOne(collection, func(c *mgo.Collection) (map[string]interface{}, error) {
		var r map[string]interface{}
		err := c.Find(bson.M{"nickname": nickname, "password": password}).Select(bson.M{"nickname": 1, "tags": 1}).One(&r)
		return r, err
	})
	if err != nil {
		println("error: loginCheck", err.Error())
		return nil
	}
	if user != nil {
		return user
	}
	return nil
}

func register(mail string, nickname string, password string) bool {
	user := User{
		Mail:     mail,
		NickName: nickname,
		Password: password,
	}
	err := mgoHelper.MgoInsert(collection, func(c *mgo.Collection) error {
		err := c.Insert(user)
		return err
	})
	if err != nil {
		println("error: Insert", collection, err.Error())
		return false
	}
	return true
}
