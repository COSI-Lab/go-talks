package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Saves weekly backups of the database going back 4 weeks
func backup() {
	dbLock.Lock()
	defer dbLock.Unlock()

	// Get the current date
	t := time.Now().In(tz)
	date := t.Format("2006-01-02")

	// File name for the backup
	fileName := fmt.Sprintf("backups/backup-%s.db", date)

	// Copy the database to the backup file
	cmd := fmt.Sprintf("cp talks.db %s", fileName)
	err := exec.Command("sh", "-c", cmd).Run()

	if err != nil {
		log.Println("[ERROR] Failed to backup database:", err)
	}

	// Remove backups older than 4 weeks
	cmd = "find backups -type f -mtime +28 -exec rm {} \\;"
	err = exec.Command("sh", "-c", cmd).Run()

	if err != nil {
		log.Println("[ERROR] Failed to remove old backups:", err)
	}
}
