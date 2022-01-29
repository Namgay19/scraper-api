package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Set up environment variables and connection to postgres
	setEnvironmentVariables()
	setup_db()

	// Handlers
	notifHandler := http.HandlerFunc(notificationHandler)

	http.Handle("/api/notifications", loggingMiddleware(notifHandler))

	// Start server
	err := http.ListenAndServe(":8081", nil)
	log.Fatal(err)
}

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Link string `json:"link"`
	Source string `json:"source"`
	Date sql.NullTime `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

func checkError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	sql := "SELECT * FROM posts LIMIT 25"
	var posts []Post

	rows, err := db.Query(sql);
	checkError(err, w)

	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Link, &post.Source, &post.Date, &post.CreatedAt,)
		checkError(err, w)

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response, err := json.Marshal(posts)
	checkError(err, w)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
