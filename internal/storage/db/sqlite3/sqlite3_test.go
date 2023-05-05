package sqlite3

import (
	"context"
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"
	"url-shortener/internal/schema"
	prep "url-shortener/internal/storage/db/queries"
)

var TestDB *Sqlite3

const pathToMigrations = "file://../../../../migrations/sqlite3"

func TestMain(m *testing.M) {
	// Write code here to run before tests

	db, err := sql.Open("sqlite3", "test-db")
	if err != nil {
		log.Fatal(err)
	}

	TestDB = New(db, pathToMigrations).(*Sqlite3)

	err = prep.Prepare(db, "sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	// Run tests
	c := m.Run()
	db.Close()

	func() {
		err = os.Remove("test-db")
		if err != nil {
			log.Fatalf("Err temp file was not removed: %v", err)
		}
	}()

	os.Exit(c)
}

func Test_FindMaxID(t *testing.T) {
	ctx := context.Background()
	want := 0
	got, err := TestDB.FindMaxID(ctx)
	if got != want {
		t.Errorf("FindMaxID() got = %v, want %v", got, want)
	} else if err != nil {
		t.Error(err)
	}

	_, err = TestDB.AddLink(ctx, "dqwdqwd", "qhwdfhqfh", "hqfhvqhv")
	if err != nil {
		t.Error(err)
	}

	want = 1
	got, err = TestDB.FindMaxID(ctx)
	if got != want {
		t.Errorf("FindMaxID() got = %v, want %v", got, want)
	} else if err != nil {
		t.Error(err)
	}

}

func Test_AddLink(t *testing.T) {
	ctx := context.Background()
	type args struct {
		longURL  string
		shortURL string
		cookie   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				longURL:  "google.com/",
				shortURL: "g.com/dqw",
				cookie:   "qwdqgfqedq",
			},
			want: "g.com/dqw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TestDB.AddLink(ctx, tt.args.longURL, tt.args.shortURL, tt.args.cookie)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetAllLinksByCookie(t *testing.T) {
	ctx := context.Background()
	_, err := TestDB.AddLink(ctx, "dqw3dqwd", "q3hwdfhqfh", "3hqfhvqhv")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		cookie  string
		baseURL string
	}
	tests := []struct {
		name    string
		args    args
		want    []schema.URL
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				cookie:  "3hqfhvqhv",
				baseURL: "127.0.0.1/",
			},
			want: []schema.URL{
				{
					LongURL:  "dqw3dqwd",
					ShortURL: "127.0.0.1/q3hwdfhqfh",
				},
			},
		},
		{
			name: "Empty",
			args: args{
				cookie:  "3hqf3hvqhv",
				baseURL: "127.0.0.1/",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TestDB.GetAllLinksByCookie(ctx, tt.args.cookie, tt.args.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllLinksByCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllLinksByCookie() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetLongLink(t *testing.T) {
	longLink := "dqwdqq"
	shortLink := "f"
	ctx := context.Background()

	_, err := TestDB.AddLink(ctx, longLink, shortLink, "wd")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		shortURL string
	}
	tests := []struct {
		name        string
		args        args
		wantLongURL string
		wantErr     bool
	}{
		{
			name:        "Ok",
			args:        args{shortURL: shortLink},
			wantLongURL: longLink,
		},
		{
			name:    "Not found",
			args:    args{shortURL: "f1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLongURL, err := TestDB.GetLongLink(ctx, tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLongLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLongURL != tt.wantLongURL {
				t.Errorf("GetLongLink() gotLongURL = %v, want %v", gotLongURL, tt.wantLongURL)
			}
		})
	}
}

func Test_MarkAsDeleted(t *testing.T) {
	ShortURL := "qwe"
	cookie := "qwsa"
	TestDB.MarkAsDeleted(ShortURL, cookie)
	ctx := context.Background()

	_, err := TestDB.GetLongLink(ctx, ShortURL)
	if err == nil {
		t.Error("The MarkAsDeleted() job was not completed")
		return
	}
}

func Test_Ping(t *testing.T) {
	ctx := context.Background()
	if err := TestDB.Ping(ctx); err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}
