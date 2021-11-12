package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global db connection to be shared
var DB *gorm.DB

// checkExists checks if a value exists and fails if it doesn't
func checkExists(exists bool, msg string) {
	if !exists {
		log.Fatal(msg)
	}
}

// ConnectDB sets up the initial connection to the database along with retrying attempts
func ConnectDB(dbType string) error {
	// Get user credentials
	dbTypeUpper := strings.ToUpper(dbType)

	if dbType == "sqlite" {
		dbName, exists := os.LookupEnv(dbTypeUpper + "_DB")
		checkExists(exists, "Couldn't find database")

		var err error
		DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
		return err
	} else {
		user, exists := os.LookupEnv(dbTypeUpper + "_USER")
		checkExists(exists, "Couldn't find database user")
		password, exists := os.LookupEnv(dbTypeUpper + "_PASSWORD")
		checkExists(exists, "Couldn't find database password")

		// Get database params
		dbServer, exists := os.LookupEnv(dbTypeUpper + "_SERVER")
		checkExists(exists, "Couldn't find database server")
		dbPort, exists := os.LookupEnv(dbTypeUpper + "_PORT")
		checkExists(exists, "Couldn't find database port")
		dbName, exists := os.LookupEnv(dbTypeUpper + "_DB")
		checkExists(exists, "Couldn't find database name")
		connectionString := fmt.Sprintf(
			"sslmode=disable host=%s port=%s dbname=%s user=%s password=%s",
			dbServer,
			dbPort,
			dbName,
			user,
			password,
		)

		// Check how many times to try the db before quitting
		attemptsStr, exists := os.LookupEnv("DB_ATTEMPTS")
		if !exists {
			attemptsStr = "5"
		}

		attempts, err := strconv.Atoi(attemptsStr)
		if err != nil {
			attempts = 5
		}

		timeoutStr, exists := os.LookupEnv("DB_CONNECTION_TIMEOUT")
		if !exists {
			timeoutStr = "5"
		}
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			timeout = 5
		}

		for i := 1; i <= attempts; i++ {
			DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
			if err != nil {
				if i != attempts {
					log.Printf(
						"WARNING: Could not connect to db on attempt %d. Trying again in %d seconds.\n",
						i,
						timeout,
					)
				} else {
					return fmt.Errorf("could not connect to db after %d attempts", attempts)
				}
				time.Sleep(time.Duration(timeout) * time.Second)
			} else {
				// No error to worry about
				break
			}
		}
		log.Println("Connection to db succeeded!")
	}
	return nil
}

// MakeDB sets up the db
func MakeDB() {
	// Create all regular tables
	DB.AutoMigrate(
		&Talk{},
	)
}

// DropTables drops everything in the db
func DropTables() {
	// Drop tables in an order that won't invoke errors from foreign key constraints
	DB.Migrator().DropTable(&Talk{})
}
