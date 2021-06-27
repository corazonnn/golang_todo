package controllers

import (
	"fmt"
	"go_todo/app/models"
	"go_todo/config"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

//ここにはサーバーの立ち上げのコードを書いていく

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	// fmt.Println("dataの中身は", data)
	templates.ExecuteTemplate(w, "layout", data) //(...,実行するテンプレート,渡すデータ)

	/*
		第二引数の明示的に示したテンプレートの中だけ{{template ""}}が使える.
		要するに、layout.htmlからは{{define "content"}}に飛ばせるけど、index.html内では{{define "content"}}は使えない

	*/
	/*
		generateHTMLの引数では〇〇にアクセスがあった時に、どのページがみたいのかを引数として読み込んでいたけど、
		ExecuteTemplateでは実際その読み込んだものをどこで使いたいのかを示す。
	*/

}

/*
どのページに飛んでもまずこの関数を実行してログインしているかどうかの確認
cookieを取得する関数を使ってアクセス制限
Checksession：いるかいないかを返すメソッド、session:どんな人がいるかを返す関数
*/
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	/*
		cookieから値を受け取って
		受け取ったcookieのUUIDがDBに存在するUUIDと同じか
	*/
	cookie, err := r.Cookie("_cookie") //rにはCookieメソッドがある。そもそも今来たやつはcookieを持っているのか？（２回目以降の客なのか）を確認
	if err == nil {
		sess = models.Session{UUID: cookie.Value} //ここで新しくsession構造体を作成
		if ok, _ := sess.CheckSession(); !ok {    //そのsessがDB内にあるかどうかを確認
			err = fmt.Errorf("invalid session")
		}
	}
	return sess, err
}

var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

//リクエストがあったらそのidを取得する関数が必要
//パターンとして覚えちゃっていい
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	fmt.Println("parseURL通過")
	return func(w http.ResponseWriter, r *http.Request) {
		//todos/edit/1。この１の部分を取得したい
		//validpathとurlがマッチした日にっmwppしら遅イソタオ
		q := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Println("rの中身は", r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		// fmt.Println("qiの中身は", qi)
		fn(w, r, qi)
	}
}

func StartMainServer() error {
	//静的ファイルを読み込みたい
	files := http.FileServer(http.Dir(config.Config.Static))

	//http.handleは特定のURLとhandler(serverHTTP...HTTPリクエストに対してレスポンスを返す)を紐づける(DefaultServerMuxに登録する)
	http.Handle("/static/", http.StripPrefix("/static/", files))

	//第一引数の場所にアクセスしたらtopにいくようにする
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", auhtenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	//   save　と edit/ の違いはURLが完全に一致するかどうか.
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))

	return http.ListenAndServe(":"+config.Config.Port, nil)
}
