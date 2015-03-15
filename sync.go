package draftmark

import (
	"draftmark/dropbox_client"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
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

//var db *gorm.DB
var db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")

func init() {
	setupDotEnv()
	setupDB()
}

func Sync(user User, prefix string, client dropbox_client.Client) {
	entries := client.GetChanges(user.DropboxAccessToken, user.DropboxCursor, "/notes")
	for entry := range entries {
		log.Println(entry)
	}

	// for _, entry := range dp.Entries {
	// 	log.Println(entry.Entry)
	// 	if entry.Entry == nil {
	// 		deleteEntry(&entry)
	// 	} else {
	// 		createOrUpdateEntry(&entry)
	// 	}
	// }
}

func deleteEntry(entry *dropbox_client.DropboxEntry) {
}

func createOrUpdateEntry(entry *dropbox_client.DropboxEntry) {
}

func setupDotEnv() {
	log.Println("Setting up dotenv")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupDB() {
	db.DB()
}
