/*
 * Copyright (c) 2018. Philipp Kalytta
 */

package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"html/template"
	"io/ioutil"
)

func cssHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	content, err := ioutil.ReadFile("static/css" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/css")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(content))
}
func jsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	content, err := ioutil.ReadFile("static/js" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/javascript")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(content))
}

func indexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "index.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "login")
}

func cartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "cart")
}

func adminHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "admin")
}

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/index.html", indexHandler)
	router.GET("/login/*suburl", loginHandler)
	router.GET("/cart/*suburl", cartHandler)
	//Find closest match for /admin, so that all admin sites are handled by adminHandler
	router.GET("/admin/*suburl", adminHandler)
	router.GET("/css/*suburl", cssHandler)
	router.GET("/js/*suburl", jsHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
