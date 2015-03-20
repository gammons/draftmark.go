package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")

func init() {
	db.DB()
	db.LogMode(true)
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

type DBClient interface {
	DeleteNote(note Note) bool
	SaveNote(note Note) bool
}

func (c *Client) DeleteNote(note Note) bool {
	return true
}

func (c *Client) SaveNote(note Note) bool {
	//log.Println(db.First(&note))
	db.Where(Note{Path: note.Path}).FirstOrCreate(&note)
	//db.Create(&note)

	return true
}
