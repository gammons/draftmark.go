package dropbox_client

import (
	"github.com/stacktic/dropbox"
	"log"
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

type DropboxClient interface {
	GetChanges(cursor, prefix string) []DropboxEntry
}

type Client struct {
	AccessToken string
	Cursor      string
	Prefix      string
	Dbox        *dropbox.Dropbox
}

func (c *Client) InitDropbox() {
	c.Dbox = dropbox.NewDropbox()
	c.Dbox.SetAppInfo(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"))
	c.Dbox.SetAccessToken(c.AccessToken)
}

func (c *Client) GetChanges(cursor string, prefix string) []DropboxEntry {
	allEntries := make([]DropboxEntry, 0)

	for {
		delta, _ := c.Dbox.Delta(cursor, prefix)
		for _, entry := range delta.Entries {
			allEntries = append(allEntries, DropboxEntry{entry.Entry.Path, "", entry.Entry.IsDeleted, time.Time(entry.Entry.Modified)})
		}
		log.Println("In entries, hasmore = ", delta.HasMore)
		if !delta.HasMore {
			break
		}
	}
	return allEntries
}
