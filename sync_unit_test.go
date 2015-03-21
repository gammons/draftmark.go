package draftmark

import (
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Fake client with new entry

type FakeClientWithNewEntry struct{}

func (c *FakeClientWithNewEntry) GetChanges() []dropbox.DropboxEntry {
	r := []dropbox.DropboxEntry{{Path: "/notes/test.md", IsDeleted: false, Modified: time.Now()}}
	return r
}

// Fake client with deleted entry

type FakeClientWithDeletedEntry struct{}

func (c *FakeClientWithDeletedEntry) GetChanges() []dropbox.DropboxEntry {
	r := []dropbox.DropboxEntry{{Path: "/notes/test.md", IsDeleted: true, Modified: time.Now()}}
	return r
}

// Fake DB

type FakeDb struct {
}

func (d *FakeDb) DeleteNote(note db.Note) bool {
	return true
}

func (d *FakeDb) SaveNote(note db.Note) bool {
	saveNoteCount++
	return true
}

var user = User{1, "gammons@gmail.com", "asdf", "asdf", "123", time.Now(), time.Now()}
var fdb = new(FakeDb)
var saveNoteCount = 0

func TestSync_NewFileAdded(t *testing.T) {
	client := &Sync{Dropbox: &FakeClientWithNewEntry{}, Db: fdb}

	client.DoSync(user, "/notes")
	assert.Equal(t, saveNoteCount, 1)
}

func TestSync_NoteDeleted(t *testing.T) {
	//client := &Sync{Dropbox: &FakeClientWithNewEntry{}, Db: fdb}

}
