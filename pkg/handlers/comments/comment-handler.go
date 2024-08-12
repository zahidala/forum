package comments

import (
	"encoding/json"
	"fmt"
	"forum/pkg/db"
	"io"
	"net/http"
)

type CommentBody struct {
	Content string `json:"content"`
	UserID  string `json:"userId"`
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

	fmt.Println("Creating comment for post with id: ", postId)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var commentBody CommentBody

	jsonErr := json.Unmarshal(reqBody, &commentBody)
	if jsonErr != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO comments (content, postId, authorId) VALUES (?, ?, ?)`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentBody.Content, postId, commentBody.UserID)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}
