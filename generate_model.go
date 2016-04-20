package main

import(
	"os"
	"path"
	"strings"
)

func generateModel(mname, fields, crupath string) {
	p, f := path.Split(mname)
	modelName := strings.Title(f)
	packageName := "models"
	if p != "" {
		i := strings.LastIndex(p[:len(p)-1], "/")
		packageName = p[i+1 : len(p)-1]
	}
	modelStruct, err := GetStruct(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not genrate models struct: %s\n", err)
		os.Exit(2)
	}
	ColorLog("[INFO] Using '%s' as model name\n", modelName)
	ColorLog("[INFO] Using '%s' as package name\n", packageName)


	filePath := path.Join(crupath ,"app", "models", p)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(filePath, 0777); err != nil {
			ColorLog("[ERRO] Could not create models directory: %s\n", err)
			os.Exit(2)
		}
	}
	fpath := path.Join(filePath, strings.ToLower(modelName)+".go")
	if f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer f.Close()
		paths := strings.Split(crupath, "/")
		projectName := paths[len(paths) - 1:][0]
		mongoPkg := path.Join(projectName, "app", "models", "mongodb")
		collectionFuncName := "new" + string(modelName) + "Collection()"
		updatedStr, err := GetAttrs(fields)
		if err != nil {
			ColorLog("[ERRO] Could not genrate models struct: %s\n", err)
			os.Exit(2)
		}
		content := strings.Replace(modelTpl, "{{packageName}}", packageName, -1)
		content = strings.Replace(content, "{{mongoPkg}}", mongoPkg, -1)
		content = strings.Replace(content, "{{modelName}}", modelName, -1)
		content = strings.Replace(content, "{{modelStruct}}", modelStruct, -1)
		content = strings.Replace(content, "{{collectionFuncName}}", collectionFuncName, -1)
		content = strings.Replace(content, "{{collectionName}}", strings.ToLower(modelName) + "s", -1)
		content = strings.Replace(content, "{{modelObject}}", strings.ToLower(modelName), -1)
		content = strings.Replace(content, "{{listModelName}}",  strings.ToLower(modelName) + "s", -1)
		content = strings.Replace(content, "{{sortFields}}",  "\"-createdAt\"", -1)
		content = strings.Replace(content, "{{updatedData}}",  updatedStr, -1)
		
		
		f.WriteString(content)
		// gofmt generated source code
		FormatSourceCode(fpath)
		ColorLog("[INFO] model file generated: %s\n", fpath)
	} else {
		// error creating file
		ColorLog("[ERRO] Could not create model file: %s\n", err)
		os.Exit(2)
	}
}

func deleteModel(mname, crupath string) {
	_, f := path.Split(mname)
	modelName := strings.Title(f)
	filePath := path.Join(crupath, "app", "models", modelName + ".go")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err = os.Remove(filePath)
		if err != nil{
			ColorLog("[ERRO] Could not delete model struct: %s\n", err)
			os.Exit(2)	
		}
		ColorLog("[INFO] model file deleted: %s\n", filePath)
		
	}

}

var modelTpl = `package {{packageName}}

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"{{mongoPkg}}"
)

{{modelStruct}}

func {{collectionFuncName}}  *mongodb.Collection  {
   return mongodb.NewCollectionSession("{{collectionName}}")
}

// Add{{modelName}} insert a new {{modelName}} into database and returns
// last inserted {{modelObject}} on success.
func Add{{modelName}}(m {{modelName}}) ({{modelObject}} {{modelName}}, err error) {
	c := {{collectionFuncName}}
	defer c.Close()
	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()
	return m, c.Session.Insert(m)
}

// Update{{modelName}} update a {{modelName}} into database and returns
// last nil on success.
func (m {{modelName}}) Update{{modelName}}() error{
	c := {{collectionFuncName}}
	defer c.Close()
	
	err := c.Session.Update(bson.M{
		"_id": m.ID,
	}, bson.M{
		"$set": {{updatedData}}
	})
	return err
}

// Delete{{modelName}} Delete {{modelName}} from database and returns
// last nil on success.
func (m {{modelName}}) Delete{{modelName}}() error{
	c := {{collectionFuncName}}
	defer c.Close()

	err := c.Session.Remove(bson.M{"_id": m.ID})
	return err
}

// Get{{modelName}}s Get all {{modelName}} from database and returns
// list of {{modelName}} on success
func Get{{modelName}}s() ([]{{modelName}}, error) {
	var (
		{{listModelName}} []{{modelName}}
		err   error
	)

	c := {{collectionFuncName}}
	defer c.Close()

	err = c.Session.Find(nil).Sort({{sortFields}}).All(&{{listModelName}})
	return {{listModelName}}, err
}

// Get{{modelName}} Get a {{modelName}} from database and returns
// a {{modelName}} on success
func Get{{modelName}}(id bson.ObjectId) ({{modelName}}, error) {
	var (
		{{modelObject}} {{modelName}}
		err   error
	)

	c := {{collectionFuncName}}
	defer c.Close()


	err = c.Session.Find(bson.M{"_id": id}).One(&{{modelObject}})
	return {{modelObject}}, err
}
`

