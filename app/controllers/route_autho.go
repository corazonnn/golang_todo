package controllers

import (
	"go_todo/app/models"
	"log"
	"net/http"
)

//いろんなハンドラを記述していく

//signupのテンプレートファイルの出力だけをしたい

//(w,渡すデータ、なんでlayoutとsignup???)
//第3引数には、このアクセスがあったときに使いたい、表示したいページ
//最終的にgenerateHTMLの中で第3引数のファイルを読み込んでいるtemplate.ParseFiles()で

func signup(w http.ResponseWriter, r *http.Request) {
	//同じ/signupでもGETでくる場合とPOSTでくる場合がある
	if r.Method == "GET" {
		/*
			zoroが１回目に来たときは、cookieを持ってないので、帰る時にコンピュータ側がcookieを渡してくれる
			zoroが２回目以降訪れるときは、そのcookieをもって行き、コンピュータがそのcookieが
		*/
		_, err := session(w, r) //今アクセスして来たやつがログインしているのか確認(そもそもcookie持ってるんか(2回目以降か)＋＋そのcookieはsession_idとして保存されてるか(ログイン状態なのか))
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}

	} else if r.Method == "POST" {
		//新しいユーザを作成したい

		//ParseForm()を行うことでデータをFormから取得できるようになる
		err := r.ParseForm()
		if err != nil {
			log.Print()
		}
		user := models.User{
			Name:     r.PostFormValue("name"), //formの中のname="name"から取り出してる
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}
		/*
			ユーザの登録が成功したらトップページにリダイレクトしたい
			http.Redirect(ResponseWriter, HttpRequest, どこに飛ばしたいか, ステータスコード)
		*/
		http.Redirect(w, r, "/", 302)

	}
}

//ログインページを表示したい
func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}

}

//現在ログインしているかどうかの確認（ログインフォームから入力してPOSTでここにくる）
func auhtenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}
	//保存されているパスワードは暗号化されているからフォームから来たやつも暗号化して比較する(encryptするらしい
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		//そういうものって覚えるらしい。Cookieへの保存の仕方
		cookie := http.Cookie{
			Name:     "_cookie",    //cookie自体の名前
			Value:    session.UUID, //cookieの値
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		//ここでやっとログインに成功した

		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

//ログアウトしたいーーーーって言われたらこのハンドラに飛んでくる
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie") //cookieがそもそもあるのか確認
	if err != nil {
		log.Println(err)
	}
	//http.ErrNoCookieは「その名前のCookieが存在しない」エラーに限定した処理(今回はCookieが存在する処理)
	if err != http.ErrNoCookie { //もしcookieがあるなら、
		session := models.Session{UUID: cookie.Value} //そのcookieの情報を使って新しいsessionを生成

		session.DeleteSessionByUUID() //それと一致するものをsessionから削除する
	}
	http.Redirect(w, r, "/login", 302)
}
