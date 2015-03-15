package persistence

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("DB_TABLENAME", "remark_test")
}

func TestSaveNote(t *testing.T) {
	note := Note{Path: "/notes/test.md"}
	client := new(Client)
	client.SaveNote(note)

}
