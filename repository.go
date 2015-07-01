package main

import (
	"strings"
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
	infoFromDB := new(RemoteJobModel)
	rjr.Collection.FindId(url).One(infoFromDB)

	infoFromDB.Company = company
	infoFromDB.Url = url
	infoFromDB.Timestamp = timestamp
	infoFromDB.UpdateDate = update

	addFlag := true
	for i := 0; infoFromDB != nil && i < len(infoFromDB.UpdateHistory); i++ {
		if strings.EqualFold(infoFromDB.UpdateHistory[i], update) {
			addFlag = false
			break
		}
	}
	if addFlag {
		infoFromDB.UpdateHistory = append(infoFromDB.UpdateHistory, update)
	}

	rjr.Collection.UpsertId(url, infoFromDB)
}

func (rjr *RemoteJobRepository) Close() {
	rjr.session.Close()
}

func (rjr *RemoteJobRepository) FindAll() []RemoteJobModel {
	var result []RemoteJobModel
	rjr.Collection.Find(nil).Sort("-timestamp", "-updatedate").All(&result)
	return result
}

type RemoteJobModel struct {
	Url           string `bson:"_id,omitempty"`
	Company       string
	UpdateDate    string
	Timestamp     time.Time
	UpdateHistory []string
}
