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
	logInfo("import", "Importing data")
	zapsiUsers, downloadedFromZapsi := DownloadUsersFromZapsi()
	heliosUsers, downloadedFromHelios := DownloadUsersFromHelios()
	sort.Slice(zapsiUsers, func(i, j int) bool {
		return zapsiUsers[i].Barcode <= zapsiUsers[j].Barcode
	})
	sort.Slice(heliosUsers, func(i, j int) bool {
		return heliosUsers[i].Cislo <= heliosUsers[j].Cislo
	})

	if downloadedFromZapsi && downloadedFromHelios {
		logInfo("import", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		logInfo("import", "Helios Users: "+strconv.Itoa(len(heliosUsers)))
		UpdateUsers(heliosUsers, zapsiUsers)
	}
	logInfo("import", "Data imported, time elapsed: "+time.Since(timer).String())
}

func UpdateUsers(heliosUsers []hvw_Zamestnanci, zapsiUsers []user) {
	timer := time.Now()
	logInfo("import", "Updating Users")
	for _, heliosUser := range heliosUsers {
		if serviceRunning {
			index, userInZapsi := BinarySearchUser(zapsiUsers, heliosUser)
			if userInZapsi {
				logInfo("import", heliosUser.Prijmeni+": User exists, just updating...")
				UpdateUserInZapsi(heliosUser, zapsiUsers[index])

			} else {
				CreateZapsiUserFrom(heliosUser)
			}
		}
	}
	logInfo("import", "Users updated, time elapsed: "+time.Since(timer).String())
}

func UpdateUserInZapsi(heliosUser hvw_Zamestnanci, zapsiUser user) {
	timer := time.Now()
	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User exists, updating...")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var userTypeIdToInsert int32 = 1
	if heliosUser.Serizovac {
		userTypeIdToInsert = 2
	}
	db.Model(&user{}).Where(user{Login: zapsiUser.Login}).Updates(user{
		Login:      heliosUser.Cislo,
		Name:       heliosUser.Prijmeni,
		FirstName:  heliosUser.Jmeno,
		Rfid:       heliosUser.Cislo,
		Barcode:    heliosUser.Cislo,
		Pin:        heliosUser.Cislo,
		UserTypeID: sql.NullInt32{Int32: userTypeIdToInsert, Valid: true},
	})
	logInfo("import", "User updated, time elapsed: "+time.Since(timer).String())
}

func CreateZapsiUserFrom(heliosUser hvw_Zamestnanci) {
	timer := time.Now()
	logInfo("MAIN", heliosUser.Jmeno+" "+heliosUser.Prijmeni+": User does not exist, creating...")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var user user
	user.FirstName = heliosUser.Jmeno
	user.Name = heliosUser.Prijmeni
	user.Rfid = heliosUser.Cislo
	user.Barcode = heliosUser.Cislo
	user.Pin = heliosUser.Cislo
	user.Login = heliosUser.Cislo
	user.Role = "user"
	user.Function = "Operator"
	if heliosUser.Serizovac {
		user.UserTypeID = sql.NullInt32{
			Int32: 2,
			Valid: true,
		}
	} else {
		user.UserTypeID = sql.NullInt32{
			Int32: 1,
			Valid: true,
		}
	}
	db.Save(&user)
	logInfo("import", "User created, time elapsed: "+time.Since(timer).String())
	return
}

func BinarySearchUser(zapsiUsers []user, heliosUser hvw_Zamestnanci) (int, bool) {
	index := sort.Search(len(zapsiUsers), func(i int) bool { return zapsiUsers[i].Barcode >= heliosUser.Cislo })
	userInZapsi := index < len(zapsiUsers) && zapsiUsers[index].Barcode == heliosUser.Cislo
	return index, userInZapsi
}

func DownloadUsersFromHelios() ([]hvw_Zamestnanci, bool) {
	timer := time.Now()
	logInfo("import", "Downloading data from Helios")
	db, err := gorm.Open(sqlserver.Open(heliosConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []hvw_Zamestnanci{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	logInfo("import", "Helias database connected")
	var users []hvw_Zamestnanci
	db.Table("user").Find(&users)
	logInfo("import", "Helios users downloaded, time elapsed: "+time.Since(timer).String())
	return users, true
}

func DownloadUsersFromZapsi() ([]user, bool) {
	timer := time.Now()
	logInfo("import", "Downloading data from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []user{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	logInfo("import", "Zapsi database connected")
	var users []user
	db.Table("user").Find(&users)
	logInfo("import", "Zapsi users downloaded, time elapsed: "+time.Since(timer).String())
	return users, true
}
