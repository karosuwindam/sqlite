package sqlite

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {

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

	t.Log("-----------delete check id=1-20 ---------------")

	rand.Seed(time.Now().UnixNano())
	deleteid := rand.Intn(20)
	for {
		if deleteid != 0 {
			break
		}
		deleteid = rand.Intn(20)
	}
	t.Logf("-----------delete id=%v ---------------", deleteid)
	if err := sql.Delete(testtablename, deleteid); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	rdata := []TableTest{}
	_ = sql.Read(testtablename, TableTest{}, &rdata, map[string]string{"id": strconv.Itoa(deleteid)}, AND)

	if len(rdata) != 0 {
		t.Errorf("No Delete data id=%v", deleteid)
	}
	t.Log("-----------delete CHECK END ---------------")

}
