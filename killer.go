package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type User struct {
	FirstName string `json:"firstname""`
	LastName  string `json:"lastname""`
	UID       string `json:"uid"`
}

type KilledList struct {
	PageTitle   string
	killedUsers []User
}

var allUsers = []User{}

func handler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		request.URL.Path = "./static/index.html"
	}
}

func getProfile(user User) string {
	ret := "<tr>"
	ret += "<th>" + user.FirstName + "</th>"
	ret += "<th>" + user.LastName + "</th>"
	ret += "<th>" + user.UID + "</th>"
	ret += "</tr>"
	return ret
}

func getAllUsers() []User {
	if len(allUsers) == 0 {
		file, err := os.Open("./data/user_data.json")
		if err != nil {
			log.Fatal(err)
		}
		decoder := json.NewDecoder(strings.NewReader(file))
		for {
			var user User
			if err := decoder.Decode(&user); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			allUsers = append(allUsers, user)
		}
	}
	return allUsers
}

func checkHandler(writer http.ResponseWriter, request *http.Request) {
	requestedUid := request.FormValue("unique_id")

}

func main() {
	router := mux.NewRouter()
	//router.HandleFunc("/", handler).Methods("GET")
	//router.HandleFi
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/"))).Methods("GET")
	router.HandleFunc("/check", checkHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
