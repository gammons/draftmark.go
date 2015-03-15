package dropbox_client

import (
	"github.com/stacktic/dropbox"
	"os"
	"time"
)

type ChangeGetter func(accessToken string, cursor string, prefix string) []DropboxEntry

type DropboxEntry struct {
	Path      string
	IsDeleted bool
	Modified  time.Time
}

var dbox = dropbox.NewDropbox()

func init() {
	setupDropbox()
}

type Client struct {
	GetChanges ChangeGetter
}

func GetChanges(accessToken string, cursor string, prefix string) []DropboxEntry {
	dbox.SetAccessToken(accessToken)
	allEntries := make([]DropboxEntry, 0)

	for {
		delta, _ := dbox.Delta(cursor, prefix)
		for _, entry := range delta.Entries {
			allEntries = append(allEntries, DropboxEntry{entry.Entry.Path, entry.Entry.IsDeleted, time.Time(entry.Entry.Modified)})
		}
		if !delta.HasMore {
			break
		}
	}
	return allEntries
}

func setupDropbox() {
	var dbox *dropbox.Dropbox
	dbox = dropbox.NewDropbox()
	dbox.SetAppInfo(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"))
}
