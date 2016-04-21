# **revel_mgo**

## About
  **revel_mgo** is a code generator for **Revel web framework** https://revel.github.io/ with **Mongodb** database.
  It is mainly focus for **RESTFul API**. 
  revel_mgo can setup mongodb in your revel project easily and it is using "gopkg.in/mgo.v2" for mongo databse driver.
  It can also generate models and controllers for your RESTFul API with revel framework. 

## Installation
    go get github.com/kyawmyintthein/revel_mgo
    cd $GOPATH/src/path/to/revel_mgo
    go install revel_mgo

  

## Usage
#### Setup mongo driver
    revel_mgo mgo:setup

#### Update some code to startup mongo in your revel application   
  *Add following code in your app.conf file and change your database configuration.*
   
      mongo.database = database_name
	  mongo.path = localhost
	  mongo.maxPool = 20
	  
*Add following code in your init.go*

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

*Add this code under your init function of init.go file*

	revel.OnAppStart(initApp)

