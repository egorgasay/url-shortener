package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"strings"
)

type urls struct {
	longName  string
	shortName string
	id        int
}

const (
	alphabet    string = "AB1CDEFG2HIJKLM3NOPQRS4TUVW5XYZabc6defgh7ijklmn8opqrs9tuvw0xyz"
	lenAlphabet int    = 62
	domain      string = "http://127.0.0.1:8080/"
)

// lastIdentificator канал в котором хранится последнее известное id
var lastIdentificator = make(chan int, 1)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	shrt := r.URL.Path
	fmt.Println(shrt)
	stm := DB.QueryRow("SELECT long, id FROM urls WHERE short = ?", shrt[1:])
	u := urls{}
	err := stm.Scan(&u.longName, &u.id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	select {
	case lastIdentificator <- u.id:
	default:
		<-lastIdentificator
		lastIdentificator <- u.id
	}
	//http.Redirect(w, r, u.longName, http.StatusTemporaryRedirect)
	w.Header().Set("Location", "")
	//w.Header().Add("Location", u.longName)
	w.Header().Add("Location", u.longName)
	w.WriteHeader(307)
	w.Header().Del("Location")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	longURL := string(b)
	li := 0
	select {
	case li = <-lastIdentificator:
		//fmt.Println("Забрал старый", li)
		lastIdentificator <- li + 1
		break
	default:
		stm := DB.QueryRow("SELECT MAX(id) FROM urls")
		err := stm.Scan(&li)
		if err != nil {
			li++
		}
		//fmt.Println("Забрать старый не получилось", li)
		lastIdentificator <- li + 1 // код дублируется - исправить
		break
	}
	shrtURL := getShortName(li)
	valueStrings := fmt.Sprintf("('%s','%s')", longURL, shrtURL)
	stmt := fmt.Sprintf("INSERT INTO urls (long, short) VALUES %s", valueStrings)
	_, err = DB.Exec(stmt, valueStrings)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	_, err = w.Write([]byte(domain + shrtURL))
	if err != nil {
		return
	}
}

func getShortName(lastID int) (shrtURL string) {
	allNums := []int{}
	if lastID < 100000 {
		lastID = 10000 * lastID
	}
	for lastID > 0 {
		allNums = append(allNums, lastID%lenAlphabet)
		lastID = lastID / lenAlphabet
	}
	//fmt.Println(allNums)
	// разворачиваем слайс
	for i, j := 0, len(allNums)-1; i < j; i, j = i+1, j-1 {
		allNums[i], allNums[j] = allNums[j], allNums[i]
	}

	chars := []string{}
	for _, el := range allNums {
		chars = append(chars, string(alphabet[el]))
	}
	shrtURL = strings.Join(chars, "")
	return
}

var DB *sql.DB

func main() {
	db, errWhileOpenDB := sql.Open("sqlite3", "urlshortener.db")
	DB = db
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			return
		}
	}(DB)
	if errWhileOpenDB != nil {
		log.Fatal(errWhileOpenDB)
	}
	router := mux.NewRouter()
	router.HandleFunc("/{id}", GetHandler)
	router.HandleFunc("/", PostHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
