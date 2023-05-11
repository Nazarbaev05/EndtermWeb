package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type Character struct {
	Id     int
	Fn     string
	Ln     string
	Age    int
	Height int
}

type Server struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "awerny.2003"
	dbname   = "bluelock"
)

func DBconnect() *Server {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		log.Fatal(err)
	}
	return &Server{db: db}
}

func (s *Server) charactersPage(w http.ResponseWriter, r *http.Request) {
	var chrs []Character
	res, _ := s.db.Query("select * from characters;")
	for res.Next() {
		var chr Character
		res.Scan(&chr.Id, &chr.Fn, &chr.Ln, &chr.Age, &chr.Height)
		chrs = append(chrs, chr)
	}
	t, _ := template.ParseFiles("static/html/Characters.html")
	t.Execute(w, chrs)
}

func (s *Server) addPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		age := r.FormValue("age")
		height := r.FormValue("height")
		if _, err := s.db.Exec("insert into characters(firstname, lastname, age, height) values($1, $2, $3, $4)", fn, ln, age, height); err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/characters", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/form.html")
	t.Execute(w, nil)
}

func (s *Server) deletePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		if _, err := s.db.Exec("delete from characters where id=$1", id); err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/characters", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/formDelete.html")
	t.Execute(w, nil)
}

func (s *Server) updatePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		age := r.FormValue("age")
		height := r.FormValue("height")
		if _, err := s.db.Exec("update characters set firstname=$1, lastname=$2, age=$3, height=$4 where id=$5", fn, ln, age, height, id); err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/characters", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/formUpdate.html")
	t.Execute(w, nil)
}

func main() {
	s := DBconnect()
	defer s.db.Close()
	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/", fileServer)
	http.HandleFunc("/characters", s.charactersPage)
	http.HandleFunc("/add", s.addPage)
	http.HandleFunc("/delete", s.deletePage)
	http.HandleFunc("/update", s.updatePage)
	http.ListenAndServe(":1234", nil)
}
