package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Order struct {
	OrderID      uint   `json:"order_id" gorm:"primary_key"`
	CustomerName string `json:"customer_name"`
	OrderedAt    string `json:"ordered_at"`
	Items        []Item `json:"items" gorm:"foreignkey:OrderID"`
}

type Item struct {
	ItemID      uint   `json:"item_id" gorm:"primary_key"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"order_id"`
}

var db *gorm.DB

func SetupDB() {
	var err error

	USER := "root"
	PASS := ""
	HOST := "localhost"
	PORT := "3306"
	DBNAME := "assignment2"

	CONNECT := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT)

	db, err = gorm.Open("mysql", CONNECT)

	if err != nil {
		fmt.Println(err)
		panic("Gagal terhubung ke database")
	}

	sql := "CREATE DATABASE " + DBNAME

	db.Exec(sql)

	sql2 := "USE " + DBNAME

	db.Exec(sql2)
	db.AutoMigrate(&Order{}, &Item{})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input Order

	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&input)

	db.Create(&input)

	json.NewEncoder(w).Encode(input)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []Order

	w.Header().Set("Content-Type", "application/json")

	db.Preload("Items").Find(&orders)

	json.NewEncoder(w).Encode(orders)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var input Order

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	ParamOrderID, _ := strconv.Atoi(params["orderId"])

	ID := uint(ParamOrderID)

	db.First(&input, ID)

	json.NewDecoder(r.Body).Decode(&input)

	db.Save(&input)

	json.NewEncoder(w).Encode(input)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ParamOrderID, _ := strconv.Atoi(params["orderId"])

	ID := uint(ParamOrderID)

	db.Where("order_id = ?", ID).Delete(&Item{})

	db.Where("order_id = ?", ID).Delete(&Order{})

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	SetupDB()
	router := mux.NewRouter()

	router.HandleFunc("/orders", CreateOrder).Methods("POST")
	router.HandleFunc("/orders", GetOrders).Methods("GET")
	router.HandleFunc("/orders/{orderId}", UpdateOrder).Methods("PUT")
	router.HandleFunc("/orders/{orderId}", DeleteOrder).Methods("DELETE")

	fmt.Println("Running")
	http.ListenAndServe(":8080", router)
}
