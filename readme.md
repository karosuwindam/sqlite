## sqlite機能のモジュール化

SQLite3の制御を別で定義して作成を便利にする関数を用意する。

Ubuntuで使用する場合はgccが必要なので、以下のコマンドで導入をしておく
```
sudo apt install gcc
```

## 使用方法

* 基本
以下の通り、記載することで、以下手順を実施します。
1. SQLiteファイルを開く
2. 不具合が発生したら、ファイルを閉じる
3. 構造体からテーブルの作成
4. 20個のデータを追加
5. ランダムで一つのデータを削除
6. 削除したIDを指定してを読み取り

```go:main.go
package main

import (
	"fmt"
	"math/rand"
	"github.com/karosuwindam/sqlite"
	"time"
)

func main() {

	type TableTest struct {
		Id  int    `db:"id"`
		Str string `db:"str"`
		I   int    `db:"i"`
	}

	testtablename := "test"
	testdbname := "test.db"

	sql := sqlite.Setup(testdbname)
	if err := sql.Open(); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sql.Close()
	if err := sql.CreateTable(testtablename, TableTest{}); err != nil {
		fmt.Println(err.Error())
		return
	}

	ckdata := []TableTest{}

	wdata := TableTest{Id: 1, Str: "data", I: 500}
	for i := 0; i < 20; i++ {
		wdata.I += i
		if err := sql.Add(testtablename, &wdata); err != nil {
			fmt.Println(err.Error())
		}
		ckdata = append(ckdata, wdata)
	}

	rand.Seed(time.Now().UnixNano())
	deleteid := rand.Intn(20)
	for {
		if deleteid != 0 {
			break
		}
		deleteid = rand.Intn(20)
	}
	fmt.Printf("-----------delete id=%v ---------------", deleteid)
	if err := sql.Delete(testtablename, deleteid); err != nil {
		fmt.Println(err.Error())
	}

	rdata := []TableTest{}
	if err := sql.Read(testtablename, &rdata, map[string]string{"id": strconv.Itoa(deleteid)}, sqlite.AND); err != nil {
		fmt.Println(err.Error())
	}
	if len(rdata) == 0 {
		fmt.Printf("No Delete data id=%v", deleteid)
	}

}
```


* その他
ほかの使い方は各テストファイルを参照にしてください。

