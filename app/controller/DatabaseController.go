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
	logDatabase("RENTED-1 SELECT * FROM rentlist r WHERE r.userid = "+strconv.FormatInt(UserID, 10), fmt.Sprint(invs))
	for invs.Next() {

		var inRent rentData
		invs.Scan(&(inRent.RentID), &(inRent.UserID), &(inRent.InvID), &(inRent.RentDate), &(inRent.ReturnDate),
			&(inRent.Bookmarked), &(inRent.Amount), &(inRent.Repair))
		rentlist = append(rentlist,inRent)

	}

	//Get all the InvIDs that are rented or bookmarked by a user
	if (!bookmarked) {
		invs, err = db.Query("SELECT r.invid FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10) + " AND r.Bookmarked = false;")

		logDatabase("RENTED-2 SELECT r.invid FROM rentlist r WHERE r.userid = "+strconv.FormatInt(UserID, 10), fmt.Sprint(invs))

		checkErr(err)
	} else {

		invs, err = db.Query("SELECT r.invid FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID, 10) + " AND r.Bookmarked = true;")

		logDatabase("RENTED-3 SELECT r.invid FROM rentlist r WHERE r.userid = "+strconv.FormatInt(UserID, 10), fmt.Sprint(invs))

		checkErr(err)
	}

	var ids []int

	//Save alle the IDs
	for invs.Next() {

		var inIDs int
		rows.Scan(inIDs)
		ids = append(ids, inIDs)

	}

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

		}
		logDatabase("SELECT * FROM equipment e WHERE e.invid ="+strconv.Itoa(element), fmt.Sprint(equipment))

	}

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

		res, err := db.Query("SELECT * FROM rentlist r WHERE r.invid = " + strconv.FormatInt(invID,10)+ ";")
		checkErr(err)
		for res.Next(){

			var inRent rentData
			res.Scan(&(inRent.RentID),&(inRent.UserID),&(inRent.InvID),&(inRent.RentDate),&(inRent.ReturnDate),&(inRent.Bookmarked),&(inRent.Amount),&(inRent.Repair))

			//Get the User Name
			var user string
			username, err := db.Query("SELECT name FROM users WHERE userid = " + strconv.FormatInt(inRent.UserID,10)+ ";")
			checkErr(err)
			username.Next()
			username.Scan(user)
			logDatabase("SELECT name FROM users WHERE userid = " + strconv.FormatInt(inRent.UserID,10)+ ";", user)
			inRent.RentedByUserName = user

			rentlist = append(rentlist,inRent)

		}
		logDatabase("SELECT * FROM rentlist r WHERE r.invid = " + strconv.FormatInt(invID,10)+ ";",fmt.Sprint(rentlist))
		return &rentlist

	//get all the rows with the given userid
	} else if UserID > 0 {

		res, err := db.Query("SELECT * FROM rentlist r WHERE r.userid = " + strconv.FormatInt(UserID,10)+ ";")
		checkErr(err)
		for res.Next(){

			var inRent rentData
			res.Scan(&(inRent.RentID),&(inRent.UserID),&(inRent.InvID),&(inRent.RentDate),&(inRent.ReturnDate),&(inRent.Bookmarked),&(inRent.Amount),&(inRent.Repair))

			//Get the User Name
			var user string
			username, err := db.Query("SELECT name FROM users WHERE userid = " + strconv.FormatInt(UserID,10)+ ";")
			checkErr(err)
			username.Next()
			username.Scan(user)
			logDatabase("SELECT name FROM users WHERE userid = " + strconv.FormatInt(UserID,10)+ ";", user)
			inRent.RentedByUserName = user

			rentlist = append(rentlist,inRent)

		}
		logDatabase("SELECT * FROM rentlist r WHERE r.invid = " + strconv.FormatInt(UserID,10)+ ";",fmt.Sprint(rentlist))
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
		if (validatePassword && inUser.Password == password) {
			return &inUser
		} else {
			log.Print("Error: The User could not be logged in")
			//TODO: What to do now?
		}
		logDatabase("SELECT * FROM users u WHERE u.Name LIKE '"+username+"';", fmt.Sprint(inUser))
		return &inUser
	}

	log.Print("Error: No User with the given username")
	return nil
}

//TODO: Gets the user, that the admin wants to edit, has to be used with session cookies or a flag in the database or query string on the URL
func getEditUser() *user {

	logDatabase("!!!DUMMY!!!", fmt.Sprint(nadmuser))

	return &nadmuser
}

//TODO: Still dummy, has to be used with session cookies
func getLoggedInUser() *user {

	//Get Cookie with UserID
	//THEN:
	//query := "SELECT * FROM users u WHERE u.userid = " + userid + ";"

	logDatabase("!!!DUMMY!!!", fmt.Sprint(nadmuser))

	return &nadmuser
}

//Gets the Cart Items, that the given User has in his cart
//TODO: Still dummy, has to be used with session cookies
func getCartItemsForUser(db *sql.DB, UserSessionCookie string) (*[]equipmentData, *user) {

	var eq []equipmentData
	eq = append(eq, nequip3)
	eq = append(eq, nequip1)
	eq = append(eq, nequip2)

	logDatabase("!!!DUMMY!!!", fmt.Sprint(eq)+"|"+fmt.Sprint(nadmuser))

	return &eq, &nadmuser
}

//TODO: Fill dummy, make use of hashing
func updateUser(db *sql.DB, user *user) bool {

	var test []byte;
	bcrypt.GenerateFromPassword(test, 0)

	return false
}
