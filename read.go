package sqlite

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type KeyWordOption string //検索オプション

//検索オプションの値
//
//AND keyword=data and
//OR keyword=data or
//AND_Like keyword like %keyword% and
//OR_LIKE keyword like %keyword% or
const (
	AND      KeyWordOption = "and"
	OR       KeyWordOption = "or"
	AND_Like KeyWordOption = "and_like"
	OR_Like  KeyWordOption = "or_like"
)

func (t *sqliteConfig) Read(tname string, stu interface{}, slice interface{}, v map[string]string, keytype KeyWordOption) error {
	cmd, err := createReadCmd(tname, stu, v, keytype)
	if err != nil {
		return err
	}
	rows, err := t.db.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		data, err := sqlite3RowsReadData(slice)
		if err != nil {
			return err
		}
		err = rows.Scan(data...)
		if err != nil {
			return err
		}
		tmpdata, err1 := silceToMap(data, slice)
		if err1 != nil {
			return err1
		}
		mapToStruct(tmpdata, slice)
	}

	return err
}

func createReadCmd(tname string, stu interface{}, keyword map[string]string, keytype KeyWordOption) (string, error) {
	rt := reflect.TypeOf(stu)
	if rt.Kind() == reflect.UnsafePointer {
		return "", errors.New("This input stu data is pointer")
	}
	cmd := "SELECT * FROM" + " " + tname
	if len(keyword) == 0 {
		// return cmd, nil
	} else {
		cmd += " " + "WHERE" + " " + convertCmd(stu, keyword, keytype)

	}

	return cmd, nil

}

func sqlite3RowsReadData(slice interface{}) ([]interface{}, error) {
	sv := reflect.ValueOf(slice)
	if sv.Type().Kind() != reflect.Ptr {
		fmt.Println(sv.Type().Kind())
		return []interface{}(nil), errors.New("Don't pointer data input data = " + sv.Type().Kind().String())
	}
	var output []interface{}
	tStruct := reflect.TypeOf(sv.Elem().Interface()).Elem()
	vStruct := reflect.New(tStruct)
	loaddata := vStruct.Elem().Interface()
	if loaddata == nil {
		return output, errors.New("Don't input tablemap data for" + string(""))
	}
	rt := reflect.TypeOf(loaddata)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field((i))
		switch ft.Type.Kind() {
		case reflect.Int:
			i := int64(0)
			output = append(output, &i)
		case reflect.String:
			str := string("")
			output = append(output, &str)
		case reflect.Struct:
			// str := string("")
			var str time.Time
			output = append(output, &str)
		default:
			return output, errors.New("Error data type " + strconv.Itoa(int(ft.Type.Kind())))
		}
	}
	if len(output) != 0 {
		for i := 0; i < 2; i++ {
			var str time.Time
			output = append(output, &str)
		}

	}
	return output, nil
}

func silceToMap(silce []interface{}, stu interface{}) (map[string]interface{}, error) {
	output := map[string]interface{}{}
	if len(silce) == 0 {
		return output, nil
	}

	sv := reflect.ValueOf(stu)
	if sv.Type().Kind() != reflect.Ptr {
		return output, errors.New("Don't struct pointer input i=" + sv.Type().Kind().String())
	}
	ii := sv.Elem().Interface()
	tStruct := reflect.TypeOf(ii).Elem()
	vStruct := reflect.New(tStruct)
	ckStruct := reflect.TypeOf(vStruct.Elem().Interface())
	for i, data := range silce {
		tmp := reflect.ValueOf(data).Elem().Interface()
		if i < ckStruct.NumField() {
			f := ckStruct.Field(i)
			switch tmp.(type) {
			case int64:
				output[f.Name] = int(tmp.(int64))
			case int:
				output[f.Name] = tmp.(int)
			case string:
				output[f.Name] = tmp.(string)
			case time.Time:
				output[f.Name] = tmp.(time.Time)
			}
		} else {
			if i == len(silce)-2 {
				output["create_at"] = tmp.(time.Time)
			} else if i == len(silce)-1 {
				output["update_at"] = tmp.(time.Time)

			}
		}
	}
	// endcount := 0
	// for i := 0; i < ckStruct.NumField(); i++ {
	// 	f := ckStruct.Field(i)

	// 	if i >= len(silce) {
	// 		break
	// 	}
	// 	if f.Type.Kind() == reflect.Int64 {
	// 		output[f.Name] = silce[i].(int)

	// 	} else {
	// 		output[f.Name] = silce[i]

	// 	}
	// 	endcount = i
	// }
	// endcount++
	// if len(silce[endcount:]) == 2 {
	// 	output["create_at"] = silce[len(silce)-2]
	// 	output["update_at"] = silce[len(silce)-1]
	// }
	return output, nil
}

func convertCmd(stu interface{}, keyword map[string]string, keytype KeyWordOption) string {
	output := ""
	if stu == nil {
		return output
	}
	st := reflect.TypeOf(stu)
	count := 0
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if keyword[f.Tag.Get("db")] != "" {
			if count != 0 {
				switch keytype {
				case AND_Like:
					output += " " + string(AND) + " "
				case OR_Like:
					output += " " + string(OR) + " "
				default:
					output += " " + string(keytype) + " "
				}

			}
			if keytype == AND || keytype == OR {
				output += f.Tag.Get("db") + "="
				switch f.Type.Kind() {
				case reflect.Int:
					output += keyword[f.Tag.Get("db")]
				case reflect.String:
					output += "'" + keyword[f.Tag.Get("db")] + "'"
				}
			} else {
				output += f.Tag.Get("db") + " like "
				switch f.Type.Kind() {
				case reflect.Int:
					output += "'%" + keyword[f.Tag.Get("db")] + "%'"
				case reflect.String:
					output += "'%" + keyword[f.Tag.Get("db")] + "%'"
				}
			}
			count++
		}
	}

	return output
}

func mapToStruct(s map[string]interface{}, i interface{}) error {

	sv := reflect.ValueOf(i)
	if sv.Type().Kind() != reflect.Ptr {
		return errors.New("Don't struct pointer input i=" + sv.Type().Kind().String())
	}
	if len(s) == 0 {
		return nil
	}
	ii := sv.Elem().Interface()
	tStruct := reflect.TypeOf(ii).Elem()
	vStruct := reflect.New(tStruct)
	ckStruct := reflect.TypeOf(vStruct.Elem().Interface())
	for i := 0; i < ckStruct.NumField(); i++ {
		f := ckStruct.Field(i)
		v := vStruct.Elem().FieldByName(f.Name)
		ss := s[f.Name]
		switch f.Type.Kind() {
		case reflect.Int & reflect.TypeOf(ss).Kind():
			v.SetInt(int64(ss.(int)))
		case reflect.String & reflect.TypeOf(ss).Kind():
			v.SetString(ss.(string))
		}
	}
	out := vStruct.Elem()
	v := sv.Elem()
	v.Set(reflect.Append(v, out))
	return nil
}
