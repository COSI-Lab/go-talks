package main

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
	Id          uint32
	Name        string
	Type        TalkType
	Description string
	IsHidden    bool
	Week        uint32
	Order       uint32
}
