package categories

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"log"
	"net/http"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) []Types.Category {
	query := `SELECT json_group_array(
								json_object(
										'id', id,
										'name', name,
										'description', description,
										'icon', icon
								)
						) AS categories
						FROM Categories;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	defer rows.Close()

	var results []Types.Category

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return nil
		}

		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &results)
		if errJsonUnmarshal != nil {
			log.Println("Error unmarshaling JSON:", errJsonUnmarshal)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return nil
		}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Println("Error iterating rows:", rowsErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}

	return results
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) Types.Category {
	categoryID := r.PathValue("id")

	query := `SELECT json_object(
								'id', id,
								'name', name,
								'description', description,
								'icon', icon
						) AS category
						FROM Categories
						WHERE id = ?;`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return Types.Category{}
	}
	defer stmt.Close()

	rows, err := stmt.Query(categoryID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return Types.Category{}
	}
	defer rows.Close()

	var category Types.Category

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return Types.Category{}
		}

		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &category)
		if errJsonUnmarshal != nil {
			log.Println("Error unmarshaling JSON:", errJsonUnmarshal)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return Types.Category{}
		}
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Println("Error iterating rows:", rowsErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return Types.Category{}
	}

	return category
}
