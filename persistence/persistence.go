package persistence

import (
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

var db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")

func init() {
	db.DB()
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

type Client struct {
}

func (c *Client) DeleteNote(note Note) bool {
	return true
}

func (c *Client) SaveNote(note Note) bool {
	return true
}
