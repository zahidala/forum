package posts

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	utils "forum/pkg/utils"
	"log"
	"net/http"
	"strconv"
	// "fmt"
)

type Post struct {
	Types.Post
	Author Types.User `json:"author"`
}

func GetPostsFromCategoryHandler(w http.ResponseWriter, r *http.Request) []Post {
	categoryId := r.PathValue("id")

	query := `SELECT json_group_array(
	json_object(
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
	)
) AS posts
FROM posts p
LEFT JOIN users u ON p.authorId = u.id
LEFT JOIN PostCategories pc ON p.id = pc.postId
LEFT JOIN categories c ON pc.categoryId = c.id
WHERE pc.categoryId = ?
GROUP BY p.id
ORDER BY p.createdAt DESC;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil
	}

	rows, err := stmt.Query(categoryId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}

	var results []Post

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result []Post
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result...)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}

type PostWithMoreDetails struct {
	Types.Post
	Author Types.User `json:"author"`
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) (PostWithMoreDetails, Types.ErrorPageProps) {
	postId := r.PathValue("id")

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
	'attachments', p.attachments
) AS post
FROM Posts p
LEFT JOIN Users u ON p.authorId = u.id
WHERE p.id = ?;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		return PostWithMoreDetails{}, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		}
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		log.Println("Error executing query:", err)

		return PostWithMoreDetails{}, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		}
	}
	defer rows.Close()

	var result PostWithMoreDetails
	var found bool

	for rows.Next() {
		found = true
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			log.Println("Error scanning row:", err)

			return PostWithMoreDetails{}, Types.ErrorPageProps{
				Error: Types.Error{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				},
				Title: "Internal Server Error",
			}
		}

		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			log.Println("Error unmarshaling json:", errJsonUnmarshal)

			return PostWithMoreDetails{}, Types.ErrorPageProps{
				Error: Types.Error{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				},
				Title: "Internal Server Error",
			}
		}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Println("Error iterating rows:", rowsErr)

		return PostWithMoreDetails{}, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			},
			Title: "Internal Server Error",
		}
	}

	if !found {
		log.Println("Post not found for ID:", postId)

		return PostWithMoreDetails{}, Types.ErrorPageProps{
			Error: Types.Error{
				Code:    http.StatusNotFound,
				Message: "Post not found",
			},
			Title: "Post not found",
		}
	}

	return result, Types.ErrorPageProps{}
}

type PostLikeWithAuthor struct {
	Types.PostLike
	Author Types.User `json:"author"`
}

func GetPostLikesHandler(w http.ResponseWriter, r *http.Request) []PostLikeWithAuthor {
	postId := r.PathValue("id")

	// Prepare the SQL statement
	query := `SELECT json_group_array(
					json_object(
									'id', pl.id,
									'postId', pl.postId,
									'userId', pl.userId,
									'author', json_object(
													'id', u.id,
													'name', u.name,
													'username', u.username,
													'profilePicture', u.profilePicture
									),
									'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', pl.createdAt),
									'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', pl.updatedAt)
					)
	)
	FROM PostLikes pl
	LEFT JOIN Users u ON pl.userId = u.id
	WHERE pl.postId = ? AND pl.isLike = 1;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}

	var results []PostLikeWithAuthor

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result []PostLikeWithAuthor
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result...)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}

type PostDislikeWithAuthor struct {
	Types.PostDisLike
	Author Types.User `json:"author"`
}

func GetPostDislikesHandler(w http.ResponseWriter, r *http.Request) []PostDislikeWithAuthor {
	postId := r.PathValue("id")

	// Prepare the SQL statement
	query := `SELECT json_group_array(
					json_object(
									'id', pd.id,
									'postId', pd.postId,
									'userId', pd.userId,
									'author', json_object(
													'id', u.id,
													'name', u.name,
													'username', u.username,
													'profilePicture', u.profilePicture
									),
									'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', pd.createdAt),
									'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', pd.updatedAt)
					)
	)
	FROM PostDislikes pd
	LEFT JOIN Users u ON pd.userId = u.id
	WHERE pd.postId = ? AND pd.isDislike = 1;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}

	var results []PostDislikeWithAuthor

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result []PostDislikeWithAuthor
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result...)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}

type PostCreateBody struct {
	UserId             string `json:"userId"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	Images             string `json:"images,omitempty"`
	SelectedCategories []int  `json:"selectedCategories"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var body PostCreateBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// Prepare the SQL statement
	query := `INSERT INTO Posts (authorId, title, content, attachments) VALUES (?, ?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.UserId, body.Title, body.Content, body.Images)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting last insert id:", err)
		return
	}

	// Prepare the SQL statement
	categoryQuery := `INSERT INTO PostCategories (postId, categoryId) VALUES (?, ?);`

	categoryStmt, err := db.GetDB().Prepare(categoryQuery)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}

	for _, categoryId := range body.SelectedCategories {
		_, err = categoryStmt.Exec(postId, categoryId)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error executing query:", err)
			return
		}
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(int(postId)), http.StatusSeeOther)
}

type PostCreateByCategoryBody struct {
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  string `json:"images,omitempty"`
}

func CreatePostByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryId := r.PathValue("id")

	var body PostCreateByCategoryBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error decoding body:", err)
		return
	}

	// Prepare the SQL statement
	query := `INSERT INTO Posts (authorId, title, content, attachments) VALUES (?, ?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.UserId, body.Title, body.Content, body.Images)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error getting last insert id:", err)
		return
	}

	// Prepare the SQL statement
	categoryQuery := `INSERT INTO PostCategories (postId, categoryId) VALUES (?, ?);`

	categoryStmt, err := db.GetDB().Prepare(categoryQuery)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}

	_, err = categoryStmt.Exec(postId, categoryId)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(int(postId)), http.StatusSeeOther)
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

type PostWithCategory struct {
	Types.Post
	Categories []Types.Category `json:"categories"`
	Author     Types.User       `json:"author"`
}

func GetNewPostsHandler(w http.ResponseWriter, r *http.Request) []PostWithCategory {
	// Prepare the SQL statement
	query := `SELECT json_group_array(
    json_object(
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
        'categories', (
            SELECT json_group_array(
                json_object(
                    'id', c.id,
                    'name', c.name,
                    'description', c.description,
                    'icon', c.icon
                )
            )
            FROM categories c
            JOIN PostCategories pc ON c.id = pc.categoryId
            WHERE pc.postId = p.id
        )
    )
) AS posts
FROM posts p
LEFT JOIN users u ON p.authorId = u.id
GROUP BY p.id
ORDER BY p.createdAt DESC
LIMIT 3;
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

	var results []PostWithCategory

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result []PostWithCategory
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", errJsonUnmarshal)
			return nil
		}

		results = append(results, result...)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil
	}

	return results
}

type PostWithFilter struct {
	PostID         int
	Title          string
	CreatedAt      string
	UpdatedAt      string
	UserID         int
	Username       string
	ProfilePicture string
	Categories     []Category `json:"categories"`
}

type Category struct {
	CategoryID   int    `json:"categoryID"`
	CategoryName string `json:"categoryName"`
}

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) ([]PostWithFilter, utils.Filters) {

	// query := ` SELECT * FROM (
	// SELECT p.id AS postID , p.title , p.createdAt , p.updatedAt , u.id AS userID , u.username , u.profilePicture , pc.categoryId, c.name
	// FROM Posts p
	// JOIN Users u ON p.authorId = u.id
	// JOIN PostCategories pc ON p.id  = pc.postId
	// JOIN Categories c ON pc.categoryId = c.id
	// ) `

	query := `SELECT * FROM (
	SELECT p.id AS postID , p.title , p.createdAt , p.updatedAt , u.id AS userID , u.username , u.profilePicture ,
	JSON_GROUP_ARRAY(
		JSON_OBJECT(
		'categoryID', pc.categoryId,
		'categoryName', c.name 
		)
		) AS categories
	FROM Posts p
	JOIN Users u ON p.authorId = u.id
	JOIN PostCategories pc ON p.id  = pc.postId
	JOIN Categories c ON pc.categoryId = c.id
	GROUP BY p.id
	)`

	filterQuery, filters := utils.GetFilteredPosts(w, r)
	query += filterQuery
	// query += "AS subquery\nORDER BY createdAt DESC"

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return nil, filters
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil, filters
	}
	defer rows.Close()

	var posts []PostWithFilter

	for rows.Next() {
		var post PostWithFilter
		var categoriesJSON string

		err := rows.Scan(
			&post.PostID,
			&post.Title,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UserID,
			&post.Username,
			&post.ProfilePicture,
			&categoriesJSON,
		)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil, filters
		}

		err = json.Unmarshal([]byte(categoriesJSON), &post.Categories)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, post)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return nil, filters
	}

	return posts, filters
}
