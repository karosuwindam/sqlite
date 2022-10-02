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
	str, err := createTableCmd(testtablename, TableTest{}, ifnotOn)
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

//未実装のためコメントアウト
// func TestTableCreateCmd(t *testing.T) {

// 	type TableTest struct {
// 		Id   int       `db:"id"`
// 		Str  string    `db:"str"`
// 		I    int       `db:"i"`
// 		time time.Time `db:"time"`
// 	}

// 	type TableTest1 struct {
// 		Id   int       `db:"id"`
// 		Str  string    `db:"str"`
// 		Strb string    `db:"strb"`
// 		time time.Time `db:"time"`
// 		I    int       `db:"i"`
// 		B    int       `db:"b"`
// 		C    int       `db:"c"`
// 	}

// 	tname := "test"
// 	testdbname := "test.db"

// 	sql := Setup(testdbname)
// 	_ = sql.Open()
// 	defer sql.Close()

// 	t.Log("----------- table cmd check ---------------")
// 	cmdA, err := createTableCmd(tname, TableTest{}, ifnotOn)
// 	if err != nil {
// 		t.Error("Do not created cmd")
// 		t.FailNow()
// 	}
// 	cmdB, err := createTableCmd(tname, TableTest1{}, ifnotOn)
// 	if err != nil {
// 		t.Error("Do not created cmd")
// 		t.FailNow()
// 	}
// 	output := altertableCmd(cmdA, cmdB)
// 	if len(output) == 0 && cmdA == cmdB {
// 		t.Error("")
// 		t.FailNow()
// 	}
// 	t.Log("----------- table cmd check OK ---------------")
// 	t.Log("----------- table create ---------------")
// 	err = sql.CreateTable(tname, TableTest{})
// 	if err != nil {
// 		t.Errorf("%v", err.Error())
// 		t.FailNow()
// 	}
// 	defer os.Remove(testdbname)

// 	t.Log("----------- table read ---------------")
// 	cmd, err1 := sql.ReadCreateTableCmd(tname)
// 	if err1 != nil {
// 		t.Errorf("%v", err1.Error())
// 		t.FailNow()
// 	}
// 	t.Logf("cmd:%v", cmd)
// 	t.Log("----------- input data ---------------")
// 	wdata := TableTest{Id: 1, Str: "data", I: 500, time: time.Now()}
// 	ckdata := []TableTest{}
// 	for i := 0; i < 20; i++ {
// 		wdata.I += i
// 		if err := sql.Add(tname, &wdata); err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		ckdata = append(ckdata, wdata)
// 	}
// 	sql.Add(tname, &ckdata)
// 	ckdata = []TableTest{}
// 	sql.Read(tname, &ckdata)
// 	t.Log(ckdata)

// 	t.Log("----------- different table create NG---------------")
// 	err = sql.CreateTable(tname, TableTest1{})
// 	if err == nil {
// 		t.FailNow()
// 	}
// 	t.Logf("%v", err.Error())

// 	t.Log("----------- different table update ---------------")
// 	err = sql.UpdateTable(tname, TableTest1{})
// 	if err != nil {
// 		t.Errorf("%v", err.Error())
// 		t.FailNow()
// 	}

// 	t.Log("----------- table read ---------------")
// 	cmd, err1 = sql.ReadCreateTableCmd(tname + "_tmp")
// 	if err1 != nil {
// 		t.Errorf("%v", err1.Error())
// 		t.FailNow()
// 	}
// 	t.Logf("cmd:%v", cmd)
// }
