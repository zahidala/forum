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

	// Get posts from subcategory
	query := `SELECT json_object(
		'id', p.id,
		'title', p.title,
		'content', p.content,
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

	rows, err := db.GetDB().Query(query, subcategoryId)
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
