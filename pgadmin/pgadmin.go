package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"io/ioutil"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "192.168.1.203"
	port     = 5432
	userdb   = "mas"
	password = "12345678"
	dbname   = "vanda"
	sslmode  = "disable"
)

type User struct {
	Userid   int           `json:"user_id"`
	Tenantid int64         `json:"tenant_id"`
	Email    string        `json:"email"`
	Fullname string        `json:"fullname"`
	Salt     string        `json:"salt"`
	Password string        `json:"password"`
	Locked   bool          `json:"locked"`
	Created  time.Time     `json:"created"`
	Modified time.Time     `json:"modified"`
	Avatar   sql.NullInt64 `json:"avatar"`
}

type userHandler struct {
	db *sql.DB
}

var outArr []User

//Show All User
func (u *userHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {

	rows, err := u.db.Query("SELECT * FROM account.user")
	checkErr(err)

	var a User

	for rows.Next() {
		err = rows.Scan(&a.Userid, &a.Tenantid, &a.Email, &a.Fullname, &a.Salt, &a.Password, &a.Locked, &a.Created, &a.Modified, &a.Avatar)
		checkErr(err)
		outArr = append(outArr, a)
	}
	json.NewEncoder(w).Encode(outArr)
}

//Get User by ID
func (u *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	uID, err := strconv.Atoi(params["user_id"])
	checkErr(err)

	rows, err := u.db.Query("SELECT * FROM account.user WHERE user_id=$1", uID)
	checkErr(err)

	var a User
	var check []User

	for rows.Next() {
		err = rows.Scan(&a.Userid, &a.Tenantid, &a.Email, &a.Fullname, &a.Salt, &a.Password, &a.Locked, &a.Created, &a.Modified, &a.Avatar)
		checkErr(err)
		check = append(check, a)
	}
	json.NewEncoder(w).Encode(check)
}

//Insert User
func (u *userHandler) InsertUser(w http.ResponseWriter, r *http.Request) {

	var a User

	a = UnwrapJson(r)

	insert, err := u.db.Query("INSERT INTO account.user VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", a.Userid, a.Tenantid, a.Email, a.Fullname, a.Salt, a.Password, a.Locked, a.Created, a.Modified, a.Avatar)
	checkErr(err)
	fmt.Println(insert)
	
}

//Update User by ID
func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	uID, err := strconv.Atoi(params["user_id"])
	checkErr(err)

	var a User
	
	a = UnwrapJson(r)

	update, err := u.db.Query("UPDATE account.user SET fullname=$1 WHERE user_id=$2", a.Fullname, uID)
	checkErr(err)
	fmt.Print(update)
}

//Delete User by ID
func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	uID, err := strconv.Atoi(params["user_id"])

	deluser, err := u.db.Query("DELETE FROM account.user WHERE user_id=$1", uID)
	checkErr(err)

	json.NewEncoder(w).Encode(deluser)
}

func UnwrapJson(r *http.Request) (a User){
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	
	data := json.Unmarshal(body, &a)
	if data != nil {
		fmt.Println("error:", data)
	}
	fmt.Println(a)
	return a
}

func main() {

	db := conndb()
	u := &userHandler{db}

	router := mux.NewRouter()
	router.HandleFunc("/user", u.GetAllUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", u.GetUserByID).Methods("GET")
	router.HandleFunc("/user", u.InsertUser).Methods("POST")
	router.HandleFunc("/user/{user_id}", u.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/{user_id}", u.UpdateUser).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", router))

	defer u.db.Close()
}

func conndb() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host = %s  port = %d  user = %s  password = %s  dbname = %s  sslmode = %s", host, port, userdb, password, dbname, sslmode)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	fmt.Println("Succesfully Connected")
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
