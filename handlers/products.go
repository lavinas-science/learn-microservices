package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/lavinas-science/learn-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marchal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	// Getting var id
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// Checking product on db and updating
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	// Update
	err := data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	// return handle funcrtion
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarchal json", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
