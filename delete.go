package main

import(
	"os"
)

var cmdDelete = &Command{
	UsageLine: "delete [Command] [filename]",
	Short:     "remove source code",
	Long: `
revel_mgo delete [command] [filename] 

revel_mgo delete model [model_name] 

revel_mgo delete controller [controller_name]           
`,
}

func init() {
	cmdDelete.Run = deleteCode
}

func deleteCode(cmd *Command, args []string){
	// get current path 
	curpath, _ := os.Getwd()

	// check command args exist
	if len(args) < 1 {
		ColorLog("[ERRO] Command is missing\n")
		ColorLog("[HINT] Usage: revel_mgo delete [model or controller] [name]\n")
		os.Exit(2)
	}

	//check gopath in local machine
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ColorLog("[ERRO] $GOPATH not found\n")
		ColorLog("[HINT] Set $GOPATH in your environment vairables\n")
		os.Exit(2)
	}

	// get file type model or controller
	dcmd := args[0]

	switch dcmd {
	case "controller":
		if len(args)<  2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo delete controller [controllername]\n")
			os.Exit(2)
		}

	case "model":
		if len(args) < 2 {
			ColorLog("[ERRO] Wrong number of arguments\n")
			ColorLog("[HINT] Usage: revel_mgo delete model [modelname]\n")
			os.Exit(2)
		}

		sname := args[1]
		ColorLog("[INFO] Removing '%s' model\n", sname)
		deleteModel(sname, curpath)
	default:
		errorf("[ERRO] command is missing\n")
	}
	ColorLog("[SUCC] successfully deleted!\n")
}
