package draftmark

import (
	"draftmark/dropbox_client"
	"testing"
	"time"
)

func TestSync_NewFileAdded(t *testing.T) {
	changes := func(accessToken string, cursor string, prefix string) []dropbox_client.DropboxEntry {
		r := []dropbox_client.DropboxEntry{{Path: "/notes/test.md", IsDeleted: false, Modified: time.Now()}}
		return r
	}

	client := dropbox_client.Client{GetChanges: changes}
	user := User{1, "gammons@gmail.com", "asdf", "asdf", "123", time.Now(), time.Now()}
	Sync(user, "/notes", client)
}
