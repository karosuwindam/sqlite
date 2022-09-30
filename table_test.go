package sqlite

import (
	"os"
	"testing"
)

func createchcmdTableTest(name string) string {
	cmd := "CREATE TABLE IF NOT EXISTS \"" + name + "\" "
	cmd += "(" + "\"id\" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL"
	cmd += ", \"str\" varchar"
	cmd += ", \"i\" INTEGER"
	cmd += ", \"created_at\" datetime, \"updated_at\" datetime)"
	return cmd
}

func TestCreateTable(t *testing.T) {

	type TableTest struct {
		id  int    `db:"id"`
		str string `db:"str"`
		i   int    `db:"i"`
	}

	testtablename := "test"
	testdbname := "test.db"

	sql := Setup(testdbname)
	_ = sql.Open()
	defer sql.Close()

	t.Log("----------- table create cmd ---------------")
	str, err := createTableCmd(testtablename, TableTest{}, ifnot_on)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	if str != createchcmdTableTest(testtablename) {
		t.Errorf("output:%v\ncheck:%v", str, createchcmdTableTest(testtablename))
	}
	t.Logf("run sql cmd:%v", str)
	t.Log("----------- table create ---------------")
	err = sql.CreateTable(testtablename, TableTest{})
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	t.Log("create test.db")

	t.Log("----------- table read ---------------")
	cmd, err1 := sql.ReadCreateTableCmd(testtablename)
	if err1 != nil {
		t.Errorf("%v", err1.Error())
	}
	t.Logf("cmd:%v", cmd)

	t.Log("----------- table list read ---------------")
	stu, err1 := sql.ReadTableList()
	if err1 != nil {
		t.Errorf("%v", err1.Error())
	}
	tableckoff := []string{"sqlite_sequence"}
	for _, name := range stu {
		if testtablename == name {
			t.Logf("created table = %v", name)
		} else {
			ck := true
			for _, offname := range tableckoff {
				if name == offname {
					ck = false
					break
				}
			}
			if ck {
				t.Errorf("%v", name)
			}
		}
	}
	t.Log("----------- table drop ---------------")
	err1 = sql.DropTable(testtablename)
	if err != nil {
		t.Errorf("%v", err1.Error())
	}
	ckdatabase, _ := sql.ReadCreateTableCmd(testtablename)
	if ckdatabase != "" {
		t.Errorf("Don't delete table cmd:%v", ckdatabase)
	}

	t.Logf("%v table deleted", testtablename)
	os.Remove(testdbname)

}
