package draftmark_test

import (
	. "draftmark"
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

type FakeDropboxClient struct {
	getChanges []*dropbox.DropboxEntry
	getContent string
}

func (c *FakeDropboxClient) GetChanges(cursor *string, prefix string) (string, []*dropbox.DropboxEntry) {
	return "asdf", c.getChanges
}

func (c *FakeDropboxClient) GetContent(path string) string {
	return c.getContent
}

type FakeDb struct {
	saveNoteCount   int
	deleteNoteCount int
}

func (d *FakeDb) DeleteNote(note db.Note) bool {
	d.deleteNoteCount++
	return true
}

func (d *FakeDb) SaveNote(note db.Note) bool {
	d.saveNoteCount++
	return true
}

func (d *FakeDb) UpdateUserCursor(user db.User, cursor string) bool {
	return true
}

var _ = Describe("Sync Unit tests", func() {
	var user = db.User{ID: 1, Email: "gammons@gmail.com", DropboxCursor: "", DropboxAccessToken: "", DropboxUserId: "123", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	var fakedb FakeDb
	var fakeDropboxClient FakeDropboxClient
	var client = &Sync{Dropbox: &fakeDropboxClient, Db: &fakedb}

	BeforeEach(func() {
		fakedb.saveNoteCount = 0
		fakedb.deleteNoteCount = 0
	})

	Context("A file was added or changed", func() {
		It("Saves the file to the db", func() {
			fakeDropboxClient.getChanges = []*dropbox.DropboxEntry{{Path: "/notes/test.md", IsDeleted: false, Modified: time.Now()}}
			client.DoSync(user, "/notes")
			Expect(fakedb.saveNoteCount).To(Equal(1))
			Expect(fakedb.deleteNoteCount).To(Equal(0))
		})
	})

	Context("A file was deleted", func() {
		It("Deletes the file to the db", func() {
			fakeDropboxClient.getChanges = []*dropbox.DropboxEntry{{Path: "/notes/test.md", IsDeleted: true, Modified: time.Now()}}
			client.DoSync(user, "/notes")
			Expect(fakedb.deleteNoteCount).To(Equal(1))
			Expect(fakedb.saveNoteCount).To(Equal(0))
		})
	})

	Describe("Ignoring files we don't care about", func() {
		It("ignores files outside of the directory we care about", func() {
			fakeDropboxClient.getChanges = []*dropbox.DropboxEntry{{Path: "/Pictures/test.md", IsDeleted: true, Modified: time.Now()}}
			client.DoSync(user, "/notes")
			Expect(fakedb.deleteNoteCount).To(Equal(0))
			Expect(fakedb.saveNoteCount).To(Equal(0))
		})

		It("ignores files that do not end in .md", func() {
			fakeDropboxClient.getChanges = []*dropbox.DropboxEntry{{Path: "/notes/test.txt", IsDeleted: true, Modified: time.Now()}}
			client.DoSync(user, "/notes")
			Expect(fakedb.deleteNoteCount).To(Equal(0))
			Expect(fakedb.saveNoteCount).To(Equal(0))
		})
	})
})
