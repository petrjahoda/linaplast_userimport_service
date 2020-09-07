package main

import (
	"github.com/kardianos/service"
	"sort"
	"strconv"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

)

const version = "2020.3.2.29"
const serviceName = "Linaplast User Import Service"
const serviceDescription = "Download users from Helios database"
const zapsiConfig = "user=zapsi_uzivatek password=zapsi dbname=zapsi2 host=zapsidatabase port=3306 sslmode=disable"
const heliosConfig = "user=postgres password=Zps05..... dbname=postgres host=database port=5432 sslmode=disable"
const downloadInSeconds = 86400

var serviceRunning = false
var processRunning = false

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	program := &program{}
	s, err := service.New(program, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	serviceRunning = true
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceRunning = false
	if processRunning {
		logInfo("MAIN", serviceName+" ["+version+"] stopping...")
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	for serviceRunning {
		processRunning = true
		start := time.Now()
		logInfo("MAIN", serviceName+" ["+version+"] running")
		ImportUsersFromHelios()
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		logInfo("MAIN", "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
		processRunning = false
	}
}

func ImportUsersFromHelios() {
	timer := time.Now()
	logInfo("import", "Importing data")
	zapsiUsers, downloadedFromZapsi := DownloadUsersFromZapsi()
	heliosUsers, downloadedFromHelios := DownloadUsersFromHelios()
	sort.Slice(zapsiUsers, func(i, j int) bool {
		return zapsiUsers[i].Barcode <= zapsiUsers[j].Barcode
	})
		sort.Slice(heliosUsers, func(i, j int) bool {
		return heliosUsers[i].OsC <= heliosUsers[j].OsC
	})

	if downloadedFromZapsi && downloadedFromHelios {
		logInfo("import", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		logInfo("import", "Helios Users: "+strconv.Itoa(len(heliosUsers)))
		UpdateUsers(k2Users, zapsiUsers)
		UpdateUsersFour(k2UsersFour, zapsiUsers)
		UpdateOrders(k2Orders, zapsiOrders)
	}
	logInfo("import", "Data imported, time elapsed: "+time.Since(timer).String())
}

func DownloadUsersFromHelios() (interface{}, interface{}) {
	timer := time.Now()
	logInfo("import", "Downloading data from Helios")
	db, err := gorm.Open(sqlserver.Open(heliosConfig), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return []user{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	logInfo("import", "Helias database connected")
	var users []user
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
