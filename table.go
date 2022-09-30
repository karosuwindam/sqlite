package sqlite

import (
	"errors"
	"reflect"
	"strings"
)

type ifnot bool

const (
	ifnot_off ifnot = false
	ifnot_on  ifnot = true
)

const (
	TimeLayout  = "2006-01-02 15:04:05.999999999"
	TimeLayout2 = "2006-01-02 15:04:05.999999 +0000 UTC"
)

func (t *sqliteConfig) CreateTable(tname string, stu interface{}) error {
	var cmd string
	backcmd, err := t.ReadCreateTableCmd(tname)
	if err != nil {
		return err
	}
	if backcmd != "" { //tableが作成済み
		cmd, err = createTableCmd(tname, stu, ifnot_off)
		if err != nil {
			return err
		}
		if cmd == backcmd { //登録コマンドと実行コマンドが同じ

		} else { //登録コマンドと実行コマンドが異なる
			a := strings.Split(cmd, ",")
			b := strings.Split(backcmd, ",")
			if len(a) > len(b) { //カラム増やす変更を増やす

			} else {
				return errors.New("Don't change table for delete column")
			}

		}

	} else { //tableが作成していない
		cmd, err = createTableCmd(tname, stu, ifnot_on)
		if err != nil {
			return err
		}
		_, err = t.db.Exec(cmd)

	}

	return err
}

func (t *sqliteConfig) ReadTableList() ([]string, error) {
	var output []string
	cmd, err := readTableAllCmd()
	if err != nil {
		return output, err
	}
	rows, err := t.db.Query(cmd)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		str := ""
		err = rows.Scan(&str)
		if err != nil {
			return []string{}, err
		}
		output = append(output, str)
	}

	return output, err
}

func (t *sqliteConfig) ReadCreateTableCmd(tname string) (string, error) {
	var output string
	cmd, err := readCreateTableCmd(tname)
	if err != nil {
		return output, err
	}
	rows, err := t.db.Query(cmd)
	if err != nil {
		return output, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&output)
		if err != nil {
			return "", err
		}
	}

	return output, err

}

func (t *sqliteConfig) DropTable(tname string) error {
	cmd, err := dropTableCmd(tname)
	if err != nil {
		return err
	}
	_, err = t.db.Exec(cmd)
	return err

}

func createTableCmd(tname string, stu interface{}, flag ifnot) (string, error) {
	cmd := "CREATE TABLE" + " "
	if flag == ifnot_on {
		cmd += "IF NOT EXISTS" + " "
	}
	if tname == "" {
		return "", errors.New("Don't input name data")
	}
	cmd += "\"" + tname + "\""
	cmd += " ("
	if reflect.TypeOf(stu).Kind() != reflect.Struct {
		return "", errors.New("Don't input st data")
	}
	rt := reflect.TypeOf(stu)
	count := 0
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tmp := ""
		if i != 0 {
			cmd += ", "
		}
		switch f.Type.Kind() {
		case reflect.Int:
			tmp = f.Tag.Get("db")
			cmd += "\"" + tmp + "\" INTEGER"
		case reflect.String:
			tmp = f.Tag.Get("db")
			cmd += "\"" + tmp + "\" varchar"
		}
		if tmp == "id" {
			cmd += " PRIMARY KEY AUTOINCREMENT NOT NULL"
			count++
		} else if tmp == "" {
			return "", errors.New("Don't tag setup for " + f.Name)
		}
	}
	if count == 0 {
		return "", errors.New("Don't Struct data for \"id\"")
	}
	cmd += ", \"created_at\" datetime"
	cmd += ", \"updated_at\" datetime"
	cmd += ")"
	return cmd, nil

}

func altertableCmd(tname, cmd_a, cmd_b string) []string {
	var output []string
	//ALTER TABLE tbl_name ADD COLUMN new_col VARCHAR(10) AFTER col1;

	return output
}

func dropTableCmd(tname string) (string, error) {
	cmd := "DROP TABLE IF EXISTS" + " '" + tname + "'"
	return cmd, nil

}

func readTableAllCmd() (string, error) {
	cmd := "SELECT name FROM sqlite_master WHERE type='table'"
	return cmd, nil

}

func readCreateTableCmd(tname string) (string, error) {
	cmd := "SELECT sql FROM sqlite_master WHERE type='table' AND name='" + tname + "'"
	return cmd, nil
}
