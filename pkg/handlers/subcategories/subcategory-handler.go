package subcategories

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"log"
	"net/http"
)

type SubcategoryWithCategory struct {
	Types.Subcategory
	Category Types.Category `json:"category"`
}

func GetSubcategoryWithCategoryHandler(w http.ResponseWriter, r *http.Request) SubcategoryWithCategory {
	subcategoryId := r.PathValue("id")

	query := `SELECT json_object(
										'id', s.id,
										'name', s.name,
										'description', s.description,
										'categoryId', s.categoryId,
										'category', json_object(
											'id', c.id,
											'name', c.name,
											'description', c.description,
											'icon', c.icon
										)
									) AS subcategory_with_category
							FROM Subcategories s
							LEFT JOIN Categories c ON s.categoryId = c.id
							WHERE s.id = ?;
	`

	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing query:", err)
		return SubcategoryWithCategory{}
	}
	defer stmt.Close()

	rows, err := stmt.Query(subcategoryId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return SubcategoryWithCategory{}
	}
	defer rows.Close()

	var result SubcategoryWithCategory

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning rows:", err)
			return SubcategoryWithCategory{}
		}

		err = json.Unmarshal([]byte(jsonString), &result)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshaling json:", err)
			return SubcategoryWithCategory{}
		}
	}

	return result
}
