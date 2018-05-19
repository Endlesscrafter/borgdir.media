/*
 * Copyright (c) 2018. Philipp Kalytta
 */

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"html/template"
	"io/ioutil"
	"fmt"
)

func cssHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	content, err := ioutil.ReadFile("static/css" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/css")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
}
func jsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	content, err := ioutil.ReadFile("static/js" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/javascript")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
	//fmt.Fprintf(w, string(content))
}

func indexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tmpl, err := template.ParseFiles("template/index.html")
	//Content-Type text/html doesnt fix malformatted css
	w.Header().Set("Content-Type", "text/html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "index.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tmpl, err := template.ParseFiles("template/login.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "login.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func registerHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tmpl, err := template.ParseFiles("template/register.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "register.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func cartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tmpl, err := template.ParseFiles("template/cart.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "cart.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func equipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tmpl, err := template.ParseFiles("template/equipment.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "equipment.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func myEquipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tmpl, err := template.ParseFiles("template/my-equipment.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "my-equipment.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}
func profileHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tmpl, err := template.ParseFiles("template/profile.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "profile.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func adminHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	fmt.Println(params.ByName("suburl"))
	if params.ByName("suburl") == "/index.html" || params.ByName("suburl") == "/" {
		tmpl, err := template.ParseFiles("template/admin/index.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "index.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/add.html"{
		tmpl, err := template.ParseFiles("template/admin/add.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "add.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/clients.html"{
		tmpl, err := template.ParseFiles("template/admin/clients.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "clients.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/edit-client.html"{
		tmpl, err := template.ParseFiles("template/admin/edit-client.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "edit-client.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/equipment.html"{
		tmpl, err := template.ParseFiles("template/admin/equipment.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "equipment.html", "DATATATATATA")
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/index.html", indexHandler)
	router.GET("/login.html", loginHandler)
	router.GET("/register.html", registerHandler)
	router.GET("/cart.html", cartHandler)
	router.GET("/equipment.html", equipHandler)
	router.GET("/my-equipment.html", myEquipHandler)
	router.GET("/profile.html", profileHandler)
	//Find closest match for /admin, so that all admin sites are handled by adminHandler
	router.GET("/admin/*suburl", adminHandler)
	router.GET("/css/*suburl", cssHandler)
	router.GET("/js/*suburl", jsHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
