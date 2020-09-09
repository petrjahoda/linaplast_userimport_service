package main

import (
	"github.com/kardianos/service"
	"time"
)

const version = "2020.3.3.9"
const serviceName = "Linaplast User Import Service"
const serviceDescription = "Download users from Helios database and imports them into Zapsi database"
const zapsiConfig = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
const heliosConfig = "sqlserver://zapsi:Zapsi@172.16.1.17/SQLHEO?database=Helios002"
const downloadInSeconds = 60

var serviceRunning = false
var processRunning = false

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	logInfo("MAIN", serviceDescription)
	logInfo("MAIN", heliosConfig)
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
