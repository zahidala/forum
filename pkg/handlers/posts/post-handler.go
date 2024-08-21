package posts

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	Types.Post
	Author Types.User `json:"author"`
}

func GetPostsFromSubCategoryHandler(w http.ResponseWriter, r *http.Request) []Post {
	subcategoryId := r.PathValue("id")

	// Prepare the SQL statement
	query := `SELECT json_object(
			'id', p.id,
			'title', p.title,
			'content', p.content,
			'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.createdAt),
			'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.updatedAt),
			'author', json_object(
											'id', u.id,
											'name', u.name,
											'username', u.username,
											'profilePicture', u.profilePicture
			)
) AS post

FROM Posts p

LEFT JOIN Users u ON p.authorId = u.id

WHERE p.subcategoryId = ?;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(subcategoryId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var results []Post

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result Post
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}

type PostLikeWithAuthor struct {
	Types.PostLike
	Author Types.User `json:"author"`
}

type PostDislikeWithAuthor struct {
	Types.PostDisLike
	Author Types.User `json:"author"`
}

type PostWithMoreDetails struct {
	Types.Post
	Likes    []PostLikeWithAuthor    `json:"likes,omitempty"`
	Dislikes []PostDislikeWithAuthor `json:"dislikes,omitempty"`
	Author   Types.User              `json:"author"`
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) PostWithMoreDetails {
	postId := r.PathValue("id")

	// Prepare the SQL statement
	query := `SELECT json_object(
	'id', p.id,
	'title', p.title,
	'content', p.content,
	'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.createdAt),
	'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.updatedAt),
	'likes', (
			SELECT json_group_array(
					json_object(
							'id', pl.id,
							'author', json_object(
									'id', u1.id,
									'name', u1.name,
									'username', u1.username,
									'profilePicture', u1.profilePicture
							),
							'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', pl.createdAt),
							'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', pl.updatedAt)
					)
			)
			FROM PostLikes pl
			LEFT JOIN Users u1 ON pl.userId = u1.id
			WHERE pl.postId = p.id AND pl.isLike = 1
	),
	'dislikes', (
			SELECT json_group_array(
					json_object(
							'id', pdl.id,
							'author', json_object(
									'id', u2.id,
									'name', u2.name,
									'username', u2.username,
									'profilePicture', u2.profilePicture
							),
							'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', pdl.createdAt),
							'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', pdl.updatedAt)
					)
			)
			FROM PostDislikes pdl
			LEFT JOIN Users u2 ON pdl.userId = u2.id
			WHERE pdl.postId = p.id AND pdl.isDislike = 1
	),
	'author', json_object(
			'id', u.id,
			'name', u.name,
			'username', u.username,
			'profilePicture', u.profilePicture
	),
	'attachments', p.attachments
) AS post
FROM Posts p
LEFT JOIN Users u ON p.authorId = u.id
WHERE p.id = ?;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return PostWithMoreDetails{}
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return PostWithMoreDetails{}
	}
	defer rows.Close()

	var result PostWithMoreDetails

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return PostWithMoreDetails{}
		}

		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return PostWithMoreDetails{}
		}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return PostWithMoreDetails{}
	}

	return result
}

type PostCreateBody struct {
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  string `json:"images,omitempty"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	subcategoryId := r.PathValue("id")

	var body PostCreateBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// Prepare the SQL statement
	query := `INSERT INTO Posts (authorId, subcategoryId, title, content, attachments) VALUES (?, ?, ?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.UserId, subcategoryId, body.Title, body.Content, body.Images)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting last insert ID:", err)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(int(postID)), http.StatusSeeOther)
}

type PostReactionBody struct {
	UserId string `json:"userId"`
}

func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

	var body PostReactionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// check if the user has already disliked the post and remove the dislike

	dislikeQuery := `DELETE FROM PostDislikes WHERE postId = ? AND userId = ?;`

	dislikeStmt, err := db.GetDB().Prepare(dislikeQuery)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}

	_, err = dislikeStmt.Exec(postId, body.UserId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	// Prepare the SQL statement
	query := `INSERT INTO PostLikes (postId, userId, isLike) VALUES (?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, body.UserId, 1)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

func PostRemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

	var body PostReactionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// Prepare the SQL statement
	query := `DELETE FROM PostLikes WHERE postId = ? AND userId = ?;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, body.UserId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

func PostDislikeHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

	var body PostReactionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// check if the user has already liked the post and remove the like

	likeQuery := `DELETE FROM PostLikes WHERE postId = ? AND userId = ?;`

	likeStmt, err := db.GetDB().Prepare(likeQuery)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}

	_, err = likeStmt.Exec(postId, body.UserId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	// Prepare the SQL statement
	query := `INSERT INTO PostDislikes (postId, userId, isDislike) VALUES (?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, body.UserId, 1)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

func PostRemoveDislikeHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")

	var body PostReactionBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// Prepare the SQL statement
	query := `DELETE FROM PostDislikes WHERE postId = ? AND userId = ?;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postId, body.UserId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
}

func IsPostLikedByCurrentUserHandler(w http.ResponseWriter, r *http.Request, postId int, userId int) bool {
	// Prepare the SQL statement
	query := `SELECT COUNT(*) FROM PostLikes WHERE postId = ? AND userId = ?;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(postId, userId).Scan(&count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return false
	}

	return count > 0
}

func IsPostDisLikedByCurrentUserHandler(w http.ResponseWriter, r *http.Request, postId int, userId int) bool {
	// Prepare the SQL statement
	query := `SELECT COUNT(*) FROM PostDislikes WHERE postId = ? AND userId = ?;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(postId, userId).Scan(&count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return false
	}

	return count > 0
}

type PostWithSubcategory struct {
	Types.Post
	Author      Types.User        `json:"author"`
	Subcategory Types.Subcategory `json:"subcategory"`
}

func GetNewPostsHandler(w http.ResponseWriter, r *http.Request) []PostWithSubcategory {
	// Prepare the SQL statement
	query := `SELECT json_object(
			'id', p.id,
			'title', p.title,
			'content', p.content,
			'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.createdAt),
			'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', p.updatedAt),
			'author', json_object(
											'id', u.id,
											'name', u.name,
											'username', u.username,
											'profilePicture', u.profilePicture
			),
			'subcategory', json_object(
											'id', s.id,
											'name', s.name,
											'description', s.description,
											'category', json_object(
																			'id', c.id,
																			'name', c.name,
																			'description', c.description
											)
			)
) AS post

FROM Posts p
LEFT JOIN Users u ON p.authorId = u.id
LEFT JOIN Subcategories s ON p.subcategoryId = s.id
LEFT JOIN Categories c ON s.categoryId = c.id
ORDER BY p.createdAt DESC
LIMIT 10;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil
	}

	rows, err := stmt.Query()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}

	var results []PostWithSubcategory

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result PostWithSubcategory
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}
