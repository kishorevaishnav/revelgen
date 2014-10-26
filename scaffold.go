package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	SCAFFOLD_NAME int = 2
)

var base_app_folder string

type sf_gorp_contStruct struct {
	Base_App_Folder  string
	ScaffoldName     string
	Lcase_sfname     string
	Ol               string // One Letter
	Pfn              string // Primary Field Name
	Fields           []Fields
	ValidationNeeded int
	RequiredArray    []string
	MinimumArray     []string
	MaximumArray     []string
}

func scaffoldRevel() {
	println("you are in scaffold Revel")
	cwd, _ := os.Getwd()
	base_app_folder = path.Base(cwd)
	log.Println(base_app_folder)
	lf, ra, mia, maa, pfn := gnModel_return(SCAFFOLD_NAME)

	abc := &sf_gorp_contStruct{
		Base_App_Folder:  base_app_folder,
		ScaffoldName:     strings.Title(os.Args[SCAFFOLD_NAME]),
		Lcase_sfname:     strings.ToLower(os.Args[SCAFFOLD_NAME]),
		Ol:               fmt.Sprintf(strings.ToLower(string(os.Args[SCAFFOLD_NAME][0]))),
		Pfn:              strings.ToLower(pfn),
		Fields:           lf,
		ValidationNeeded: len(ra),
		RequiredArray:    ra,
		MinimumArray:     mia,
		MaximumArray:     maa,
	}
	funcMap := template.FuncMap{
		"title": strings.Title,
		// "firstLetterLower": func(a string) string { return fmt.Sprintf(strings.ToLower(string(a[0]))) },
		"formatFieldName": func(a string) string { return fmt.Sprintf("%s", strings.Title(a)) },
		"formatDataType":  func(a string) string { return fmt.Sprintf("%-*s", 8, a) },
	}

	// Start Writing Controller File
	temp_data, _ := template_sf_gorp_controller_rvltpl()
	t, err := template.New("sf_gorp_controller.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	writeFile(os.Args[SCAFFOLD_NAME], buf, CONTR_FOL_PATH)

	// Start Writing Model File
	temp_data, _ = template_sf_gorp_model_rvltpl()
	t, err = template.New("sf_gorp_model.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	os.MkdirAll(MODEL_FOL_PATH, 0755)
	writeFile(os.Args[SCAFFOLD_NAME], buf, MODEL_FOL_PATH)

	// Start Writing Views File
	view_dir := VIEW_FOL_PATH + strings.ToLower(os.Args[SCAFFOLD_NAME]) + "/"
	log.Println(view_dir)
	os.MkdirAll(view_dir, 0775)
	temp_data, _ = template_sf_form_rvltpl()
	t, err = template.New("sf_form.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	log.Println(view_dir + "_form.html")
	if !fileExists(view_dir+"_form.html") || OVERWRITE_FILES {
		ioutil.WriteFile(view_dir+"_form.html", buf.Bytes(), 0644)
	} else {
		fmt.Println("view file already exists")
	}

	temp_data, _ = template_sf_edit_rvltpl()
	t, err = template.New("sf_edit.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	log.Println(view_dir + "edit.html")
	if !fileExists(view_dir+"edit.html") || OVERWRITE_FILES {
		ioutil.WriteFile(view_dir+"edit.html", buf.Bytes(), 0644)
	} else {
		fmt.Println("view file already exists")
	}

	temp_data, _ = template_sf_index_rvltpl()
	t, err = template.New("sf_index.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	log.Println(view_dir + "index.html")
	if !fileExists(view_dir+"index.html") || OVERWRITE_FILES {
		ioutil.WriteFile(view_dir+"index.html", buf.Bytes(), 0644)
	} else {
		fmt.Println("view file already exists")
	}

	temp_data, _ = template_sf_new_rvltpl()
	t, err = template.New("sf_new.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	log.Println(view_dir + "new.html")
	if !fileExists(view_dir+"new.html") || OVERWRITE_FILES {
		ioutil.WriteFile(view_dir+"new.html", buf.Bytes(), 0644)
	} else {
		fmt.Println("view file already exists")
	}

	temp_data, _ = template_sf_show_rvltpl()
	t, err = template.New("sf_show.rvltpl").Funcs(funcMap).Parse(string(temp_data))
	checkError(err)
	buf = new(bytes.Buffer)
	err = t.Execute(buf, abc)
	checkError(err)
	log.Println(view_dir + "show.html")
	if !fileExists(view_dir+"show.html") || OVERWRITE_FILES {
		ioutil.WriteFile(view_dir+"show.html", buf.Bytes(), 0644)
	} else {
		fmt.Println("view file already exists")
	}

	log.Println("Scaffold completed.")
}
