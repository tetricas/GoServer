package internal

import (
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	db *sqlx.DB

	dataInsert = `INSERT INTO users (name, email, isAdmin) VALUES (:name, :email, :isAdmin)`
	/*dataDelete = `DELETE FROM container WHERE data = (?)`
	dataSelect = `SELECT data FROM container`*/
)

type UserInternal struct {
	Name    string `db:"name"`
	Email   string `db:"email"`
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
