package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
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
	Userid   int			`json:"user_id"`
	Tenantid int64			`json:"tenant_id"`
	Email    string        	`json:"email"`
	Fullname string        	`json:"full_name"`
	Salt     string        	`json:"-"`
	Password string        	`json:"-"`
	Locked   bool          	`json:"locked"`
	Created  time.Time     	`json:"created"`
	Modified time.Time     	`json:"modified"`
	Avatar   sql.NullInt64 	`json:"avatar"`
}

var db *sql.DB
var outArr []User

func GetUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(outArr)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range outArr {
		s := strconv.Itoa(item.Userid)
        if s == params["user_id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(nil)
}

// func GetPerson(w http.ResponseWriter, r *http.Request) {
// 	nId := r.URL.Query().Get("user_id")
// 	selPerson, err := db.Query("SELECT * FROM account.user WHERE user_id=?", nId)
// 	checkErr(err)

	
// }

func CreatePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
	var person User
	_ = json.NewDecoder(r.Body).Decode(&person)
	userid,_ := strconv.Atoi(params["user_id"])
    person.Userid = userid
    outArr = append(outArr, person)
    json.NewEncoder(w).Encode(outArr)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range outArr {
		s := strconv.Itoa(item.Userid)
        if s == params["user_id"] {
            outArr = append(outArr[:index], outArr[index+1:]...)
            break
        }
        json.NewEncoder(w).Encode(outArr)
    }
}

func main() {

	psqlInfo := fmt.Sprintf("host = %s  port = %d  user = %s  password = %s  dbname = %s  sslmode = %s", host, port, userdb, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	fmt.Println("Succesfully Connected")

	rows, _ := db.Query("SELECT * FROM account.user")

	var a User

	for rows.Next() {
		if err = rows.Scan(&a.Userid, &a.Tenantid, &a.Email, &a.Fullname, &a.Salt, &a.Password, &a.Locked, &a.Created, &a.Modified, &a.Avatar); err != nil {
			fmt.Println("Scanning failed.....")
			fmt.Println(err.Error())
			return
		}
		outArr = append(outArr, a)
	}

	router := mux.NewRouter()
	router.HandleFunc("/user", GetUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", GetPerson).Methods("GET")
	router.HandleFunc("/user/{user_id}", CreatePerson).Methods("POST")
	router.HandleFunc("/user/{user_id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
