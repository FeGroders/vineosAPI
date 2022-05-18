package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
    DB_USER     = "postgres"
    DB_PASSWORD = "postgres"
    DB_NAME     = "postgres"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

// DB set up
func setupDB() *sql.DB {
    // dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
    dbinfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

    checkErr(err)

    return db
}

// Main function
func main() {
    // Init the mux router
    router := mux.NewRouter()

	// Route handles & endpoints
    // Get all vinhos
    router.HandleFunc("/vinhos/", GetVinhos).Methods("GET")

    // Create a movie
    router.HandleFunc("/vinhos/", CreateVinho).Methods("POST")

    // Delete a specific movie by the movieID
    router.HandleFunc("/vinhos/{movieid}", DeleteVinho).Methods("DELETE")

    // serve the app
    fmt.Println("Server at 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

type Vinho struct {
    vinhoID         int `json:"id"`
    vinhoNome       string `json:"nome"`
	vinhoDescricao  string `json:"descricao"`
	vinhoAno        int `json:"ano"`
	vinhoPreco      float32 `json:"preco"`
	vinhoImagem     string `json:"imagem"`
	vinhoDisponivel bool `json:"disponivel"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Vinho `json:"data"`
    Message string `json:"message"`
}

// Function for handling messages
func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

// Get all vinhos
// response and request handlers
func GetVinhos(w http.ResponseWriter, r *http.Request) {
    db := setupDB()

    printMessage("Getting vinhos...")

    // Get all vinhos from vinhos table that don't have movieID = "1"
    rows, err := db.Query("SELECT * FROM vinhos")

    // check errors
    checkErr(err)

    // var response []JsonResponse
    var vinhos []Vinho

    // Foreach movie
    for rows.Next() {
        var id int
        var nome string
		var descricao string
		var ano int
		var preco float32
		var imagem string
		var disponivel bool

        err = rows.Scan(&id, &nome, &descricao, &ano, &preco, &imagem, &disponivel)

        // check errors
        checkErr(err)

        vinhos = append(vinhos, Vinho{vinhoID: id, vinhoNome: nome, vinhoDescricao: descricao, vinhoAno: ano, 
			vinhoPreco: preco, vinhoImagem: imagem, vinhoDisponivel: disponivel})
    }

    var response = JsonResponse{Type: "success", Data: vinhos}

    json.NewEncoder(w).Encode(response)
}

// Create a movie
// response and request handlers
func CreateVinho(w http.ResponseWriter, r *http.Request) {
    vinhoID := r.FormValue("vinhoID")
    vinhoNome := r.FormValue("vinhoNome")

    var response = JsonResponse{}

    if vinhoID == "" || vinhoNome == "" {
        response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
    } else {
        db := setupDB()

        printMessage("Inserting vinho into DB")

        fmt.Println("Inserting new vinho with ID: " + vinhoID + " and name: " + vinhoNome)

        var lastInsertID int
    err := db.QueryRow("INSERT INTO vinhos(movieID, movieName) VALUES($1, $2) returning id;", vinhoID, vinhoNome).Scan(&lastInsertID)

    // check errors
    checkErr(err)

    response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
    }

    json.NewEncoder(w).Encode(response)
}

// Delete a movie
// response and request handlers
func DeleteVinho(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    vinhoID := params["vinhoID"]

    var response = JsonResponse{}

    if vinhoID == "" {
        response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
    } else {
        db := setupDB()

        printMessage("Deleting movie from DB")

        _, err := db.Exec("DELETE FROM vinhos where movieID = $1", vinhoID)

        // check errors
        checkErr(err)

        response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
    }

    json.NewEncoder(w).Encode(response)
}