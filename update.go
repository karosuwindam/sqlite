package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// (*cfg)Update(tname str) = error
//
// SQLiteのデータベースから登録してあるデータを書き換える
//
// tname(string) : 対象のテーブル名
// str(interface{}) : データを書き換えるための構造体
func (cfg *sqliteConfig) Update(tname string, str interface{}) error {
	cmd := createUpdateCmd(tname, str)
	if cmd == "" {
		return errors.New("Don't create updata cmd")
	}
	_, err := cfg.db.Exec(cmd)
	return err
}

// createUpdateCmd(tname, str) = string
//
// SQLite用の更新コマンドを作る
//
// tname(string) : 対象のテーブル名
// str(interface{}) : データを書き換えるための構造体
func createUpdateCmd(tname string, str interface{}) string {
	cmd := "UPDATE " + tname + " SET"
	rv := reflect.ValueOf(str)
	rt := reflect.TypeOf(str)
	id := 0
	flag := false
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := rv.FieldByName(f.Name).Interface()
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
		return ""
	}
	now := time.Now()
	cmd += " ," + "updated_at" + "=" + "'" + "" + now.Format(TimeLayout) + "'"
	cmd += " " + "WHERE" + " " + "id=" + strconv.Itoa(id)

	return cmd
}
