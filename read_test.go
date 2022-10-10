package sqlite

import (
	"os"
	"testing"
	"time"
)

func TestTableCreateReadCmd(t *testing.T) {

	type TableTest struct {
		Id   int       `db:"id"`
		Str  string    `db:"str"`
		I    int       `db:"i"`
		time time.Time `db:"time"`
	}

	testtablename := "test"
	keyword := map[string]string{"id": "1", "str": "bb"}
	cmd, err := createReadCmd(testtablename, &[]TableTest{})
	if err != nil {
		t.Error(err)
	}
	t.Log(cmd)
	cmd, err = createReadCmd(testtablename, &[]TableTest{}, keyword)
	if err != nil {
		t.Error(err)
	}
	cmd1, err := createReadCmd(testtablename, &[]TableTest{}, keyword, AND)
	if err != nil {
		t.Error(err)
	}
	cmd2, err := createReadCmd(testtablename, &[]TableTest{}, keyword, AND, AND)
	if err != nil {
		t.Error(err)
	}
	if cmd != cmd1 && cmd1 != cmd2 {
		t.Error(cmd1)
		t.Error(cmd2)
		t.FailNow()
	}
	t.Log(cmd)
	cmd, err = createReadCmd(testtablename, &[]TableTest{}, OR, keyword)
	if err != nil {
		t.Error(err)
	}
	t.Log(cmd)

}

func TestTableRead(t *testing.T) {

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

	ckdata := []TableTest{}

	wdata := TableTest{Id: 1, Str: "data", I: 500}
	for i := 0; i < 20; i++ {
		wdata.I += i
		wdata.time = time.Now()
		sql.Add(testtablename, &wdata)
		ckdata = append(ckdata, wdata)
	}

	t.Log("-----------Read data ---------------")
	rdata := []TableTest{}
	err := sql.Read(testtablename, &rdata, map[string]string{}, AND)
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
		if tmp.time.Format(TimeLayout) != tmp2.time.Format(TimeLayout) {
			t.Errorf("read ng for time %v=%v", tmp.time.Format(TimeLayout), tmp2.time.Format(TimeLayout))
		}
	}
	t.Log(rdata)
	t.Log("read all OK")
	t.Log("-----------Read data Serch by ID---------------")
	rdata1 := []TableTest{}
	err1 := sql.Read(testtablename, &rdata1, map[string]string{"id": "2"}, AND)
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

func TestReadWhileTime(t *testing.T) {
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
	nowtime := time.Now()
	cmdbase := "INSERT INTO " + testtablename
	cmdback := " (" + "id,str,i,created_at,updated_at) VALUES (" + "1,'a',0,'" + nowtime.Format(TimeLayout) + "','" + nowtime.Format(TimeLayout) + "'" + ")"

	if _, err := sql.db.Exec(cmdbase + cmdback); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	cmdback = " (" + "id,str,i,created_at,updated_at) VALUES (" + "2,'b',0,'" + nowtime.Add(-24*time.Hour*6).Format(TimeLayout) + "','" + nowtime.Add(-24*time.Hour*6).Format(TimeLayout) + "'" + ")"
	if _, err := sql.db.Exec(cmdbase + cmdback); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	cmdback = " (" + "id,str,i,created_at,updated_at) VALUES (" + "3,'c',0,'" + nowtime.Add(-24*time.Hour*29).Format(TimeLayout) + "','" + nowtime.Add(-24*time.Hour*29).Format(TimeLayout) + "'" + ")"
	if _, err := sql.db.Exec(cmdbase + cmdback); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	cmdback = " (" + "id,str,i,created_at,updated_at) VALUES (" + "4,'d',0,'" + nowtime.Add(-24*time.Hour*31).Format(TimeLayout) + "','" + nowtime.Add(-24*time.Hour*31).Format(TimeLayout) + "'" + ")"
	if _, err := sql.db.Exec(cmdbase + cmdback); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	rdata := []TableTest{}
	if err := sql.ReadToday(testtablename, &rdata); err == nil {
		if len(rdata) == 1 {
			t.Log(rdata)
		} else {
			t.Fail()
		}
	} else {
		t.Fatal(err.Error())
	}
	rdata = []TableTest{}
	if err := sql.ReadToWeek(testtablename, &rdata); err == nil {
		if len(rdata) == 2 {
			t.Log(rdata)
		} else {
			t.Fail()
		}
	} else {
		t.Fatal(err.Error())
	}
	rdata = []TableTest{}
	if err := sql.ReadToMonth(testtablename, &rdata); err == nil {
		if len(rdata) == 3 {
			t.Log(rdata)
		} else {
			t.Fail()
		}
	} else {
		t.Fatal(err.Error())
	}

}
