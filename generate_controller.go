package main

import(
	"os"
	"path"
	"strings"
)

// generate controller
func generateController(cname, crupath string) {
	// get controller name and package 
	p, f := path.Split(cname)

	// set controller name to uppercase
	controllerName := strings.Title(f)

	//set default package
	packageName := "controllers"
	if p != "" {
		i := strings.LastIndex(p[:len(p)-1], "/")
		packageName = p[i+1 : len(p)-1]
	}

	// get struct for controller 
	controllerStruct, err := GetControllerStruct(controllerName)
	if err != nil {
		ColorLog("[ERRO] Could not genrate controllers struct: %s\n", err)
		os.Exit(2)
	}

	ColorLog("[INFO] Using '%s' as controller name\n", controllerName)
	ColorLog("[INFO] Using '%s' as package name\n", packageName)

	// create controller folders
	filePath := path.Join(crupath ,"app", "controllers", p)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(filePath, 0777); err != nil {
			ColorLog("[ERRO] Could not create controllers directory: %s\n", err)
			os.Exit(2)
		}
	}

	// create common controller.go
	commonCtrFp := path.Join(crupath, "app", "controllers", "controller.go")
	if _, err := os.Stat(commonCtrFp); os.IsNotExist(err) {
		if cf, err := os.OpenFile(commonCtrFp, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
			defer cf.Close()	
			content := strings.Replace(commonTpl, "{{packageName}}", packageName, -1)		
			cf.WriteString(content)
			// gofmt generated source code
			FormatSourceCode(commonCtrFp)
			ColorLog("[INFO] controller file generated: %s\n", commonCtrFp)
		} else {
			// error creating file
			ColorLog("[ERRO] Could not create controller file: %s\n", err)
			os.Exit(2)
		}
	}

	mPath := path.Join(crupath, "app", "models", strings.ToLower(controllerName)+".go")
	if _, err := os.Stat(mPath); os.IsNotExist(err) {
		ColorLog("[ERRO] Could not find model file: %s\n", err)
		os.Exit(2)
	}
	// create controller file
	fpath := path.Join(filePath, strings.ToLower(controllerName)+".go")
	if f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer f.Close()

		paths := strings.Split(crupath, "/")
		projectName := paths[len(paths) - 1:][0]
		modelsPkg := path.Join(projectName, "app", "models")

		content := strings.Replace(controllerTpl, "{{packageName}}", packageName, -1)
		content = strings.Replace(content, "{{modelsPkg}}", modelsPkg, -1)
		content = strings.Replace(content, "{{controllerStruct}}", controllerStruct, -1)
		content = strings.Replace(content, "{{contorllerStructName}}", controllerName + "Controller", -1)
		content = strings.Replace(content, "{{modelObjects}}", strings.ToLower(controllerName+"s"), -1)
		content = strings.Replace(content, "{{modelObject}}", strings.ToLower(controllerName), -1)
		content = strings.Replace(content, "{{modelStruct}}", controllerName, -1)
		content = strings.Replace(content, "{{modelStructs}}", controllerName + "s", -1)
		
		
		f.WriteString(content)
		// gofmt generated source code
		FormatSourceCode(fpath)
		ColorLog("[INFO] model file generated: %s\n", fpath)
	} else {
		// error creating file
		ColorLog("[ERRO] Could not create controller file: %s\n", err)
		os.Exit(2)
	}
}

// delete controller
func deleteController(cname, crupath string) {
	_, f := path.Split(cname)
	controllerName := strings.Title(f)
	filePath := path.Join(crupath, "app", "controllers", controllerName + ".go")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err = os.Remove(filePath)
		if err != nil{
			ColorLog("[ERRO] Could not delete controller struct: %s\n", err)
			os.Exit(2)	
		}
		ColorLog("[INFO] controller file deleted: %s\n", filePath)
		
	}

}


var commonTpl = `package {{packageName}}
import (
	"strconv"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type CtrlErr map[string]interface{}

func parseUintOrDefault(intStr string, _default uint64) uint64 {
    if value, err := strconv.ParseUint(intStr, 0, 64); err != nil {
        return _default
    } else {
        return value
    }
}

func parseIntOrDefault(intStr string, _default int64) int64 {
    if value, err := strconv.ParseInt(intStr, 0, 64); err != nil {
        return _default
    } else {
        return value
    }
}

func convertToObjectIdHex(id string) (result bson.ObjectId, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Unable to convert %v to object id", id)
		}
	}()

	return bson.ObjectIdHex(id), err
}

func buildErrResponse(err error, errorCode string) CtrlErr {
	ctrlErr := CtrlErr{}
	ctrlErr["error_message"] = err.Error()
	ctrlErr["error_code"] = errorCode
	return ctrlErr
}


`
var controllerTpl = `package {{packageName}}

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"github.com/revel/revel"
	"encoding/json"
	"{{modelsPkg}}"
)

{{controllerStruct}}

func (c {{contorllerStructName}}) Index() revel.Result {  
	var (
		{{modelObjects}} []models.{{modelStruct}}
		err error
	)
	{{modelObjects}}, err = models.Get{{modelStructs}}()
	if err != nil{
		errResp := buildErrResponse(err,"500")
		c.Response.Status = 500
		return c.RenderJSON(errResp)
	}
	c.Response.Status = 200
    return c.RenderJSON({{modelObjects}})
}

func (c {{contorllerStructName}}) Show(id string) revel.Result {  
    var (
    	{{modelObject}} models.{{modelStruct}}
    	err error
    	{{modelObject}}ID bson.ObjectId 
    )

    if id == ""{
    	errResp := buildErrResponse(errors.New("Invalid {{modelObject}} id format"),"400")
    	c.Response.Status = 400
    	return c.RenderJSON(errResp)
    }

    {{modelObject}}ID, err = convertToObjectIdHex(id)
    if err != nil{
    	errResp := buildErrResponse(errors.New("Invalid {{modelObject}} id format"),"400")
    	c.Response.Status = 400
    	return c.RenderJSON(errResp)
    }

    {{modelObject}}, err = models.Get{{modelStruct}}({{modelObject}}ID)
    if err != nil{
    	errResp := buildErrResponse(err,"500")
    	c.Response.Status = 500
    	return c.RenderJSON(errResp)
    }
  
    c.Response.Status = 200
    return c.RenderJSON({{modelObject}})
}

func (c {{contorllerStructName}}) Create() revel.Result {  
    var (
    	{{modelObject}} models.{{modelStruct}}
    	err error
    )

    err = c.Params.BindJSON(&{{modelObject}})
	if err != nil {
		errResp := buildErrResponse(err, "403")
		c.Response.Status = 403
		return c.RenderJSON(errResp)
	}

	{{modelObject}}, err = models.Add{{modelStruct}}({{modelObject}})
	if err != nil{
		errResp := buildErrResponse(err,"500")
    	c.Response.Status = 500
    	return c.RenderJSON(errResp)
	}
    c.Response.Status = 201
    return c.RenderJSON({{modelObject}})
}

func (c {{contorllerStructName}}) Update() revel.Result {  
	var (
    	{{modelObject}} models.{{modelStruct}}
    	err error
    )
    err = c.Params.BindJSON(&{{modelObject}})
	if err != nil{
		errResp := buildErrResponse(err,"400")
    	c.Response.Status = 400
    	return c.RenderJSON(errResp)
	}

	err = {{modelObject}}.Update{{modelStruct}}()
	if err != nil{
		errResp := buildErrResponse(err,"500")
    	c.Response.Status = 500
    	return c.RenderJSON(errResp)
	}
    return c.RenderJSON({{modelObject}})
}

func (c {{contorllerStructName}}) Delete(id string) revel.Result { 
	var (
    	err error
    	{{modelObject}} models.{{modelStruct}}
    	{{modelObject}}ID bson.ObjectId 
    )
     if id == ""{
    	errResp := buildErrResponse(errors.New("Invalid {{modelObject}} id format"),"400")
    	c.Response.Status = 400
    	return c.RenderJSON(errResp)
    }

    {{modelObject}}ID, err = convertToObjectIdHex(id)
    if err != nil{
    	errResp := buildErrResponse(errors.New("Invalid {{modelObject}} id format"),"400")
    	c.Response.Status = 400
    	return c.RenderJSON(errResp)
    }

    {{modelObject}}, err = models.Get{{modelStruct}}({{modelObject}}ID)
    if err != nil{
    	errResp := buildErrResponse(err,"500")
    	c.Response.Status = 500
    	return c.RenderJSON(errResp)
    }
	err = {{modelObject}}.Delete{{modelStruct}}()
	if err != nil{
		errResp := buildErrResponse(err,"500")
    	c.Response.Status = 500
    	return c.RenderJSON(errResp)
	} 
	c.Response.Status = 204
    return c.RenderJSON(nil)
}
`

