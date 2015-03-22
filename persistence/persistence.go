package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type User struct {
	email                string
	encrypted_password   string
	dropbox_access_token string
	dropbox_cursor       string
	dropbox_user_id      string
	created_at           time.Time
	updated_at           time.Time
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
}

type Client struct {
	Db gorm.DB
}

func (c *Client) InitDB() {
	c.Db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")
	c.Db.LogMode(true)
}

func NewClient() DBClient {
	return &Client{}
}

func (c *Client) NoteCount() int {
	var count int
	c.Db.Table("notes").Count(&count)
	return count
}

func (c *Client) DeleteNote(note Note) bool {
	c.Db.Delete(&note)
	return true
}

func (c *Client) SaveNote(note Note) bool {
	c.Db.Where(Note{Path: note.Path}).FirstOrCreate(&note)
	return true
}
