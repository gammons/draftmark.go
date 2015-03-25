package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type User struct {
	ID                 int
	Email              string
	EncryptedPassword  string
	DropboxAccessToken string
	DropboxCursor      string
	DropboxUserId      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Note struct {
	ID        int
	Title     string
	Content   string `sql:"type:text"`
	UserId    int
	Path      string
	Mtime     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DBClient interface {
	DeleteNote(note Note) bool
	SaveNote(note Note) bool
	UpdateUserCursor(user User, cursor string) bool
}

type Client struct {
	Db gorm.DB
}

func (c *Client) InitDB() {
	c.Db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")
}

func NewClient() DBClient {
	client := &Client{}
	client.InitDB()
	return client
}

func (c *Client) NoteCount() int {
	var count int
	c.Db.Table("notes").Count(&count)
	return count
}

func (c *Client) SaveNoteContents(note Note, content string) {
	note.Content = content
	c.Db.Update(&note)
}

func (c *Client) DeleteNote(note Note) bool {
	c.Db.Delete(&note)
	return true
}

func (c *Client) SaveNote(note Note) bool {
	c.Db.Create(&note)
	return true
}

func (c *Client) UpdateUserCursor(user User, cursor string) bool {
	user.DropboxCursor = cursor
	c.Db.Save(&user)
	return true
}
