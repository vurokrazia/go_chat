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
	WebSocket *websocket.Conn
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

func CreateUser(user_name string, ws *websocket.Conn) User {
	return User{user_name, ws}
}

func AddUser(user User){
	Users.Lock()
	defer Users.Unlock()
	Users.m[user.User_Name] = user 
}

func RemoveUser(user_name string)  {
	Users.Lock()
	defer Users.Unlock()
	log.Println("Bye user")
	delete(Users.m, user_name)
}

func SendMessage(type_message int, message [] byte)  {
	Users.RLock()
	defer Users.RUnlock()
	for _, user := range Users.m{
		user.WebSocket.WriteMessage(type_message, message)
		if err := user.WebSocket.WriteMessage(type_message, message); err != nil {
			log.Println("Error message")
			return
		} else {
			log.Println(user.User_Name)
			log.Println(string(message))
		}
	}
}

func ToArrayByte(value string) [] byte {
	return []byte(value)
}

func ConcatMessage(user_name string, array []byte) string {
	return user_name + ":" + string(array)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_name := vars["user_name"]

	ws, err := websocket.Upgrade(w,r,nil,1024,1024)
	if err != nil {
		log.Println("Error")
	}

	current_user := CreateUser(user_name,ws)
	AddUser(current_user)
	log.Println("New user registered")
	for{
		type_message, message, err := ws.ReadMessage()
		if err != nil {
			RemoveUser(user_name)
			return
		}
		final_message := ConcatMessage(user_name, message)
		SendMessage(type_message, ToArrayByte(final_message))
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