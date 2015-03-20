package dropbox_client

import (
	"github.com/stacktic/dropbox"
	"os"
	"time"
)

type DropboxEntry struct {
	Path      string
	Content   string
	IsDeleted bool
	Modified  time.Time
}

var dbox = dropbox.NewDropbox()

func init() {
	setupDropbox()
}

type DropboxClient interface {
	GetChanges() []DropboxEntry
}

type Client struct {
	AccessToken string
	Cursor      string
	Prefix      string
}

func (c *Client) GetChanges() []DropboxEntry {
	dbox.SetAccessToken(c.AccessToken)
	allEntries := make([]DropboxEntry, 0)

	for {
		delta, _ := dbox.Delta(c.Cursor, c.Prefix)
		for _, entry := range delta.Entries {
			allEntries = append(allEntries, DropboxEntry{entry.Entry.Path, "", entry.Entry.IsDeleted, time.Time(entry.Entry.Modified)})
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
