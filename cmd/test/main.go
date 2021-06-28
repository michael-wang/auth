package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michael-wang/auth"
)

func main() {
	go getUserProfile()

	startHttpServer(":8080")
}

func getUserProfile() {
	a := auth.New("faf9d0195aed6907ee12", "2e1e0b8b51d9dcb8c8ef9c7def16383632f20b23", "http://localhost:8080/auth/callback")
	url := a.GetLoginURL()
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		panic(err)
	}

	user, _ := a.GetUserProfile(code)
	bb, _ := json.MarshalIndent(user, "", "\t")
	fmt.Println(string(bb))
	fmt.Println("test done, you can exit by 'Ctrl+C'")
}

func startHttpServer(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/auth/callback", http.HandlerFunc(auth.CallbackHandler))
	fmt.Println("http server listen on", addr)
	http.ListenAndServe(addr, mux)
}
