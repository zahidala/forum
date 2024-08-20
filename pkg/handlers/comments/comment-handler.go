package comments

import (
	"encoding/json"
	"forum/pkg/db"
	"io"
	"log"
	"net/http"
)

type CommentBody struct {
	Content string `json:"content"`
	UserID  string `json:"userId"`
	Images  string `json:"images,omitempty"`
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

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

	query := `INSERT INTO comments (content, postId, authorId, attachments) VALUES (?, ?, ?, ?)`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentBody.Content, postId, commentBody.UserID, commentBody.Images)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

type CommentReactionBody struct {
	UserId string `json:"userId"`
	PostId string `json:"postId"`
}

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	commentId := r.PathValue("id")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var commentLikeBody CommentReactionBody

	jsonErr := json.Unmarshal(reqBody, &commentLikeBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// check if user already disliked the comment and remove the like

	dislikeQuery := `DELETE FROM commentdislikes WHERE commentId = ? AND userId = ?`

	dislikeStmt, err := db.GetDB().Prepare(dislikeQuery)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = dislikeStmt.Exec(commentId, commentLikeBody.UserId)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO commentlikes (commentId, userId, isLike) VALUES (?, ?, ?)`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentId, commentLikeBody.UserId, 1)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+commentLikeBody.PostId, http.StatusSeeOther)
}

func CommentRemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	commentId := r.PathValue("id")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var commentLikeBody CommentReactionBody

	jsonErr := json.Unmarshal(reqBody, &commentLikeBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM commentlikes WHERE commentId = ? AND userId = ?`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentId, commentLikeBody.UserId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+commentLikeBody.PostId, http.StatusSeeOther)
}

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	commentId := r.PathValue("id")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var commentDislikeBody CommentReactionBody

	jsonErr := json.Unmarshal(reqBody, &commentDislikeBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// check if user already liked the comment and remove the like

	likeQuery := `DELETE FROM commentlikes WHERE commentId = ? AND userId = ?`

	likeStmt, err := db.GetDB().Prepare(likeQuery)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = likeStmt.Exec(commentId, commentDislikeBody.UserId)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO commentdislikes (commentId, userId, isDislike) VALUES (?, ?, ?)`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentId, commentDislikeBody.UserId, 1)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+commentDislikeBody.PostId, http.StatusSeeOther)
}

func CommentRemoveDislikeHandler(w http.ResponseWriter, r *http.Request) {
	commentId := r.PathValue("id")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var commentDislikeBody CommentReactionBody

	jsonErr := json.Unmarshal(reqBody, &commentDislikeBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM commentdislikes WHERE commentId = ? AND userId = ?`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(commentId, commentDislikeBody.UserId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+commentDislikeBody.PostId, http.StatusSeeOther)
}
