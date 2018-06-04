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
	noDefaultValues = true
)

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
var equip1 = equipmentData{"Kamera Obscura", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Entliehen", 0, "None", true, 1, "/img/equipment/gandalf.gif", true, true, false, 1, "Max Mustermann", "25.05.18", "25.05.18", 1245, "Baungasse", 2}
var equip2 = equipmentData{"Stativ", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Verfügbar", 2, "None", true, 2, "/img/equipment/gandalf.gif", true, false, false, 1, "Karl Karstens", "25.05.18", "25.05.18", 13452, "Schrank", 2}
var equip3 = equipmentData{"Mikrophon", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Vorgemerkt", 1, "None", true, 3, "/img/equipment/gandalf.gif", true, false, false, 1, "Max Mustermann", "25.05.18", "25.05.18", 2374, "Regal 12", 2}

var user1 = user{1, "Max Mustermann", "max@muster.de", "asdf", "img/equipment/gandalf.gif", "Benutzer", false, "25.07.18"}
var user2 = user{2, "Peter Müller", "peter@mueller.de", "asdf", "img/equipment/gandalf.gif", "Administrator", false, "25.07.22"}

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

		data := siteData{}
		data.Equipment = append(data.Equipment, equip1)
		data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)

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
	logAccess(r, params, "")
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
	logAccess(r, params, "")
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
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/profile.html")
	if err == nil {

		data := siteData{}
		data.User = user1
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

			data := siteData{}
			data.User = user2
			data.AdminUserList = append(data.AdminUserList, user1)
			data.AdminUserList = append(data.AdminUserList, user2)
			data.AdminUserList = append(data.AdminUserList, user1)

			tmpl.ExecuteTemplate(w, "clients.html", data)
		}
	}
	if params.ByName("suburl") == "/edit-client.html" {
		tmpl, err := template.ParseFiles("template/admin/edit-client.html")
		if err == nil {

			data := siteData{}
			data.User = user2
			data.AdminUserList = append(data.AdminUserList, user1)

			//AdminUser list contais the one user that should be edited
			tmpl.ExecuteTemplate(w, "edit-client.html", data)
		}
	}
	if params.ByName("suburl") == "/equipment.html" {
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

	fmt.Println(time.Now().Format(time.RFC3339) + " " + r.Header.Get("User-Agent") + " " + r.Method + " /" + fileDir + params.ByName("suburl"))

}

func connectDatabase() *sql.DB {

	//CREATE DATABASE IF NOT EXISTS borgdirmedia OWNER goserver ENCODING 'UTF-8'

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		//It fails here TODO: #6 BUG
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

	_, err := db.Exec("INSERT INTO users VALUES (" +
		"DEFAULT," +
		"'Max Mustermann'," +
		"'max@muster.de'," +
		"'asdf'," +
		"'img/equipment/gandalf.gif'," +
		"'Benutzer'," +
		"false," +
		"'2018-07-25'," +
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
		"'2022-08-11'," +
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
		"'2019-02-15'," +
		");")

	if err3 != nil {
		log.Fatal(err)
	}

	//TODO: Insert equipment 1,2,3 #7

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
		"RentedByUserID bigint REFERENCES users(UserID)," +
		"RentDate date," +
		"ReturnDate date," +
		"InvID bigserial," +
		"StorageLocation varchar(60)," +
		"EquipmentOwnerID bigint REFERENCES users(UserID)," +
		"PRIMARY KEY (InvID)" +
		");")

	if err2 != nil {
		log.Fatal(err2)
	}

}

//Gets a User form the Database and returns a pointer to it, if wanted, you can specify a password an let it test against the database one
func getUserFromName(db *sql.DB, username string, password string, validatePassword bool) *user {

	var inUser *user;
	rows, err := db.Query("SELECT * FROM users u WHERE u.Name LIKE '" + username + "';")

	checkErr(err)
	for rows.Next() {

		rows.Scan(&(inUser.UserID), &(inUser.Name), &(inUser.Email), &(inUser.Password), &(inUser.ProfileImageSRC), &(inUser.UserLevel), &(inUser.Blocked), &(inUser.ActiveUntilDate))
		if (validatePassword && inUser.Password == password) {
			return inUser
		} else {
			log.Print("Error: The User could not be logged in")
			return nil
		}
		return inUser
	}

	log.Print("Error: No User with the given username")
	return nil
}

//Gets the featured products

//Gets products, that are not rented and therefore can be rented

//Gets products that are rented by the user UserID

//Gets the products that are owned by the user UserID
func getEquipFromOwner(db *sql.DB,  UserID int64) *[]equipmentData{

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.EquipmentOwnerID == " + string(UserID) + ";")

	checkErr(err)

	var equipment []equipmentData
	for rows.Next(){

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.Rented), &(inEquip.Bookmarked), &(inEquip.Repair),
			&(inEquip.RentedByUserID), &(inEquip.RentedByUserName), &(inEquip.RentDate), &(inEquip.ReturnDate),
			&(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		equipment = append(equipment, inEquip)


	}

	return &equipment

}

//Get one Special product with the given InvID
func getProduct(db *sql.DB, invID int64) *equipmentData{

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.InvID == " + string(invID) + ";")

	checkErr(err)
	for rows.Next() {

		var inEquip *equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.Rented), &(inEquip.Bookmarked), &(inEquip.Repair),
			&(inEquip.RentedByUserID), &(inEquip.RentedByUserName), &(inEquip.RentDate), &(inEquip.ReturnDate),
			&(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		return inEquip

	}

	log.Print("The Product with ID " + string(invID) + "wasn't found in the Database")
	return nil
}

//Gets all Users
func getAllUsers(db *sql.DB) *[]user {

	var users []user
	rows, err := db.Query("SELECT * FROM users u;")

	checkErr(err)
	for rows.Next() {

		var userN user
		rows.Scan(&(userN.UserID), &(userN.Name), &(userN.Email), &(userN.Password), &(userN.ProfileImageSRC), &(userN.UserLevel), &(userN.Blocked), &(userN.ActiveUntilDate))
		users = append(users, userN)

	}

	return &users

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	db := connectDatabase()
	if db != nil {

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

		log.Fatal(http.ListenAndServe(":8080", router))

	}

}
