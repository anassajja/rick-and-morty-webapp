package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Users struct {
	Users []User `json:"results"`
}

type User struct {
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Species  string   `json:"species"`
	Gender   string   `json:"gender"`
	Origin   Origin   `json:"origin"`
	Location Location `json:"location"`
	Image    string   `json:"image"`
}

type Origin struct {
	Origin string `json:"name"`
}

type Location struct {
	Location string `json:"name"`
}

func getCharacters(page int) (*Users, error) {
	url := fmt.Sprintf("https://rickandmortyapi.com/api/character?page=%d", page)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var results Users
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return &results, nil
}

func getAllCharacters() ([]User, error) {
	var allUsers []User
	for i := 1; i <= 42; i++ {
		users, err := getCharacters(i)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, users.Users...)
	}
	return allUsers, nil
}

func getRandomCharacters(numCharacters int) ([]User, error) {
	allUsers, err := getAllCharacters()
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().Unix())

	rand.Shuffle(len(allUsers), func(i, j int) {
		allUsers[i], allUsers[j] = allUsers[j], allUsers[i]
	})

	return allUsers[:numCharacters], nil
}

func serveWebsite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.New("siteweb.html").ParseFiles("siteweb.html")
	if err != nil {
		log.Fatal(err)
	}

	users, err := getRandomCharacters(826)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, users)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Serve static files from the "image" directory
	r.PathPrefix("/image/").Handler(http.StripPrefix("/image", http.FileServer(http.Dir("./image"))))

	// Define the route for the main page
	r.HandleFunc("/", serveWebsite).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
