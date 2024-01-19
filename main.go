package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "db_act5" // ubah dbName dengan nama database yg kalian pengen
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(localhost)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func executeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close()

	var emp Employee
	var res []Employee
	for selDB.Next() {
		var id int
		var name, city string
		err := selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	executeTemplate(w, "Index", res)
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close()

	var emp Employee
	for selDB.Next() {
		var id int
		var name, city string
		err := selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	executeTemplate(w, "Show", emp)
}

func New(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close()

	var emp Employee
	for selDB.Next() {
		var id int
		var name, city string
		err := selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	executeTemplate(w, "Edit", emp)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		defer insForm.Close()

		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		defer insForm.Close()

		insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer delForm.Close()

	delForm.Exec(emp)
	log.Println("DELETE")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	log.Println("Server started on: http://localhost:8068") // ganti port jadi npm kalian
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8068", nil) //ganti port jadi npm kalian
}
