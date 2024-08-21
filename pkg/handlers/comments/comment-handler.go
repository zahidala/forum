package comments

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"io"
	"log"
	"net/http"
)

type CommentLikeWithAuthor struct {
	Types.CommentLike
	Author Types.User `json:"author"`
}

type CommentDislikeWithAuthor struct {
	Types.CommentDisLike
	Author Types.User `json:"author"`
}

type CommentWithMoreDetails struct {
	Types.Comment
	Author   Types.User                 `json:"author"`
	Likes    []CommentLikeWithAuthor    `json:"likes,omitempty"`
	Dislikes []CommentDislikeWithAuthor `json:"dislikes,omitempty"`
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) []CommentWithMoreDetails {
	postId := r.PathValue("id")

	query := `SELECT json_group_array(
			json_object(
					'id', c.id,
					'postId', c.postId,
					'content', c.content,
					'author', json_object(
							'id', u.id,
							'name', u.name,
							'username', u.username,
							'profilePicture', u.profilePicture
					),
					'attachments', c.attachments,
					'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', c.createdAt),
					'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', c.updatedAt),
					'likes', (
							SELECT json_group_array(
									json_object(
											'id', cl.id,
											'commentId', cl.commentId,
											'author', json_object(
													'id', ul.id,
													'name', ul.name,
													'username', ul.username,
													'profilePicture', ul.profilePicture
											),
											'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', cl.createdAt),
											'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', cl.updatedAt)
									)
							) FROM commentlikes cl
							JOIN users ul ON cl.userId = ul.id
							WHERE cl.commentId = c.id
					),
					'dislikes', (
							SELECT json_group_array(
									json_object(
											'id', cd.id,
											'commentId', cd.commentId,
											'author', json_object(
													'id', ud.id,
													'name', ud.name,
													'username', ud.username,
													'profilePicture', ud.profilePicture
											),
											'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', cd.createdAt),
											'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', cd.updatedAt)
									)
							) FROM commentdislikes cd
							JOIN users ud ON cd.userId = ud.id
							WHERE cd.commentId = c.id
					)
			)
	) AS comments
	FROM comments c
	JOIN users u ON c.authorId = u.id
	WHERE c.postId = ?`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return []CommentWithMoreDetails{}
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return []CommentWithMoreDetails{}
	}
	defer rows.Close()

	var commentsJson string

	for rows.Next() {
		err := rows.Scan(&commentsJson)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return []CommentWithMoreDetails{}
		}
	}

	var comments []CommentWithMoreDetails

	errJsonUnmarshal := json.Unmarshal([]byte(commentsJson), &comments)
	if errJsonUnmarshal != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error unmarshaling JSON:", errJsonUnmarshal)
		return []CommentWithMoreDetails{}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return []CommentWithMoreDetails{}
	}

	return comments
}

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
