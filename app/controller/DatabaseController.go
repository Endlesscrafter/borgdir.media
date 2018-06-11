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

	rent, err := db.Query("SELECT * FROM rentlist r WHERE r.invid = " + strconv.FormatInt(invID, 10) + ";")

	checkErr(err)

	var inEquip equipmentData

	for rows.Next() {

		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.InvID), &(inEquip.StorageLocation))

	}
	if !rows.Next() {

		log.Print("The Product with ID " + strconv.FormatInt(invID, 10) + " wasn't found in the Database")
		return nil

	}

	if rent.Next() {

		var irr string

		rent.Scan(&(inEquip.RentedByUserID),&irr,&(inEquip.RentDate),&(inEquip.ReturnDate),&(inEquip.Bookmarked),&irr,&(inEquip.Repair))

	} else{

		inEquip.Rented = false
		inEquip.RentDate = "0"
		inEquip.ReturnDate = "0"
		inEquip.Bookmarked = false
		inEquip.Repair = false

	}

	return &inEquip

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
func getRentedEquip(db *sql.DB, UserID int64, bookmarked bool) *[]equipmentData {

	var rows *sql.Rows
	var err error

	if (!bookmarked) {
		rows, err = db.Query("SELECT * FROM equipment e WHERE e.RentedbyUserID = " + strconv.FormatInt(UserID, 10) + " AND e.Bookmarked = false;") //Works together with template
	} else {
		rows, err = db.Query("SELECT * FROM equipment e WHERE e.RentedbyUserID = " + strconv.FormatInt(UserID, 10) + " AND e.Bookmarked = true;")
	}
	checkErr(err)

	var equipment []equipmentData
	for rows.Next() {

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.Rented), &(inEquip.Bookmarked), &(inEquip.Repair),
			&(inEquip.RentedByUserID), &(inEquip.RentedByUserName), &(inEquip.RentDate), &(inEquip.ReturnDate),
			&(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		equipment = append(equipment, inEquip)

	}
	logDatabase("SELECT * FROM equipment e WHERE e.RentedbyUserID = "+strconv.FormatInt(UserID, 10)+" ((AND e.Bookmarked = true));", fmt.Sprint(equipment))
	return &equipment

}

//Gets products, that are not rented and therefore can be rented
func getAvailableEqip(db *sql.DB) *[]equipmentData {

	rows, err := db.Query("SELECT * FROM equipment e WHERE e.Rented = false AND e.Bookmarked = false;")

	checkErr(err)

	var equipment []equipmentData
	for rows.Next() {

		var inEquip equipmentData
		rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
			&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
			&(inEquip.FeaturedImageSRC), &(inEquip.Rented), &(inEquip.Bookmarked), &(inEquip.Repair),
			&(inEquip.RentedByUserID), &(inEquip.RentedByUserName), &(inEquip.RentDate), &(inEquip.ReturnDate),
			&(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
		equipment = append(equipment, inEquip)

	}
	logDatabase("SELECT * FROM equipment e WHERE e.Rented = false AND e.Bookmarked = false;", fmt.Sprint(equipment))
	return &equipment

}

//Gets the featured products, in correct order
func getFeaturedProducts(db *sql.DB) *[]equipmentData {

	var equipment []equipmentData

	for i := 1; i <= 4; i++ {
		rows, err := db.Query("SELECT * FROM equipment e WHERE e.Featured = true AND e.FeaturedID = " + strconv.Itoa(i) + ";")

		checkErr(err)

		for rows.Next() {

			var inEquip equipmentData
			rows.Scan(&(inEquip.Name), &(inEquip.Desc), &(inEquip.ImageSRC), &(inEquip.ImageAlt), &(inEquip.Stock),
				&(inEquip.StockAmount), &(inEquip.Category), &(inEquip.Featured), &(inEquip.FeaturedID),
				&(inEquip.FeaturedImageSRC), &(inEquip.Rented), &(inEquip.Bookmarked), &(inEquip.Repair),
				&(inEquip.RentedByUserID), &(inEquip.RentedByUserName), &(inEquip.RentDate), &(inEquip.ReturnDate),
				&(inEquip.InvID), &(inEquip.StorageLocation), &(inEquip.EquipmentOwnerID))
			equipment = append(equipment, inEquip)

		}
	}

	logDatabase("SELECT * FROM equipment e WHERE e.Featured = true AND e.FeaturedID = *;", fmt.Sprint(equipment))
	return &equipment

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
	eq = append(eq, nequip1)
	eq = append(eq, nequip2)
	eq = append(eq, nequip3)

	logDatabase("!!!DUMMY!!!", fmt.Sprint(eq)+"|"+fmt.Sprint(nadmuser))

	return &eq, &nadmuser
}

//TODO: Fill dummy, make use of hashing
func updateUser(db *sql.DB, user *user) bool {

	var test []byte;
	bcrypt.GenerateFromPassword(test, 0)

	return false
}
