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
}

func (c *FakeDropboxClient) GetChanges(cursor string, prefix string) []*dropbox.DropboxEntry {
	return c.getChanges
}

type FakeDb struct {
	saveNoteCount   int
	deleteNoteCount int
}

func (d *FakeDb) DeleteNote(note db.Note) bool {
	d.deleteNoteCount++
	return true
}

func (d *FakeDb) SaveNote(note *db.Note) bool {
	d.saveNoteCount++
	return true
}

var _ = Describe("Sync Unit tests", func() {
	var user = User{1, "gammons@gmail.com", "asdf", "asdf", "123", time.Now(), time.Now()}
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
