package draftmark

import (
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"strings"
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

type Sync struct {
	Db      db.DBClient
	Dropbox dropbox.DropboxClient
}

func init() {
	setupDotEnv()
}

func NewSync() *Sync {
	return &Sync{Db: &db.Client{}, Dropbox: &dropbox.Client{}}
}

func (s *Sync) DoSync(user User, prefix string) {
	entries := s.Dropbox.GetChanges(user.DropboxCursor, prefix)
	for _, entry := range entries {
		log.Println("Entry is ", entry)
		if !strings.HasPrefix(entry.Path, prefix) || !strings.HasSuffix(entry.Path, ".md") {
			log.Println("entry does not meet suffix and prefix requirements")
			continue
		}

		if entry.IsDeleted {
			s.deleteEntry(&user, entry)
		} else {
			s.createOrUpdateEntry(&user, entry)
		}
	}
}

func (s *Sync) deleteEntry(user *User, entry *dropbox.DropboxEntry) {
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified}
	s.Db.DeleteNote(note)
}

func (s *Sync) createOrUpdateEntry(user *User, entry *dropbox.DropboxEntry) {
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified}
	s.Db.SaveNote(note)
}

func setupDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
