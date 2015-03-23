package persistence

import (
	"fmt"
	"os"
	"testing"
)

func init() {
	os.Setenv("DB_TABLENAME", "remark_test")
}

func TestSaveNote(t *testing.T) {
	user := &User{Email: "test@test.com"}
	fmt.Println(user)
	// note := &Note{Path: "/notes/test.md"}
	// client := new(Client)
	// client.SaveNote(note)

}
