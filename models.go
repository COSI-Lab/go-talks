package main

import "time"

// TalkType is the type of talk
type TalkType uint32

const (
	// ForumTopic are for general lab discussion
	ForumTopic TalkType = iota
	// LightningTalk are for short (5-10 minute) talks
	LightningTalk
	// ProjectUpdate are quick updates regarding ongoing projects
	ProjectUpdate
	// Announcement are for... announcements
	Announcement
	// AfterMeetingSlot holding events that happen after the meeting
	AfterMeetingSlot
)

func (tt TalkType) String() string {
	if tt > 4 {
		return "unknown"
	}
	return []string{"forum topic", "lightning talk", "project update", "announcement", "after meeting slot"}[tt]
}

// Talk is used to represent a talk in the database
type Talk struct {
	ID          uint32    `gorm:"AUTO_INCREMENT, primary key" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Type        TalkType  `gorm:"not null" json:"type"`
	Description string    `gorm:"not null" json:"description"`
	IsHidden    bool      `gorm:"not null" json:"-"`
	Week        string    `gorm:"index, not null" json:"-"`
	Order       uint32    `gorm:"not null" json:"-"` // TODO: Talk ordering
	CreatedAt   time.Time `json:"-"`
}
