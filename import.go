package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

func ImportUsersFromHelios() {
	timer := time.Now()
	logInfo("MAIN", "Importing process started")
	zapsiUsers, downloadedFromZapsi := DownloadUsersFromZapsi()
	heliosUsers, downloadedFromHelios := DownloadUsersFromHelios()
	if downloadedFromZapsi && downloadedFromHelios {
		sort.Slice(zapsiUsers, func(i, j int) bool {
			return zapsiUsers[i].Login <= zapsiUsers[j].Login
		})
		sort.Slice(heliosUsers, func(i, j int) bool {
			return heliosUsers[i].Cislo <= heliosUsers[j].Cislo
		})
		logInfo("MAIN", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		logInfo("MAIN", "Helios Users: "+strconv.Itoa(len(heliosUsers)))
		UpdateUsers(heliosUsers, zapsiUsers)
	}
	logInfo("MAIN", "Importing process complete, time elapsed: "+time.Since(timer).String())
}

func UpdateUsers(heliosUsers []hvw_Zamestnanci, zapsiUsers []user) {
	timer := time.Now()
	logInfo("MAIN", "Updating users")
	for _, heliosUser := range heliosUsers {
		if serviceRunning {
			index, userInZapsi := BinarySearchUser(zapsiUsers, heliosUser)
			if userInZapsi {
				UpdateUserInZapsi(heliosUser, zapsiUsers[index])
			} else {
				CreateUserInZapsi(heliosUser)
			}
		}
	}
	logInfo("MAIN", "Users updated, time elapsed: "+time.Since(timer).String())
}

func UpdateUserInZapsi(heliosUser hvw_Zamestnanci, zapsiUser user) {
	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User exists in Zapsi, updating to "+heliosUser.EVOLoginZam)
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	var userTypeIdToInsert int
	updateUserType := false
	if heliosUser.Serizovac {
		userTypeIdToInsert = 2
		updateUserType = true
	}
	db.Model(&user{}).Where(user{Login: zapsiUser.Login}).Updates(user{
		Name:       heliosUser.Prijmeni,
		FirstName:  heliosUser.Jmeno,
		Rfid:       heliosUser.EVOLoginZam,
		Barcode:    heliosUser.EVOLoginZam,
		Pin:        heliosUser.EVOLoginZam,
		UserTypeID: sql.NullInt32{Int32: int32(userTypeIdToInsert), Valid: updateUserType},
	})
}

func CreateUserInZapsi(heliosUser hvw_Zamestnanci) {
	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User does not exist in Zapsi, creating...")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var user user
	user.Login = heliosUser.Cislo
	user.FirstName = heliosUser.Jmeno
	user.Name = heliosUser.Prijmeni
	user.Rfid = heliosUser.EVOLoginZam
	user.Barcode = heliosUser.EVOLoginZam
	user.Pin = heliosUser.EVOLoginZam
	if heliosUser.Serizovac {
		user.UserTypeID = sql.NullInt32{Int32: 2, Valid: true}
	} else {
		user.UserTypeID = sql.NullInt32{Int32: 1, Valid: true}
	}
	user.UserRoleID = sql.NullInt32{Int32: 2, Valid: true}
	db.Save(&user)
	return
}

func BinarySearchUser(zapsiUsers []user, heliosUser hvw_Zamestnanci) (int, bool) {
	index := sort.Search(len(zapsiUsers), func(i int) bool { return zapsiUsers[i].Login >= heliosUser.Cislo })
	userInZapsi := index < len(zapsiUsers) && zapsiUsers[index].Login == heliosUser.Cislo
	return index, userInZapsi
}

func DownloadUsersFromHelios() ([]hvw_Zamestnanci, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading users from Helios")
	db, err := gorm.Open(sqlserver.Open(heliosConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []hvw_Zamestnanci{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var users []hvw_Zamestnanci
	db.Find(&users)
	logInfo("MAIN", "Helios users downloaded, time elapsed: "+time.Since(timer).String())
	return users, true
}

func DownloadUsersFromZapsi() ([]user, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading users from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []user{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var users []user
	db.Find(&users)
	logInfo("MAIN", "Zapsi users downloaded, time elapsed: "+time.Since(timer).String())
	return users, true
}
