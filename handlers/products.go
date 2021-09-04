package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/lavinas-science/learn-microservices/data"
)

type Products struct {
	l *log.Logger
}

type KeyProduct struct{}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	// get all products
	lp := data.GetProducts()
	// to Json and show error if it not works
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marchal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	// get product from context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	// add product
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	// get var id
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// get product from context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	// update and show error if it not works
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

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	// return handle funcrtion
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// create product
		prod := &data.Product{}
		// get from json and show error if it not works
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarchal json", http.StatusBadRequest)
			return
		}
		// Validade product
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s", err), http.StatusBadRequest)
		}

		// product on context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
