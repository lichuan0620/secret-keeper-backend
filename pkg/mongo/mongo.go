package mongo

import (
	"time"

	"github.com/globalsign/mgo"
)

const (
	dbName = "secret-keeper"

	CollectionBox = "box"
)

var baseSession *mgo.Session

func Init(url string) error {
	const (
		timeout     = 5 * time.Second
		maxIdleTime = int(5 * time.Minute / time.Millisecond)
		minPoolSize = 3
	)
	info, err := mgo.ParseURL(url)
	if err != nil {
		return err
	}
	info.Timeout = timeout
	info.ReadTimeout = timeout
	info.WriteTimeout = timeout
	info.MinPoolSize = minPoolSize
	info.MaxIdleTimeMS = maxIdleTime
	baseSession, err = mgo.DialWithInfo(info)
	if err != nil {
		return err
	}
	baseSession.SetSyncTimeout(timeout)
	return nil
}

func DB() *mgo.Database {
	return baseSession.Copy().DB(dbName)
}
