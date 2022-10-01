package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// (*cfg)AddOne(tname tabledatap) = error
//
// テーブルに追加する構造体のデータによって追加するテーブルを切り替える
//
// tname(string) : テーブル名
// tabledatap(interface{}) :テーブルに追加する構造体データのポインタ
func (cfg *sqliteConfig) Add(tname string, tabledatap interface{}) error {
	pv := reflect.ValueOf(tabledatap)
	if pv.Kind() != reflect.Ptr {
		return errors.New("tabledatap input not pointer")
	}
	ppv := reflect.ValueOf(pv.Elem().Interface())
	switch ppv.Kind() {
	case reflect.Slice: //配列構造体の入力
		return nil
		pv := reflect.ValueOf(reflect.ValueOf(tabledatap).Elem().Interface())
		for i := 0; i < pv.Len(); i++ {
			f := pv.Index(i).Interface()
			//以下の分で、構造体のポインタ渡しで失敗しているので
			//ポインタの同じ構造体を作ってそっちに渡して動作するかテストを実施する。
			cfg.addOne(tname, &f)
		}
		return nil
	case reflect.Struct: //構造体の入力
		return cfg.addOne(tname, tabledatap)
	}
	return nil
}

// (*cfg)addOne(tname tabledatap) = error
//
// SQLのテーブルにレコードを一つ追加する。
//
// tname(string) : テーブル名
// tabledatap(interface{}) :テーブルに追加する構造体データのポインタ
func (cfg *sqliteConfig) addOne(tname string, tabledatap interface{}) error {
	cangeDbID(cfg.sqlite3IdMax(tname), tabledatap)
	tabledata := reflect.ValueOf(tabledatap).Elem().Interface()

	cmd, err := createaddCmdByID(tname, tabledata)
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
func createaddCmdByID(tname string, tabledata interface{}) (string, error) {
	cmd := "INSERT INTO " + tname + " "
	cmdColume := ""
	cmdVaule := ""
	now := time.Now()
	st := reflect.TypeOf(tabledata)
	sv := reflect.ValueOf(tabledata)
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fv := sv.Field(i)
		if key := ft.Tag.Get("db"); key != "" {
			fvi := fv.Interface()
			if i != 0 {
				cmdColume += ","
				cmdVaule += ","
			}
			cmdColume += key
			switch fvi.(type) {
			case int:
				cmdVaule += strconv.Itoa(fvi.(int))
			case string:
				cmdVaule += "'" + fvi.(string) + "'"
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
func (cfg *sqliteConfig) sqlite3IdMax(tname string) int {
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
