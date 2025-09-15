package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ps/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("unable to decode the request body. %v", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "stock created succesfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {

		log.Fatalf("unable to convert string into int, %v", err)

	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("unable to get stock, %v", err)

	}

	json.NewEncoder(w).Encode(stock)

}

func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("unable to getallstocks, %v", err)
	}

	json.NewEncoder(w).Encode(stocks)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("UNABLE to convert string into int, %v ", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unablr to decode the request body, %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("stock updated successfully. total rows.records affedted %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"])

	if err != nil {
		log.Fatalf("unable to convert string to int %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("stock deleted successfully. Total rows/records deleted %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func insertStock(stock models.Stock) int64 {
	db := createConnection()

	defer db.Close()
	sqlStatement := `INSERT INTO stocks(name,price,company) VALUES ($1,$2, $3) RETURNING stockid`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	fmt.Sprintf("Inserted a single record %v", id)
}

func getStock(id int64) (models.Stock, error) {

	db := createConnection()

	defer db.Close()

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row %v", err)

		return stock, err
	}
}

func getAllStocks() ([]models.Stock, error) {

}

func updateStock(id int64, stock models.Stock) int64 {

}

func deleteStock(id int64) int64 {

}
