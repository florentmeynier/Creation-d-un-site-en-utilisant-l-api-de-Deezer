package music

import (
	"encoding/json"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/utils"
	"io/ioutil"
	"net/http"
)

func SearchMusic(resp http.ResponseWriter, req *http.Request) {
	search := req.FormValue("search")
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
	for _, d := range data.Data {
		res := InsertMusic(d.ID, d.Title, d.Artist.Name, d.Album.Title)
		if res != "" {
			utils.Response(resp, http.StatusInternalServerError, `{"message":"`+res+`"}`)
		}
	}
	utils.Response(resp, http.StatusOK, `{"message":"Research completed"}`)
}

func InsertMusic(id int, title string, artist string, album string) string {
	if !ContainsMusic(id) {
		db, err := database.Connect()
		if err != nil {
			return "An error occured with database"
		}
		res, err := db.Exec("INSERT INTO Music(id, title, artist, album) VALUES (?, ?, ?, ?)",
			id, title, artist, album)
		if err != nil || db.Close() != nil {
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
