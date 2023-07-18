package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
)

// Global talks object
var talks *Talks

// Talks stores all it's data in an append-only log
type Talks struct {
	// Logs are protected by a RWMutex
	sync.RWMutex
	// The file we're writing to (opened in append mode)
	encoder *json.Encoder
	// Maps talk IDs to talks
	talks map[uint32]*Talk

	// An index of talks by week
	weeks map[string][]*Talk

	// The current ID counter
	id uint32
}

// NewTalks creates a new talks object, loading the talks from the given reader
// and writing new events to the given writer
func NewTalks(r io.Reader, w io.Writer) (*Talks, error) {
	encoder := json.NewEncoder(w)

	// Create the talks object
	t := &Talks{
		RWMutex: sync.RWMutex{},
		encoder: encoder,
		talks:   make(map[uint32]*Talk),
		weeks:   make(map[string][]*Talk),
		id:      0,
	}

	// Load the talks from the file
	if err := t.load(r); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Talks) load(r io.Reader) error {
	// json decode the file
	dec := json.NewDecoder(r)
	for {
		var event TalkEvent
		if err := dec.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := t.event(event); err != nil {
			return err
		}
	}

	return nil
}

// Writes an event to the log
func (t *Talks) write(event TalkEvent) error {
	// Write the event to the file
	return t.encoder.Encode(event)
}

// Applies an event to the in-memory state
func (t *Talks) event(event TalkEvent) error {
	// switch on the event type
	switch event.Type {
	case Create:
		t.create(event.Create)
	case Hide:
		t.hide(event.Hide)
	case Delete:
		t.delete(event.Delete)
	default:
		return fmt.Errorf("unknown event type: %v", event.Type)
	}

	return nil
}

// Write and apply an event and handle any errors
func (t *Talks) writeAndApply(event TalkEvent) {
	log.Println("[INFO] Applying event:", event.Type, "[", event, "]")

	if err := t.write(event); err != nil {
		log.Println("[ERROR] Failed to write event:", err)
	}

	if err := t.event(event); err != nil {
		log.Println("[ERROR] Failed to apply event:", err)
	}
}

// Applies a create event to the in-memory state
func (t *Talks) create(c *CreateTalkEvent) {
	if c.ID > t.id {
		t.id = c.ID
	}

	t.talks[c.ID] = &Talk{
		ID:          c.ID,
		Name:        c.Name,
		Type:        c.Type,
		Description: c.Description,
		Week:        c.Week,
		Hidden:      false,
	}

	// Add the talk to the week index
	if _, ok := t.weeks[c.Week]; !ok {
		t.weeks[c.Week] = make([]*Talk, 0)
	}
	t.weeks[c.Week] = append(t.weeks[c.Week], t.talks[c.ID])

	sort.Slice(t.weeks[c.Week], func(i, j int) bool {
		if t.weeks[c.Week][i].Type == t.weeks[c.Week][j].Type {
			return t.weeks[c.Week][i].ID < t.weeks[c.Week][j].ID
		}
		return t.weeks[c.Week][i].Type < t.weeks[c.Week][j].Type
	})
}

// Create creates a new talk
func (t *Talks) Create(name string, talkType TalkType, description string, week string) uint32 {
	t.Lock()

	// Increment the ID counter
	t.id++

	id := t.id
	event := TalkEvent{
		Time: Now(),
		Type: Create,
		Create: &CreateTalkEvent{
			ID:          id,
			Name:        name,
			Type:        talkType,
			Description: description,
			Week:        week,
		},
	}
	t.writeAndApply(event)

	t.Unlock()
	return id
}

// Applies a hide event to the in-memory state
func (t *Talks) hide(h *HideTalkEvent) {
	if _, ok := t.talks[h.ID]; !ok {
		return
	}

	t.talks[h.ID].Hidden = true
}

// Hide hides a talk
func (t *Talks) Hide(id uint32) {
	t.Lock()

	event := TalkEvent{
		Time: Now(),
		Type: Hide,
		Hide: &HideTalkEvent{
			ID: id,
		},
	}
	t.writeAndApply(event)

	t.Unlock()
}

// Applies a delete event to the in-memory state
func (t *Talks) delete(d *DeleteTalkEvent) {
	if _, ok := t.talks[d.ID]; !ok {
		return
	}

	// Remove the talk from the week index
	week := t.talks[d.ID].Week
	for i, talk := range t.weeks[week] {
		if talk.ID == d.ID {
			t.weeks[week] = append(t.weeks[week][:i], t.weeks[week][i+1:]...)
		}
	}
	delete(t.talks, d.ID)
}

// Delete deletes a talk
func (t *Talks) Delete(id uint32) {
	t.Lock()

	event := TalkEvent{
		Time: Now(),
		Type: Delete,
		Delete: &DeleteTalkEvent{
			ID: id,
		},
	}
	t.writeAndApply(event)

	t.Unlock()
}

// AllTalks returns all talks for a given week
func (t *Talks) AllTalks(week string) []*Talk {
	t.RLock()
	talks := t.weeks[week]
	t.RUnlock()

	return talks
}

// VisibleTalks returns all visible talks for a given week
func (t *Talks) VisibleTalks(week string) []*Talk {
	t.RLock()
	talks := make([]*Talk, 0)
	for _, talk := range t.weeks[week] {
		if !talk.Hidden {
			talks = append(talks, talk)
		}
	}
	t.RUnlock()

	return talks
}
