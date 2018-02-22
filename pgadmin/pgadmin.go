package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
	Fullname string        `json:"full_name"`
	Salt     string        `json:"-"`
	Password string        `json:"-"`
	Locked   bool          `json:"locked"`
	Created  time.Time     `json:"created"`
	Modified time.Time     `json:"modified"`
	Avatar   sql.NullInt64 `json:"avatar"`
}

type userHandler struct {
	db *sql.DB
}

var outArr []User

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

func main() {

	var u userHandler
	u.conndb()

	router := mux.NewRouter()
	router.HandleFunc("/user", u.GetAllUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", u.GetUserByID).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

	defer u.db.Close()
}

func (u *userHandler) conndb() {
	psqlInfo := fmt.Sprintf("host = %s  port = %d  user = %s  password = %s  dbname = %s  sslmode = %s", host, port, userdb, password, dbname, sslmode)

	var err error
	u.db, err = sql.Open("postgres", psqlInfo)
	checkErr(err)

	err = u.db.Ping()
	checkErr(err)

	fmt.Println("Succesfully Connected")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
