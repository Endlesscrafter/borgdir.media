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
	"strings"
)


type equipmentData struct {
	Name string
	Desc string
	ImageSRC string
	ImageAlt string
	Stock string
	StockAmount int
	Category string
	Featured bool
	FeaturedID int
	FeaturedImageSRC string
	Rented bool
	Bookmarked bool
	Repair bool
	RentedByUserID int64
	RentedByUserName string
	RentDate string
	ReturnDate string
	InvID int64
	StorageLocation string
	equipmentOwnerID int64
}

type user struct{
	UserID int64
	Name string
	Email string
	Password string
	ProfileImageSRC string
	UserLevel string
	Blocked bool
	ActiveUntilDate string
}

type siteData struct {
	Equipment []equipmentData
	User user
	AdminUserList []user
}

//Test Data
var equip1 = equipmentData{"Kamera Obscura", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder" ,"Entliehen",0 ,"None", true, 1,"/img/equipment/gandalf.gif",true, true, false ,1,"Max Mustermann" ,"25.05.18", "25.05.18", 1245, "Baungasse",2}
var equip2 = equipmentData{"Stativ", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder" ,"Verfügbar",2 , "None", true, 2,"/img/equipment/gandalf.gif", true, false, false ,1,"Karl Karstens" ,"25.05.18", "25.05.18", 13452 ,"Schrank",2}
var equip3 = equipmentData{"Mikrophon", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder" ,"Vorgemerkt",1, "None" , true, 3, "/img/equipment/gandalf.gif", true, false, false,1,"Max Mustermann" ,"25.05.18", "25.05.18", 2374, "Regal 12",2}

var user1 = user{1,"Max Mustermann", "max@muster.de", "asdf", "img/equipment/gandalf.gif", "Benutzer", false, "25.07.18"}
var user2 = user{2,"Peter Müller", "peter@mueller.de", "asdf", "img/equipment/gandalf.gif", "Administrator" ,false, "25.07.22"}


func cssHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("static/css" + params.ByName("suburl"))
	content, err := ioutil.ReadFile("static/css" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/css")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
}
func jsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("static/js" + params.ByName("suburl"))
	content, err := ioutil.ReadFile("static/js" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/javascript")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
	//fmt.Fprintf(w, string(content))
}

func imgHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("app/model/images/" + params.ByName("suburl"))
	content, err := ioutil.ReadFile("app/model/images/" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "image/*")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
	//fmt.Fprintf(w, string(content))
}

func indexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/index.html")
	//Content-Type text/html doesnt fix malformatted css
	w.Header().Set("Content-Type", "text/html")
	if err == nil {

		data := siteData{}
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)

		tmpl.ExecuteTemplate(w, "index.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/login.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "login.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func registerHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/register.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "register.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func cartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/cart.html")
	if err == nil {

		data := siteData{}
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)
		data.Equipment = append(data.Equipment, equip3)
		data.User = user1

		tmpl.ExecuteTemplate(w, "cart.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func equipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/equipment.html")
	if err == nil {

		data := siteData{}
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)
		data.Equipment = append(data.Equipment, equip3)
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip1)
		data.User = user1

		//Always give a even number of data-Sets (0, 2, 4, 6 etc): Otherwise it breaks the Design
		tmpl.ExecuteTemplate(w, "equipment.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func myEquipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/my-equipment.html")
	if err == nil {

		data := siteData{}
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)
		data.User = user1

		tmpl.ExecuteTemplate(w, "my-equipment.html", data)
	}
	//fmt.Fprintf(w, "index")
}
func profileHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println(params.ByName("suburl"))
	tmpl, err := template.ParseFiles("template/profile.html")
	if err == nil {

		data := siteData{}
		data.User = user1
		tmpl.ExecuteTemplate(w, "profile.html", data)
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

			data := siteData{}
			data.User = user2
			data.AdminUserList = append(data.AdminUserList, user1)
			data.AdminUserList = append(data.AdminUserList, user2)
			data.AdminUserList = append(data.AdminUserList, user1)

			tmpl.ExecuteTemplate(w, "clients.html", data)
		}
	}
	if params.ByName("suburl") == "/edit-client.html"{
		tmpl, err := template.ParseFiles("template/admin/edit-client.html")
		if err == nil {

			data := siteData{}
			data.User = user2
			data.AdminUserList = append(data.AdminUserList, user1)

			//AdminUser list contais the one user that should be edited
			tmpl.ExecuteTemplate(w, "edit-client.html", data)
		}
	}
	if params.ByName("suburl") == "/equipment.html"{
		tmpl, err := template.ParseFiles("template/admin/equipment.html")
		if err == nil {

			data := siteData{}
			data.Equipment = append(data.Equipment, equip1)
			data.Equipment = append(data.Equipment, equip2)
			data.Equipment = append(data.Equipment, equip3)
			data.User = user2

			tmpl.ExecuteTemplate(w, "equipment.html", data)
		}
	}
	if strings.Contains(params.ByName("suburl"), "img/"){
		fmt.Println(params.ByName("suburl"))
		content, err := ioutil.ReadFile("app/model/images/" + strings.Trim(params.ByName("suburl"),"/img"))
		w.Header().Set("Content-Type", "image/*")
		if err != nil {
			log.Fatal(err)
		} else {
			w.Write(content)
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
	router.GET("/img/*suburl", imgHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
