package main

import (
	"github.com/bouchenakihabib/PC3R_DEEZER/src/comment"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/database"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/like"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/music"
	"github.com/bouchenakihabib/PC3R_DEEZER/src/user"
	"log"
	"net/http"
)

func handleUser(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	switch req.Method {
	case "GET":
		user.GetUser(resp, req)
	case "POST":
		user.AddUser(resp, req)
	case "DELETE":
		user.DeleteUser(resp, req)
	}
}

func handleMusic(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	switch req.Method {
	case "GET":
		music.SearchMusic(resp, req)
	}
}

func handleComment(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	switch req.Method {
	case "GET":
		comment.GetComment(resp, req)
	case "POST":
		comment.AddComment(resp, req)
	case "DELETE":
		comment.DeleteComment(resp, req)
	}
}

func handleLikeMusic(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	switch req.Method {
	case "GET":
		like.GetLikeMusic(resp, req)
	case "POST":
		like.AddLikeMusic(resp, req)
	case "DELETE":
		like.DeleteLikeMusic(resp, req)
	}
}

func handleLikeComment(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	switch req.Method {
	case "GET":
		like.GetLikeComment(resp, req)
	case "POST":
		like.AddLikeComment(resp, req)
	case "DELETE":
		like.DeleteLikeComment(resp, req)
	}
}

func main() {
	log.Printf("test")

	database.Create()

	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/music", handleMusic)
	http.HandleFunc("/comment", handleComment)
	http.HandleFunc("/like_music", handleLikeMusic)
	http.HandleFunc("/like_comment", handleLikeComment)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err.Error())
	}

}
