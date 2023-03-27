package mapstorage

import (
	"os"
	"reflect"
	"testing"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
)

// TestDB database for testing.
var TestDB storage.IStorage

// TestMain ...
func TestMain(m *testing.M) {
	// Write code here to run before tests
	TestDB = NewMapStorage()
	// Run tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

// TestPostgres_FindMaxID ...
func TestPostgres_FindMaxID(t *testing.T) {
	want := 0
	got, err := TestDB.FindMaxID()
	if got != want {
		t.Errorf("FindMaxID() got = %v, want %v", got, want)
	} else if err != nil {
		t.Error(err)
	}

	_, err = TestDB.AddLink("dqwdqwd", "qhwdfhqfh", "hqfhvqhv")
	if err != nil {
		t.Error(err)
	}

	want = 1
	got, err = TestDB.FindMaxID()
	if got != want {
		t.Errorf("FindMaxID() got = %v, want %v", got, want)
	} else if err != nil {
		t.Error(err)
	}
}

// TestPostgres_AddLink ...
func TestPostgres_AddLink(t *testing.T) {
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
		{
			name: "Already Exists",
			args: args{
				longURL:  "google.com/",
				shortURL: "g.com/dqw",
				cookie:   "qwdqgfqedq",
			},
			wantErr: true,
			want:    "g.com/dqw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TestDB.AddLink(tt.args.longURL, tt.args.shortURL, tt.args.cookie)
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

// TestPostgres_GetAllLinksByCookie ...
func TestPostgres_GetAllLinksByCookie(t *testing.T) {
	_, err := TestDB.AddLink("dqw3dqwd", "q3hwdfhqfh", "3hqfhvqhv")
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
			got, err := TestDB.GetAllLinksByCookie(tt.args.cookie, tt.args.baseURL)
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

// TestPostgres_GetLongLink ...
func TestPostgres_GetLongLink(t *testing.T) {
	_, err := TestDB.AddLink("dqwdqq", "f", "wd")
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
			gotLongURL, err := TestDB.GetLongLink(tt.args.shortURL)
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

// TestPostgres_MarkAsDeleted ...
func TestPostgres_MarkAsDeleted(t *testing.T) {
	ShortURL := "qwe"
	cookie := "qwsa"
	TestDB.MarkAsDeleted(ShortURL, cookie)

	_, err := TestDB.GetLongLink(ShortURL)
	if err == nil {
		t.Error("The MarkAsDeleted() job was not completed")
		return
	}
}

// TestPostgres_Ping ...
func TestPostgres_Ping(t *testing.T) {
	if err := TestDB.Ping(); err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}
