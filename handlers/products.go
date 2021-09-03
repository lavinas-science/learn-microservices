package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/lavinas-science/learn-microservices/data"
)

type Products struct {
	l *log.Logger
}


func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP (rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
	}

	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)

	}

	// Catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}


func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marchal json", http.StatusInternalServerError)
		return
	}
}


func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(rw, "Unable to unmarchal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)

}


func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	rg := regexp.MustCompile(`/([0-9]+)`)
	g := rg.FindAllStringSubmatch(r.URL.Path, -1)
	if len(g) != 1 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}
	if len(g[0]) != 2 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}
	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(rw, "Unable to unmarchal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} 

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}


}