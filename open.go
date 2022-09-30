package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //sqliteを使用しているので
)

// sqliteConfig SQLiteのテーブル設定
type sqliteConfig struct {
	filepass string  // ファイルパス
	db       *sql.DB // 開いたデータベースの値
}

// Setup sqliteConfig=Setup(filepath)
//
// 基本セットアップ
//
// filepath(string): sqlite3のデータベースパス
func Setup(filepath string) sqliteConfig {
	return sqliteConfig{filepass: filepath}
}

// Open()
//
// SQLiteのファイルを開く
func (t *sqliteConfig) Open() error {
	var err error
	t.db, err = sql.Open("sqlite3", t.filepass)
	return err
}

// Close()
//
// SQLiteのファイルを閉じる
func (t *sqliteConfig) Close() error {
	return t.db.Close()

}

//ファイルパスの読み取り
func (t *sqliteConfig) ReturnFilePass() string { return t.filepass }
