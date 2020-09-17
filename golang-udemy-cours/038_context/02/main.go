package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ctx = context.WithValue(ctx, "userID", 777)
	ctx = context.WithValue(ctx, "fname", "Bond")

	result := dbAccess(ctx)
	fmt.Fprintln(w, result)

	result2 := getName(ctx)
	fmt.Fprintln(w, result2)
}

func getName(ctx context.Context) string {
	uName := ctx.Value("fname").(string)
	return uName
}

func dbAccess(ctx context.Context) int {
	uid := ctx.Value("userID").(int)
	return uid
}

func bar(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println(ctx)
	fmt.Fprintln(w, ctx)
}
