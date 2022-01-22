package main

import "time"

type TalkType uint32

const (
	FORUM_TOPIC TalkType = iota
	LIGHTNING_TALK
	PROJECT_UPDATE
	ANNOUNCEMENT
	AFTER_MEETING_SLOT
)

func (tt TalkType) String() string {
	if tt > 4 {
		return "unknown"
	}
	return []string{"forum topic", "lightning talk", "project update", "announcment", "after meeting slot"}[tt]
}

type Talk struct {
	Id          uint32   `gorm:"AUTO_INCREMENT, primary key"`
	Name        string   `gorm:"not null"`
	Type        TalkType `gorm:"not null"`
	Description string   `gorm:"not null"`
	IsHidden    bool     `gorm:"not null"`
	Week        uint32   `gorm:"index, not null"`
	Order       uint32   `gorm:"not null"`
	CreatedAt   time.Time
}
