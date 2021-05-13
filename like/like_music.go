package like

import (
	"encoding/json"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"net/http"
)

func GetLikeMusic(resp http.ResponseWriter, req *http.Request) {
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	if idMusic == "" && idUser == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	request := ""
	if idMusic != "" {
		request += "id_Music = " + idMusic
		if idUser != "" {
			request += "AND id_User = " + idUser
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
	res, err := db.Query("SELECT * FROM music_like WHERE " + request)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while searching music likes"}`)
		return
	}
	likes := make([]utils.LikeMusic, 0)
	for res.Next() {
		l := utils.LikeMusic{}
		err := res.Scan(&l.IdMusic, &l.IdUser)
		if err != nil {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while collecting music like"}`)
			return
		}
		likes = append(likes, l)
	}
	jlikes, err := json.Marshal(likes)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"An error occured with the result"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"likes found", "result":`+string(jlikes)+`}`)
}

func AddLikeMusic(resp http.ResponseWriter, req *http.Request) {
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	if idUser == "" && idMusic == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	db, err := database.Connect()
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	res, err := db.Exec("INSERT INTO music_like(id_Music, id_User) VALUES (?, ?)", idMusic, idUser)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while adding a music like"}`)
		return
	}
	r, err := res.RowsAffected()
	if r == 0 {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"User has already liked this music}`)
		return
	}
	if err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Music liked"}`)
}

func DeleteLikeMusic(resp http.ResponseWriter, req *http.Request) {
	idMusic := req.FormValue("id_Music")
	idUser := req.FormValue("id_User")
	if idUser == "" && idMusic == "" {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Argument missing"}`)
		return
	}
	request := ""
	if idMusic != "" {
		request += "id_Music = " + idMusic
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
	res, err := db.Exec("DELETE FROM music_like WHERE " + request)
	if err != nil || db.Close() != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured while deleting music likes"}`)
		return
	}
	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"An error occured with database"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Music likes deleted"}`)
}
