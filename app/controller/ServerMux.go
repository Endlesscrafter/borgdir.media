/*
 * Copyright (c) 2018. Philipp Kalytta
 */

package main

import (
	"net/http"
	_ "github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"html/template"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	"time"
)

const (
	DB_USER         = "goserver"
	DB_PASSWORD     = "c58WvoedyiVRmPjaEoEi"
	DB_NAME         = "goserver"
	noDefaultValues = false
	debug           = true
)

var GLOBALDB *sql.DB

type equipmentData struct {
	Name             string
	Desc             string
	ImageSRC         string
	ImageAlt         string
	Stock            string
	StockAmount      int
	Category         string
	Featured         bool
	FeaturedID       int
	FeaturedImageSRC string
	Rented           bool
	Bookmarked       bool
	Repair           bool
	RentedByUserID   int64
	RentedByUserName string
	RentDate         string
	ReturnDate       string
	InvID            int64
	StorageLocation  string
	EquipmentOwnerID int64
}

type user struct {
	UserID          int64
	Name            string
	Email           string
	Password        string
	ProfileImageSRC string
	UserLevel       string
	Blocked         bool
	ActiveUntilDate string
}

type siteData struct {
	Equipment     []equipmentData
	User          user
	AdminUserList []user
}

//Test Data
var nequip1 = equipmentData{"Kamera Obscura", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Entliehen", 0, "None", true, 1, "/img/equipment/gandalf.gif", true, true, false, 1, "Max Mustermann", "25.05.18", "25.05.18", 1245, "Baungasse", 2}
var nequip2 = equipmentData{"Stativ", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Verfügbar", 2, "None", true, 2, "/img/equipment/gandalf.gif", true, false, false, 1, "Karl Karstens", "25.05.18", "25.05.18", 13452, "Schrank", 2}
var nequip3 = equipmentData{"Mikrophon", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Vorgemerkt", 1, "None", true, 3, "/img/equipment/gandalf.gif", true, false, false, 1, "Max Mustermann", "25.05.18", "25.05.18", 2374, "Regal 12", 2}

var nstduser = user{1, "Max Mustermann", "max@muster.de", "asdf", "img/equipment/gandalf.gif", "Benutzer", false, "25.07.18"}
var nadmuser = user{2, "Peter Müller", "peter@mueller.de", "asdf", "img/equipment/gandalf.gif", "Administrator", false, "25.07.22"}

func cssHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "static/css")
	content, err := ioutil.ReadFile("static/css" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/css")
	if err != nil {
		log.Fatal(err)
	} else {
		w.Write(content)
	}
}
func jsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "static/js")
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
	logAccess(r, params, "app/model/images")
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
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/index.html")
	//Content-Type text/html doesnt fix malformatted css
	w.Header().Set("Content-Type", "text/html")
	if err == nil {

		eq := getFeaturedProducts(GLOBALDB)

		data := siteData{}
		data.Equipment = append(data.Equipment, (*eq)[0])
		data.Equipment = append(data.Equipment, (*eq)[1])
		data.Equipment = append(data.Equipment, (*eq)[2])

		tmpl.ExecuteTemplate(w, "index.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func loginHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/login.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "login.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func registerHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/register.html")
	if err == nil {
		tmpl.ExecuteTemplate(w, "register.html", "DATATATATATA")
	}
	//fmt.Fprintf(w, "index")
}

func cartHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/cart.html")
	if err == nil {

		eq, user := getCartItemsForUser(GLOBALDB, "TEST")

		data := siteData{}

		for _, element := range *eq {
			data.Equipment = append(data.Equipment, element)
		}
		/*data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)
		data.Equipment = append(data.Equipment, equip3)*/
		data.User = *user

		tmpl.ExecuteTemplate(w, "cart.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func equipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/equipment.html")
	if err == nil {

		eq := getAvailableEqip(GLOBALDB)
		user := getLoggedInUser()

		data := siteData{}
		for _, element := range *eq {
			data.Equipment = append(data.Equipment, element)
		}
		data.User = *user

		//Always give a even number of data-Sets (0, 2, 4, 6 etc): Otherwise it breaks the Design
		tmpl.ExecuteTemplate(w, "equipment.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func myEquipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/my-equipment.html")
	if err == nil {

		user := getLoggedInUser()
		//eq := getEquipFromOwner(GLOBALDB, user.UserID)
		eq := getRentedEquip(GLOBALDB, user.UserID,false)
		eqb := getRentedEquip(GLOBALDB, user.UserID, true)

		data := siteData{}
		for _, element := range *eq {
			data.Equipment = append(data.Equipment, element)
		}
		for _, element := range *eqb {
			data.Equipment = append(data.Equipment, element)
		}
		data.User = *user

		tmpl.ExecuteTemplate(w, "my-equipment.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func profileHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/profile.html")
	if err == nil {

		user := getLoggedInUser()

		data := siteData{}
		data.User = *user
		tmpl.ExecuteTemplate(w, "profile.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func adminHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	logAccess(r, params, "admin")
	if params.ByName("suburl") == "/index.html" || params.ByName("suburl") == "/" {
		tmpl, err := template.ParseFiles("template/admin/index.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "index.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/add.html" {
		tmpl, err := template.ParseFiles("template/admin/add.html")
		if err == nil {
			tmpl.ExecuteTemplate(w, "add.html", "DATATATATATA")
		}
	}
	if params.ByName("suburl") == "/clients.html" {
		tmpl, err := template.ParseFiles("template/admin/clients.html")
		if err == nil {

			user := getLoggedInUser()
			users := getAllUsers(GLOBALDB)

			data := siteData{}
			data.User = *user
			for _, element := range *users {
				data.AdminUserList = append(data.AdminUserList, element)
			}

			tmpl.ExecuteTemplate(w, "clients.html", data)
		}
	}
	if params.ByName("suburl") == "/edit-client.html" {
		tmpl, err := template.ParseFiles("template/admin/edit-client.html")
		if err == nil {

			user := getLoggedInUser()
			edituser := getEditUser()

			data := siteData{}
			data.User = *user
			data.AdminUserList = append(data.AdminUserList, *edituser)

			//AdminUser list contais the one user that should be edited
			tmpl.ExecuteTemplate(w, "edit-client.html", data)
		}
	}
	if params.ByName("suburl") == "/equipment.html" {
		tmpl, err := template.ParseFiles("template/admin/equipment.html")
		if err == nil {

			user := getLoggedInUser()
			eq := getEquipFromOwner(GLOBALDB, user.UserID)

			data := siteData{}
			for _, element := range *eq {
				data.Equipment = append(data.Equipment, element)
			}
			data.User = *user

			tmpl.ExecuteTemplate(w, "equipment.html", data)
		}
	}
	if params.ByName("suburl") == "/profile.html" {

		http.Redirect(w, r, "/profile.html", http.StatusFound)

	}
	if strings.Contains(params.ByName("suburl"), "img/") {
		content, err := ioutil.ReadFile("app/model/images/" + strings.Trim(params.ByName("suburl"), "/img"))
		w.Header().Set("Content-Type", "image/*")
		if err != nil {
			log.Fatal(err)
		} else {
			w.Write(content)
		}
	}
}

func logAccess(r *http.Request, params httprouter.Params, fileDir string) {

	fmt.Println(time.Now().Format(time.RFC3339) + " ACCESS: " + r.Header.Get("User-Agent") + " " + r.Method + " /" + fileDir + params.ByName("suburl"))

}

func logDatabase(query string, data string) {
	if (debug) {
		fmt.Println(time.Now().Format(time.RFC3339) + " " + "QUERY: " + query + " | RESULT: " + data)
	} else {
		fmt.Println(time.Now().Format(time.RFC3339) + " " + "QUERY: " + query + " | RESULT OMITTED")
	}
}

func connectDatabase() *sql.DB {

	//CREATE DATABASE IF NOT EXISTS borgdirmedia OWNER goserver ENCODING 'UTF-8'

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		//It fails here TODO: #6 BUG
		log.Fatal(err)
		log.Fatal("Error: Could not establish a connection with the database")
	}

	//Create the tables
	createTables(db)

	//Create default Admin user, a distributor, a standard user and some equipment
	if noDefaultValues {
		createDummyValues(db)
	}

	return db;
}

//Creates some dummy values
func createDummyValues(db *sql.DB) {

	//Delete everything beforehand
	/*_, e := db.Exec("DELETE FROM users; DELETE FROM equipment;")
	checkErr(e)

	_, err := db.Exec("INSERT INTO users VALUES (" +
		"DEFAULT," +
		"'Max Mustermann'," +
		"'max@muster.de'," +
		"'asdf'," +
		"'img/equipment/gandalf.gif'," +
		"'Benutzer'," +
		"false," +
		"'2018-07-25'" +
		");")

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := db.Exec("INSERT INTO users VALUES (" +
		"DEFAULT," +
		"'Peter Müller'," +
		"'peter@mueller.de'," +
		"'asdf'," +
		"'img/equipment/gandalf.gif'," +
		"'Administrator'," +
		"false," +
		"'2022-08-11'" +
		");")

	if err2 != nil {
		log.Fatal(err)
	}

	_, err3 := db.Exec("INSERT INTO users VALUES (" +
		"DEFAULT," +
		"'Distri Butor'," +
		"'distri@butor.de'," +
		"'asdf'," +
		"'img/equipment/gandalf.gif'," +
		"'Distributor'," +
		"false," +
		"'2019-02-15'" +
		");")

	if err3 != nil {
		log.Fatal(err)
	}*/

	_, err4 := db.Exec("INSERT INTO equipment VALUES(" +
		"'Kamera Obscura'," +
		"'Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit'," +
		"'img/equipment/generic.gif'," +
		"'Generic Placeholder'," +
		"'Verfügbar'," +
		"1," +
		"'None'," +
		"true," +
		"1," +
		"'/img/equipment/gandalf.gif'," +
		"false," +
		"false," +
		"false," +
		"NULL," +
		"NULL," +
		"NULL," +
		"NULL," +
		"DEFAULT," +
		"'Baungasse'," +
		"3" +
		");")
	checkErr(err4)

	_, err5 := db.Exec("INSERT INTO equipment VALUES(" +
		"'Stativ'," +
		"'Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit'," +
		"'img/equipment/generic.gif'," +
		"'Generic Placeholder'," +
		"'Verfügbar'," +
		"2," +
		"'None'," +
		"true," +
		"2," +
		"'/img/equipment/gandalf.gif'," +
		"false," +
		"false," +
		"false," +
		"NULL," +
		"NULL," +
		"NULL," +
		"NULL," +
		"DEFAULT," +
		"'Schrank'," +
		"3" +
		");")
	checkErr(err5)

	_, err6 := db.Exec("INSERT INTO equipment VALUES(" +
		"'Mikrophon'," +
		"'Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit'," +
		"'img/equipment/generic.gif'," +
		"'Generic Placeholder'," +
		"'Vorgemerkt'," +
		"1," +
		"'None'," +
		"true," +
		"3," +
		"'/img/equipment/gandalf.gif'," +
		"false," +
		"false," +
		"false," +
		"NULL," +
		"NULL," +
		"NULL," +
		"NULL," +
		"DEFAULT," +
		"'Regal 12'," +
		"3" +
		");")
	checkErr(err6)

}

//Creates the necessary tables in the Database
func createTables(db *sql.DB) {

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"UserID bigserial NOT NULL," +
		"Name varchar(60) NOT NULL," +
		"Email varchar(128) NOT NULL," +
		"Password varchar(128) NOT NULL," +
		"ProfileImageSRC text," +
		"UserLevel varchar(30)," +
		"Blocked boolean," +
		"ActiveUntilDate date," +
		"PRIMARY KEY (UserID)" +
		");")

	if err != nil {
		log.Fatal("Fehler beim Anlegen der users-DB")
		log.Fatal(err)
	}

	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS equipment (" +
		"Name varchar(128)," +
		"Description text," +
		"ImageSRC text," +
		"ImageAlt varchar(128)," +
		"Stock varchar(30)," +
		"StockAmount int," +
		"Category varchar(60)," +
		"Featured boolean," +
		"FeaturedID int," +
		"FeaturedImageSRC text," +
		"Rented boolean," +
		"Bookmarked boolean," +
		"Repair boolean," +
		"RentedByUserID bigint REFERENCES users(UserID)," +
		"RentedByUserName varchar(60)," +
		"RentDate date," +
		"ReturnDate date," +
		"InvID bigserial," +
		"StorageLocation varchar(60)," +
		"EquipmentOwnerID bigint REFERENCES users(UserID)," +
		"PRIMARY KEY (InvID)" +
		");")

	if err2 != nil {
		log.Fatal("Fehler beim Anlegen der equip-DB")
		log.Fatal(err2)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	GLOBALDB = connectDatabase()
	if GLOBALDB != nil {

		//fmt.Print("Connection to Database successful\n")
		//fmt.Print(getAllUsers(GLOBALDB))
		//fmt.Print(getUserFromName(GLOBALDB, "Max Mustermann", "", false))
		//fmt.Print(getFeaturedProducts(GLOBALDB))
		//fmt.Print(getAvailableEqip(GLOBALDB))
		//fmt.Print(getRentedEquip(GLOBALDB, 1, false))
		//fmt.Print(getEquipFromOwner(GLOBALDB, 2))
		//fmt.Print(getEquip(GLOBALDB, 1))
		//Start Routing the Information
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

		log.Print("Server started successfully")
		log.Fatal(http.ListenAndServe(":80", router))

	}
	defer GLOBALDB.Close()
}
