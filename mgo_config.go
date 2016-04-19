package main

import (
	"log"
	"os"
	"path"
	"strings"
)

var cmdMgoConfig = &Command{
	UsageLine: "mgo:config",
	Short:     "create new config for Revel application",
	Long: `
Create new database.conf file in app/conf of your application.

It puts all necessary files as import. 

--driver is required. The configuration will be changed based on --driver parameter.

For example:

    revel db:config --driver=[mysql]
`,
}

var connTplStr = `[dev]
		mongo.path = {{path}}
		mongo.name = {{name}}
		mongo.maxPool = {{maxPool}}`

func init() {
	cmdMgoConfig.Run = mgoConfig
}

func mgoConfig(cmd *Command, args []string) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("$GOPATH not found.\nRun 'revel help db' for usage.\n")
		os.Exit(2)
	}
	pwd, _ := os.Getwd()
	destFile := path.Join(pwd, "conf", "database.conf")
	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		if err = os.Remove(destFile); err != nil {
			log.Fatalln("Failed to remove existing file database.go:", err)
		}
	}
	if file, err := os.OpenFile(destFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer file.Close()
		content := strings.Replace(connTplStr, "{{path}}", "localhost", -1)
		content = strings.Replace(content, "{{name}}", "test", -1)
		content = strings.Replace(content, "{{maxPool}}", "30", -1)
		file.WriteString(content)
	} else {
		errorf("No driver given.\nRun 'revel help db' for usage.\n")
	}
}
