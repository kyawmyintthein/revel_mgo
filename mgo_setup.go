package main

import(
	"os"
	"path"
	// "strings"
	"log"
)

var cmdMgoSetup = &Command{
	UsageLine: "mgo:setup",
	Short:     "create new config for Revel application",
	Long: `
Create new database.conf file in app/conf of your application.

It puts all necessary files as import. 

--driver is required. The configuration will be changed based on --driver parameter.

For example:
    revel db:setup
`,
}

var collectionTpl = `package mongodb

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Collection struct {
	db      *Database
	name    string
	Session *mgo.Collection
}

func (c *Collection) Connect() {
	session := *c.db.session.C(c.name)
	c.Session = &session
}

func NewCollectionSession(name string) *Collection {
	log.Println(name)
	var c = Collection{
		db:   newDBSession(DBNAME),
		name: name,
	}
	log.Printf("%+v",c)
	c.Connect()
	return &c
}

func (c *Collection) Close() {
	service.Close(c)
}

`

var databaseTpl = `package mongodb

import "gopkg.in/mgo.v2"

type Database struct {
	s       *mgo.Session
	name    string
	session *mgo.Database
}

func (db *Database) Connect() {

	db.s = service.Session()
	session := *db.s.DB(db.name)
	db.session = &session

}

func newDBSession(name string) *Database {

	var db = Database{
		name: name,
	}
	db.Connect()
	return &db
}
`

var driverTpl = `package mongodb

var MaxPool int
var PATH    string
var DBNAME  string


func CheckAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = PATH
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}


`

var serviceTpl = `package mongodb

import "gopkg.in/mgo.v2"

type Service struct {
	baseSession *mgo.Session
	queue       chan int
	URL         string
	Open        int
}

var service Service

func (s *Service) New() error {
	var err error
	s.queue = make(chan int, MaxPool)
	for i := 0; i < MaxPool; i = i + 1 {
		s.queue <- 1
	}
	s.Open = 0
	s.baseSession, err = mgo.Dial(s.URL)
	return err
}

func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}

`


func init() {
	cmdMgoSetup.Run = mgoSetup
}

func mgoSetup(cmd *Command, args []string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("$GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}
	pwd, _ := os.Getwd()
	databaseFolder := path.Join(pwd, "app", "models", "mongodb")
	files := []string{"database","collection","driver","service"}

	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		os.MkdirAll(databaseFolder, 0777)
	}

	for _, filename := range files{
		databaseFile := path.Join(databaseFolder, filename+".go")
		if file, err := os.OpenFile(databaseFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
			defer file.Close()
			content := "";
			switch filename{
			case "database":
				content = databaseTpl
			case "collection": 
			    content = collectionTpl
			case "driver":
				content = driverTpl
			case "service":
				content = serviceTpl
			}
			file.WriteString(content)
		} else {
			log.Println(err)
			errorf("Missing database.go.\nRun 'revel help db' for usage.\n")
		}	
	}
}
