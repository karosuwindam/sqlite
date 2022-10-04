package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //sqliteを使用しているので
)

// SqliteConfig SQLiteのテーブル設定
type SqliteConfig struct {
	filepass string  // ファイルパス
	db       *sql.DB // 開いたデータベースの値
}

// Setup SqliteConfig=Setup(filepath)
//
// 基本セットアップ
//
// filepath(string): sqlite3のデータベースパス
func Setup(filepath string) SqliteConfig {
	return SqliteConfig{filepass: filepath}
}

// (*cfg)Open()
//
// SQLiteのファイルを開く
func (cfg *SqliteConfig) Open() error {
	var err error
	cfg.db, err = sql.Open("sqlite3", cfg.filepass)
	return err
}

// (*cfg)Close()
//
// SQLiteのファイルを閉じる
func (cfg *SqliteConfig) Close() error {
	return cfg.db.Close()

}

// (*cfg)ReturnFilePass()
//
// ファイルパスの読み取り
func (cfg *SqliteConfig) ReturnFilePass() string { return cfg.filepass }
