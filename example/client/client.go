package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// Authorization Code Grant
	http.HandleFunc("/test/authorization_code", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/authorize?client_id=1&redirect_uri=http://localhost:8080/test/code_to_token&response_type=code", 302)
	})

	http.HandleFunc("/test/code_to_token", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")

		log.Println("code:" + code)

		url := fmt.Sprintf("http://localhost:8081/token?client_id=1&client_secret=1he5k5ZUrHFjznxN&grant_type=authorization_code&code=%s", code)

		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}

		body, _ := ioutil.ReadAll(res.Body)

		log.Println(fmt.Sprintf("POST url: %s repsonse %s", url, body))

		w.Write(body)
	})

	// Implicit Grant
	http.HandleFunc("/test/implicit", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/authorize?client_id=1&redirect_uri=http://localhost:8080/test/implicit_token&response_type=token", 302)
	})

	http.HandleFunc("/test/implicit_token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("access_token in uri (only for frontend)"))
	})

	// Password Credentials Grant
	http.HandleFunc("/test/password", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8081/token?client_id=1&client_secret=1he5k5ZUrHFjznxN&grant_type=password&username=admin&password=123456"

		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}

		body, _ := ioutil.ReadAll(res.Body)

		log.Println(fmt.Sprintf("POST url: %s repsonse %s", url, body))

		w.Write(body)
	})

	// Client Credentials Grant
	http.HandleFunc("/test/client_credentials", func(w http.ResponseWriter, r *http.Request) {
		url := "http://localhost:8081/token?client_id=1&client_secret=1he5k5ZUrHFjznxN&grant_type=client_credentials"

		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}

		body, _ := ioutil.ReadAll(res.Body)

		log.Println(fmt.Sprintf("POST url: %s repsonse %s", url, body))

		w.Write(body)
	})

	// Refreshing an access token
	http.HandleFunc("/test/refresh_token", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		refresh_token := query.Get("refresh_token")
		if refresh_token == "" {
			w.Write([]byte("missing params:refresh_token"))
			return
		}
		url := fmt.Sprintf("http://localhost:8081/token?client_id=1&client_secret=1he5k5ZUrHFjznxN&grant_type=refresh_token&refresh_token=%s", refresh_token)

		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}

		body, _ := ioutil.ReadAll(res.Body)

		log.Println(fmt.Sprintf("POST url: %s repsonse %s", url, body))

		w.Write(body)
	})

	log.Println("Start Listen :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
