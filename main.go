package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time" 

	_ "modernc.org/sqlite"
)

type Skill struct {
	Name string
	Icon string
}

type SocialLink struct {
	Platform string
	URL string
	Icon string
}

type Profile struct {
	Name string
	Title string
	Languages []Skill
	Databases []Skill
	Entries []Entry
	ContactEmail string
}

type Entry struct {
	Name string
	Message string
	Timestamp time.Time
}

var profile = Profile{
	Name: "Luch28",
	Title: "Frontend and Backend Developer",
	ContactEmail: "Luch28NET@gmail.com",
	Languages: []Skill{
		{"Go", ""},
		{"Python", ""},
		{"C#", ""},
		{"HTML", ""},
		{"CSS", ""},
		{"Bash Script", ""},
		{"JavaScript",""},
		{"Java", ""},
	},
	Databases: []Skill{
		{"MariaDB", ""},
		{"SQLite", ""},
		{"MySQL", ""},
	},
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("CREATE TABLE IF NOT EXISTS entries (name TEXT, message TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")

	fmt.Println("Done")
	handeReqest()
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	rows, err := db.Query("SELECT name, message, timestamp FROM entries ORDER BY timestamp DESC")
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Name, &e.Message, &e.Timestamp); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		entries = append(entries, e)
	}

	currentProfile := profile
	currentProfile.Entries = entries
	tmpl, err := template.ParseFiles("templates/index.html", "templates/head.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "index", currentProfile)
}


func guest(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name, message, timestamp FROM entries ORDER BY timestamp DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var e Entry
		rows.Scan(&e.Name, &e.Message, &e.Timestamp)
		entries = append(entries, e)
	}
	tmpl, err := template.ParseFiles("templates/guest.html", "templates/head.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	tmpl.ExecuteTemplate(w, "guest", entries)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		message := r.FormValue("message")
		if name != "" && message != "" {
			db.Exec("INSERT INTO entries (name, message) VALUES (?, ?)", name, message)
		}
	}
	http.Redirect(w, r, "guest", http.StatusSeeOther)
}

func handeReqest() {
	fsStatic := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fsStatic))
	fsImages := http.FileServer(http.Dir("images"))
    http.Handle("/images/", http.StripPrefix("/images/", fsImages))
	http.HandleFunc("/", index)
	http.HandleFunc("/guest", guest)
	http.HandleFunc("/submit", handleSubmit)
	log.Println("Running on http://localhost:28")
	log.Fatal(http.ListenAndServe(":28", nil))
}
