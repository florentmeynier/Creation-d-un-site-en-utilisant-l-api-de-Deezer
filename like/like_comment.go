package like

import (
	"encoding/json"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"net/http"
)

func GetLikeComment(resp http.ResponseWriter, req *http.Request) {
	idComment := req.FormValue("id_Comment")
	idUser := req.FormValue("id_User")
	if idComment == "" && idUser == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	request := ""
	if idComment != "" {
		request += "id_Comment = " + idComment
		if idUser != "" {
			request += "AND id_User = " + idUser
		}
	} else {
		if idUser != "" {
			idUser += "id_User = " + idUser
		}
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Query("SELECT * FROM comment_like WHERE " + request + ";")
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while searching music likes"}`)
		return
	}
	likes := make([]utils.LikeComment, 0)
	for res.Next() {
		l := utils.LikeComment{}
		err := res.Scan(&l.IdComment, &l.IdUser)
		if err != nil {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while collecting comment like"}`)
			return
		}
		likes = append(likes, l)
	}
	jlikes, err := json.Marshal(likes)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"An error occured with the result"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"likes found", "result":"`+string(jlikes)+`"}`)
}

func AddLikeComment(resp http.ResponseWriter, req *http.Request) {
	idComment := req.FormValue("id_Comment")
	idUser := req.FormValue("id_User")
	if idUser == "" && idComment == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Exec("INSERT INTO comment_like(id_Comment, id_User) VALUES (?, ?)", idComment, idUser)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while adding a comment like"}`)
		return
	}
	r, err := res.RowsAffected()
	if r == 0 {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"User has already liked this comment}`)
		return
	}
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Comment liked"}`)
}

func DeleteLikeComment(resp http.ResponseWriter, req *http.Request) {
	idComment := req.FormValue("id_Comment")
	idUser := req.FormValue("id_User")
	if idUser == "" && idComment == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	request := ""
	if idComment != "" {
		request += "id_Comment = " + idComment
		if idUser != "" {
			request += " AND id_User = " + idUser
		}
	} else {
		if idUser != "" {
			request += "id_User = " + idUser
		}
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Exec("DELETE FROM comment_like WHERE " + request)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while deleting comment likes"}`)
		return
	}
	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Comment likes deleted"}`)
}

func GetNbLikesFromIdComment(id string) int {
	db, err := database.Connect()
	if err != nil {
		return -1
	}
	exist := -1
	err = db.QueryRow("SELECT COUNT(*) FROM comment_like WHERE id_Comment = ?", id).Scan(&exist)
	if err != nil {
		return -1
	}
	err = db.Close()
	if err != nil {
		return -1
	}
	return exist
}
