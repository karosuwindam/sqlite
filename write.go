package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// Add
func (cfg *sqliteConfig) Add(tname string, tabledatap interface{}) error {
	cangeDbID(cfg.sqlite3IdMax(tname), tabledatap)
	tabledata := reflect.ValueOf(tabledatap).Elem().Interface()

	cmd, err := createaddCmdByID(tname, tabledata)
	if err != nil {
		return err
	}
	_, err = cfg.db.Exec(cmd)
	return err
}

// createaddCmdByID()
// 挿入するコマンドを作成
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

//idを探して値を置き換える
func cangeDbID(id int, tabledatap interface{}) {
	if reflect.TypeOf(tabledatap).Kind() != reflect.Ptr || id < 0 {
		return
	}
	sv := reflect.ValueOf(tabledatap)
	svi := sv.Elem().Interface()
	st := reflect.TypeOf(svi)
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		if key := ft.Tag.Get("db"); key != "" {
			if key == "id" {
				sv.Elem().FieldByName(ft.Name).SetInt(int64(id))
				// fv.SetInt(int64(id))
			}
		}
	}
}

// idの値から最大値+1の値を設定
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
