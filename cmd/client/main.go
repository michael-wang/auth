package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	loginGithubEndpoint = "login/github"
)

func main() {
	flagHost := flag.String("host", "localhost", "host name of the server")
	flagPort := flag.Int("port", 8080, "port number of the server")
	flag.Parse()

	url := fmt.Sprintf("http://%s:%d/%s", *flagHost, *flagPort, loginGithubEndpoint)
	fmt.Printf("HTTP GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response headers:\n%v\n", resp.Header)
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	fmt.Printf("Response body:\n%s\n", bb)
}
