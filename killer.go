package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	UID       string   `json:"uid"`
	Killed    []string `json:"killed"`
}

type KilledList struct {
	PageTitle   string
	KilledUsers []User
}

type NoSuchUser struct {
	PageTitle string
	UID       string
}

type UserList struct {
	AllUsers []User `json:"knownusers"`
}

var allUsers UserList

func handler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		request.URL.Path = "./static/index.html"
	}
}

func getAllUsers() []User {
	if len(allUsers.AllUsers) == 0 {
		file, err := ioutil.ReadFile("./data/user_data.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &allUsers)
	}
	fmt.Println(allUsers.AllUsers)
	return allUsers.AllUsers
}

func isUidExist(uid string) (User, bool) {
	for _, elem := range getAllUsers() {
		if elem.UID == uid {
			return elem, true
		}
	}
	return User{}, false
}

func checkHandler(writer http.ResponseWriter, request *http.Request) {
	requestedUid := request.FormValue("unique_id")
	userData, exist := isUidExist(requestedUid)
	killedList := KilledList{
		PageTitle:   "Мои жертвы",
		KilledUsers: []User{},
	}
	if exist {
		for _, target := range userData.Killed {
			killedUserData, exists := isUidExist(target)
			if exists {
				killedList.KilledUsers = append(killedList.KilledUsers, killedUserData)
			}
		}
	} else {
		nosuchUserHandler(writer, request)
		return
	}
	tmpl := template.Must(template.ParseFiles("./static/profile.html"))
	tmpl.Execute(writer, killedList)
}

func nosuchUserHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/nosuchuser.html"))
	tmpl.Execute(writer, NoSuchUser{
		PageTitle: "No such user id:",
		UID:       request.FormValue("unique_id"),
	})
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/"))).Methods("GET")
	router.HandleFunc("/check", checkHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
