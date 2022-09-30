package sqlite

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func (cfg *sqliteConfig) Add(tname string, tabledatap interface{}) error {
	cangeDbId(cfg.sqlite3IdMax(tname), tabledatap)
	tabledata := reflect.ValueOf(tabledatap).Elem().Interface()

	cmd, err := createaddCmdById(tname, tabledata)
	if err != nil {
		return err
	}
	_, err = cfg.db.Exec(cmd)
	return err
}

// 挿入するコマンドを作成
func createaddCmdById(tname string, tabledata interface{}) (string, error) {
	cmd := "INSERT INTO " + tname + " "
	cmd_colume := ""
	cmd_vaule := ""
	now := time.Now()
	st := reflect.TypeOf(tabledata)
	sv := reflect.ValueOf(tabledata)
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fv := sv.Field(i)
		if key := ft.Tag.Get("db"); key != "" {
			fvi := fv.Interface()
			if i != 0 {
				cmd_colume += ","
				cmd_vaule += ","
			}
			cmd_colume += key
			switch fvi.(type) {
			case int:
				cmd_vaule += strconv.Itoa(fvi.(int))
			case string:
				cmd_vaule += "'" + fvi.(string) + "'"
			case time.Time:
				cmd_vaule += "'" + fvi.(time.Time).Format(TimeLayout) + "'"
			}
		}
	}
	if cmd_colume == "" || cmd_vaule == "" {
		return "", errors.New("Don't created command")
	} else {
		cmd_colume += "," + "created_at"
		cmd_vaule += "," + "'" + now.Format(TimeLayout) + "'"
		cmd_colume += "," + "updated_at"
		cmd_vaule += "," + "'" + now.Format(TimeLayout) + "'"
		cmd += " (" + cmd_colume + ") " + "VALUES" + " (" + cmd_vaule + ")"
	}

	return cmd, nil
}

//idを探して値を置き換える
func cangeDbId(id int, tabledatap interface{}) {
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
