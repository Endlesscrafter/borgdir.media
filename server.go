/*
 * Copyright (c) 2018. Philipp Kalytta
 */

package main

import (
	"net/http"
	"log"
	"html/template"
	"io/ioutil"
)

//Defines a Webpage with a Title and a Body
/*type Page struct {
	Title string
	Body  []byte
}

//Saves a given Page to a TXT-Document
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//Loads a given Page by Name from a TXT-Document
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}*/

//Handles all view requests after /view/, offers the TXT-Document as HTML
/*func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	//Redirects non-existing pages to th edit site:
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

//Handles all requests after /edit/, offers a form to the HTTPClient
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}*/

//Renders the Webpage from a given template file
/*func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

type WebPage struct {
	htmlDocument []byte
}

//Entry Point for the Caller
func main() {
	//http.HandleFunc("/view/", viewHandler)
	//http.HandleFunc("/edit/", editHandler)

	//Not jet implemented
	//http.HandleFunc("/save/", saveHandler)

	http.HandleFunc("/admin/",adminHandler)
	http.HandleFunc("/", standardHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func adminHandler(writer http.ResponseWriter, request *http.Request) {
	
}

func standardHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path
	println("IS TITLE INDEX?: " + title=="/")
	if(title == "/"){
		title = "index"
	}
	println("TITLE: " + title)
	webpage, err := loadPage(title)
	//Redirects non-existing pages to the index site:
	if err != nil {
		http.Redirect(writer, request, "index", http.StatusFound)
		return
	}
	println("RENDER TEMPLATE: " + "pulbic/" + title)
	renderTemplate(writer, "public/"+title , webpage)
}

func renderTemplate(w http.ResponseWriter, templatePath string, wp *WebPage){


	println("RENDER TEMPLATE: " + templatePath + ".html")
	t, err := template.ParseFiles(templatePath + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, wp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}


//Loads a page from the given path and returns it as a pointer to a WebPage variable
func loadPage(path string) (*WebPage, error) {
	filename := "public/" + path + ".html"
	document, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &WebPage{htmlDocument:document}, nil
}