package user

import (
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"net/http"
)

func GetUser(resp http.ResponseWriter, req *http.Request) {
	// TODO
}

func AddUser(resp http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	mail := req.FormValue("mail")
	pwd := req.FormValue("pwd")
	if IsValidUserInformation(login, mail, pwd) {
		res := InsertUser(login, mail, pwd)
		if res == "" {
			utils.Response(resp, http.StatusOK, `{"message":"New user created"}`)
		} else {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
		}
	} else {
		utils.Response(resp, http.StatusNotFound, `{"message":"Wrong information"}`)
	}
}

func DeleteUser(resp http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	if ExistUser(login) {
		res := RemoveUser(login)
		if res == "" {
			utils.Response(resp, http.StatusOK, `{"message":"User deleted"}`)
		} else {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
		}
	} else {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured"}`)
	}
}

func IsValidUserInformation(login string, mail string, pwd string) bool {
	return login != "" && mail != "" && pwd != ""
}

func InsertUser(login string, mail string, pwd string) string {
	db, err := database.Connect()
	if err != nil {
		return "An error occured with database" + err.Error()
	}
	if ExistUser(login) {
		return "Login already exist"
	}
	res, err := db.Exec("INSERT INTO User(login, mail, pwd) VALUES (?, ?, ?);", login, mail, pwd)
	if err != nil || db.Close() != nil {
		return "An error occured while adding user"
	}
	r, err := res.RowsAffected()
	if r == 0 {
		return "Login usermail already exist"
	}
	if err != nil {
		return "An error occured"
	}
	return ""
}

func RemoveUser(login string) string {
	db, err := database.Connect()
	if err != nil {
		return "An error occured" + err.Error()
	}
	res, err := db.Exec("DELETE FROM User WHERE login=?", login)
	if err != nil || db.Close() != nil {
		return "An error occured" + err.Error()
	}
	r, err := res.RowsAffected()
	if r == 0 {
		return "User has not been deleted"
	}
	if err != nil {
		return "An error occured" + err.Error()
	}
	return ""
}

func ExistUser(login string) bool {
	db, err := database.Connect()
	if err != nil {
		return false
	}
	exist := 0
	err = db.QueryRow("SELECT COUNT(*) FROM User WHERE login = ?;", login).Scan(&exist)
	if err != nil || db.Close() != nil {
		return false
	}
	return exist == 1
}
