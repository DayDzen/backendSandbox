package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/user/:id", getUser)
	http.ListenAndServe("http://localhost:8080/", r)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func getUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{
		Name:   "Ilya",
		Gender: "male",
		Age:    "24",
		Id:     p.ByName("id"),
	}
}
