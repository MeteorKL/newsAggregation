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

func loginCheck(nickname string, password string) bool {
	n, err := mgoHelper.MgoCount(collection, func(c *mgo.Collection) (int, error) {
		return c.Find(bson.M{"nickname": nickname, "password": password}).Count()
	})
	if err != nil {
		println("error: loginCheck", err.Error())
		return false
	}
	if n == 1 {
		return true
	}
	return false
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
