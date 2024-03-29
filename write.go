package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

// (*cfg)AddOne(tname tabledatap) = error
//
// テーブルに追加する構造体のデータによって追加するテーブルを切り替える
//
// tname(string) : テーブル名
// tabledatap(interface{}) :テーブルに追加する構造体データのポインタ
func (cfg *SqliteConfig) Add(tname string, tabledatap interface{}) error {
	pv := reflect.ValueOf(tabledatap)
	if pv.Kind() != reflect.Ptr {
		return errors.New("tabledatap input not pointer")
	}
	ppv := reflect.ValueOf(pv.Elem().Interface())
	switch ppv.Kind() {
	case reflect.Slice: //配列構造体の入力
		// return nil
		pv := reflect.ValueOf(reflect.ValueOf(tabledatap).Elem().Interface())
		for i := 0; i < pv.Len(); i++ {
			fi := pv.Index(i)
			fi = reflect.NewAt(fi.Type(), unsafe.Pointer(fi.UnsafeAddr()))
			f := fi.Interface()
			if err := cfg.addOne(tname, f); err != nil {
				return err
			}

		}
		return nil
	case reflect.Struct: //構造体の入力
		return cfg.addOne(tname, tabledatap)
	}
	return nil
}

// (*cfg)add(tname tabledatap) = error
//
// SQLのテーブルにレコードを一つ追加する。
//
// tname(string) : テーブル名
// tabledatap(interface{}) :テーブルに追加する構造体データのポインタ
func (cfg *SqliteConfig) addOne(tname string, tabledatap interface{}) error {
	cangeDbID(cfg.sqlite3IdMax(tname), tabledatap)
	// tabledata := reflect.ValueOf(tabledatap).Elem().Interface()

	cmd, err := createaddCmdByID(tname, tabledatap)
	if err != nil {
		return err
	}
	_, err = cfg.db.Exec(cmd)
	return err
}

// createaddCmdByID(tname, tabledata) = string, error
//
// 挿入するコマンドを作成
//
// tname(string) : テーブル名
// tabledata(interface{}) : データを作成する構造体
func createaddCmdByID(tname string, ptabledata interface{}) (string, error) {
	if reflect.TypeOf(ptabledata).Kind() != reflect.Ptr {
		return "", errors.New("input data not Pointer")
	}
	cmd := "INSERT INTO " + tname + " "
	cmdColume := ""
	cmdVaule := ""
	now := time.Now()
	tabledata := reflect.ValueOf(ptabledata).Elem()

	st := reflect.TypeOf(tabledata.Interface())
	sv := reflect.ValueOf(tabledata.Interface())
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fv := sv.Field(i)
		if key := ft.Tag.Get("db"); key != "" {
			if i != 0 {
				cmdColume += ","
				cmdVaule += ","
			}
			if fv.Kind() == timeKind { //時刻処理
				fv = tabledata.FieldByName(ft.Name)
				fv = reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()
				fvt := fv.Interface()
				cmdColume += key
				cmdVaule += "'" + fvt.(time.Time).Format(TimeLayout) + "'"
				continue
			}
			fvi := fv.Interface()
			cmdColume += key
			switch fvi.(type) {
			case int:
				cmdVaule += strconv.Itoa(fvi.(int))
			case string:
				cmdVaule += "'" + fvi.(string) + "'"
			case float32:
				cmdVaule += strconv.FormatFloat(float64(fvi.(float32)), 'f', -1, 32)
			case float64:
				cmdVaule += strconv.FormatFloat(fvi.(float64), 'f', -1, 64)
			case time.Time:
				cmdVaule += "'" + fvi.(time.Time).Format(TimeLayout) + "'"
			}
		}
	}
	if cmdColume == "" || cmdVaule == "" {
		return "", errors.New("Don't created command")
	}
	cmdColume += "," + "created_at"
	cmdVaule += "," + "'" + now.Format(TimeLayout) + "'"
	cmdColume += "," + "updated_at"
	cmdVaule += "," + "'" + now.Format(TimeLayout) + "'"
	cmd += " (" + cmdColume + ") " + "VALUES" + " (" + cmdVaule + ")"

	return cmd, nil
}

// (*cfg)sqlite3IdMax(tname) = int
//
// 対象のSQLテーブルからKey名idの最大値に+1した値を返す
//
// tname : 対象のテーブル
func (cfg *SqliteConfig) sqlite3IdMax(tname string) int {
	id := 0
	cmd := "select max(id) from " + string(tname)
	rows, err := cfg.db.Query(cmd)
	if err != nil {
		return -1
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&id)
	id++
	return id
}
