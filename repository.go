package main

import (
	"time"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

type RemoteJobRepository struct {
	session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
}

func NewRemoteJobRepository() *RemoteJobRepository {
	return new(RemoteJobRepository)
}

func (rjr *RemoteJobRepository) Open() {
	session, err := mgo.Dial("mongodb://bluemirr.synology.me:8003/myHome")
	if err != nil {
		panic("mongodb not connected")
	}
	rjr.session = session
	rjr.Db = session.DB("myHome")
	rjr.Collection = rjr.Db.C("remotejobinfo")
}

func (rjr *RemoteJobRepository) Save(url, company, update string, timestamp time.Time) {
	model := new(RemoteJobModel)
	model.Company = company
	model.Url = url
	model.Timestamp = timestamp
	model.UpdateDate = update
	rjr.Collection.UpsertId(url, model)
}

func (rjr *RemoteJobRepository) Close() {
	rjr.session.Close()
}

func (rjr *RemoteJobRepository) FindAll() []RemoteJobModel {
	var result []RemoteJobModel
	rjr.Collection.Find(nil).Sort("-updatedate").All(&result)
	return result
}

type RemoteJobModel struct {
	Url       string `bson:"_id,omitempty"`
	Company   string
	UpdateDate string
	Timestamp time.Time
}