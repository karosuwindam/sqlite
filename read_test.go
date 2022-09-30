package sqlite

import (
	"os"
	"testing"
)

func TestTableRead(t *testing.T) {

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

	ckdata := []TableTest{}

	wdata := TableTest{Id: 1, Str: "data", I: 500}
	for i := 0; i < 20; i++ {
		wdata.I += i
		sql.Add(testtablename, &wdata)
		ckdata = append(ckdata, wdata)
	}

	t.Log("-----------Read data ---------------")
	rdata := []TableTest{}
	err := sql.Read(testtablename, TableTest{}, &rdata, map[string]string{}, AND)
	if err != nil {
		t.Errorf("read err :%v", err.Error())
		t.FailNow()
	}
	if len(ckdata) != len(rdata) {
		t.Errorf("read count ng %v=%v", len(ckdata), len(rdata))
		t.FailNow()
	}
	for i, tmp := range rdata {
		tmp2 := ckdata[i]
		if tmp.I != tmp2.I {
			t.Errorf("read ng for i %v=%v", tmp.I, tmp2.I)
		}
		if tmp.Str != tmp2.Str {
			t.Errorf("read ng for str %v=%v", tmp.Str, tmp2.Str)
		}
	}
	t.Log(rdata)
	t.Log("read all OK")
	t.Log("-----------Read data Serch by ID---------------")
	rdata1 := []TableTest{}
	err1 := sql.Read(testtablename, TableTest{}, &rdata1, map[string]string{"id": "2"}, AND)
	if err1 != nil {
		t.Errorf("read err :%v", err1.Error())
	}
	if len(rdata1) != 1 {
		t.Errorf("read count ng %v=%v", 1, len(rdata1))
		t.FailNow()
	}
	for _, tmp := range rdata1 {
		for _, tmp2 := range ckdata {
			if tmp2.Id == tmp.Id {
				if tmp.I != tmp2.I {
					t.Errorf("read ng for i %v=%v", tmp.I, tmp2.I)
				}
				if tmp.Str != tmp2.Str {
					t.Errorf("read ng for str %v=%v", tmp.Str, tmp2.Str)
				}
			}
		}
	}
	t.Log("ID serch OK")
	t.Log("table data read OK")

}
