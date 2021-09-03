package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lavinas-science/learn-microservices/handlers"
)

func main_old() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("Opsss"))
			http.Error(rw, "Ops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(rw, "Hello %s\n", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}

func main() {

	// Starting log
	l := log.New(os.Stdout, "product-api-", log.LstdFlags)

	// Handling - injecting log
	sm := http.NewServeMux()
	// sm.Handle("/", handlers.NewHello(l))
	// sm.Handle("/bye", handlers.NewGoodbye(l))
	sm.Handle("/", handlers.NewProducts(l))
	
	
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
