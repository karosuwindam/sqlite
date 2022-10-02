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

// (*cfg)Open()
//
// SQLiteのファイルを開く
func (cfg *sqliteConfig) Open() error {
	var err error
	cfg.db, err = sql.Open("sqlite3", cfg.filepass)
	return err
}

// (*cfg)Close()
//
// SQLiteのファイルを閉じる
func (cfg *sqliteConfig) Close() error {
	return cfg.db.Close()

}

// (*cfg)ReturnFilePass()
//
// ファイルパスの読み取り
func (cfg *sqliteConfig) ReturnFilePass() string { return cfg.filepass }
