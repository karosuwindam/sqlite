package sqlite

import "testing"

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
