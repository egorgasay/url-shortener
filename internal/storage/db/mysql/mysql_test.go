package mysql

import (
	"context"
	"github.com/egorgasay/dockerdb"
	"log"
	"os"
	"reflect"
	"testing"
	prep "url-shortener/internal/storage/db/queries"
	shortener "url-shortener/pkg/api"
)

var TestDB *MySQL

const pathToMigrations = "file://../../../../migrations/mysql"

func TestMain(m *testing.M) {
	// Write code here to run before tests
	ctx := context.TODO()
	cfg := dockerdb.CustomDB{
		DB: dockerdb.DB{
			Name:     "admin",
			User:     "admin",
			Password: "admin",
		},
		Port: "31188",
		Vendor: dockerdb.Vendor{
			Name:  "mysql",
			Image: "mysql:5.7",
		},
	}

	err := dockerdb.Pull(ctx, "mysql:5.7")
	if err != nil {
		log.Fatal(err)
		return
	}

	vdb, err := dockerdb.New(ctx, cfg)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	TestDB = New(vdb.DB, pathToMigrations).(*MySQL)

	queries := []string{
		"SET foreign_key_checks = 0;",
		"TRUNCATE links;",
		"SET foreign_key_checks = 1;",
	}

	tx, err := TestDB.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	TestDB = New(vdb.DB, pathToMigrations).(*MySQL)
	// Run tests

	err = prep.Prepare(TestDB.DB.DB, "mysql")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
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
		want    []*shortener.UserURL
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				cookie:  "3hqfhvqhv",
				baseURL: "127.0.0.1/",
			},
			want: []*shortener.UserURL{
				{
					OriginalUrl: "dqw3dqwd",
					ShortUrl:    "127.0.0.1/q3hwdfhqfh",
				},
			},
		},
		{
			name: "Empty",
			args: args{
				cookie:  "3hqf3hvqhv",
				baseURL: "127.0.0.1/",
			},
			want: []*shortener.UserURL{},
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
	ctx := context.Background()
	_, err := TestDB.AddLink(ctx, "dqwdqq", "f", "wd")
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
			args:        args{shortURL: "f"},
			wantLongURL: "dqwdqq",
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
	ctx := context.Background()
	ShortURL := "qwe"
	cookie := "qwsa"
	TestDB.MarkAsDeleted(ShortURL, cookie)

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
