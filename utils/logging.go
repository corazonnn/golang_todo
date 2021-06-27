package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	//渡されたlogfileを条件をつけて呼び出す
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //読み書き,ファイルがなければ作成,追記
	if err != nil {
		log.Fatalln(err)
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)   //logの書き込み先を、標準出力とlogfileに指定している(まだここでは設定の段階)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //logのフォーマットを設定
	log.SetOutput(multiLogFile)                          //logの出力先を設定(ここでやっと実行)
}
