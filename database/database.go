package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test")

	return db, err
}

func Create() {
	db, err := Connect()
	if err != nil {
		panic(err.Error())
	}
	db.Exec("CREATE TABLE User (" +
		"\n id INT PRIMARY KEY NOT NULL AUTO_INCREMENT," +
		"\n login VARCHAR(255)," +
		"\n mail VARCHAR(255)," +
		"\n pwd VARCHAR(255)," +
		"\n CONSTRAINT Unique_login UNIQUE (login)," +
		"\n CONSTRAINT Unique_mail UNIQUE (mail)" +
		"\n )")
	db.Exec("CREATE TABLE Music (" +
		"\n id INT PRIMARY KEY NOT NULL" +
		"\n )")
	db.Exec("CREATE TABLE Comment (" +
		"\n id INT PRIMARY KEY NOT NULL AUTO_INCREMENT," +
		"\n id_Music INT NOT NULL," +
		"\n id_User INT NOT NULL," +
		"\n date_p DATETIME DEFAULT CURRENT_TIMESTAMP," +
		"\n msg VARCHAR(1000) NOT NULL," +
		"\n likes INT DEFAULT 0," +
		"\n FOREIGN KEY (id_Music) REFERENCES Music(id)," +
		"\n FOREIGN KEY (id_User) REFERENCES User(id)" +
		"\n )")
	db.Exec("CREATE TABLE Music_Like (" +
		"\n id_Music INT NOT NULL," +
		"\n id_User INT NOT NULL," +
		"\n PRIMARY KEY(id_User, id_Music)," +
		"\n FOREIGN KEY (id_Music) REFERENCES Music(id)," +
		"\n FOREIGN KEY (id_User) REFERENCES User(id)" +
		"\n )")
	db.Exec("CREATE TABLE Comment_Like (" +
		"\n id_Comment INT NOT NULL," +
		"\n id_User INT NOT NULL," +
		"\n PRIMARY KEY (id_Comment, id_User)," +
		"\n FOREIGN KEY (id_Comment) REFERENCES Comment(id)," +
		"\n FOREIGN KEY (id_User) REFERENCES User(id)" +
		"\n )")
	db.Exec("CREATE TABLE Connection (" +
		"\n id INT NOT NULL," +
		"\n idSession VARCHAR(32) NOT NULL," +
		"\n PRIMARY KEY (id)," +
		"\n FOREIGN KEY (id) REFERENCES User(id)" +
		"\n )")
}
