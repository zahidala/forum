package posts

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"log"
	"net/http"
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

type PostWithComments struct {
	Types.Post
	Author      Types.User        `json:"author"`
	Subcategory Types.Subcategory `json:"subcategory"`
	Comments    []Types.Comment   `json:"comments"`
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
		'comments', json_group_array(
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
						'attachments', c.attachments
				)
		)
) AS post

FROM Posts p

LEFT JOIN Users u ON p.authorId = u.id
LEFT JOIN Subcategories s ON p.subcategoryId = s.id
LEFT JOIN Categories c2 ON s.categoryId = c2.id
LEFT JOIN Comments c ON p.id = c.postId
LEFT JOIN Users u2 ON c.authorId = u2.id

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
