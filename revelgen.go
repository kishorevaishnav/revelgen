package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	ACTION          int = 1 // 2nd Argument
	CONTROLLER_NAME int = 2
	MODEL_NAME      int = 2 // 3rd Argument

	CONTR_FOL_PATH string = "app/controllers/"
	MODEL_FOL_PATH string = "app/models/"
	ROUTE_FOL_PATH string = "conf/"

	OVERWRITE_FILES bool = true
	DEBUG           bool = true
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
	fmt.Println("")
	generateViews()
}

func generateViews() {
	view_dir := "app/views/" + strings.ToLower(os.Args[CONTROLLER_NAME])
	os.Mkdir(view_dir, 0775)
	for _, v := range os.Args[CONTROLLER_NAME+1 : len(os.Args)] {
		os.Create(view_dir + "/" + v + ".html")
	}
}
func load_parse_ControllerTemplate(title string, contValue *contStruct) (*bytes.Buffer, error) {
	temp_data, _ := template_controller_rvltpl()
	t, err := template.New("controller.rvltpl").Parse(string(temp_data))
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
	var createdField bool = false
	var updatedField bool = false
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
		switch strings.ToLower(fieldSplit[0]) {
		case strings.ToLower(os.Args[MODEL_NAME] + "id"), "id":
			primaryField = true
		case "updated":
			updatedField = true
		case "created":
			createdField = true
		}
		// if strings.ToLower(os.Args[MODEL_NAME])+"id" == strings.ToLower(fieldSplit[0]) {
		// 	primaryField = true
		// }
	}

	// check if the primary_key can be added or not.
	if false == primaryField {
		a := Fields{Name: strings.Title(os.Args[MODEL_NAME]) + "Id", Datatype: "int", Db_data: strings.Title(os.Args[MODEL_NAME]) + "id", Json_data: strings.Title(os.Args[MODEL_NAME]) + "id"}
		lineFields = append([]Fields{a}, lineFields...)
		if max_field_name_length < (strings.Count(os.Args[MODEL_NAME]+"Id", "") - 1) {
			max_field_name_length = strings.Count(os.Args[MODEL_NAME]+"Id", "") - 1
		}
	}

	// check if the created field can be added or not.
	if false == createdField {
		a := Fields{Name: "Created", Datatype: "int64", Db_data: "created", Json_data: "created"}
		lineFields = append(lineFields, a)
		if max_field_name_length < 7 {
			max_field_name_length = 7
		}
	}

	// check if the updated field can be added or not.
	if false == updatedField {
		a := Fields{Name: "Updated", Datatype: "int64", Db_data: "updated", Json_data: "updated"}
		lineFields = append(lineFields, a)
		if max_field_name_length < 7 {
			max_field_name_length = 7
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
	temp_data, _ := template_model_rvltpl()
	t, err := template.New("model.rvltpl").Funcs(funcMap).Parse(string(temp_data))
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
	filename := strings.ToLower(write_path + name + ".go")
	if !fileExists(filename) || OVERWRITE_FILES {
		ioutil.WriteFile(filename, content.Bytes(), 0644)
		exec.Command("go", "fmt", filename).Output()
		exec.Command("goimports", "-w=true", filename).Output()
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
