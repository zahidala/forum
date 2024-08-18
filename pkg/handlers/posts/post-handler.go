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
	Author      Types.User        `json:"author"`
	Subcategory Types.Subcategory `json:"subcategory"`
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
			log.Println("Error unmarshalling json:", errJsonUnmarshal)
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

type PostWithComments struct {
	Types.Post
	Author      Types.User               `json:"author"`
	Subcategory Types.Subcategory        `json:"subcategory"`
	Comments    []CommentWithMoreDetails `json:"comments"`
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) PostWithComments {
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
        'subcategory', json_object(
            'id', s.id,
            'name', s.name,
            'description', s.description,
            'category', json_object(
                'id', c2.id,
                'name', c2.name,
                'description', c2.description
            )
        ),
        'comments', (
            SELECT json_group_array(
                json_object(
                    'id', c.id,
                    'postId', c.postId,
                    'content', c.content,
                    'author', json_object(
                        'id', u2.id,
                        'name', u2.name,
                        'username', u2.username,
                        'profilePicture', u2.profilePicture
                    ),
                    'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', c.createdAt),
                    'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', c.updatedAt),
                    'attachments', c.attachments,
                    'likes', (
                        SELECT json_group_array(
                            json_object(
                                'id', cl.id,
                                'author', json_object(
                                    'id', u3.id,
                                    'name', u3.name,
                                    'username', u3.username,
                                    'profilePicture', u3.profilePicture
                                ),
                                'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', cl.createdAt),
                                'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', cl.updatedAt)
                            )
                        )
                        FROM CommentLikes cl
                        LEFT JOIN Users u3 ON cl.userId = u3.id
                        WHERE cl.commentId = c.id AND cl.isLike = 1
                    ),
                    'dislikes', (
                        SELECT json_group_array(
                            json_object(
                                'id', cdl.id,
                                'author', json_object(
                                    'id', u4.id,
                                    'name', u4.name,
                                    'username', u4.username,
                                    'profilePicture', u4.profilePicture
                                ),
                                'createdAt', strftime('%Y-%m-%dT%H:%M:%SZ', cdl.createdAt),
                                'updatedAt', strftime('%Y-%m-%dT%H:%M:%SZ', cdl.updatedAt)
                            )
                        )
                        FROM CommentDislikes cdl
                        LEFT JOIN Users u4 ON cdl.userId = u4.id
                        WHERE cdl.commentId = c.id AND cdl.isDislike = 1
                    )
                )
            )
            FROM Comments c
            LEFT JOIN Users u2 ON c.authorId = u2.id
            WHERE c.postId = p.id
        )
) AS post
FROM Posts p
LEFT JOIN Users u ON p.authorId = u.id
LEFT JOIN Subcategories s ON p.subcategoryId = s.id
LEFT JOIN Categories c2 ON s.categoryId = c2.id
WHERE p.id = ?;
`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return PostWithComments{}
	}
	defer stmt.Close()

	rows, err := stmt.Query(postId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return PostWithComments{}
	}
	defer rows.Close()

	var result PostWithComments

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return PostWithComments{}
		}

		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshalling json:", errJsonUnmarshal)
			return PostWithComments{}
		}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error iterating rows:", rowsErr)
		return PostWithComments{}
	}

	return result
}

type PostCreateBody struct {
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
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
	query := `INSERT INTO Posts (authorId, subcategoryId, title, content) VALUES (?, ?, ?, ?);`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(body.UserId, subcategoryId, body.Title, body.Content)
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
