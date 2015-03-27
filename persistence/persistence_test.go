package persistence_test

import (
	database "draftmark/persistence"
	//"fmt"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var db = &database.Client{}
var _ = BeforeSuite(func() {
	err := godotenv.Load("../test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	//db.Db.LogMode(true)
	db.Db.CreateTable(&database.User{})
	db.Db.CreateTable(&database.Note{})
})

var _ = AfterSuite(func() {
	db.Db.DropTable(&database.User{})
	db.Db.DropTable(&database.Note{})
})

var _ = Describe("ListNotes", func() {
	var user = database.User{Email: "gammons@gmail.com"}
	var note1 = &database.Note{Path: "/notes/test.md", UserId: 1}
	var note2 = &database.Note{Path: "/notes/test2.md", UserId: 1}
	var _ = &database.Note{Path: "/notes/test3.md", UserId: 2}

	It("Lists notes for a user", func() {
		db.Db.Create(&user)
		db.Db.Create(note1)
		db.Db.Create(note2)
		notes := db.ListNotes(&user)

		var paths []string
		for _, note := range notes {
			paths = append(paths, note.Path)
		}
		var expected = []string{"/notes/test.md", "/notes/test2.md"}
		Expect(paths).To(Equal(expected))
	})

})
