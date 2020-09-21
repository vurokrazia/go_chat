package main

import(
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"sync"
)

type Response struct {
	Message string  `json: "message"`
	Status  int     `json: "status"`
	IsValid bool    `json: "is_valid"`
}

type User struct {
	User_Name string
}

var Users = struct {
	m map[string] User
	sync.RWMutex
}{m: make(map[string] User)}

func UserExit(user_name string) bool {
	Users.RLock()
	defer Users.RUnlock()
	if _, ok := Users.m[user_name]; ok {
		return true
	}
	return false
}

func Hello(w http.ResponseWriter, r *http.Request)  {
	log.Println("hi")
	w.Write([]byte("Hello world from go"))
}

func Home(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Hello world from home"))
}

func LoadStatic(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "./front/index.html")
}

func HelloJson(w http.ResponseWriter, r *http.Request)  {
	response := CreateResponse("Mensaje",0,true)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, valid bool) Response {
	return Response{message,0,valid}
}

func Validate(w http.ResponseWriter, r *http.Request)  {

	r.ParseForm()
	user_name := r.FormValue("user_name")
	response := Response{}

	if UserExit(user_name) {
		response.IsValid = false
		// response := CreateResponse("Dont is valid",false)
	} else {
		response.IsValid = true
		// response := CreateResponse("Is valid",false)
	}
	json.NewEncoder(w).Encode(response)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w,r,nil,1024,1024)
	if err != nil {
		log.Println("Error")
	}
}

func main()  {

	cssHandle := http.FileServer(http.Dir("./front/css"))
	jsHandle := http.FileServer(http.Dir("./front/js"))

	mux := mux.NewRouter()
	mux.HandleFunc("/hello", Hello).Methods("GET")
	mux.HandleFunc("/hello_json", HelloJson).Methods("GET")
	mux.HandleFunc("/static", LoadStatic).Methods("GET")
	mux.HandleFunc("/validate", Validate).Methods("POST")
	mux.HandleFunc("/", LoadStatic).Methods("GET")
	mux.HandleFunc("/chat/{user_name}", WebSocket).Methods("GET")

	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))

	log.Println("Server success!! Listened port:8000")
	log.Fatal(http.ListenAndServe(":8000",nil))

}