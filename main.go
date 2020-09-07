package main

import (
	"github.com/kardianos/service"
	"time"
)

const version = "2020.3.3.7"
const serviceName = "Linaplast User Import Service"
const serviceDescription = "Download users from Helios database"
const zapsiConfig = "user=zapsi_uzivatel password=zapsi dbname=zapsi2 host=zapsidatabase port=3306 sslmode=disable"
const heliosConfig = "sqlserver://zapsi:Zapsi@sql14.linaplast.local\\sql2014:1433?database=Helios002"
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
