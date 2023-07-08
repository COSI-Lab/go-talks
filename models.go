package main

import (
	"encoding/json"
	"time"
)

// TalkEventType is an enum of the different types of events
type TalkEventType int

const (
	// UnknownEvent is for when we can't parse the event type
	UnknownEvent TalkEventType = iota
	// Create is the talk creation event {type: "create", create: {id: 1, name: "foo", ...}
	Create
	// Hide is the talk hiding event {type: "hide", hide: {id: 1}}
	Hide
	// Delete is the talk deletion event {type: "delete", delete: {id: 1}}
	Delete
)

// MarshalJSON implements the json.Marshaler interface
func (t TalkEventType) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case Create:
		s = "create"
	case Hide:
		s = "hide"
	case Delete:
		s = "delete"
	default:
		s = "unknown"
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *TalkEventType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "create":
		*t = Create
	case "hide":
		*t = Hide
	case "delete":
		*t = Delete
	default:
		*t = UnknownEvent
	}
	return nil
}

// TalkType is the type of talk
type TalkType int

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

// MarshalJSON implements the json.Marshaler interface
func (t TalkType) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case ForumTopic:
		s = "forum topic"
	case LightningTalk:
		s = "lightning talk"
	case ProjectUpdate:
		s = "project update"
	case Announcement:
		s = "announcement"
	case AfterMeetingSlot:
		s = "after meeting slot"
	default:
		s = "unknown"
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *TalkType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "forum topic":
		*t = ForumTopic
	case "lightning talk":
		*t = LightningTalk
	case "project update":
		*t = ProjectUpdate
	case "announcement":
		*t = Announcement
	case "after meeting slot":
		*t = AfterMeetingSlot
	}
	return nil
}

// String implements the fmt.Stringer interface. This is used when templating
func (t TalkType) String() string {
	switch t {
	case ForumTopic:
		return "forum topic"
	case LightningTalk:
		return "lightning talk"
	case ProjectUpdate:
		return "project update"
	case Announcement:
		return "announcement"
	case AfterMeetingSlot:
		return "after meeting slot"
	default:
		return "unknown"
	}
}

// JSONTime is a Wrapper for time.Time struct with custom JSON behavior
// matching sqlite's datetime format
type JSONTime struct {
	time.Time
}

var format = "2006-01-02 15:04:05.999999-07:00"

// MarshalJSON implements the json.Marshaler interface
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(format))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	tt, err := time.Parse(format, s)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// Now returns a JSONTime with the time from time.Now()
func Now() JSONTime {
	return JSONTime{time.Now()}
}

// TalkEvent is stored in the database. Since go doesn't have a good way to
// represent a union type, we'll just use a large struct with a type field
// and a field for each event type
type TalkEvent struct {
	Time   JSONTime         `json:"time"`
	Type   TalkEventType    `json:"type"`
	Create *CreateTalkEvent `json:"create,omitempty"`
	Hide   *HideTalkEvent   `json:"hide,omitempty"`
	Delete *DeleteTalkEvent `json:"delete,omitempty"`
}

// CreateTalkEvent is created when a talk is created
type CreateTalkEvent struct {
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	Type        TalkType `json:"type"`
	Description string   `json:"description"`
	Week        string   `json:"week"`
}

// HideTalkEvent is created when a talk is hidden
type HideTalkEvent struct {
	ID uint32 `json:"id"`
}

// DeleteTalkEvent is created when a talk is deleted
type DeleteTalkEvent struct {
	ID uint32 `json:"id"`
}

// Talk is the resulting type produced by the database
type Talk struct {
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	Type        TalkType `json:"type"`
	Description string   `json:"description"`
	Week        string   `json:"week"`
	Hidden      bool     `json:"hidden"`
}
