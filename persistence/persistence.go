package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type User struct {
	ID                 int `gorm:"primary_key"`
	Email              string
	EncryptedPassword  string
	DropboxAccessToken string
	DropboxCursor      string
	DropboxUserId      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Note struct {
	ID             int       `gorm:"primary_key" json:"id"`
	Title          string    `json:"title"`
	Content        string    `sql:"type:text" json:"-"`
	ContentPreview string    `json:"preview"`
	UserId         int       `json:"user_id"`
	Path           string    `json:"path"`
	Mtime          time.Time `json:"mtime"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DBClient interface {
	DeleteNote(note *Note) bool
	SaveNote(note *Note) bool
	UpdateUserCursor(user *User, cursor string) bool
	ListNotes(user *User) []Note
	GetNoteContents(id int) string
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

func (c *Client) ListNotes(user *User) []Note {
	var notes []Note
	c.Db.Select("id, title, left(content, 255) as content_preview, user_id, path, mtime, created_at, updated_at").Where("user_id = ?", user.ID).Limit(100).Order("mtime desc").Find(&notes)
	return notes
}

func (c *Client) GetNoteContents(id int) string {
	var note Note
	c.Db.Select("content").Where("id = ?", id).First(&note)
	return note.Content
}

func (c *Client) DeleteNote(note *Note) bool {
	c.Db.Delete(note)
	return true
}

func (c *Client) SaveNote(note *Note) bool {
	c.Db.Create(note)
	return true
}

func (c *Client) UpdateUserCursor(user *User, cursor string) bool {
	user.DropboxCursor = cursor
	c.Db.Save(user)
	return true
}
