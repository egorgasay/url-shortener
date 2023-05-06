package grpchandler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/usecase"
	shortener "url-shortener/pkg/api"
)

func TestHandler_Create(t *testing.T) {
	cfg := config.Config{Key: []byte("test-key"), DBConfig: &repository.Config{DriverName: "map"}, Host: ":789",
		BaseURL: "http://localhost:789/"}
	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		t.Fatal(err)
	}

	uc := usecase.New(storage)

	h := NewHandler(&cfg, uc)

	grpcServer := grpc.NewServer()

	log.Println("Starting Shortener ...")
	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	shortener.RegisterShortenerServer(grpcServer, h)

	go func() {
		log.Println("Starting GRPC", cfg.Host)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	conn, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	var header metadata.MD // variable to store header and trailer

	cl := shortener.NewShortenerClient(conn)
	token := "529967c34009a2fc523e597248c32dbf95166c3a49092d4223cfb60d4a1324cd-31363833333634373635313831383032363030"

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	badCtx := context.Background()

	type args struct {
		ctx context.Context
		req *shortener.CreateRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *shortener.CreateResponse
		wantToken bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: &shortener.CreateRequest{
					Url: "https://ya.ru",
				},
			},
			want: &shortener.CreateResponse{
				Shortened: "http://localhost:789/zE",
			},
		},
		{
			name: "no token",
			args: args{
				ctx: badCtx,
				req: &shortener.CreateRequest{
					Url: "https://ya.ru",
				},
			},
			want: &shortener.CreateResponse{
				Shortened: "http://localhost:789/Xz",
			},
		},
		{
			name: "no token #2",
			args: args{
				ctx: badCtx,
				req: &shortener.CreateRequest{
					Url: "https://ya.ru",
				},
			},
			want: &shortener.CreateResponse{
				Shortened: "http://localhost:789/rx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cl.Create(tt.args.ctx, tt.args.req, grpc.Header(&header))
			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}

			if got.Shortened != tt.want.Shortened {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if tt.wantToken && len(header.Get("token")) == 0 {
				t.Errorf("Create() header = %v, want %v", header, "token")
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	cfg := config.Config{Key: []byte("test-key"), DBConfig: &repository.Config{DriverName: "map"}, Host: ":788",
		BaseURL: "http://localhost:788/"}
	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		t.Fatal(err)
	}

	uc := usecase.New(storage)

	h := NewHandler(&cfg, uc)

	grpcServer := grpc.NewServer()

	log.Println("Starting Shortener ...")
	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	shortener.RegisterShortenerServer(grpcServer, h)

	go func() {
		log.Println("Starting GRPC", cfg.Host)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	conn, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	cl := shortener.NewShortenerClient(conn)

	token := "529967c34009a2fc523e597248c32dbf95166c3a49092d4223cfb60d4a1324cd-31363833333634373635313831383032363030"

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// SUCCESS

	// PREPARE

	_, err = cl.Create(ctx, &shortener.CreateRequest{Url: "http://ya.ru"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// TEST

	get, err := cl.Get(ctx, &shortener.GetRequest{Shortened: "zE"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if get.OriginalUrl != "http://ya.ru" {
		t.Errorf("Get() got = %v, want %v", get, "http://ya.ru")
	}

	// DELETED

	// PREPARE

	_, err = cl.Delete(ctx, &shortener.DeleteRequest{ShortenedUrls: []string{"zE"}})
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// TEST

	get, err = cl.Get(ctx, &shortener.GetRequest{Shortened: "zE"})
	if err == nil {
		t.Fatal("No error was returned (deleted url)")
	}

	if status.Code(err) != codes.Unavailable {
		t.Fatal("Wrong error was returned (deleted url)")
	}

	// NOT FOUND

	// TEST

	get, err = cl.Get(ctx, &shortener.GetRequest{Shortened: "qwdfqdfqwsqwdqfqfew"})
	if err == nil {
		t.Fatal("No error was returned (deleted url)")
	}
}

func TestHandler_Delete(t *testing.T) {
	cfg := config.Config{Key: []byte("test-key"), DBConfig: &repository.Config{DriverName: "map"}, Host: ":788",
		BaseURL: "http://localhost:787/"}
	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		t.Fatal(err)
	}

	uc := usecase.New(storage)

	h := NewHandler(&cfg, uc)

	grpcServer := grpc.NewServer()

	log.Println("Starting Shortener ...")
	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	shortener.RegisterShortenerServer(grpcServer, h)

	go func() {
		log.Println("Starting GRPC", cfg.Host)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}

	}()

	conn, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	cl := shortener.NewShortenerClient(conn)

	token := "529967c34009a2fc523e597248c32dbf95166c3a49092d4223cfb60d4a1324cd-31363833333634373635313831383032363030"

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// SUCCESS

	// PREPARE

	_, err = cl.Create(ctx, &shortener.CreateRequest{Url: "http://ya.ru"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// TEST

	get, err := cl.Get(ctx, &shortener.GetRequest{Shortened: "zE"})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if get.OriginalUrl != "http://ya.ru" {
		t.Errorf("Get() got = %v, want %v", get, "http://ya.ru")
	}

	// DELETED

	// PREPARE

	_, err = cl.Delete(ctx, &shortener.DeleteRequest{ShortenedUrls: []string{"zE"}})
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// TEST

	get, err = cl.Get(ctx, &shortener.GetRequest{Shortened: "zE"})
	if err == nil {
		t.Fatal("No error was returned (deleted url)")
	}

	if status.Code(err) != codes.Unavailable {
		t.Fatal("Wrong error was returned (deleted url)")
	}

	// NOT FOUND

	// TEST

	get, err = cl.Get(ctx, &shortener.GetRequest{Shortened: "qwdfqdfqwsqwdqfqfew"})
	if err == nil {
		t.Fatal("No error was returned (deleted url)")
	}
}
