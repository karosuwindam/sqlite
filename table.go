package sqlite

import (
	"errors"
	"reflect"
	"strings"
)

// ifnot
//
// SQLコマンド作成時でＩＦ　EXISTSをつけるフラグ
type ifnot bool

const (
	ifnotOff ifnot = false
	ifnotOn  ifnot = true
)

const (
	// TimeLayout String変換用テンプレート
	TimeLayout = "2006-01-02 15:04:05.999999999"
	// TimeLayout2 String変換用テンプレート
	TimeLayout2 = "2006-01-02 15:04:05.999999 +0000 UTC"
)

// CreateTable(tname, stu) = error
//
// SQL内にテーブルを作成する
//
// tname(string) : 作成するテーブル名
// stu(interface{}) : 作成するテーブル内の構造体
func (t *sqliteConfig) CreateTable(tname string, stu interface{}) error {
	var cmd string
	backcmd, err := t.ReadCreateTableCmd(tname)
	if err != nil {
		return err
	}
	if backcmd != "" { //tableが作成済み
		cmd, err = createTableCmd(tname, stu, ifnotOff)
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
		cmd, err = createTableCmd(tname, stu, ifnotOn)
		if err != nil {
			return err
		}
		_, err = t.db.Exec(cmd)

	}

	return err
}

// ReadTableList() = []string, error
//
// SQL内のテーブル名を取得
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

// ReadCreateTableCmd(tname) = string, error
//
// SQL内の作成したテーブルのコマンド情報を取得
//
// tname(string) : 読み取り対象のテーブル
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

// DropTable(tname) = error
//
// SQL内のテーブルを削除する
//
// tname(string) : 削除対象のテーブル
func (t *sqliteConfig) DropTable(tname string) error {
	cmd, err := dropTableCmd(tname)
	if err != nil {
		return err
	}
	_, err = t.db.Exec(cmd)
	return err

}

// createTableCmd(tname, stu, flag) = string, error
//
// テーブルを作成するSQLコマンドを作成
// 構造体データにcreated_atとupdated_atを追加して、作成と更新のタイムスタンプをつける
//
// tname(string) : テーブル名
// stu(interface{}) : テーブル内のデータ精製オプション
// flag(ifnot) : IF NOT EXISTSをつけるオプション
func createTableCmd(tname string, stu interface{}, flag ifnot) (string, error) {
	cmd := "CREATE TABLE" + " "
	if flag == ifnotOn {
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

// altertableCmd(tname, cmdA, cmdB) = []string
//
// 追加作成用のテーブルを作るコマンド
// :ToDo
//
// cmdA : 登録してあるコマンド
// cmdB : これから設定するコマンド
func altertableCmd(tname, cmdA, cmdB string) []string {
	var output []string
	//ALTER TABLE tbl_name ADD COLUMN new_col VARCHAR(10) AFTER col1;

	return output
}

// dropTableCmd(tname) = string, error
//
// テーブルを削除するSQLコマンドを作る
//
// tname(string) : 削除対象のテーブル
func dropTableCmd(tname string) (string, error) {
	cmd := "DROP TABLE IF EXISTS" + " '" + tname + "'"
	return cmd, nil

}

// readTableAllCmd() = string, error
//
// SQLiteに登録してあるテーブルを取得するコマンドを作る
func readTableAllCmd() (string, error) {
	cmd := "SELECT name FROM sqlite_master WHERE type='table'"
	return cmd, nil

}

// readCreateTableCmd() = string, error
//
// SQLiteに登録してあるテーブルを作成したコマンドを読み取るSQLコマンドを作る
//
// tname(string) : 読み取り対象となるテーブル
func readCreateTableCmd(tname string) (string, error) {
	cmd := "SELECT sql FROM sqlite_master WHERE type='table' AND name='" + tname + "'"
	return cmd, nil
}
