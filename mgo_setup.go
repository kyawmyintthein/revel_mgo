package main

import(
	"os"
	"path"
	// "strings"
	"log"
	"github.com/robfig/config"
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

var collectionTpl = `package db

import (
	"gopkg.in/mgo.v2"
	"errors"
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
	dbname,flag := dbConfig.String("mongo.name")
	if flag == false{
		panic(errors.New("Missing databse name in config."))
	}
	var c = Collection{
		db:   newDBSession(dbname),
		name: name,
	}
	c.Connect()
	return &c
}

func (c *Collection) Close() {
	service.Close(c)
}
`

var databaseTpl = `package db

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

var driverTpl = `package db

import "github.com/revel/revel"

const (
	dbConfigPath string = "conf/"
	dbConfigName string = "database.conf"
)

var maxPool int
var dbConfig *revel.MergedConfig

func init() {
	dbConfig = loadDatabaseConfig()
	maxPool = dbConfig.IntDefault("mongo.maxPool", 30)
	checkAndInitServiceConnection()
}

func checkAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL, _ = dbConfig.String("mongo.maxPool")
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}

func loadDatabaseConfig() *revel.MergedConfig {
	revel.ConfPaths = append(revel.ConfPaths, dbConfigPath)
	Config, err := revel.LoadConfig(dbConfigName)
	if err != nil{
		panic(err)
	}
	return Config;
}  
`

var serviceTpl = `package db

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
	s.queue = make(chan int, maxPool)
	for i := 0; i < maxPool; i = i + 1 {
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

	config, err := config.ReadDefault(path.Join(pwd, "conf", "database.conf"))
	if err != nil || config == nil {
		log.Fatalln("Failed to load database.conf:", err)
	}

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
