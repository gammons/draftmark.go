package draftmark_test

import (
	. "draftmark"
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"time"
)

var database = &db.Client{}
var dbox = &dropbox.Client{}
var tmpFolder = "/__draftmarktest__"

var _ = BeforeSuite(func() {
	_ = godotenv.Load("test.env")
	setupDatabase()
	setupDropbox()
})

var _ = AfterSuite(func() {
	log.Println(dbox.Dbox.Delete(tmpFolder))
	dbox.Dbox.Delete(tmpFolder)
	database.Db.Exec("drop table users")
	database.Db.DropTable(&db.User{})
	database.Db.DropTable(&db.Note{})
})

var _ = Describe("Sync", func() {
	var user = User{1, "gammons@gmail.com", "asdf", "asdf", "123", time.Now(), time.Now()}
	var sync = &Sync{Db: database, Dropbox: dbox}
	Context("A new file is added in dropbox", func() {
		It("Adds the file to the db", func() {
			sync.DoSync(user, tmpFolder)
			log.Println("note count = ", database.NoteCount())
			Expect(database.NoteCount()).To(Equal(1))
			//Expect(database.Db.Table("notes").Count(&count)).To(Equal(1))
		})
	})
})

func setupDatabase() {
	database.InitDB()
	database.Db.CreateTable(&db.User{})
	database.Db.CreateTable(&db.Note{})
}

func setupDropbox() {
	dbox.AccessToken = os.Getenv("DROPBOX_ACCESS_TOKEN")
	dbox.InitDropbox()
	dbox.Dbox.CreateFolder(tmpFolder)
}
