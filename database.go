package main

import "database/sql"

type user struct {
	OID        int           `gorm:"primary_key;column:OID"`
	Login      string        `gorm:"column:Login"`
	Password   string        `gorm:"column:Password"`
	Name       string        `gorm:"column:Name"`
	FirstName  string        `gorm:"column:FirstName"`
	Rfid       string        `gorm:"column:Rfid"`
	Role       string        `gorm:"column:Role"`
	Barcode    string        `gorm:"column:Barcode"`
	Pin        string        `gorm:"column:Pin"`
	Function   string        `gorm:"column:Function"`
	UserTypeID sql.NullInt32 `gorm:"column:UserTypeID"`
	Email      string        `gorm:"column:Email"`
	Phone      string        `gorm:"column:Phone"`
}

func (user) TableName() string {
	return "user"
}

type hvw_Zamestnanci struct {
	Cislo        string `gorm:"primary_key;column:Cislo"`
	Prijmeni     string `gorm:"column:Prijmeni"`
	Jmeno        string `gorm:"column:Jmeno"`
	_EVOLoginZam string `gorm:"column:_EVOLoginZam"`
	Serizovac    bool   `gorm:"column:Serizovac"`
}

func (hvw_Zamestnanci) TableName() string {
	return "hvw_Zamestnanci"
}
