package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const (
	ACTION          int    = 1 // 2nd Argument
	CONTROLLER_NAME int    = 2
	MODEL_NAME      int    = 2 // 3rd Argument
	CONTR_FOL_PATH  string = "app/controllers/"
	MODEL_FOL_PATH  string = "app/models/"
	ROUTE_FOL_PATH  string = "conf/"
	OVERWRITE_FILES bool   = true
	DEBUG           bool   = true
)

var (
	max_field_name_length int
	// validationNeeded      bool = false
)

func main() {
	if os.Args[0] != "revelgen" {
		fmt.Println("wrong usage")
		os.Exit(1)
	}
	switch os.Args[ACTION] {
	case "controller", "c":
		generateController()
	case "model", "m":
		generateModel()
	case "route", "r":
		updateRoute()
	case "scaffold", "s":
		scaffoldRevel()
	default:
		panic("No actions provided")
	}
	os.Exit(1)
}

type contStruct struct {
	ControllerName string
	MethodNames    []string
}

func generateController() {
	fmt.Println("you are in generateController")
	contValue := &contStruct{
		ControllerName: strings.Title(os.Args[CONTROLLER_NAME]),
		MethodNames:    os.Args[CONTROLLER_NAME+1 : len(os.Args)],
	}
	p, err := load_parse_ControllerTemplate("controller", contValue)
	checkError(err)
	writeFile(os.Args[CONTROLLER_NAME], p, CONTR_FOL_PATH)
}

func load_parse_ControllerTemplate(title string, contValue *contStruct) (*bytes.Buffer, error) {
	filename := "./template/" + title + ".rvltpl"
	t, err := template.ParseFiles(filename)
	checkError(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, contValue)
	checkError(err)
	return buf, nil
}

type modelStruct struct {
	ModelName        string
	Fields           []Fields
	ValidationNeeded int
	ValidationArray  []string
}

type Fields struct {
	Name      string
	Datatype  string
	Db_data   string
	Json_data string
}

func generateModel() {
	fmt.Println("you are in generateModel")
	var primaryField bool = false
	var requiredArray []string
	var lineFields []Fields
	fieldArray := os.Args[MODEL_NAME+1 : len(os.Args)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ")
		fieldSplit := strings.Split(fieldArray[key], ":")
		switch fieldSplit[1] {
		case "int", "int64", "string", "varchar", "text":
		default:
			fmt.Println("wrong data type", fieldSplit[0], ":", fieldSplit[1])
			os.Exit(1)
		}
		if string(fieldSplit[0][len(fieldSplit[0])-1]) == "*" {
			fieldSplit[0] = fieldSplit[0][0 : len(fieldSplit[0])-1]
			requiredArray = append(requiredArray, fieldSplit[0])
		}
		if max_field_name_length < (strings.Count(fieldSplit[0], "") - 1) {
			max_field_name_length = strings.Count(fieldSplit[0], "") - 1
		}
		lineFields = append(lineFields, Fields{Name: fieldSplit[0], Datatype: fieldSplit[1], Db_data: fieldSplit[0], Json_data: fieldSplit[0]})
		if strings.ToLower(os.Args[MODEL_NAME])+"id" == strings.ToLower(fieldSplit[0]) {
			primaryField = true
		}
	}

	// check if the primary_key can be added or not.
	if false == primaryField {
		a := Fields{Name: strings.Title(os.Args[MODEL_NAME]) + "Id", Datatype: "int", Db_data: strings.Title(os.Args[MODEL_NAME]) + "Id", Json_data: strings.Title(os.Args[MODEL_NAME]) + "Id"}
		lineFields = append([]Fields{a}, lineFields...)
		if max_field_name_length < (strings.Count(os.Args[MODEL_NAME]+"Id", "") - 1) {
			max_field_name_length = strings.Count(os.Args[MODEL_NAME]+"Id", "") - 1
		}
	}

	// fmt.Println(max_field_name_length)
	// for _, v := range lineFields {
	// 	fmt.Println(v.Name, v.Datatype, v.Db_data, v.Json_data)
	// }
	modelValue := &modelStruct{
		ModelName:        strings.Title(os.Args[MODEL_NAME]),
		Fields:           lineFields,
		ValidationNeeded: len(requiredArray),
		ValidationArray:  requiredArray,
	}
	p, err := load_parse_ModelTemplate("model", modelValue)
	checkError(err)
	writeFile(os.Args[MODEL_NAME], p, MODEL_FOL_PATH)
}

func load_parse_ModelTemplate(title string, modelValue *modelStruct) (*bytes.Buffer, error) {
	funcMap := template.FuncMap{
		"title":            strings.Title,
		"formatFieldName":  func(a string) string { return fmt.Sprintf("%-*s", max_field_name_length, strings.Title(a)) },
		"formatDataType":   func(a string) string { return fmt.Sprintf("%-*s", 8, a) },
		"firstLetterLower": func(a string) string { return fmt.Sprintf(strings.ToLower(string(a[0]))) },
	}
	filename := "./template/" + title + ".rvltpl"
	t, err := template.New("model.rvltpl").Funcs(funcMap).ParseFiles(filename)
	checkError(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, modelValue)
	checkError(err)
	return buf, nil
}

func updateRoute() {
	fmt.Println("you are in updateRoute")
}

func scaffoldRevel() {
	generateController()
	generateModel()
	updateRoute()
	println("you are in scaffold Revel")
}

func writeFile(name string, content *bytes.Buffer, write_path string) {
	filename := write_path + name + ".go"
	if !fileExists(filename) || OVERWRITE_FILES {
		ioutil.WriteFile(filename, content.Bytes(), 0644)
		fmt.Println("...completed.")
	} else {
		fmt.Println("file already exists")
	}
	if DEBUG {
		fmt.Println("------------ DEBUG ------------")
		dat, err := ioutil.ReadFile(filename)
		checkError(err)
		fmt.Println(string(dat))
		fmt.Println("============ DEBUG ============")
	}
	return
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}