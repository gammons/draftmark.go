package draftmark

import (
	"draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type User struct {
	ID                 int
	Email              string
	DropboxCursor      string
	DropboxAccessToken string
	DropboxUserId      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Note struct {
	ID        int
	Title     string
	content   string `sql:"type:text"`
	UserId    int
	Path      string
	Mtime     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	setupDotEnv()
}

var dbClient = new(db.Client)

func Sync(user User, prefix string, client dropbox_client.Client) {
	entries := client.GetChanges(user.DropboxAccessToken, user.DropboxCursor, "/notes")
	for _, entry := range entries {
		if entry.IsDeleted {
			deleteEntry(&entry)
		} else {
			createOrUpdateEntry(&entry)
		}
	}
}

func deleteEntry(entry *dropbox_client.DropboxEntry) {
	note := db.Note{Path: entry.Path, Mtime: entry.Modified}
	dbClient.DeleteNote(note)
}

func createOrUpdateEntry(entry *dropbox_client.DropboxEntry) {
	note := db.Note{Path: entry.Path, Mtime: entry.Modified}
	dbClient.SaveNote(note)
}

func setupDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
