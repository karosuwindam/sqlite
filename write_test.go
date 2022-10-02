package sqlite

import (
	"os"
	"testing"
	"time"
)

func TestTableWrite(t *testing.T) {

	type TableTest struct {
		Id  int    `db:"id"`
		Str string `db:"str"`
		I   int    `db:"i"`
	}

	testtablename := "test"
	testdbname := "test.db"

	sql := Setup(testdbname)
	_ = sql.Open()
	defer sql.Close()
	sql.CreateTable(testtablename, TableTest{})
	defer os.Remove(testdbname)

	t.Log("-----------Add data ---------------")
	wdata := TableTest{Id: 10, Str: "data", I: 500}
	err := sql.Add(testtablename, &wdata)
	if err != nil {
		t.Errorf("Don't Added %v Table", testtablename)
	}
	t.Log("Add recode")
	t.Log("----------- check data --------------")
	i := sql.sqlite3IdMax(testtablename)
	if i != 2 {
		t.Errorf("%v table count+1 = %v", testtablename, i)
	} else {
		t.Log("table data add OK")

	}
}

func TestTableWriteMulti(t *testing.T) {

	type TableTest struct {
		Id   int       `db:"id"`
		Str  string    `db:"str"`
		I    int       `db:"i"`
		time time.Time `db:"time"`
	}

	testtablename := "test"
	testdbname := "test.db"

	sql := Setup(testdbname)
	_ = sql.Open()
	defer sql.Close()
	sql.CreateTable(testtablename, TableTest{})
	defer os.Remove(testdbname)

	t.Log("-----------Add Multi data ---------------")
	mwdata := []TableTest{}
	wdata := TableTest{Id: 10, Str: "data", I: 500}
	for i := 0; i < 10; i++ {
		wdata.time = time.Now()
		mwdata = append(mwdata, wdata)
	}
	err := sql.Add(testtablename, &mwdata)
	if err != nil {
		t.Errorf("Don't Added %v Table", testtablename)
	}
	t.Log("Add recode")
	t.Log("----------- check data --------------")
	i := sql.sqlite3IdMax(testtablename)
	if i != 11 {
		t.Errorf("%v table count+1 = %v", testtablename, i)
	} else {
		t.Log("table data add OK")

	}
}
