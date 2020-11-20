package main

import (
	"flag"
	"github.com/mdapathy/url-shortener/database"
	server2 "github.com/mdapathy/url-shortener/server"
	"github.com/mdapathy/url-shortener/tools"
	"github.com/mdapathy/url-shortener/url/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	configFile     = flag.String("configFile", "config.json", "Config file")
	httpPortNumber = flag.Int("p", 8080, "HTTP port number")
)

func main() {
	db := database.NewDBConfig(*configFile)
	defer db.Close()

	mid := middleware.New(db, tools.NewCacheStorage())

	server := &server2.ApiServer{
		Port:       *httpPortNumber,
		Controller: mid.NewController(),
	}

	go func() {

		log.Println("Starting the server...")

		err := server.Start()
		if err == http.ErrServerClosed {
			log.Printf("HTTP server stopped")
		} else {
			log.Fatalf("Cannot start HTTP server: %s", err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<-sigChannel

	if err := server.Stop(); err != nil && err != http.ErrServerClosed {
		log.Printf("Error stopping the server: %s", err)
	}
}
