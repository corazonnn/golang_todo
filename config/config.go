package config

import (
	"go_todo/utils"
	"log"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port      string //ポート番号
	SQLDriver string //使用するSQL
	DbName    string //DBの名前
	LogFile   string //logを残すファイル
	Static    string //bootstrapとjqueryがあるファイルの階層を設定する
}

var Config ConfigList //グローバル変数を宣言

func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile) //mainより前に読み込んでおきたい
}

func LoadConfig() {
	cfg, err := ini.Load("config.ini") //iniファイルを読み込む
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{ //structの中に代入
		Port:      cfg.Section("web").Key("port").MustString("8080"), //MustString:iniファイルに対象の情報がなければ8080を入れる
		SQLDriver: cfg.Section("db").Key("driver").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		LogFile:   cfg.Section("web").Key("logfile").String(),
		Static:    cfg.Section("web").Key("static").String(),
	}
}
