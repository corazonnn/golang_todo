package models

import (
	"log"
	"time"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

//タスクを作成する！！！！！！
func (u *User) CreateTodo(content string) (err error) {
	cmd := `insert into todos (
		content, 
		user_id,
		created_at) values (?,?,?)`

	_, err = Db.Exec(cmd, content, u.ID, time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetTodo(id int) (todo Todo, err error) {
	cmd := `select id ,content, user_id, created_at from todos where id= ?`
	todo = Todo{}
	err = Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt)

	return todo, err
}

//全todoを取得したい！！！！！！！！！！
func GetTodos() (todos []Todo, err error) {
	cmd := `select id, content, user_id, created_at from todos`
	rows, err := Db.Query(cmd)
	//rowsはまだ意味不明な値:&{0xc000148090 0x40f6520 0xc00011c540 <nil> <nil> {{0 0} 0 0 0 0} false <nil> []}こんな感じ

	if err != nil {
		log.Fatalln(err)
	}
	//とってきた情報をloopで回す
	for rows.Next() {
		/*
						todo = Todo{}とvar todo Todoの宣言の違いは何か？？
			構造体の初期化方法は複数存在する。
			①変数定義後にフィールドを設定する方法 var user User
			②{}で順番にフィールドの値を渡す方法 user := User{}
			③
			色々あるから今回の場合はどっちでもいい。
		*/
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)

		if err != nil {
			log.Fatalln(err)
		}
		//append(追加元となる[]Todo型のtodos,追加する[]Todo型のtodo)
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

//特定のユーザのタスクを取得したい！！！！！！(ユーザのメソッドとして定義する)
func (u *User) GetTodoByUser() (todos []Todo, err error) {
	cmd := `select id, content, user_id, created_at from todos where user_id = ?`
	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	/*
		Closeの必要性は？？なんとなくは理解していてもどんな時にすればいいのかわからない
	*/
	rows.Close()

	return todos, err
}

func (t *Todo) UpdateTodo() error {
	cmd := `update todos set content = ?, user_id = ? where id = ?`
	_, err = Db.Exec(cmd, t.Content, t.UserID, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (t *Todo) DeleteTodo() error {
	cmd := `delete from todos where id = ?`
	_, err = Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
