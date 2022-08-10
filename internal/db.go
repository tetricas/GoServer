package internal

import (
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	db *sqlx.DB

	dataInsert = `INSERT INTO users (name, email, secret, isAdmin) VALUES (:name, :email, :secret, :isAdmin)`
	dataDelete = `DELETE FROM users WHERE email = (?)`
	dataSelect = `SELECT * FROM users WHERE email = (?)`
)

type UserInternal struct {
	Name    string `db:"name"`
	Email   string `db:"email"`
	Secret  string `db:"secret"`
	IsAdmin bool   `db:"isAdmin"`
}

func ConnectToDB() {
	db, _ = sqlx.Connect("sqlite3", "testOne.db")
	path := filepath.Join("db", "create_db.sql")
	file, _ := ioutil.ReadFile(path)
	schema := string(file)
	db.MustExec(schema)
}

func AddUserToDB(user *UserInternal) {
	tx := db.MustBegin()
	_, err := tx.NamedExec(dataInsert, user)
	if err != nil {
		log.Println(err)
		return
	}
	_ = tx.Commit()
}

func GetUserFromDB(email string) (user *UserInternal) {
	tx := db.MustBegin()
	user = &UserInternal{}
	err := db.Get(user, dataSelect, email)
	if err != nil {
		log.Println(err)
		return
	}
	_ = tx.Commit()

	return user
}

func DeleteUserFromDB(email string) {
	tx := db.MustBegin()
	tx.MustExec(dataDelete, email)
	_ = tx.Commit()
}
