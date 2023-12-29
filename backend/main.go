package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

const SecretKey = "secret"

var db *sql.DB

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Role   string `json:"role"`
	Group  string `json:"group"`
}

type User struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    []byte `json:"-"`
	Deactivated int    `json:"deactivated"`
	Role        string `json:"role"`
	Group       string `json:"group"`
}

func main() {
	fmt.Println("Starting")
	var err error
	db, err = sql.Open("mysql", "root:@root123#123@tcp(localhost:3306)/task_management")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/user", getUser).Methods("GET")
	router.HandleFunc("/deactivate", DeactivateAccount).Methods("DELETE")
	router.HandleFunc("/delete", DeleteAccount).Methods("DELETE")
	router.HandleFunc("/upload/users", uploadUsersCSV).Methods("POST")
	router.HandleFunc("/upload/tasks", uploadTasksCSV).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := c.Handler(router)

	log.Println("Starting the HTTP server on port 3001")
	log.Fatal(http.ListenAndServe(":3001", handler))
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{}
	rows, err := db.Query("SELECT id, title, status, roleId, groupId FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Status, &task.Role, &task.Group)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	_, err := db.Exec("INSERT INTO tasks (title, status, roleId, groupId) VALUES (?, ?, ?, ?)", task.Title, task.Status, task.Role, task.Group)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	_, err := db.Exec("UPDATE tasks SET title = ?, status = ? WHERE id = ?", task.Title, task.Status, id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := User{
		Name:        data["name"],
		Email:       data["email"],
		Password:    password,
		Role:        data["role"],
		Group:       data["group"],
		Deactivated: 0,
	}

	_, err := db.Exec("INSERT INTO users (name, email, password, roleId, groupId, deactivated) VALUES (?, ?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Password, user.Role, user.Group, user.Deactivated)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT id, name, email, password, deactivated FROM users WHERE email = ?", data["email"]).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Deactivated)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Deactivated == 1 {
		http.Error(w, "User is deactivated", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))

	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		http.Error(w, "Could not login", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to extract claims from token", http.StatusInternalServerError)
		return
	}
	userIDStr, ok := claims["iss"].(string)
	if !ok {
		http.Error(w, "User ID not found in token claims", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID in token claims", http.StatusInternalServerError)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, name, email, roleId, groupId FROM users WHERE id = ?", userID).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Group)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to extract claims from token", http.StatusInternalServerError)
		return
	}
	userIDStr, ok := claims["iss"].(string)
	if !ok {
		http.Error(w, "User ID not found in token claims", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID in token claims", http.StatusInternalServerError)
		return
	}

	_, updateErr := db.Exec("UPDATE users SET deactivated = ? WHERE id = ?", 1, userID)
	if updateErr != nil {
		http.Error(w, "Failed to deactivate the account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account deactivated"))
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to extract claims from token", http.StatusInternalServerError)
		return
	}

	userIDStr, ok := claims["iss"].(string)
	if !ok {
		http.Error(w, "User ID not found in token claims", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID in token claims", http.StatusInternalServerError)
		return
	}

	_, deleteErr := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if deleteErr != nil {
		http.Error(w, "Failed to delete user account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account deleted"))
}

func uploadTasksCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Failed to retrieve the CSV file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading CSV record:", err)
			continue
		}

		if len(record) != 4 {
			log.Println("Invalid CSV record format")
			continue
		}

		title := record[0]
		status := record[1]
		role := record[2]
		group := record[3]

		_, err = db.Exec("INSERT INTO tasks (title, status, roleId, groupId) VALUES (?, ?, ?, ?)", title, status, role, group)
		if err != nil {
			log.Println("Failed to insert task:", err)
			continue
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tasks uploaded successfully"))
}

func uploadUsersCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Failed to retrieve the CSV file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading CSV record:", err)
			continue
		}

		if len(record) != 4 {
			log.Println("Invalid CSV record format")
			continue
		}

		name := record[0]
		email := record[1]
		role := record[2]
		group := record[3]

		_, err = db.Exec("INSERT INTO users (name, email, roleId, groupId) VALUES (?, ?, ?, ?)", name, email, role, group)
		if err != nil {
			log.Println("Failed to insert user:", err)
			continue
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users uploaded successfully"))
}
