package main

import(
	"os"
	"path"
)

var cmdMgoSetup = &Command{
	UsageLine: "mgo:setup",
	Short:     "setup mongodb in your Revel application",
	Long: `
Create new files such database.go, driver.go, collection.go, service.go in mongodb package under app/models path. Congiuration is required. 
Add following in your app.conf under [env]. 
	mongo.database = revel_sample_dev
	mongo.path = localhost
	mongo.maxPool = 20

Add following code in your init.go
   	
   	revel.OnAppStart(initApp)

	// initApp contains all application level initialization
	func initApp() {
		Config, err := revel.LoadConfig("app.conf")
		if err != nil || Config == nil {
			log.Fatalf("%+v",err)
		}
		mongodb.MaxPool = revel.Config.IntDefault("mongo.maxPool", 0)
		mongodb.PATH,_ = revel.Config.String("mongo.path")
		mongodb.DBNAME, _ = revel.Config.String("mongo.database")
		mongodb.CheckAndInitServiceConnection()
	}

For example:
    revel db:setup
`,
}

var collectionTpl = `package mongodb

import (
	"gopkg.in/mgo.v2"
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
	var c = Collection{
		db:   newDBSession(DBNAME),
		name: name,
	}
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

	// check #GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ColorLog("$GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}

	// get current path
	pwd, _ := os.Getwd()

	// get mongodb package path
	databaseFolder := path.Join(pwd, "app", "models", "mongodb")
	
	// mongodb setup files
	files := []string{"database","collection","driver","service"}

	//create mongodb package folder
	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		os.MkdirAll(databaseFolder, 0777)
	}

	// Write mongo file based on template
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
			ColorLog("[SUCC] mongo file generated as '%s' .\n", databaseFile)
		} else {
			ColorLog("Missing database.go.\n '%s' \nRun 'revel help db' for usage.\n", err)
		}	
	}
	ColorLog("[SUCC] mongodb package now in your project.\n")
	ColorLog("[SUCC] Please add database config in your 'app.conf'.\n")
	ColorLog("[SUCC] Please add mongo init code in your 'init.go'.\n")
}
