package storage

//type Storage *sql.DB
//
//type Repositories interface {
//	Close() error
//	Exec(query string, args ...any) (sql.Result, error)
//	QueryRow(query string, args ...any) *sql.Row
//}
//
//func NewStorage() Repositories {
//	db, errWhileOpenDB := sql.Open("sqlite3", "urlshortener.db")
//	if errWhileOpenDB != nil {
//		log.Fatal(errWhileOpenDB)
//	}
//	return db
//}
