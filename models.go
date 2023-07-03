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
	return []string{"forum topic", "lightning talk", "project update", "announcement", "after meeting slot"}[tt]
}

type Talk struct {
	Id          uint32    `gorm:"AUTO_INCREMENT, primary key" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Type        TalkType  `gorm:"not null" json:"type"`
	Description string    `gorm:"not null" json:"description"`
	IsHidden    bool      `gorm:"not null" json:"-"`
	Week        string    `gorm:"index, not null" json:"-"`
	Order       uint32    `gorm:"not null" json:"-"` // TODO: Talk ordering
	CreatedAt   time.Time `json:"-"`
}
