package main

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global db connection to be shared
var DB *gorm.DB
var DB_LOCK sync.Mutex

// ConnectDB sets up the initial connection to the database along with retrying attempts
func ConnectDB(config *Config) error {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	var err error
	DB, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	return err
}

// MakeDB sets up the db
func MakeDB() {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	// Create all regular tables
	DB.AutoMigrate(
		&Talk{},
	)
}

// DropTables drops everything in the db
func DropTables() {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	// Drop tables in an order that won't invoke errors from foreign key constraints
	DB.Migrator().DropTable(&Talk{})
}

func VisibleTalks(week string) []Talk {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	if week == "" {
		week = nextWednesday()
	}

	var talks []Talk
	result := DB.Where("is_hidden = false").Where("week = ?", week).Order("type").Find(&talks)

	if result.Error != nil {
		log.Println("[WARN] could not get visible talks:", result)
	}

	return talks
}

func AllTalks(week string) []Talk {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	if week == "" {
		week = nextWednesday()
	}

	var talks []Talk
	result := DB.Where("week = ?", week).Order("type").Find(&talks)

	if result.Error != nil {
		log.Println("[WARN] could not get all talks:", result)
	}

	return talks
}

func CreateTalk(talk *Talk) uint32 {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	result := DB.Create(talk)

	if result.Error != nil {
		log.Println("[WARN] could not create talk:", result)
	}

	log.Println("[INFO] Created talk {", talk.Name, talk.Description, talk.Type, talk.Week, talk.Id, "}")
	return talk.Id
}

func HideTalk(id uint32) {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	talk := Talk{}
	result := DB.First(&talk, id)

	if result.Error != nil {
		log.Println("[WARN] could not find talk:", result)
	}

	talk.IsHidden = true
	result = DB.Save(&talk)

	if result.Error != nil {
		log.Println("[WARN] could not hide talk:", result)
	}
}

func DeleteTalk(id uint32) {
	DB_LOCK.Lock()
	defer DB_LOCK.Unlock()

	talk := Talk{}
	result := DB.First(&talk, id)

	if result.Error != nil {
		log.Println("[WARN] could not find talk:", result)
	}

	result = DB.Delete(&talk)

	if result.Error != nil {
		log.Println("[WARN] could not delete talk:", result)
	}
}
