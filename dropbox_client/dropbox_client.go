package dropbox_client

import (
	"github.com/stacktic/dropbox"
	"io/ioutil"
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
	GetChanges(cursor *string, prefix string) (string, []*DropboxEntry)
	GetContent(path string) (string, error)
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

func (c *Client) GetChanges(cursor *string, prefix string) (string, []*DropboxEntry) {
	allEntries := []*DropboxEntry{}
	var nextCursor string

	for {
		delta, _ := c.Dbox.Delta(*cursor, prefix)
		nextCursor = delta.Cursor.Cursor
		for _, entry := range delta.Entries {
			if entry.Entry == nil {
				newEntry := &DropboxEntry{entry.Path, "", true, time.Now()}
				allEntries = append(allEntries, newEntry)
			} else {
				newEntry := &DropboxEntry{entry.Path, "", entry.Entry.IsDeleted, time.Time(entry.Entry.Modified)}
				allEntries = append(allEntries, newEntry)
			}
		}
		if !delta.HasMore {
			break
		}
	}
	return nextCursor, allEntries
}

func (c *Client) GetContent(path string) (string, error) {
	noteio, length, err := c.Dbox.Download(path, "", 0)
	if length != 0 {
		note, _ := ioutil.ReadAll(noteio)
		return string(note), nil
	} else {
		return "", err
	}
}
