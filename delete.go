package sqlite

import "strconv"

// (*cfg)Delete(tname,id)
//
// SQLのテーブルからid指定で削除する
//
// tname:id指定するための参照テーブル
// id:削除指定を出すID
func (cfg *SqliteConfig) Delete(tname string, id int) error {
	cmd, err := createDelCmd(tname, id)
	if err != nil {
		return err
	}
	_, err = cfg.db.Exec(cmd)
	return err

}

//createDelCmd(tname,id)
//
// SQL用の削除コマンドを作成
//
// tname(string):削除対応のテーブル指定
// id:削除指定を出すID
func createDelCmd(tname string, id int) (string, error) {
	cmd := "DELETE" + " " + "FROM" + " " + tname
	cmd += " " + "WHERE" + " " + "id=" + strconv.Itoa(id)
	return cmd, nil
}
