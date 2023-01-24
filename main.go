package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	UserId   int
	UserName string
}

var db *gorm.DB

func main() {
	var err error

	dsn := "host=host.docker.internal user=postgres password=12345 dbname=postgres port=9000 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Database is connected sucessfully")

	dbinstance, _ := db.DB()

	defer dbinstance.Close()

	db.AutoMigrate(User{})

	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/products/{id}", getProducts).Methods("GET")
	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/orders", createOrders).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", router)
	fmt.Println("server is running on 8080 port")

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user []User
	db.Find(&user)
	json.NewEncoder(w).Encode(user)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	params := mux.Vars(r)

	db.First(&user, params["id"])

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("%#v\n", user)
	db.Save(&user)

	json.NewEncoder(w).Encode(&user)

}

func getProducts(w http.ResponseWriter, r *http.Request) {

	client := resty.New()

	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8001/products")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.Body())

}

func getOrders(w http.ResponseWriter, r *http.Request) {

	client := resty.New()

	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8003/orders")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.Body())

}

func createOrders(w http.ResponseWriter, r *http.Request) {

	client := resty.New()

	resp, _ := client.R().
		Get("http://localhost:8002/stocks")

	if resp.StatusCode() == 200 {

		resp, err := client.R().Post("http://localhost:8003/orders")

		if err != nil {
			panic(err)
		}
		w.Write([]byte(resp.Body()))
	} else {
		w.Write([]byte("Product is not available in memory"))
	}

}
