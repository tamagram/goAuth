package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/github"
)

func main() {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
		github.New(os.Getenv("GITHUB_CLIENT"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
	)

	m := make(map[string]string)
	m["google"] = "Google"
	m["github"] = "Github"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("../template/user.html")
		t.Execute(res, user)
	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			t, _ := template.ParseFiles("../template/user.html")
			t.Execute(res, gothUser)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request){
		t, _ := template.ParseFiles("../template/index.html")
		t.Execute(res, providerIndex)
	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

