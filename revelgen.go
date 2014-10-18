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
		MethodNames:    os.Args[4:len(os.Args)],
	}
	println(os.Args[4], os.Args[5], os.Args[6])
	p, err := load_parse_ControllerTemplate("controller", contValue)
	checkError(err)
	writeFile(os.Args[CONTROLLER_NAME], p)
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

func generateModel() {
	fmt.Println("you are in generateModel")
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

func writeFile(name, content) {
	file_name := CONTR_FOL_PATH + name + ".go"
	if !fileExists(file_name) {
		ioutil.WriteFile(file_name, content.Bytes(), 0644)
		fmt.Println("...completed.")
	} else {
		fmt.Println("file already exists")
	}

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
