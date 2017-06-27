package mgoHelper

import (
	"time"

	"github.com/MeteorKL/koala"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

var dialInfo map[string]string

func getSession() *mgo.Session {
	if dialInfo == nil {
		dialInfo, _ = koala.ReadJSONFile("./mgo.conf")
	}
	if mgoSession == nil {
		var err error
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:     []string{dialInfo["addr"]},
			Username:  dialInfo["username"],
			Password:  dialInfo["password"],
			Database:  dialInfo["database"],
			Source:    dialInfo["source"],
			Mechanism: dialInfo["mechanism"],
			Timeout:   60 * time.Second,
		}
		mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo) //连接数据库
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func MgoInsert(collection string, query func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoUpdateOne(collection string, query func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoUpdateAll(collection string, query func(*mgo.Collection) (*mgo.ChangeInfo, error)) (*mgo.ChangeInfo, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoCount(collection string, query func(*mgo.Collection) (int, error)) (int, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoSelectOne(collection string, query func(*mgo.Collection) (map[string]interface{}, error)) (map[string]interface{}, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoSelectAll(collection string, query func(*mgo.Collection) ([]map[string]interface{}, error)) ([]map[string]interface{}, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func MgoQuery(collection string, query func(*mgo.Collection) (interface{}, error)) (interface{}, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}

func selectFromCollection(collection string, query func(*mgo.Collection) (map[string]interface{}, error)) (map[string]interface{}, error) {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo["database"]).C(collection)
	return query(c)
}
