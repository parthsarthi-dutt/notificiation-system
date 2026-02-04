package main

import (
	"encoding/json"
	"fmt"
	redisclient "internals/redis"
	api "internals/api"
	"log"
	"net/http"
	"strconv"
	"sync"
	
)

type User struct{
	Name string `json:"name"`
}

var userCache=make(map[int]User)

var cacheMutex sync.RWMutex 



func main() {

log.Println("Starting service...")

	redisclient.InitRedis()

	mux := http.NewServeMux()
	mux.HandleFunc("/set", api.SetHandler)
	mux.HandleFunc("/get", api.GetHandler)

	log.Println("HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	mux:= http.NewServeMux()
// 	mux.HandleFunc("/",HandleRoot)
// 	mux.HandleFunc("POST /users",CreateUser)
// 	mux.HandleFunc("GET /users/{id}",GetUser)
// 	mux.HandleFunc("DELETE /users/{id}",DeleteUser)
// 	fmt.Println("Listening to server at port :8080")
// 	http.ListenAndServe(":8080",mux)
// }
func HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello world!")
}

func CreateUser(w http.ResponseWriter, r *http.Request){
	var user User
	err:=json.NewDecoder(r.Body).Decode(&user)
	if(err!=nil){
		http.Error(
			w,
			err.Error(), 
			http.StatusBadRequest,
		)
		return 
	}
	if user.Name==""{
			http.Error(
			w,
			"name is required", 
			http.StatusBadRequest,
		)
		return 
	}
	cacheMutex.Lock()
	userCache[len(userCache)+1]=user
	cacheMutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)

json.NewEncoder(w).Encode(map[string]string{
	"message": "User created",
})
}

func GetUser(w http.ResponseWriter, r *http.Request){
	id,err:=strconv.Atoi(r.PathValue("id"))
	if(err!=nil){
		http.Error(
			w,
			err.Error(), 
			http.StatusBadRequest,
		)
		return 
	}
	cacheMutex.RLock()
	user,ok:=userCache[id]
	cacheMutex.RUnlock()
	if !ok{
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return 
	}
	w.Header().Set("Content-Type","application/json")
	j,err:=json.Marshal(user)
	if(err!=nil){
		http.Error(
			w,
			err.Error(), 
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func DeleteUser(w http.ResponseWriter, r *http.Request){
	id,err:=strconv.Atoi(r.PathValue("id"))
	if(err!=nil){
		http.Error(
			w,
			err.Error(), 
			http.StatusBadRequest,
		)
		return 
	}
	cacheMutex.RLock()
	_,ok:=userCache[id]
	cacheMutex.RUnlock()
	if !ok{
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return 
	}
	cacheMutex.RLock()
	delete(userCache,id)
	cacheMutex.RUnlock()

	w.WriteHeader(http.StatusNoContent)


}