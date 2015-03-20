package draftmark

import (
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"log"
	"testing"
	"time"
)

type FakeClient struct {
}

func (c *FakeClient) GetChanges() []dropbox.DropboxEntry {
	r := []dropbox.DropboxEntry{{Path: "/notes/test.md", IsDeleted: false, Modified: time.Now()}}
	return r
}

type FakeDb struct {
}

func (d *FakeDb) DeleteNote(note db.Note) bool {
	return true
}

func (d *FakeDb) SaveNote(note db.Note) bool {
	log.Println("I'm doing the fake SaveNote!")
	return true
}

func TestSync_NewFileAdded(t *testing.T) {
	client := &Sync{Dropbox: &FakeClient{}, Db: &FakeDb{}}

	user := User{1, "gammons@gmail.com", "asdf", "asdf", "123", time.Now(), time.Now()}
	client.DoSync(user, "/notes")
}
