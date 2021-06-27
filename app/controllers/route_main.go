package controllers

// "/"にアクセスがあった時にどんな処理をするかを書いてある
import (
	"fmt"
	"go_todo/app/models"
	"log"
	"net/http"
)

func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
	// fmt.Println("topにはきてる")

}
func index(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("topにはきてる")
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		todos, _ := user.GetTodoByUser()
		user.Todos = todos
		// fmt.Println("userの中身は", user.Name)
		generateHTML(w, user, "layout", "index", "private_navbar")
	}
}
func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "login", 302)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

func todoSave(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)

	if err != nil {
		http.Redirect(w, r, "login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}

		content := r.PostFormValue("content")
		if err := user.CreateTodo(content); err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/todos", 302)
	}
}

/*
edit/:idにアクセス
アクセスされたtodoを特定(ここがわからない)
そのtodoの編集ページに飛ばす
*/
func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	fmt.Println("きてますか？？")
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

/*
【流れ】
0)引数で変更したいTODOのidを渡す
1)今ログインしてるユーザの情報を取得
取得できたら,
2)フォームの中身を解析
3)その中からname=contentを抜き出す
4)Sessionユーザ情報からUserユーザの情報を取得
5)t= 新しいTODOのstruct(引数のid,フォームからのcontent,sessionから取ったuser_id)
6)t.UpdateTodo()で実際に作成
7)成功したらリダイレクト
*/
func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r) //ログインしてるユーザの情報を取ってくる
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm() //TODOを更新するためにフォームにどんなのが入っているのかparseする
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content") //parseした中からname=contentの情報が欲しい
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := t.UpdateTodo(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos", 302)
	}
}

/*
【流れ】
0)引数で削除したいTODOのidを渡す
1)今ログインしてるユーザの情報を取得
取得できたら,
2)Sessionユーザ情報からUserユーザの情報を取得
3)引数のidから削除したいTODO情報を取得
4)t.DeleteTodo()で削除
5)
6)
7)成功したらリダイレクト
*/
func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(id) //指定したidからTODOを取得
		if err != nil {
			log.Println(err)
		}
		if err := t.DeleteTodo(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/todos", 302)
	}
}
