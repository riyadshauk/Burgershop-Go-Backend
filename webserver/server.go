package main

import (
	// "golang.org/x/net/html"
	// "context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/riyadshauk/burgershopdb"
	"encoding/json"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/compute/v1"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello riyad!") // send data to client side
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client hit endpoint 'insertUser'")

	decoder := json.NewDecoder(r.Body)
	var userGo burgershopdb.User
	err := decoder.Decode(&userGo)
	if err != nil {
		log.Fatal("error:")
		panic(err)
	}
	username := userGo.Username
	password := userGo.Password
	name := userGo.Name
	email := userGo.Email
	mobile := userGo.Mobile
	fmt.Printf("insertUser: received user JSON with userGo.Username: %s, userGo.Password: %s, userGo.Name: %s, userGo.Email: %s, userGo.Mobile: %s\n", userGo.Username, userGo.Password, userGo.Name, userGo.Email, userGo.Mobile)
	burgershopdb.InsertUser(username, password, name, email, mobile)
	fmt.Fprintf(w, "Hello %s, your account (username: %s) has been (...hopefully...) added to the database.\n", name, username)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client has hit endpoint 'getUser'")
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Printf("http request... username: %s, password: %s\n", username, password)
	userJSONAsData := burgershopdb.GetUser(username, password)
	fmt.Fprint(w, string(userJSONAsData))
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client hit endpoint 'getOrders'")
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Printf("http request... username: %s, password: %s\n", username, password)
	ordersJSONAsData := burgershopdb.GetOrders(username, password)
	fmt.Fprint(w, string(ordersJSONAsData))
}

func insertOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client hit endpoint 'insertOrder'")
	decoder := json.NewDecoder(r.Body)
	var orderRequest burgershopdb.OrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		log.Fatal("error:")
		panic(err)
	}
	burgershopdb.InsertOrder(orderRequest)
	fmt.Printf("insertOrder: received orderRequest with orderRequest.username: %s, orderRequest.password: %s\n", orderRequest.Username, orderRequest.Password)
	fmt.Fprintf(w, "insertOrder: received orderRequest with orderRequest.username: %s, orderRequest.password: %s\n", orderRequest.Username, orderRequest.Password)
	// fmt.Fprintf(w, "Hello %s, your order has been (...hopefully...) added to the database.\n", username)
}

// the mobile client won't have a valid authcode, thus there will be no in-app interface for them to use this endpoint
func insertOrderItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client hit endpoint 'insertOrderItem'")
	r.ParseForm()
	authcode := r.PostFormValue("authcode")
	name := r.PostFormValue("name")
	price := r.PostFormValue("price")
	fmt.Printf("http request... name: %s, price: %s\n", name, price)
	burgershopdb.InsertOrderItem(authcode, name, price)
	// fmt.Fprintf(w, "Hello %s, your account (username: %s) has been (...hopefully...) added to the database.\n", name, username)
}

func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/insertUser", insertUser)
	http.HandleFunc("/insertOrder", insertOrder)
	http.HandleFunc("/insertOrderItem", insertOrderItem)

	http.HandleFunc("/getUser", getUser)
	http.HandleFunc("/getOrders", getOrders)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// // https://github.com/google/google-api-go-client/tree/master/#application-default-credentials-example
	// // Use oath2.NoContext if there isn't a good context to pass in.
	// ctx := context.Background()
	// client, err := google.DefaultClient(ctx, compute.ComputeScope)
	// if err != nil {
	// 	// ...
	// }
	// computeService, err := compute.New(client)
	// if err != nil {
	// 	// ...
	// }
	// ts, err := google.DefaultTokenSource(ctx)
	// if err != nil {
	// 	// ...
	// }
	// client := oath2.NewClient(ctx, ts)

	// // also see: https://godoc.org/golang.org/x/oauth2/google

}
