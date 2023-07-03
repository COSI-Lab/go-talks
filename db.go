package main

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// db is a global db connection to be shared
var db *gorm.DB
var dbLock sync.Mutex

// ConnectDB sets up the initial connection to the database along with retrying attempts
func ConnectDB(config *Config) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	var err error
	db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	return err
}

// MakeDB sets up the db
func MakeDB() {
	dbLock.Lock()
	defer dbLock.Unlock()

	// Create all regular tables
	db.AutoMigrate(
		&Talk{},
	)
}

// DropTables drops everything in the db
func DropTables() {
	dbLock.Lock()
	defer dbLock.Unlock()

	// Drop tables in an order that won't invoke errors from foreign key constraints
	db.Migrator().DropTable(&Talk{})
}

// VisibleTalks returns all visible talks for a given week
// If week is empty, it will default to this week
func VisibleTalks(week string) []Talk {
	dbLock.Lock()
	defer dbLock.Unlock()

	if week == "" {
		week = nextWednesday()
	}

	var talks []Talk
	result := db.Where("is_hidden = false").Where("week = ?", week).Order("type").Find(&talks)

	if result.Error != nil {
		log.Println("[WARN] could not get visible talks:", result)
	}

	return talks
}

// AllTalks returns all talks for a given week
// If week is empty, it will default to this week
func AllTalks(week string) []Talk {
	dbLock.Lock()
	defer dbLock.Unlock()

	if week == "" {
		week = nextWednesday()
	}

	var talks []Talk
	result := db.Where("week = ?", week).Order("type").Find(&talks)

	if result.Error != nil {
		log.Println("[WARN] could not get all talks:", result)
	}

	return talks
}

// CreateTalk inserts a new talk into the db
func CreateTalk(talk *Talk) uint32 {
	dbLock.Lock()
	defer dbLock.Unlock()

	result := db.Create(talk)

	if result.Error != nil {
		log.Println("[WARN] could not create talk:", result)
	}

	log.Println("[INFO] Created talk {", talk.Name, talk.Description, talk.Type, talk.Week, talk.ID, "}")
	return talk.ID
}

// HideTalk updates a talk, setting its isHidden field to true
func HideTalk(id uint32) {
	dbLock.Lock()
	defer dbLock.Unlock()

	talk := Talk{}
	result := db.First(&talk, id)

	if result.Error != nil {
		log.Println("[WARN] could not find talk:", result)
	}

	talk.IsHidden = true
	result = db.Save(&talk)

	if result.Error != nil {
		log.Println("[WARN] could not hide talk:", result)
	}
}

// DeleteTalk deletes a talk from the db
func DeleteTalk(id uint32) {
	dbLock.Lock()
	defer dbLock.Unlock()

	talk := Talk{}
	result := db.First(&talk, id)

	if result.Error != nil {
		log.Println("[WARN] could not find talk:", result)
	}

	result = db.Delete(&talk)

	if result.Error != nil {
		log.Println("[WARN] could not delete talk:", result)
	}
}
