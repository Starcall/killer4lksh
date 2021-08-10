package main

import (
    "github.com/gorilla/mux"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "os"
    "io"
)

type User struct {
    FirstName string `json:"firstname""`
    LastName string `json:"lastname""`
    UID string `json:"uid"`
}

var allUsers = []User {
    User {
         FirstName: "Sasha",
         LastName: "Maneev",
         UID: "1",
    },
    User {
         FirstName: "Bogdan",
         LastName: "Trubetskoy",
         UID: "12",
    },
}

func handler(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(writer, "<h1>Personal Place</h1>" +
        "write you unique ID" +
        "<form action=\"/check\" method=\"POST\">" +
        "<textarea name=\"body\"></textarea><br>" +
        "<input type=\"submit\" value=\"Login\">" +
        "</form>")
}

func getProfile(user User) string {
    ret := "<tr>"
    ret += "<th>" + user.FirstName + "</th>"
    ret += "<th>" + user.LastName + "</th>"
    ret += "<th>" + user.UID + "</th>"
    ret += "</tr>"
    return ret
}

func checkHandler(writer http.ResponseWriter, request *http.Request) {
    requestedUid := request.FormValue("body")
    for _, elem := range allUsers {
        if elem.UID == requestedUid {
            file, err := os.Open("./data/" + elem.UID + ".json")
            if err != nil {
                file, err = os.Create("./data/" + elem.UID + ".json")
                if err != nil {
                    writer.WriteHeader(401)
                    writer.Write([]byte("Something went wrong"))
                    log.Fatal(err)
                    return
                }
            }
            decoder := json.NewDecoder(file)
            writer.Write([]byte("<table style=\"width:100%\">"))
            writer.Write([]byte("<tr><th>Firstname</th><th>Lastname</th><th>UID</th></tr>"))
            for {
                var killed User
                if err := decoder.Decode(&killed); err == io.EOF {
                    break
                } else if err != nil {
                    log.Fatal(err)
                }
                writer.Write([]byte(getProfile(killed)))
            }
            writer.Write([]byte("</table>"))
            return
        }
    }
    writer.Write([]byte("UID is incorrect"))

}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", handler).Methods("GET")
    router.HandleFunc("/check", checkHandler).Methods("POST")
    log.Fatal(http.ListenAndServe(":8080", router))
}

