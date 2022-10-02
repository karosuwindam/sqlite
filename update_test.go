package sqlite

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
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

	t.Log("-----------Update data ---------------")
	rand.Seed(time.Now().UnixNano())
	updateid := rand.Intn(20)
	dataInt := rand.Intn(99999)
	for {
		if updateid != 0 {
			break
		}
		updateid = rand.Intn(20)
	}
	udata := TableTest{Id: updateid, Str: "databadse", I: dataInt}
	if err := sql.Update(testtablename, &udata); err != nil {
		t.Errorf("data update err")
		t.FailNow()
	}
	t.Log("----------- Update data OK---------------")
	t.Log("----------- Update data Read---------------")
	rdata := []TableTest{}
	sql.Read(testtablename, &rdata, map[string]string{}, AND)
	if len(rdata) == 0 {
		t.Errorf("read data not id=%v", updateid)
		t.FailNow()
	}
	for i := 0; i < len(rdata); i++ {
		ck := ckdata[i]
		if rdata[i].Id == updateid {
			ck = udata
			t.Logf("chage id=%v", updateid)
		} else {
			t.Logf("ck id=%v", rdata[i].Id)
		}
		if rdata[i].Id != ck.Id {
			t.Errorf("read data not id %v,%v", rdata[i].Id, ck.Id)
			t.FailNow()
		}
		if rdata[i].Str != ck.Str {
			t.Errorf("read data not str %v,%v", rdata[i].Str, ck.Str)
			t.FailNow()
		}
		if rdata[i].I != ck.I {
			t.Errorf("read data not i %v,%v", rdata[i].I, ck.I)
			t.FailNow()
		}
	}
	t.Log("----------- Update data Read OK ---------------")

}
