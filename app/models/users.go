package models

import (
	"fmt"
	"log"
	"time"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo
}
type Session struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

//userの作成. User structを持ってないとCreateUserは実行できない.
func (u *User) CreateUser() (err error) {
	/*
		uにはUser内のIDやpasswordがもうすでに入っている
		あとはその情報を元にuserを作成するだけ(コマンド作成→Db.Exec())
	*/
	//コマンド作成.まだ実行してない
	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values (?, ?, ?, ?, ?)`

	/*
		Dbはusers.go内には存在しないが,package models内にあるから使うことができる
		ファイル単位でスコープがあるのではなく,package単位でスコープがあるのか
	*/
	_, err = Db.Exec(cmd,
		createUUID(), //base.goにかいてある
		u.Name,
		u.Email,
		Encrypt(u.PassWord), //base.goにかいてある
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

/*
ユーザを取得したい
探したいユーザのidを引数として探す。そのidを手がかりとして、id,uuid,name,email,password,created_atの情報を探す

*/
func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `select id,uuid,name,email,password,created_at from users where id = ?`
	/*
		QueryRow:１レコードを取得、Query：複数レコードを取得
		Scan:sql文でcmdに情報が入っただけでアプリケーションの方にはきてない。
		　　　user = User{}で初期化しているuserの中にscan（データを移す）ため
	*/
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

//ユーザ更新
func (u *User) UpdateUser() (err error) {
	cmd := `update users set name= ?,email = ? where id = ?`
	//あるuさんの情報を変更したいのに、Db.Exec(uさんのName...)じゃ変更できてなくね??
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

//ユーザ削除
/*
なぜ u *Userで渡すのか??u Userでもよくないか??
もし後者で渡そうとすると、もともとあった吉祥寺のyuyaを八王子のyuyaに持ち出して変更しているから、
八王子では問題ないかもしれないが、本来解決したかった吉祥寺のyuyaは変更できてない
じゃあどうすればいいかというと、吉祥寺のyuyaを変更すればいい。そのために必要なのがポインタ
だからUserではなく*Userでないといけない。
*/
func (u *User) DeleteUser() (err error) {
	cmd := `delete from users where id = ?`
	//レコードを取得する必要のない、クエリはExecメソッドを使う
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

//ログインの際にemailを入力してもらい、そのemailを使ってDBからユーザを取得する
func GetUserByEmail(email string) (user User, err error) {
	fmt.Println("email:", email)
	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where email = ?`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)
	fmt.Println("userの中身は", user.ID)
	return user, err

}

//セッションを作るメソッド
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `insert into sessions (
		uuid,
		email,
		user_id,
		created_at) values (?,?,?,?)`
	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	//cmd1で作成したsessionをそのままcmd2で取得する
	cmd2 := `select id,uuid, email,user_id,created_at
	from sessions where user_id = ? and email = ?`
	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt)

	return session, err
}

//ユーザがログイン中かどうかを確かめるメソッド(渡したsessionがDBの中にあるかどうか)
func (sess *Session) CheckSession() (valid bool, err error) {
	cmd := `select id, uuid, email, user_id, created_at
	from sessions where uuid = ?`

	//queryrowで取得したらscanでsessに渡す
	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt)
	//sessが存在するかどうか
	if err != nil {
		valid = false
		return
	}
	//IDが初期値でないなら(なにかしらDBから見つけることができたなら)
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

func (sess *Session) DeleteSessionByUUID() (err error) {
	cmd := `delete from sessions where uuid=?`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `select id,uuid,name,email,created_at from users
	where id = ?`
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt)

	return user, err
}
