/*
 * Copyright (c) 2018. Philipp Kalytta
 */

package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/sessions"
	"time"
)

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
	logDatabase("SELECT * FROM users u;", fmt.Sprint(users))
	return &users

}

//Get one Special product with the given InvID
func getEquip(db *sql.DB, invID int64) *equipmentData {

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.InvID = " + strconv.FormatInt(invID, 10) + ";")

	checkErr(err)
	for rows.Next() {

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		logDatabase("SELECT * FROM equipment e WHERE e.InvID = "+strconv.FormatInt(invID, 10)+";", fmt.Sprint(inEquip))
		return &inEquip
	}

	log.Print("The Product with ID " + strconv.FormatInt(invID, 10) + " wasn't found in the Database")
	return nil
}

//Gets the products that are owned by the user UserID
func getEquipFromOwner(db *sql.DB, UserID int64) *[]equipmentData {

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.EquipmentOwnerID = " + strconv.FormatInt(UserID, 10) + ";")

	checkErr(err)

	var equipment []equipmentData
	for rows.Next() {

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		equipment = append(equipment, inEquip)

	}
	logDatabase("SELECT * FROM equipment e WHERE e.EquipmentOwnerID = "+strconv.FormatInt(UserID, 10)+";", fmt.Sprint(equipment))
	return &equipment

}

//Gets products that are rented by the user UserID (or that are bookmarked)
func getRentedEquip(db *sql.DB, UserID int64, bookmarked bool) (*[]equipmentData, *[]rentData) {

	var rows *sql.Rows
	var err error

	var invs *sql.Rows
	var rentlist []rentData

	//Fill the Rentlist (All Rents in the List for user UserID)
	invs, err = db.Query("SELECT * FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10))
	checkErr(err)

	for invs.Next() {

		var inRent rentData
		invs.Scan(&(inRent.RentID), &(inRent.UserID), &(inRent.InvID), &(inRent.RentDate), &(inRent.ReturnDate),
			&(inRent.Bookmarked), &(inRent.Amount), &(inRent.Repair))
		rentlist = append(rentlist, inRent)

	}
	logDatabase("RENTED-1 SELECT * FROM rentlist r WHERE r.userid = "+strconv.FormatInt(UserID, 10), fmt.Sprint(rentlist))

	//Get all the InvIDs that are rented or bookmarked by a user
	if (!bookmarked) {
		invs, err = db.Query("SELECT r.invid FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10) + " AND r.Bookmarked = false;")

		checkErr(err)

	} else {

		invs, err = db.Query("SELECT r.invid FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10) + " AND r.Bookmarked = true;")

		checkErr(err)

	}

	var ids []int

	//Save alle the IDs
	for invs.Next() {

		var inIDs int
		invs.Scan(&inIDs)
		ids = append(ids, inIDs)

	}

	logDatabase("RENTED-2 SELECT r.invid FROM rentlist r WHERE r.userid = "+strconv.FormatInt(UserID, 10)+" AND r.Bookmarked = "+fmt.Sprint(bookmarked), fmt.Sprint(ids))

	var equipment []equipmentData

	//select every item from the database
	for _, element := range ids {

		rows, err = db.Query("SELECT * FROM equipment e WHERE e.invid =" + strconv.Itoa(element))
		checkErr(err)

		for rows.Next() {

			var inEquip equipmentData
			rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
				&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
				&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
			equipment = append(equipment, inEquip)
			logDatabase("SELECT * FROM equipment e WHERE e.invid ="+strconv.Itoa(element), fmt.Sprint(inEquip))

		}

	}
	logDatabase("SELECT * FROM equipment e WHERE e.invid =*", fmt.Sprint(equipment))
	return &equipment, &rentlist

}

//Gets products, that are not rented and therefore can be rented, measured by StockAmount, if thats zero, every item is
//rented
func getAvailableEquip(db *sql.DB) *[]equipmentData {

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.StockAmount > 0;")

	checkErr(err)

	var equipment []equipmentData
	for rows.Next() {

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		equipment = append(equipment, inEquip)

	}
	logDatabase("SELECT * FROM equipment e WHERE e.StockAmount > 0;", fmt.Sprint(equipment))
	return &equipment

}

//Gets the featured products, in correct order
func getFeaturedEquip(db *sql.DB) *[]equipmentData {

	var equipment []equipmentData

	for i := 1; i <= 4; i++ {
		rows, err := db.Query("SELECT * FROM equipment e WHERE e.Featured = true AND e.FeaturedID = " + strconv.Itoa(i) + ";")

		checkErr(err)

		for rows.Next() {

			var inEquip equipmentData
			rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
				&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
				&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
			equipment = append(equipment, inEquip)

		}
	}

	logDatabase("SELECT * FROM equipment e WHERE e.Featured = true AND e.FeaturedID = *;", fmt.Sprint(equipment))
	return &equipment

}

//Gets a list of Renting information for a InvID or a UserID (one of the two can be passed, invID preferred)
func getRentList(db *sql.DB, invID int64, UserID int64) *[]rentData {

	var rentlist []rentData

	//Get all the rows with the given invid
	if invID > 0 {

		res, err := db.Query("SELECT * FROM rentlist r WHERE r.invid = " + strconv.FormatInt(invID, 10) + ";")
		checkErr(err)
		for res.Next() {

			var inRent rentData
			res.Scan(&(inRent.RentID), &(inRent.UserID), &(inRent.InvID), &(inRent.RentDate), &(inRent.ReturnDate), &(inRent.Bookmarked), &(inRent.Amount), &(inRent.Repair))

			//Get the User Name
			var user string
			username, err := db.Query("SELECT name FROM users WHERE userid = " + strconv.FormatInt(inRent.UserID, 10) + ";")
			checkErr(err)
			username.Next()
			username.Scan(user)
			logDatabase("SELECT name FROM users WHERE userid = "+strconv.FormatInt(inRent.UserID, 10)+";", user)
			inRent.RentedByUserName = user

			rentlist = append(rentlist, inRent)

		}
		logDatabase("SELECT * FROM rentlist r WHERE r.invid = "+strconv.FormatInt(invID, 10)+";", fmt.Sprint(rentlist))
		return &rentlist

		//get all the rows with the given userid
	} else if UserID > 0 {

		res, err := db.Query("SELECT * FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10) + ";")
		checkErr(err)
		for res.Next() {

			var inRent rentData
			res.Scan(&(inRent.RentID), &(inRent.UserID), &(inRent.InvID), &(inRent.RentDate), &(inRent.ReturnDate), &(inRent.Bookmarked), &(inRent.Amount), &(inRent.Repair))

			//Get the User Name
			var user string
			username, err := db.Query("SELECT name FROM users WHERE userid = " + strconv.FormatInt(UserID, 10) + ";")
			checkErr(err)
			username.Next()
			username.Scan(user)
			logDatabase("SELECT name FROM users WHERE userid = "+strconv.FormatInt(UserID, 10)+";", user)
			inRent.RentedByUserName = user

			rentlist = append(rentlist, inRent)

		}
		logDatabase("SELECT * FROM rentlist r WHERE r.invid = "+strconv.FormatInt(UserID, 10)+";", fmt.Sprint(rentlist))
		return &rentlist

	}

	log.Fatal("getRentList got no valid invID or UserID, it needs one of the two to work")
	return nil
}

//Gets a User form the Database and returns a pointer to it, if wanted, you can specify a password an let it test against the database one
func getUserFromName(db *sql.DB, username string, password string, validatePassword bool) *user {

	rows, err := db.Query("SELECT * FROM users u WHERE u.Name LIKE '" + username + "';")

	checkErr(err)
	for rows.Next() {
		var inUser user

		rows.Scan(&(inUser.UserID), &(inUser.Name), &(inUser.Email), &(inUser.Password), &(inUser.ProfileImageSRC), &(inUser.UserLevel), &(inUser.Blocked), &(inUser.ActiveUntilDate))
		logDatabase("SELECT * FROM users u WHERE u.Name LIKE '"+username+"';", fmt.Sprint(inUser))
		if (validatePassword && CheckPasswordHash(password,inUser.Password)) {
			return &inUser
		} else {
			log.Print("Error: The User could not be logged in")
			//TODO: What to do now?
		}

		return &inUser
	}

	log.Print("Error: No User with the given username")
	return nil
}

func getUserFromID(db *sql.DB, userid int) *user {

	rows, err := db.Query("SELECT * FROM users u WHERE u.userid =" + fmt.Sprint(userid) + ";")

	checkErr(err)
	for rows.Next() {

		var inUser user

		rows.Scan(&(inUser.UserID), &(inUser.Name), &(inUser.Email), &(inUser.Password), &(inUser.ProfileImageSRC), &(inUser.UserLevel), &(inUser.Blocked), &(inUser.ActiveUntilDate))
		logDatabase("SELECT * FROM users u WHERE u.useridLIKE '"+fmt.Sprint(userid)+"';", fmt.Sprint(inUser))

		return &inUser

	}

	return nil

}

//TODO: Gets the user, that the admin wants to edit, has to be used with session cookies or a flag in the database or query string on the URL
func getEditUser() *user {

	log.Println("Kein Datenbankzugriff nötig")

	return editUser
}

func getLoggedInUser(db *sql.DB, w http.ResponseWriter, r *http.Request, params httprouter.Params) *user {

	//Get Cookie with UserID
	session, _ := store.Get(r, "session")
	userid := session.Values["userid"]
	//THEN:
	//query := "SELECT * FROM users u WHERE u.userid = " + userid + ";"
	log.Print(fmt.Sprint(userid))
	id, _ := strconv.Atoi(fmt.Sprint(userid))
	if (id > 0) {
		rows, err := db.Query("SELECT * FROM users u WHERE u.userid = " + fmt.Sprint(userid) + ";")
		var userN user

		checkErr(err)
		for rows.Next() {

			rows.Scan(&(userN.UserID), &(userN.Name), &(userN.Email), &(userN.Password), &(userN.ProfileImageSRC), &(userN.UserLevel), &(userN.Blocked), &(userN.ActiveUntilDate))

		}
		logDatabase("SELECT * FROM users u WHERE u.userid = "+fmt.Sprint(userid)+";", fmt.Sprint(userN))
		if (userN.UserID != 0) {
			return &userN
		} else {
			http.Redirect(w, r, "/login.html", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/login.html", http.StatusFound)
	}
	//Gast zurückgeben
	return getUserFromName(db, "Gast", "NONE", false)
}

//Gets the Cart Items, that the logged-in User has in his cart
func getCartItemsForUser(db *sql.DB, session *sessions.Session) (*[]equipmentData) {

	var eq []equipmentData

	if (getExistingKey(session.Values["cart"]) == nil) {

		log.Println("There is no Cart, do nothing in the Database")

	} else {

		log.Println("There is a Cart")

		//Get the Items
		cartids := getExistingKey(session.Values["cart"])

		for _, element := range cartids {

			rows, err := db.Query("SELECT * FROM equipment e WHERE e.invid = " + fmt.Sprint(element) + ";")
			checkErr(err)

			for rows.Next() {
				var inEquip equipmentData
				rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
					&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
					&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
				eq = append(eq, inEquip)

				logDatabase("SELECT * FROM equipment e WHERE e.invid = "+fmt.Sprint(element)+";", fmt.Sprint(inEquip))
			}

		}

	}

	return &eq

}

func addEquipment(db *sql.DB, eq equipmentData) {

	_, err4 := db.Exec("INSERT INTO equipment VALUES(" +
		"'" + eq.Name + "'," +
		"'" + eq.Desc + "'," +
		"'" + eq.ImageSRC + "'," +
		"'Generic Placeholder'," +
		"'" + eq.Stock + "'," +
		"" + fmt.Sprint(eq.StockAmount) + "," +
		"'" + eq.Category + "'," +
		"false," +
		"-1," +
		"'NONE'," +
		"DEFAULT," +
		"'" + eq.StorageLocation + "'," +
		"" + strconv.FormatInt(eq.EquipmentOwnerID, 10) + "" +
		");")
	checkErr(err4)
	logDatabase("INSERT INTO equipment VALUES("+
		"'"+ eq.Name+ "',"+
		"'"+ eq.Desc+ "',"+
		"'"+ eq.ImageSRC+ "',"+
		"'Generic Placeholder',"+
		"'"+ eq.Stock+ "',"+
		""+ fmt.Sprint(eq.StockAmount)+ ","+
		"'"+ eq.Category+ "',"+
		"false,"+
		"-1,"+
		"'NONE',"+
		"DEFAULT,"+
		"'"+ eq.StorageLocation+ "',"+
		""+ strconv.FormatInt(eq.EquipmentOwnerID, 10)+ ""+
		");", "")

}

//TODO: Fill dummy, make use of hashing
func updateUser(db *sql.DB, user *user) bool {

	db.Exec("UPDATE users SET name='" + user.Name + "', email='" + user.Email + "',password=" + user.Password + " WHERE userid=" + strconv.FormatInt(user.UserID, 10) + ";")
	logDatabase("UPDATE users SET name='"+user.Name+"', email='"+user.Email+"',password="+user.Password+" WHERE userid="+strconv.FormatInt(user.UserID, 10)+";", "")
	var test []byte;
	bcrypt.GenerateFromPassword(test, 0)

	return false
}

func addUser(db *sql.DB, username string, email string, password string) {

	_, err := db.Exec("INSERT INTO users VALUES (" +
		"DEFAULT," +
		"'" + username + "'," +
		"'" + email + "'," +
		"'" + password + "'," +
		"'img/equipment/generic.gif'," +
		"'Benutzer'," +
		"false," +
		"'" + fmt.Sprint(time.Now().AddDate(1, 0, 0).Format(time.RFC822)) + "'" +
		");")

	if err != nil {
		log.Fatal(err)
	}
	logDatabase("INSERT INTO users VALUES ("+
		"DEFAULT,"+
		"'"+ username+ "',"+
		"'"+ email+ "',"+
		"'"+ password+ "',"+
		"'img/equipment/generic.gif',"+
		"'Benutzer',"+
		"false,"+
		"'"+ fmt.Sprint(time.Now().AddDate(1, 0, 0).Format(time.RFC822))+ "'"+
		");", "")

}

func createRent(db *sql.DB, userid int64, invid int64, bookmarked bool, amount int, repair bool) {

	_, err := db.Exec("INSERT INTO rentlist VALUES (" +
		"DEFAULT," +
		"" + fmt.Sprint(userid) + "," +
		"" + fmt.Sprint(invid) + "," +
		"'" + fmt.Sprint(time.Now().Format(time.RFC822)) + "'," +
		"'" + fmt.Sprint(time.Now().AddDate(0, 1, 0).Format(time.RFC822)) + "'," +
		"" + fmt.Sprint(bookmarked) + "," +
		"" + fmt.Sprint(amount) + "," +
		"" + fmt.Sprint(repair) + "" +
		");")

	if err != nil {
		log.Fatal(err)
	}

	//Amount reduzieren und Verfügbarkeit setzen

	eq := getEquip(db,invid)

	if(eq.StockAmount > 0) {

		_, err2 := db.Exec("UPDATE equipment SET " + " stockamount=" + fmt.Sprint(eq.StockAmount-amount)+ " WHERE invid=" + fmt.Sprint(invid))

		if err != nil {
			log.Fatal(err2)
		}

	}
	if (eq.StockAmount-amount)==0{

		_, err2 := db.Exec("UPDATE equipment SET " + " stock='Entliehen' WHERE invid=" + fmt.Sprint(invid))

		if err != nil {
			log.Fatal(err2)
		}

	}
	if(eq.StockAmount-amount)<0 || bookmarked{

		_, err2 := db.Exec("UPDATE equipment SET " + " stock='Vorgemerkt' WHERE invid=" + fmt.Sprint(invid))

		if err != nil {
			log.Fatal(err2)
		}

	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
