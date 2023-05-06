package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/config"
	grpchandler "url-shortener/internal/handler/grpc"
	resthandler "url-shortener/internal/handler/rest"
	"url-shortener/internal/repository"
	"url-shortener/internal/routes"
	"url-shortener/internal/storage/db/queries"
	"url-shortener/internal/usecase"
	shortener "url-shortener/pkg/api"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

const startText = `
Build version: %s
Build date: %s
Build commit: %s
`

func main() {
	fmt.Printf(startText, buildVersion, buildDate, buildCommit)
	cfg := config.New()

	storage, err := repository.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to initialize: %s", err.Error())
	}

	logic := usecase.New(storage)
	router := gin.Default()
	h := resthandler.NewHandler(cfg, logic)

	public := router.Group("/")
	routes.PublicRoutes(public, h)

	router.Use(gzip.Gzip(gzip.BestSpeed))

	if cfg.GRPC != "" {
		go func() {
			log.Println("Server is running on grpc://" + cfg.GRPC)
			grpcServer := grpc.NewServer()
			ghandler := grpchandler.NewHandler(cfg, logic)

			lis, err := net.Listen("tcp", cfg.Host)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			shortener.RegisterShortenerServer(grpcServer, ghandler)

			err = grpcServer.Serve(lis)
			if err != nil {
				log.Fatalf("grpcServer Serve: %v", err)
			}
		}()
	}

	go func() {
		if cfg.HTTPS {
			log.Println("Server is running on https://" + cfg.Host)
			ca := &x509.Certificate{
				SerialNumber: big.NewInt(2019),
				Subject: pkix.Name{
					Organization: []string{"GasaySecure, INC."},
					Country:      []string{"US"},
					Locality:     []string{"New York"},
				},
				NotBefore:             time.Now(),
				NotAfter:              time.Now().AddDate(10, 0, 0),
				IsCA:                  true,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
				KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				BasicConstraintsValid: true,
			}
			caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
			if err != nil {
				log.Fatalf("Failed to generate private key: %s", err.Error())
			}

			caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
			if err != nil {
				log.Fatalf("Failed to create certificate: %s", err.Error())
			}

			caPEM := new(bytes.Buffer)
			pem.Encode(caPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: caBytes,
			})

			caPrivKeyPEM := new(bytes.Buffer)
			pem.Encode(caPrivKeyPEM, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
			})

			caFile, err := os.Create("ca.crt")
			if err != nil {
				log.Fatalf("Failed to create file: %s", err.Error())
			}
			_, err = caFile.Write(caPEM.Bytes())
			if err != nil {
				log.Fatalf("Failed to write file: %s", err.Error())
			}

			defer caFile.Close()

			caFile, err = os.Create("ca.key")
			if err != nil {
				log.Fatalf("Failed to create file: %s", err.Error())
			}

			_, err = caFile.Write(caPrivKeyPEM.Bytes())
			if err != nil {
				log.Fatalf("Failed to write file: %s", err.Error())
			}

			defer caFile.Close()

			http.ListenAndServeTLS(cfg.Host, "ca.crt", "ca.key", router)
		} else {
			log.Println("Server is running on http://" + cfg.Host)
			router.Run(cfg.Host)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")

	err = storage.Shutdown()
	if err != nil {
		log.Println("Failed to shutdown storage: ", err)
	}

	err = queries.Close()
	if err != nil {
		log.Println("Failed to shutdown storage: ", err)
	}

}
