package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lavinas-science/learn-microservices/handlers"
)

func main() {

	// create log
	l := log.New(os.Stdout, "product-api-", log.LstdFlags)

	// create handlers
	ph := handlers.NewProducts(l)

	//create a new serve mux and register the handle
	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", ph.GetProducts)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/", ph.AddProduct)
	postR.Use(ph.MiddlewareProductValidation)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putR.Use(ph.MiddlewareProductValidation)

	// Run the server in a go routine with parameters
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Finishing Gracefully
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, cc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cc()
	s.Shutdown(tc)

}
