package draftmark

import (
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Sync struct {
	Db      db.DBClient
	Dropbox dropbox.DropboxClient
}

func init() {
	setupDotEnv()
}

func NewSync(accessToken string) *Sync {
	dbox := &dropbox.Client{AccessToken: accessToken}
	dbox.InitDropbox()
	dbase := db.NewClient()
	return &Sync{Db: dbase, Dropbox: dbox}
}

func (s *Sync) DoSync(user db.User, prefix string) {
	nextCursor, entries := s.Dropbox.GetChanges(&user.DropboxCursor, prefix)

	s.Db.UpdateUserCursor(&user, nextCursor)

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Path, prefix) || !strings.HasSuffix(entry.Path, ".md") {
			continue
		}

		if entry.IsDeleted {
			s.deleteEntry(&user, entry)
		} else {
			s.createOrUpdateEntry(&user, entry)
		}
	}
}

func (s *Sync) deleteEntry(user *db.User, entry *dropbox.DropboxEntry) {
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified}
	s.Db.DeleteNote(&note)
}

func (s *Sync) createOrUpdateEntry(user *db.User, entry *dropbox.DropboxEntry) {
	content, _ := s.Dropbox.GetContent(entry.Path)
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified, Content: content}
	s.Db.SaveNote(&note)
}

func setupDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
