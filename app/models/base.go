package models

//ここにはテーブルの作成のコードを書いていく
import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"go_todo/config"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

//よくわからないけどいるっぽい
var Db *sql.DB

var err error

const (
	tableNameUser    = "users"
	tableNameTodo    = "todos"
	tableNameSession = "sessions"
)

func init() {
	//ドライバの名前を指定してデーターベースに接続する
	Db, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}
	//Sprintf...(string+int)でも全体をまとめてstring型にしてくれる
	//どんなテーブルを作成したいのかidは？nameは？なんのカラムが欲しいの？
	//usertableがなければ作成する　//カラムの設定
	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		email STRING,
		password STRING,
		created_at DATETIME)`, tableNameUser)
	//実行！！！
	Db.Exec(cmdU)

	//todoテーブル
	cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME)`, tableNameTodo)

	Db.Exec(cmdT)

	cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		email STRING,
		user_id INTEGER,
		created_at DATETIME)`, tableNameSession)

	Db.Exec(cmdS)
}

//UUIDの生成.オブジェクトを識別するための一意のIDをあたえるもの
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

//パスワードの保存.パスワードはハッシュ値にする必要があるからcryptを使う
func Encrypt(plaintext string) (cryptext string) {
	/*
		sha1って何？
		SHAには種類があって,どれも暗号化してくれる.
		SHA-1は、160ビット（20バイト）のハッシュ値を生成する。
	*/
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
