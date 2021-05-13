package comment

import (
	"encoding/json"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/like"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"net/http"
	"strconv"
)

func GetComment(resp http.ResponseWriter, req *http.Request) {
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	request := ""
	if idMusic != "" {
		request += "id_Music = " + idMusic
		if idUser != "" {
			request += " AND id_User = " + idUser
		}
	} else {
		if idUser != "" {
			request += "id_User = " + idUser
		} else {
			utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
			return
		}
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Query("SELECT * FROM Comment WHERE " + request)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while searching comments"}`)
		return
	}
	comments := make([]utils.Comment, 0)
	for res.Next() {
		c := utils.Comment{}
		err := res.Scan(&c.Id, &c.IdMusic, &c.IdUser, &c.Datep, &c.Msg, &c.Likes)
		c.Likes = like.GetNbLikesFromIdComment(strconv.Itoa(c.Id))
		if err != nil {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while collecting comments"}`)
			return
		}
		comments = append(comments, c)
	}
	jcomments, err := json.Marshal(comments)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"An error occured with the result"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"comments found", "result":`+string(jcomments)+`}`)
}

func AddComment(resp http.ResponseWriter, req *http.Request) {
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	msg := req.FormValue("msg")
	if idMusic == "" || idUser == "" || msg == "" {
		return
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Exec("INSERT INTO Comment(id_Music, id_User, msg) VALUES (?, ?, ?)", idMusic, idUser, msg)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while adding comment"}`)
		return
	}
	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"comment added"}`)
	return
}

func DeleteComment(resp http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	request := ""
	if id != "" {
		request += "id = " + id
		return
	} else {
		if idMusic != "" {
			request += "id_Music = " + idMusic
			if idUser != "" {
				request += "AND id_User = " + idUser
			}
		} else {
			if idUser != "" {
				request += "id_User = " + idUser
			} else {
				utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
				return
			}
		}
	}
	if request != "" {
		res := RemoveComment(request)
		if res != "" {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
			return
		}
		utils.Response(resp, http.StatusOK, `{"message":"comment deleted"}`)
		return
	}
	utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
}

func RemoveComment(request string) string {
	db, err := database.Connect()
	if err != nil {
		return "An error occured with database"
	}
	res, err := db.Exec("DELETE FROM Comment WHERE " + request)
	if err != nil || db.Close() != nil {
		return "An error occured while deleting comment"
	}
	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		return "An error occured with database"
	}
	return ""
}
