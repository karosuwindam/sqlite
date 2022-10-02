package sqlite

import (
	"testing"
	"time"
)

func TestCommonStructToSlice(t *testing.T) {
	type TableTest struct {
		Id  int    `db:"id"`
		Str string `db:"str"`
		I   int    `db:"i"`
	}
	data := structToSlice(TableTest{}).([]TableTest)

	if len(data) != 0 {
		t.Errorf("data len = %v", len(data))
		t.FailNow()
	}
	aData := TableTest{Id: 1, Str: "data", I: 10}
	data = append(data, aData)
	if len(data) != 1 {
		t.Errorf("data len =%v", len(data))
		t.FailNow()
	}

	if data[0].Id != aData.Id {
		t.Errorf("data NG %v=%v", data[0].Id, aData.Id)
		t.FailNow()
	}
	if data[0].Str != aData.Str {
		t.Errorf("data NG %v=%v", data[0].Str, aData.Str)
		t.FailNow()
	}
	if data[0].I != aData.I {
		t.Errorf("data NG %v=%v", data[0].I, aData.I)
		t.FailNow()
	}
	t.Log("------------------ Struct To Slice OK -----------------")
}

func TestCommonMapToStruct(t *testing.T) {
	type TableTest struct {
		Id  int    `db:"id"`
		Str string `db:"str"`
		I   int    `db:"i"`
	}
	output := []TableTest{}
	input := map[string]interface{}{"Id": 1, "Str": "data", "I": 300}
	if err := mapToStruct(input, &output); err != nil {
		t.Errorf("Don't Map to struct")
		t.FailNow()
	}
	if len(output) != 1 {
		t.Errorf("Don't Map to struct len = %v", len(output))
		t.FailNow()
	}
	if output[0].Id != input["Id"] {
		t.Errorf("Don't Map to struct Id %v = %v", output[0].Id, input["Id"])
		t.FailNow()
	}
	if output[0].Str != input["Str"] {
		t.Errorf("Don't Map to struct Id %v = %v", output[0].Str, input["Str"])
		t.FailNow()
	}
	if output[0].I != input["I"] {
		t.Errorf("Don't Map to struct Id %v = %v", output[0].I, input["I"])
		t.FailNow()
	}
	t.Log("------------------ Map To Struct OK -----------------")

}

func TestCommonCangeDbID(t *testing.T) {
	type TableTest struct {
		Id  int    `db:"id"`
		Str string `db:"str"`
		I   int    `db:"i"`
	}
	data := TableTest{Id: 30, Str: "data", I: 455}
	cangeDbID(2, &data)
	if data.Id != 2 {
		t.Errorf("Don't change struct Id %v", data.Id)
		t.FailNow()
	}
	if data.Str != "data" {
		t.Errorf("Don't change struct Str %v", data.Str)
		t.FailNow()
	}
	if data.I != 455 {
		t.Errorf("Don't change struct I %v", data.I)
		t.FailNow()
	}

	t.Log("------------------ Cange Db for ID OK -----------------")
}

func TestCommonStructTopStruct(t *testing.T) {
	type TableTest struct {
		Id   int       `db:"id"`
		Str  string    `db:"str"`
		I    int       `db:"i"`
		time time.Time `db:"time"`
	}
	if structTopStruct(int(1)) != nil {
		t.Errorf("err data input not err")
		t.FailNow()
	}
	data := TableTest{Id: 30, Str: "data", I: 455}
	pdata := structTopStruct(data)
	if pdata == nil {
		t.Errorf("input data not pointer")
		t.FailNow()
	}
	switch pdata.(type) {
	case *TableTest:
		ppdata := pdata.(*TableTest)
		if ppdata.Id != data.Id {
			t.Errorf("pointer data not basedata,Id %v = %v", ppdata.Id, data.Id)
			t.FailNow()
		}
		if ppdata.Str != data.Str {
			t.Errorf("pointer data not basedata,Str %v = %v", ppdata.I, data.I)
			t.FailNow()
		}
		if ppdata.I != data.I {
			t.FailNow()
			t.Errorf("pointer data not basedata,I %v = %v", ppdata.I, data.I)
		}
	default:
		t.Errorf("not pointer")
		t.FailNow()
	}
	t.Log("------------------ Cange Struct to Point Struct OK -----------------")

}
