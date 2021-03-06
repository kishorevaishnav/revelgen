package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

const (
	ACTION          int = 1 // 2nd Argument
	CONTROLLER_NAME int = 2
	MODEL_NAME      int = 2 // 3rd Argument

	CONTR_FOL_PATH string = "app/controllers/"
	MODEL_FOL_PATH string = "app/models/"
	VIEW_FOL_PATH  string = "app/views/"
	ROUTE_FOL_PATH string = "conf/"

	OVERWRITE_FILES bool = true
	DEBUG           bool = true

	FIELD_DATATYPE_REGEXP = `^([A-Za-z]{2,15})([\*\-]{0,1}):([A-Za-z]{2,15})(\((\d{0,2}),(\d{0,2})\)){0,1}$`
)

var (
	max_field_name_length int
	// validationNeeded      bool = false
	allowed_datatype = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64", "int", "int8", "int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}
)

type contStruct struct {
	ControllerName string
	MethodNames    []string
}

type modelStruct struct {
	ModelName        string
	Fields           []Fields
	ValidationNeeded int
	RequiredArray    []string
	MinimumArray     []string
	MaximumArray     []string
}

type Fields struct {
	Name      string
	Datatype  string
	Db_data   string
	Json_data string
}

type viewStruct struct {
	Name     string
	FilePath string
}

func main() {
	if os.Args[0] != "revelgen" {
		fmt.Println("wrong usage")
		os.Exit(1)
	}
	switch os.Args[ACTION] {
	case "controller", "c":
		generateController()
		generateViews()
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

func generateController() {
	contValue := &contStruct{
		ControllerName: strings.Title(os.Args[CONTROLLER_NAME]),
		MethodNames:    os.Args[CONTROLLER_NAME+1 : len(os.Args)],
	}
	p, err := load_parse_ControllerTemplate("controller", contValue)
	checkError(err)
	writeFile(os.Args[CONTROLLER_NAME], p, CONTR_FOL_PATH)
}

func generateViews() {
	view_dir := VIEW_FOL_PATH + strings.ToLower(os.Args[CONTROLLER_NAME]) + "/"
	os.Mkdir(view_dir, 0775)
	temp_data, _ := template_view_rvltpl()
	for _, v := range os.Args[CONTROLLER_NAME+1 : len(os.Args)] {
		viewValue := &viewStruct{
			Name:     strings.Title(os.Args[CONTROLLER_NAME]) + "#" + strings.ToLower(v),
			FilePath: strings.ToLower(view_dir + v + ".html"),
		}
		t, err := template.New("view.rvltpl").Parse(string(temp_data))
		checkError(err)
		buf := new(bytes.Buffer)
		err = t.Execute(buf, viewValue)
		checkError(err)
		if !fileExists(view_dir+v+".html") || OVERWRITE_FILES {
			ioutil.WriteFile(view_dir+v+".html", buf.Bytes(), 0644)
		} else {
			fmt.Println("view file already exists")
		}
	}
}

func load_parse_ControllerTemplate(title string, contValue *contStruct) (*bytes.Buffer, error) {
	funcMap := template.FuncMap{
		"title":            strings.Title,
		"firstLetterLower": func(a string) string { return fmt.Sprintf(strings.ToLower(string(a[0]))) },
	}
	temp_data, _ := template_controller_rvltpl()
	t, err := template.New("controller.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, contValue)
	checkError(err)
	return buf, nil
}

func gnModel_return(mn int) ([]Fields, []string, []string, []string, string) {
	var primaryField bool = false
	var createdField bool = false
	var updatedField bool = false
	var requiredArray []string
	var minimumArray []string
	var maximumArray []string
	var lineFields []Fields
	var primaryFieldName string
	fieldArray := os.Args[mn+1 : len(os.Args)]
	for key, value := range fieldArray {
		fieldArray[key] = strings.Trim(value, ", ") // TODO - Need to check why I added "," instead of just " ".
		fds, err := fld_dtype_sep(fieldArray[key])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f_name, f_data_type, f_required, f_min, f_max := fds[0], fds[1], fds[2], fds[3], fds[4]

		if f_required != "" {
			requiredArray = append(requiredArray, f_name)
		}
		if f_min != "" {
			minimumArray = append(minimumArray, f_name)
		}
		if f_max != "" {
			maximumArray = append(maximumArray, f_name)
		}
		if max_field_name_length < (strings.Count(f_name, "") - 1) {
			max_field_name_length = strings.Count(f_name, "") - 1
		}
		lineFields = append(lineFields, Fields{Name: f_name, Datatype: f_data_type, Db_data: f_name, Json_data: f_name})
		switch strings.ToLower(f_name) {
		case strings.ToLower(os.Args[mn] + "id"), "id":
			primaryFieldName = strings.ToLower(f_name)
			primaryField = true
		case "updated":
			updatedField = true
		case "created":
			createdField = true
		}
	}

	// check if the primary_key can be added or not.
	if false == primaryField {
		a := Fields{Name: strings.Title(os.Args[mn]) + "id", Datatype: "int", Db_data: strings.Title(os.Args[mn]) + "id", Json_data: strings.Title(os.Args[mn]) + "id"}
		primaryFieldName = strings.ToLower(os.Args[mn] + "id")
		lineFields = append([]Fields{a}, lineFields...)
		if max_field_name_length < (strings.Count(os.Args[mn]+"id", "") - 1) {
			max_field_name_length = strings.Count(os.Args[mn]+"id", "") - 1
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
	return lineFields, requiredArray, minimumArray, maximumArray, primaryFieldName
}

func generateModel() {
	lf, ra, mia, maa, _ := gnModel_return(MODEL_NAME)
	modelValue := &modelStruct{
		ModelName:        strings.Title(os.Args[MODEL_NAME]),
		Fields:           lf,
		ValidationNeeded: len(ra),
		RequiredArray:    ra,
		MinimumArray:     mia,
		MaximumArray:     maa,
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
	os.MkdirAll(MODEL_FOL_PATH, 0755)
	err = t.Execute(buf, modelValue)
	checkError(err)
	return buf, nil
}

func updateRoute() {
	fmt.Println("you are in updateRoute")
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

func fld_dtype_sep(orig_string string) (parsed_string []string, err error) {
	r, _ := regexp.Compile(FIELD_DATATYPE_REGEXP)
	if r.MatchString(orig_string) == true {
		split := r.FindAllStringSubmatch(orig_string, -1)
		valid_datatype := false
		for _, v := range allowed_datatype {
			if v == split[0][3] {
				valid_datatype = true
			}
		}
		if valid_datatype {
			return []string{split[0][1], split[0][3], split[0][2], split[0][5], split[0][6]}, nil
		} else {
			err := errors.New("Wrong Datatype " + orig_string)
			return nil, err
		}
	} else {
		log.Println(orig_string)
		err := errors.New("Wrong Format " + orig_string)
		return nil, err
	}
}
