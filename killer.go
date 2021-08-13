package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"os"
	"github.com/gorilla/mux"
)

type User struct {
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	UID       string   `json:"uid"`
	ToKillUID string   `json:"tokilluid"`
	IsKilled  bool     `json:"iskilled"`
	Killed    []string `json:"killed"`
}

type KilledList struct {
	PageTitle   string
	UserUID     string
	KilledUsers []User
	TargetName  string
}

type NoSuchUser struct {
	PageTitle string
	UID       string
}

type UserList struct {
	Locker sync.Mutex
	AllUsers []User `json:"knownusers"`
}

var allUsers UserList

func handler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		request.URL.Path = "./static/index.html"
	}
}

func getAllUsers() []User {
	allUsers.Locker.Lock()
	if len(allUsers.AllUsers) == 0 {
		file, err := ioutil.ReadFile("./data/user_data.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &allUsers)
	}
	fmt.Println(allUsers.AllUsers)
	allUsers.Locker.Unlock()
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
		UserUID:     requestedUid,
		TargetName:  "",
		KilledUsers: []User{},
	}
	if exist {
		killedList.PageTitle = userData.FirstName + " " + userData.LastName;
		targetUser, _ := isUidExist(userData.ToKillUID)
		killedList.TargetName = targetUser.FirstName + " " + targetUser.LastName
		for _, target := range userData.Killed {
			killedUserData, exists := isUidExist(target)
			if exists {
				killedList.KilledUsers = append(killedList.KilledUsers, killedUserData)
			}
		}
	} else {
		nosuchUserHandler(writer, request, "unique_id")
		return
	}
	tmpl := template.Must(template.ParseFiles("./static/profile.html"))
	tmpl.Execute(writer, killedList)
}

func nosuchUserHandler(writer http.ResponseWriter, request *http.Request, formId string) {
	tmpl := template.Must(template.ParseFiles("./static/nosuchuser.html"))
	fmt.Println(formId)
	tmpl.Execute(writer, NoSuchUser{
		PageTitle: "No such user id:",
		UID:       request.FormValue(formId),
	})
}

func isSameUser(userA User, userB User) bool {
	return userA.UID == userB.UID
}

func saveJson() {
	json, _ := json.Marshal(allUsers)
	f, _ := os.Create("./data/user_data.json")
	defer f.Close()
	f.Write(json)
}

func processKill(killer User, target User) {
	if isSameUser(killer, target) {
		return
	}
	if killer.ToKillUID != target.UID {
		return
	}
	if target.IsKilled {
		return
	}
	if killer.IsKilled {
		return
	}
	for pos, curUser := range getAllUsers() {
		if isSameUser(curUser, target) {
			allUsers.Locker.Lock()
			allUsers.AllUsers[pos].IsKilled = true
			allUsers.Locker.Unlock()
			break
		}
	}
	for pos, curUser := range getAllUsers() {
		if isSameUser(curUser, killer) {
			allUsers.Locker.Lock()
			allUsers.AllUsers[pos].Killed = append(allUsers.AllUsers[pos].Killed, target.Killed...)
			allUsers.AllUsers[pos].Killed = append(allUsers.AllUsers[pos].Killed, target.UID)
			allUsers.AllUsers[pos].ToKillUID = target.ToKillUID
			allUsers.Locker.Unlock()
			break
		}
	}
	saveJson()
	return
}
	
func killHandler(writer http.ResponseWriter, request *http.Request) {
	requestedUID := request.FormValue("killapply")
	target, targetExist := isUidExist(requestedUID)
	if targetExist {
		userUID := request.FormValue("unique_id")
		killer, killerExist := isUidExist(userUID)
		if killerExist {
			processKill(killer, target)
			checkHandler(writer, request)
		} else {
			nosuchUserHandler(writer, request, "unique_id")
		}
	} else {
		nosuchUserHandler(writer, request, "killapply")
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/"))).Methods("GET")
	router.HandleFunc("/check", checkHandler).Methods("POST")
	router.HandleFunc("/kill", killHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
