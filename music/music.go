package music

import (
	"encoding/json"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func SearchMusic(resp http.ResponseWriter, req *http.Request) {
	search := req.FormValue("search")
	musicId := req.FormValue("music_Id")
	if search != "" {
		Search(resp, search)
	} else {
		if musicId != "" {
			SearchMusicId(resp, musicId)
		}
	}

}

func Search(resp http.ResponseWriter, search string) {
	search = strings.ReplaceAll(search, " ", "%20")
	research := "https://api.deezer.com/search?q=" + search
	apiResp, apiErr := http.Get(research)
	if apiErr != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"API Error"}`)
		return
	}
	body, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Error while collecting ressources"}`)
		return
	}

	var data utils.SearchJSON
	json.Unmarshal(body, &data)
	if len(data.Data) == 0 {
		utils.Response(resp, http.StatusBadRequest, `{"message":"No music found"}`)
		return
	}

	for _, d := range data.Data {
		res := InsertMusic(d.ID)

		if res != "" {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
			return
		}
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"An error occured with the result"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Research completed", "result":`+string(jdata)+`}`)
}

func SearchMusicId(resp http.ResponseWriter, musicId string) {
	research := "https://api.deezer.com/track/" + musicId
	apiResp, apiErr := http.Get(research)
	if apiErr != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"API Error"}`)
		return
	}
	body, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"Error while collecting ressources"}`)
		return
	}

	var track utils.Track
	json.Unmarshal(body, &track)

	res := InsertMusic(track.ID)

	if res != "" {
		utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
		return
	}

	jdata, err := json.Marshal(track)
	if err != nil {
		utils.Response(resp, http.StatusBadRequest, `{"message":"An error occured with the result"}`)
		return
	}
	utils.Response(resp, http.StatusOK, `{"message":"Research completed", "result":`+string(jdata)+`}`)
}

func InsertMusic(id int) string {
	if !ContainsMusic(id) {
		db, err := database.Connect()
		if err != nil {
			return "An error occured with database"
		}
		res, err := db.Exec("INSERT INTO Music(id) VALUES (?)", id)
		if err != nil || db.Close() != nil {
			panic(err.Error())
			return "An error occured while adding Music"
		}
		r, err := res.RowsAffected()
		if r == 0 || err != nil {
			return "An error occured"
		}
	}
	return ""
}

func ContainsMusic(id int) bool {
	db, err := database.Connect()
	if err != nil {
		panic(err.Error())
	}
	contains := 0
	err = db.QueryRow("SELECT COUNT(*) FROM Music WHERE id = ?;", id).Scan(&contains)
	if err != nil || db.Close() != nil {
		panic(err.Error())
		return true
	}
	return contains == 1
}
