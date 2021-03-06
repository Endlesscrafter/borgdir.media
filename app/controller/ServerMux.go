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
	"crypto/rand"
	"github.com/gorilla/sessions"
	"strconv"
	"github.com/satori/go.uuid"
)

const (
	DB_USER         = "goserver"
	DB_PASSWORD     = "c58WvoedyiVRmPjaEoEi"
	DB_NAME         = "goserver"
	noDefaultValues = false
	debug           = true //true
)

var store *sessions.CookieStore
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

type rentData struct {
	RentID           int64
	UserID           int64
	InvID            int64
	RentDate         string
	ReturnDate       string
	Bookmarked       bool
	Amount           int
	Repair           bool
	RentedByUserName string
}

type siteData struct {
	Equipment           []equipmentData
	Rentlist            []rentData
	User                user
	AdminUserList       []user
	EquipmentBookmarked []equipmentData
	ReturnDate          string
}

var editUser *user

//Test Data
var nequip1 = equipmentData{"Kamera Obscura", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Entliehen", 0, "None", true, 1, "/img/equipment/gandalf.gif", 1245, "Baungasse", 2}
var nequip2 = equipmentData{"Stativ", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Verfügbar", 2, "None", true, 2, "/img/equipment/gandalf.gif", 13452, "Schrank", 2}
var nequip3 = equipmentData{"Mikrophon", "Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit", "img/equipment/generic.gif", "Generic Placeholder", "Vorgemerkt", 1, "None", true, 3, "/img/equipment/gandalf.gif", 2374, "Regal 12", 2}

var nstduser = user{1, "Max Mustermann", "max@muster.de", "asdf", "img/equipment/gandalf.gif", "Benutzer", false, "25.07.18"}
var nadmuser = user{2, "Peter Müller", "peter@mueller.de", "asdf", "img/equipment/gandalf.gif", "Administrator", false, "25.07.22"}

func cssHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "static/css")
	content, err := ioutil.ReadFile("static/css" + params.ByName("suburl"))
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "max-age=3600")
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
	w.Header().Set("Cache-Control", "max-age=3600")
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
	checkErr(err)
	w.Header().Set("Content-Type", "image/*")
	w.Header().Set("Cache-Control", "max-age=3600")
	if err != nil {
		w.WriteHeader(404)
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

		eq := getFeaturedEquip(GLOBALDB)

		data := siteData{}
		data.Equipment = append(data.Equipment, (*eq)[0])
		data.Equipment = append(data.Equipment, (*eq)[1])
		data.Equipment = append(data.Equipment, (*eq)[2])
		//4th?
		//data.Equipment = append(data.Equipment, (*eq)[4])

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
	/*tmpl.Funcs(template.FuncMap{
		"Format": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	})*/
	if err == nil {

		session, _ := store.Get(r, "session")
		eq := getCartItemsForUser(GLOBALDB, session)

		data := siteData{}

		for _, element := range *eq {
			data.Equipment = append(data.Equipment, element)
		}
		/*data.Equipment = append(data.Equipment, equip2)
		data.Equipment = append(data.Equipment, equip3)
		data.Equipment = append(data.Equipment, equip3)*/
		data.User = *getLoggedInUser(GLOBALDB, w, r, params)

		data.ReturnDate = fmt.Sprint(time.Now().AddDate(0, 1, 0).Format(time.RFC822))

		tmpl.ExecuteTemplate(w, "cart.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func equipHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/equipment.html")
	if err == nil {

		//eq := getAvailableEquip(GLOBALDB)
		eq := getAvailableEquip(GLOBALDB, true)
		user := getLoggedInUser(GLOBALDB, w, r, params)

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

		user := getLoggedInUser(GLOBALDB, w, r, params)
		//eq := getEquipFromOwner(GLOBALDB, user.UserID)
		eq, rent1 := getRentedEquip(GLOBALDB, user.UserID, false)
		eqb, _ := getRentedEquip(GLOBALDB, user.UserID, true)

		data := siteData{}
		for _, element := range *eq {
			data.Equipment = append(data.Equipment, element)

		}
		for _, element := range *eqb {
			data.EquipmentBookmarked = append(data.EquipmentBookmarked, element)
		}
		data.User = *user

		for _, element := range *rent1 {
			data.Rentlist = append(data.Rentlist, element)
		}

		tmpl.ExecuteTemplate(w, "my-equipment.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func profileHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
	tmpl, err := template.ParseFiles("template/profile.html")
	if err == nil {

		user := getLoggedInUser(GLOBALDB, w, r, params)

		if (user.UserID == 4) {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		data := siteData{}
		data.User = *user
		tmpl.ExecuteTemplate(w, "profile.html", data)
	}
	//fmt.Fprintf(w, "index")
}

func adminHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	user := getLoggedInUser(GLOBALDB, w, r, params)
	if (user.UserLevel == "Benutzer") {
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}

	logAccess(r, params, "admin")
	if params.ByName("suburl") == "/index.html" || params.ByName("suburl") == "/" {
		tmpl, err := template.ParseFiles("template/admin/index.html")
		if err == nil {
			data := siteData{}
			user := getLoggedInUser(GLOBALDB, w, r, params)
			data.User = *user
			tmpl.ExecuteTemplate(w, "index.html", data)
		}
	}
	if params.ByName("suburl") == "/add.html" {
		tmpl, err := template.ParseFiles("template/admin/add.html")
		if err == nil {

			data := siteData{}
			user := getLoggedInUser(GLOBALDB, w, r, params)
			data.User = *user

			tmpl.ExecuteTemplate(w, "add.html", data)
		}
	}
	if /*strings.Contains(params.ByName("suburl"), "/edit.html")*/params.ByName("suburl") == "/edit.html" {

		value := r.URL.Query().Get("i")
		log.Println("query-wert: " + value)

		tmpl, err := template.ParseFiles("template/admin/edit.html")
		log.Println(err)
		if err == nil {

			suburl := params.ByName("suburl")

			log.Print(suburl)

			//parts := strings.SplitAfter(suburl, "?i=")

			//invid, _ := strconv.ParseInt(parts[1], 10, 64)

			invid, _ := strconv.ParseInt(value, 10, 64)

			log.Print("Zu editierendes Equip: " + fmt.Sprint(invid))
			data := siteData{}
			user := getLoggedInUser(GLOBALDB, w, r, params)
			data.User = *user
			eq := getEquip(GLOBALDB, invid)
			data.Equipment = append(data.Equipment, *eq)

			tmpl.ExecuteTemplate(w, "edit.html", data)
		}

	}
	if params.ByName("suburl") == "/clients.html" {
		tmpl, err := template.ParseFiles("template/admin/clients.html")
		if err == nil {

			user := getLoggedInUser(GLOBALDB, w, r, params)
			users := getAllUsers(GLOBALDB)
			rents := getCompleteRentList(GLOBALDB)

			data := siteData{}
			data.User = *user
			data.Rentlist = *rents
			for _, element := range *users {
				data.AdminUserList = append(data.AdminUserList, element)
			}

			tmpl.ExecuteTemplate(w, "clients.html", data)
		}
	}
	if params.ByName("suburl") == "/edit-client.html" {
		tmpl, err := template.ParseFiles("template/admin/edit-client.html")
		if err == nil {

			user := getLoggedInUser(GLOBALDB, w, r, params)
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

			user := getLoggedInUser(GLOBALDB, w, r, params)
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

		imageurl := strings.Trim(params.ByName("suburl"), "/img")
		if !strings.Contains(imageurl, ".jpg") {
			if (strings.Contains(imageurl, ".jp")) {
				imageurl += "g"
			}
		}
		content, err := ioutil.ReadFile("app/model/images/" + imageurl)
		w.Header().Set("Content-Type", "image/*")
		if err != nil {
			log.Fatal(err)
		} else {
			w.Write(content)
		}
	}
}

func loginPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	session, _ := store.Get(r, "session")

	username := r.FormValue("username")
	password := r.FormValue("password")

	user := getUserFromName(GLOBALDB, username, password, true)
	if (user != nil && !user.Blocked) {
		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Values["username"] = user.Name
		session.Values["userid"] = user.UserID
		session.Save(r, w)
		if (user.UserLevel == "Administrator") {
			http.Redirect(w, r, "/admin/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}

	} else {
		http.Redirect(w, r, "login.html", http.StatusFound)
	}

	logAccess(r, params, "")
}
func loginGuestPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	session, _ := store.Get(r, "session")

	// Set user as authenticated
	session.Values["authenticated"] = false
	session.Values["username"] = "Gast"
	session.Values["userid"] = 4
	session.Save(r, w)

	http.Redirect(w, r, "/equipment.html", http.StatusFound)

	logAccess(r, params, "")
}
func registerPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	username := r.FormValue("username")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	if (password1 != password2) {
		http.Redirect(w, r, "/register.html", http.StatusFound)
	} else {

		hash, _ := HashPassword(password1)
		addUser(GLOBALDB, username, email, hash)
		http.Redirect(w, r, "/login.html", http.StatusFound)

	}

	logAccess(r, params, "")
}

func cartPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	invid, _ := strconv.Atoi(r.FormValue("cart"))
	session, _ := store.Get(r, "session")

	if (getExistingKey(session.Values["cart"]) == nil) {

		log.Println("There is no Cart, create it")
		var cartids []int
		cartids = append(cartids, invid)

		session.Values["cart"] = cartids
		session.Save(r, w)
		http.Redirect(w, r, "/equipment.html", http.StatusFound)

	} else {

		log.Println("There already is a cart, use it")

		cartids := getExistingKey(session.Values["cart"])

		cartids = append(cartids, invid)

		log.Println("Cart: "+fmt.Sprint(cartids), "Added this time: "+fmt.Sprint(invid))
		session.Values["cart"] = cartids
		session.Save(r, w)
		http.Redirect(w, r, "/equipment.html", http.StatusFound)

	}

	logAccess(r, params, "")
}

func myEquipPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")
}

func wishPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")

	wish := r.FormValue("wish")
	invid, _ := strconv.ParseInt(wish,10,64)

	session, _ := store.Get(r, "session")

	username := session.Values["username"].(string)

	user := getUserFromName(GLOBALDB,username,"",false)

	createRent(GLOBALDB,user.UserID,invid,true,1,false)

	http.Redirect(w, r, "my-equipment.html", http.StatusFound)

}

func blockPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")

	user := r.FormValue("userid")
	userid, _ := strconv.ParseInt(user, 10, 64)

	blockUser(GLOBALDB, userid)

	http.Redirect(w, r, "clients.html", http.StatusFound)

}

func editEqPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	Name := r.FormValue("name")
	Desc := r.FormValue("description")
	StockAmount := r.FormValue("amount")
	Category := r.FormValue("category")
	StorageLocation := r.FormValue("storagelocation")
	InvID := r.FormValue("invid")

	log.Println("INVENTARID:" + InvID)

	id, _ := strconv.ParseInt(InvID, 10, 64)
	amount, _ := strconv.Atoi(StockAmount)

	file, _, err := r.FormFile("image")
	uploadedFile := false
	uuid := uuid.Must(uuid.NewV4())
	if err != nil {
		log.Println("Image could not be read correctly: " + fmt.Sprint(err))
		uploadedFile = false
	} else {
		//defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Image could not be read correctly #2")
			uploadedFile = false
		}

		err2 := ioutil.WriteFile("./app/model/images/equipment/"+fmt.Sprint(uuid)+".jpg", fileBytes, 0644)
		uploadedFile = true
		check(err2)
	}

	eq := getEquip(GLOBALDB, id)
	eq.Name = Name
	eq.Desc = Desc
	if uploadedFile {
		eq.ImageSRC = "img/equipment/" + fmt.Sprint(uuid) + ".jpg"
	}
	if(amount > 0){
		eq.Stock = "Verfügbar"
	}
	eq.Category = Category
	eq.StorageLocation = StorageLocation
	eq.StockAmount = amount

	updateEquip(GLOBALDB, eq)
	http.Redirect(w, r, "equipment.html", http.StatusFound)

	logAccess(r, params, "")
}

func profilePOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")

	username := r.FormValue("username")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	//password0 := r.FormValue("password0")

	/*Profilbild*/
	file, _, err := r.FormFile("image")
	uploadedFile := true
	uuid := uuid.Must(uuid.NewV4())
	if err != nil {
		log.Println("Image could not be read correctly: " + fmt.Sprint(err))
		uploadedFile = false
	} else {
		//defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Image could not be read correctly #2")
			uploadedFile = false
		}

		err2 := ioutil.WriteFile("./app/model/images/profile/"+fmt.Sprint(uuid)+".jpg", fileBytes, 0644)
		check(err2)
	}

	if (password1 == "" || password2 == "") {
		http.Redirect(w, r, "/profile.html", http.StatusFound)
	}
	if (password1 != password2) {
		http.Redirect(w, r, "/profile.html", http.StatusFound)
	} else {
		hash, _ := HashPassword(password1)
		user := getUserFromName(GLOBALDB, username, hash, false)
		if (username != "") {
			user.Name = username
		}
		if (email != "") {
			user.Email = email
		}
		if (password1 != "") {
			user.Password = hash
		}
		//Hoffe das klappt so, soll schauen ob die datei leer war
		if uploadedFile {
			user.ProfileImageSRC = "img/profile/" + fmt.Sprint(uuid) + ".jpg"
		}

		updateUser(GLOBALDB, user)
		http.Redirect(w, r, "/profile.html", http.StatusFound)
	}

}

func addPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	name := r.FormValue("name")
	desc := r.FormValue("description")
	cat := r.FormValue("category")
	amount := r.FormValue("amount")
	storagelocation := r.FormValue("storagelocation")
	//image := r.FormValue("image")

	file, _, err := r.FormFile("image")
	if err != nil {
		log.Fatal("Image could not be read correctly")
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Image could not be read correctly #2")
	}
	uuid := uuid.Must(uuid.NewV4())
	err2 := ioutil.WriteFile("./app/model/images/equipment/"+fmt.Sprint(uuid)+".jpg", fileBytes, 0644)
	check(err2)
	/*encoded := base64.StdEncoding.EncodeToString(fileBytes)
	var prefix = "data:image/jpg;base64,";
	encoded = prefix + encoded*/

	var eq equipmentData
	eq.Name = name
	eq.StorageLocation = storagelocation
	eq.FeaturedImageSRC = "NONE"
	eq.Featured = false
	eq.FeaturedID = -1
	eq.Category = cat
	eq.StockAmount, _ = strconv.Atoi(amount)
	eq.Stock = "Verfügbar"
	eq.ImageAlt = "NONE"
	eq.Desc = desc
	if (fmt.Sprint(uuid) != "") {
		eq.ImageSRC = "img/equipment/" + fmt.Sprint(uuid) + ".jpg"
	} else {
		eq.ImageSRC = "img/equipment/generic.gif"
	}

	session, _ := store.Get(r, "session")

	userid := session.Values["userid"]
	id := fmt.Sprint(userid)
	idi, _ := strconv.ParseInt(id, 10, 64)
	eq.EquipmentOwnerID = idi

	addEquipment(GLOBALDB, eq)

	logAccess(r, params, "")
	http.Redirect(w, r, "/admin/equipment.html", http.StatusFound)
}

//Links to edit-client, when a client should be edited
func editPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	useridstr := r.FormValue("userid")
	userid, _ := strconv.Atoi(useridstr)

	editUser = getUserFromID(GLOBALDB, userid)

	if editUser != nil {

		http.Redirect(w, r, "/admin/edit-client.html", http.StatusFound)

	}

	logAccess(r, params, "")
}

//links back to owerview. Gets the edited client as update
func editedPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	name := r.FormValue("username")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	//No use??
	//image := r.FormValue("image")

	if (password1 != password2) {
		http.Redirect(w, r, "/admin/edit-client.html", http.StatusFound)
	} else {
		hash, _ := HashPassword(password1)

		user := getUserFromName(GLOBALDB, name, hash, false)
		user.Name = name
		user.Email = email
		user.Password = hash
		updateUser(GLOBALDB, user)
		http.Redirect(w, r, "/admin/clients.html", http.StatusFound)
	}

	logAccess(r, params, "")
}
func logoutPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logAccess(r, params, "")

	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Values["username"] = "NONE"
	session.Values["userid"] = 0
	//Logout löschen des Warenkorbs
	var cartids []int
	session.Values["cart"] = cartids
	sessions.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}

func rentPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	logAccess(r, params, "")

	//User ranholen
	session, _ := store.Get(r, "session")
	useridstr := session.Values["userid"]
	userid, _ := useridstr.(int64)

	//Warenkorb holen

	eq := getCartItemsForUser(GLOBALDB, session)

	//Rents erstellen
	for _, element := range *eq {

		createRent(GLOBALDB, userid, element.InvID, false, 1, false)

	}

	//Warenkorb leeren
	var cartids []int
	session.Values["cart"] = cartids
	session.Save(r, w)
	http.Redirect(w, r, "/equipment.html", http.StatusFound)

}

//Löscht aus der Datenbank
func delEquipPOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	logAccess(r, params, "")

	invid, _ := strconv.Atoi(r.FormValue("cart"))
	log.Print("Zu löschende InventarID: " + fmt.Sprint(invid))

	delEquipment(GLOBALDB, invid)

}

//Löscht einträge aus dem Warenkorb
func deletePOSTHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	logAccess(r, params, "")

	invid, _ := strconv.Atoi(r.FormValue("cart"))
	log.Print("Zu löschende InventarID: " + fmt.Sprint(invid))
	session, _ := store.Get(r, "session")

	cartids := getExistingKey(session.Values["cart"])
	index := 0
	for p, v := range cartids {
		if (v == invid) {
			index = p
		}
	}

	cartids = remove(cartids, index)
	session.Values["cart"] = cartids
	session.Save(r, w)
	http.Redirect(w, r, "/equipment.html", http.StatusFound)

}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func logAccess(r *http.Request, params httprouter.Params, fileDir string) {

	fmt.Println(time.Now().Format(time.RFC3339) + " ACCESS: " + r.Header.Get("User-Agent") + " " + r.Method + " /" + fileDir + params.ByName("suburl"))

}

func logDatabase(query string, data string) {
	if (debug) {
		fmt.Println(time.Now().Format(time.RFC3339) + " " + "QUERY: " + query + " | RESULT: " + data)
	} else {
		//No output
		//fmt.Println(time.Now().Format(time.RFC3339) + " " + "QUERY: " + query + " | RESULT OMITTED")
	}
}

func connectDatabase() *sql.DB {

	//CREATE DATABASE IF NOT EXISTS borgdirmedia OWNER goserver ENCODING 'UTF-8'

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	//defer db.Close()

	log.Println("Trying to connect to database...")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error: Could not establish a connection with the database")
	}

	//Create the tables
	createTables(db)

	log.Println("Connection works, tables created....")

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
		"'img/equipment/gandalf.gif'," +
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
		"Repair boolean," +
		"InvID bigserial," +
		"StorageLocation varchar(60)," +
		"EquipmentOwnerID bigint REFERENCES users(UserID)," +
		"PRIMARY KEY (InvID)" +
		");")

	if err2 != nil {
		log.Fatal("Fehler beim Anlegen der equip-DB")
		log.Fatal(err2)
	}

	_, err3 := db.Exec("CREATE TABLE IF NOT EXISTS rentlist (" +
		"rentid bigserial," +
		"userid bigint references users(UserID)," +
		"invid bigint references equipment(invid)," +
		"rentdate date," +
		"returndate date," +
		"bookmarked boolean," +
		"repair boolean," +
		"amount int," +
		"PRIMARY KEY (rentid)" +
		");")

	if err3 != nil {
		log.Fatal("Fehler beim Anlegen der rentlist-DB")
		log.Fatal(err3)
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

		fmt.Print("Connection to Database successful\n")
		//fmt.Print(getAllUsers(GLOBALDB))
		//fmt.Print(getUserFromName(GLOBALDB, "Max Mustermann", "", false))
		//fmt.Print(getFeaturedEquip(GLOBALDB))
		//fmt.Print(getAvailableEquip(GLOBALDB))
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
		router.GET("/admin/*suburl", adminHandler)
		router.GET("/css/*suburl", cssHandler)
		router.GET("/js/*suburl", jsHandler)
		router.GET("/img/*suburl", imgHandler)

		router.POST("/login.html", loginPOSTHandler)
		router.POST("/loginGuest.html", loginGuestPOSTHandler)
		router.POST("/register.html", registerPOSTHandler)
		router.POST("/cart.html", cartPOSTHandler)
		//router.POST("/my-equipment.html", myEquipPOSTHandler)
		router.POST("/profile.html", profilePOSTHandler)
		router.POST("/admin/add.html", addPOSTHandler)
		router.POST("/admin/edit-client.html", editPOSTHandler)
		router.POST("/admin/edited-client.html", editedPOSTHandler)
		router.POST("/logout.html", logoutPOSTHandler)
		router.POST("/cart-rent.html", rentPOSTHandler)
		router.POST("/cart-del.html", deletePOSTHandler)
		router.POST("/delete-equip.html", delEquipPOSTHandler)
		router.POST("/admin/edit.html", editEqPOSTHandler)
		router.POST("/admin/block.html", blockPOSTHandler)
		router.POST("/wish.html", wishPOSTHandler)

		log.Print("Server started successfully")
		log.Fatal(http.ListenAndServe(":80", router))

	}
	defer GLOBALDB.Close()
}

// Is executed automatically on package load
func init() {
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewCookieStore(key)
}

func getExistingKey(f interface{}) []int {
	if f != nil {
		if key, ok := f.([]int); ok {
			return key
		}
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
