package main

import (
	"bytes"
	"testing"
)

// Creates a memory-backed DB for testing
func CreateDB(events []TalkEvent) *Talks {
	// Create a writer to save the events
	buf := bytes.NewBuffer(nil)

	// Create the DB
	db, err := NewTalks(bytes.NewReader(nil), buf)
	if err != nil {
		panic(err)
	}

	// Apply the events
	for _, event := range events {
		if err := db.event(event); err != nil {
			panic(err)
		}

		// Write the event to the buffer
		if err := db.write(event); err != nil {
			panic(err)
		}
	}

	return db
}

func TestDB(t *testing.T) {
	events := []TalkEvent{
		{Type: Create, Create: &CreateTalkEvent{
			ID:          0,
			Name:        "",
			Type:        0,
			Description: "",
			Week:        "20230705",
		}},
		{Type: Create, Create: &CreateTalkEvent{
			ID:          1,
			Name:        "Test Talk",
			Type:        0,
			Description: "This is a test talk",
			Week:        "20230705",
		}},
		{Type: Hide, Hide: &HideTalkEvent{
			ID: 1,
		}},
		{Type: Delete, Delete: &DeleteTalkEvent{
			ID: 0,
		}},
		{Type: Create, Create: &CreateTalkEvent{
			ID:          3,
			Name:        "Test Talk 2",
			Type:        0,
			Description: "This is a test talk",
			Week:        "20230705",
		}},
	}

	db := CreateDB(events)

	// Check the state
	if len(db.talks) != 2 {
		t.Fatal("Expected 2 talks")
	}

	// Verify that the talk was hidden
	if !db.talks[1].Hidden {
		t.Fatal("Expected talk to be hidden")
	}

	// Verify that the talk was deleted
	if _, ok := db.talks[0]; ok {
		t.Fatal("Expected talk to be deleted")
	}

	// Verify that talk 2 doesn't exist
	if _, ok := db.talks[2]; ok {
		t.Fatal("Expected talk to be missing")
	}
}

// Verify that deleting a talk that doesn't exist doesn't cause a panic
func TestDeleteNonExistentTalk(t *testing.T) {
	events := []TalkEvent{
		{Type: Delete, Delete: &DeleteTalkEvent{
			ID: 0,
		}},
	}

	db := CreateDB(events)

	// Check the state
	if len(db.talks) != 0 {
		t.Fatal("Expected 0 talks")
	}
}

// Verify that hiding a talk that doesn't exist doesn't cause a panic
func TestHideNonExistentTalk(t *testing.T) {
	events := []TalkEvent{
		{Type: Hide, Hide: &HideTalkEvent{
			ID: 0,
		}},
	}

	db := CreateDB(events)

	// Check the state
	if len(db.talks) != 0 {
		t.Fatal("Expected 0 talks")
	}
}

// Verify that leaving holes in the ID sequence doesn't cause a panic
func TestHolesInIDSequence(t *testing.T) {
	events := []TalkEvent{
		{Type: Create, Create: &CreateTalkEvent{
			ID:          0,
			Name:        "",
			Type:        0,
			Description: "",
			Week:        "20230705",
		}},
		{Type: Create, Create: &CreateTalkEvent{
			ID:          2,
			Name:        "Test Talk",
			Type:        0,
			Description: "This is a test talk",
			Week:        "20230705",
		}},
	}

	db := CreateDB(events)

	// Check the state
	if len(db.talks) != 2 {
		t.Fatal("Expected 2 talks")
	}

	// assert that talk 0 has the correct info
	if db.talks[0] == nil {
		t.Fatal("Expected talk 0 to exist")
	}

	if db.talks[1] != nil {
		t.Fatal("Expected talk 1 to not exist")
	}

	if db.talks[2] == nil {
		t.Fatal("Expected talk 2 to exist")
	}

	// Create a new event and expect it to have ID 3
	db.Create("Test Talk 2", ForumTopic, "This is a test talk", "20230705")

	if db.talks[3] == nil {
		t.Fatal("Expected talk 3 to exist")
	}
}
