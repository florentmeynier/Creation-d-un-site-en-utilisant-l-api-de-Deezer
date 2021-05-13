package connection

import (
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/user"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"math/rand"
	"net/http"
)

func Connect(resp http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	pwd := req.FormValue("pwd")
	userId := user.GetUserIdFromLogin(login)
	match := user.PwdMatchWithUser(userId, pwd)
	if match == "" {
		idSession := AddConnection(userId)
		if idSession == "" {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"Error with session"}`)
			return
		}
		utils.Response(resp, http.StatusOK, `{"message":"User Connected", "idSession":"`+idSession+`"}`)
	} else {
		utils.Response(resp, http.StatusBadRequest, `{"message":"`+match+`"}`)
	}
}

func Disconnect(resp http.ResponseWriter, req *http.Request) {
	idSession := req.FormValue("idSession")
	res := RemoveConnectionSession(idSession)
	if res == "" {
		utils.Response(resp, http.StatusOK, `{"message":"User disconnected"}`)
	} else {
		utils.Response(resp, http.StatusBadRequest, `{"message":"`+res+`"}`)
	}
}

func GetConnection(resp http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	idSession := req.FormValue("idSession")
	if login != "" {

	}
	if idSession != "" {
		userId := GetUserIdFromIdSession(idSession)
		if userId != "" {
			utils.Response(resp, http.StatusOK, `{"message":"User found","userId":"`+userId+`"}`)
			return
		}
	}
	utils.Response(resp, http.StatusBadRequest, `{"message":"User not connected"}`)
}

func UserIsConnected(userId string) bool {
	db, err := database.Connect()
	if err != nil {
		return false
	}
	exist := 0
	err = db.QueryRow("SELECT COUNT(*) FROM Connection WHERE id = ?", userId).Scan(&exist)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}
	return exist == 1
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func AddConnection(userId string) string {
	db, err := database.Connect()
	if err != nil {
		return ""
	}
	ru := make([]rune, 32)
	for i := range ru {
		ru[i] = letters[rand.Intn(len(letters))]
	}
	idSession := string(ru)
	if UserIsConnected(userId) {
		if RemoveConnection(userId) != "" {
			return ""
		}
	}
	res, err := db.Exec("INSERT INTO Connection (id, idSession) VALUES(?, ?)", userId, idSession)
	if err != nil {
		return ""
	}
	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		return ""
	}
	return idSession
}

func RemoveConnection(userId string) string {
	db, err := database.Connect()
	if err != nil {
		return "Error with database"
	}
	res, err := db.Exec("DELETE FROM Connection WHERE id = ?", userId)
	if err != nil || db.Close() != nil {
		return "Error while disconnecting"
	}
	r, err := res.RowsAffected()
	if r == 0 {
		return "No user disconnected"
	}
	if err != nil {
		return "Error with database"
	}
	return ""
}

func RemoveConnectionSession(idSession string) string {
	db, err := database.Connect()
	if err != nil {
		return "Error with database"
	}
	res, err := db.Exec("DELETE FROM Connection WHERE idSession = ?", idSession)
	if err != nil || db.Close() != nil {
		return "Error while disconnecting"
	}
	r, err := res.RowsAffected()
	if r == 0 {
		return "No user disconnected"
	}
	if err != nil {
		return "Error with database"
	}
	return ""
}

func GetUserIdFromIdSession(idSession string) string {
	db, err := database.Connect()
	if err != nil {
		return ""
	}
	id := ""
	err = db.QueryRow("SELECT id FROM Connection WHERE idSession = ?", idSession).Scan(&id)
	if err != nil || db.Close() != nil {
		return ""
	}
	return id
}
