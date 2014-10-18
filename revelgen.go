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
	SUB_ACTION      int    = 2
	CONTROLLER_NAME int    = 3
	MODEL_NAME      int    = 3 // 3rd Argument
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
	case "gc":
		generateController()
	case "gm":
		generateModel()
	case "gr":
		generateRoute()
	case "s":
		scaffoldRevel()
	case "generate", "g":
		switch os.Args[SUB_ACTION] {
		case "controller", "c":
			generateController()
		case "model", "m":
			generateModel()
		case "route", "r":
			generateRoute()
		case "scaffold", "s":
			scaffoldRevel()
		default:
			panic("No SubActions provided")
		}
	default:
		panic("No options provided")
	}
}

type contStruct struct {
	AppName    string
	MethodName []string
}

func generateController() {
	fmt.Println("you are in generateController")
	contValue := &contStruct{
		AppName:    strings.Title(os.Args[CONTROLLER_NAME]),
		MethodName: []string{os.Args[4]},
	}
	p, err := load_parse_ControllerTemplate("controller", contValue)
	checkError(err)
	if !fileExists(CONTR_FOL_PATH + os.Args[CONTROLLER_NAME] + ".go") {
		ioutil.WriteFile("app/controllers/"+os.Args[CONTROLLER_NAME]+".go", p.Bytes(), 0644)
		fmt.Println("...completed.")
	} else {
		fmt.Println("file already exists")
	}
}

func generateModel() {
	fmt.Println("you are in generateModel")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
func generateRoute() {
	fmt.Println("you are in generateRoute")
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func scaffoldRevel() {
	println("you are in scaffold Revel")
}
