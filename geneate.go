package main

import(
	"os"
	"log"
)

var cmdGenerate = &Command{
	UsageLine: "generate [Command]",
	Short:     "source code generator",
	Long: `
bee generate scaffold [scaffoldname] [-fields=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    The generate scaffold command will do a number of things for you.
    -fields: a list of table fields. Format: field:type, ...
    -driver: [mysql | postgres | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test
    example: bee generate scaffold post -fields="title:string,body:text"

bee generate model [modelname] [-fields=""]
    generate RESTFul model based on fields
    -fields: a list of table fields. Format: field:type, ...

bee generate controller [controllerfile]
    generate RESTFul controllers             
`,
}

var fields flagValue
func init() {
	cmdGenerate.Run = generateCode
	cmdGenerate.Flag.Var(&fields, "fields", "specify the fields want to generate.")
}

func generateCode(cmd *Command, args []string){
	curpath, _ := os.Getwd()
	if len(args) < 1 {
		errorf("[ERRO] command is missing\n")
		os.Exit(2)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("[ERRO] $GOPATH not found\n")
		errorf("[HINT] Set $GOPATH in your environment vairables\n")
		os.Exit(2)
	}

	gcmd := args[0]
	log.Println(len(args))
	switch gcmd {
	case "scaffold":
	case "controller":
		if len(args) == 2 {
			// cname := args[1]
			// generateModel(cname, curpath)
		} else {
			errorf("[ERRO] Wrong number of arguments\n")
			errorf("[HINT] Usage: revel_mgo generate model [controllername]\n")
			os.Exit(2)
		}
	case "model":
		if len(args) < 2 {
			log.Println("[ERRO] Wrong number of arguments\n")
			errorf("[HINT] Usage: revel_mgo generate model [modelname] [-fields=\"\"]\n")
			os.Exit(2)
		}
		cmd.Flag.Parse(args[2:])
		if fields == "" {
			log.Println("[ERRO] Wrong number of arguments\n")
			errorf("[HINT] Usage: bee generate model [modelname] [-fields=\"title:string,body:text\"]\n")
			os.Exit(2)
		}
		sname := args[1]
		log.Printf("[INFO] Using '%s' as model name\n", sname)
		generateModel(sname, fields.String(), curpath)
	default:
		errorf("[ERRO] command is missing\n")
	}
	errorf("[SUCC] generate successfully created!\n")
}
