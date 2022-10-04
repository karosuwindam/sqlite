package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

// (*cfg)Update(tname str) = error
//
// SQLiteのデータベースから登録してあるデータを書き換える
//
// tname(string) : 対象のテーブル名
// str(interface{}) : データを書き換えるための構造体のポインタ
func (cfg *SqliteConfig) Update(tname string, str interface{}) error {
	cmd, err := createUpdateCmd(tname, str)
	if err != nil {
		return err
	}
	_, err = cfg.db.Exec(cmd)
	return err
}

// createUpdateCmd(tname, str) = string
//
// SQLite用の更新コマンドを作る
//
// tname(string) : 対象のテーブル名
// str(interface{}) : データを書き換えるための構造体
func createUpdateCmd(tname string, ptabledata interface{}) (string, error) {
	if reflect.TypeOf(ptabledata).Kind() != reflect.Ptr {
		return "", errors.New("input data not Pointer")
	}

	cmd := "UPDATE " + tname + " SET"
	str := reflect.ValueOf(ptabledata).Elem().Interface()
	rv := reflect.ValueOf(str)
	rt := reflect.TypeOf(str)
	id := 0
	flag := false
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		var v interface{}
		if rv.FieldByName(f.Name).Kind() == reflect.Struct {
			if rv.FieldByName(f.Name).Kind() == timeKind {
				fv := reflect.ValueOf(ptabledata).Elem().FieldByName(f.Name)
				fv = reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()
				v = fv.Interface()
			}
		} else {
			v = rv.FieldByName(f.Name).Interface()
		}
		if f.Tag.Get("db") == "" {
			continue
		} else if f.Tag.Get("db") == "id" {
			id = v.(int)
			continue
		}
		if flag {
			cmd += ","
		}
		switch v.(type) {
		case int:
			cmd += " " + f.Tag.Get("db") + "=" + strconv.Itoa(v.(int))
			flag = true
		case string:
			cmd += " " + f.Tag.Get("db") + "=" + "'" + v.(string) + "'"
			flag = true
		case time.Time:
			cmd += " " + f.Tag.Get("db") + "=" + "'" + v.(time.Time).Format(TimeLayout) + "'"
			flag = true
		}
	}
	if !flag {
		return "", errors.New("Don't create updata cmd")
	}
	now := time.Now()
	cmd += " ," + "updated_at" + "=" + "'" + "" + now.Format(TimeLayout) + "'"
	cmd += " " + "WHERE" + " " + "id=" + strconv.Itoa(id)

	return cmd, nil
}
