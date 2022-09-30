package sqlite

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	testdbname := "test.db"

	t.Log("-----------sqlite setup ---------------")
	sql := Setup(testdbname)
	t.Log("-----------sqlite Open ---------------")
	err := sql.Open()
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	t.Log("-----------sqlite Close ---------------")
	err = sql.Close()
	if err != nil {
		t.Errorf("%v", err.Error())

	}
	os.Remove(testdbname)
}
