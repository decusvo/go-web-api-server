package handlers // defines the fact that this file will be part of a package called 'handlers'.

import (
	"example/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products { // a function that takes in a log.logger object and returns a handler as a reference.
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Conditional to determine the method of the HTTP request.
	if r.Method == http.MethodGet { // If this is a GET request, call getProducts.
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost { // If this is a POST request, call addProduct.
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {

		regexpattern := regexp.MustCompile("/([0-9])+")
		regexMatches := regexpattern.FindAllStringSubmatch(r.URL.Path, -1)

		if len(regexMatches) != 1 {
			p.l.Println("Regex Match 1: ", regexMatches)
			http.Error(rw, "Invalid URI item ID provided in the URL.", http.StatusBadRequest)
			return
		}

		if len(regexMatches[0]) != 2 {
			p.l.Println("Regex Match 2: ", regexMatches)
			http.Error(rw, "Invalid item ID provided in the URL.", http.StatusBadRequest)
			return
		}

		strResult := regexMatches[0][1]
		productId, err := strconv.Atoi(strResult)

		if err != nil {
			http.Error(rw, "Invalid URL: Couldn't convert ID from string to integer.", http.StatusBadRequest)
			return
		}

		p.updateProducts(productId, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed) // Return default header.

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Couldn't return JSON data", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{} // Create a new Product object. data is a package containing the definition of the Product struct.

	// r.Body implements io.Reader which we pass into the FromJSON function.
	// The reason we get an io.Reader because not all of the data in the HTTP request is read, mainly because there may be a lot of data
	// and it may not be even necessary to read the entirety of the request.
	err := prod.FromJSON(r.Body)
	p.l.Println(err)
	if err != nil {
		http.Error(rw, "Couldn't decode the incoming JSON data", http.StatusBadRequest)
	}

	p.l.Printf("Product : %#v", prod)

	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("PUT Request - Update Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	p.l.Println(err)
	if err != nil {
		http.Error(rw, "Couldn't decode the incoming JSON data", http.StatusBadRequest)

	}

	p.l.Println("Update Product Start", id)
	err2 := data.UpdateProduct(id, prod)

	p.l.Println("Update Product Error: ", err2)

	if err2 == data.ErrorProductNotFound {
		http.Error(rw, "ERROR: Product Not Found", http.StatusNotFound)
		return
	}

	if err2 != nil {
		http.Error(rw, "ERROR: Product Not Found", http.StatusInternalServerError)
		return
	}
}
