package categories

import (
	"encoding/json"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"log"
	"net/http"
)

type CategoriesWithSubcategories struct {
	Types.Category
	Subcategories []Types.Subcategory `json:"subcategories"`
}

func GetCategoriesWithSubcategoriesHandler(w http.ResponseWriter, r *http.Request) []CategoriesWithSubcategories {
	query := `SELECT 
    json_object(
				'id', c.id,
        'name', c.name,
        'icon', c.icon,
        'subcategories', json_group_array(
            json_object(
                'id', s.id,
                'name', s.name,
                'description', s.description,
								'categoryId', s.categoryId
            )
        )
    ) AS category_with_subcategories
FROM 
    Categories c
LEFT JOIN 
    Subcategories s ON c.id = s.categoryId
GROUP BY 
    c.id;
`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	var results []CategoriesWithSubcategories

	for rows.Next() {
		var jsonString string
		err := rows.Scan(&jsonString)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return nil
		}

		var result CategoriesWithSubcategories
		errJsonUnmarshal := json.Unmarshal([]byte(jsonString), &result)
		if errJsonUnmarshal != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error unmarshalling JSON:", errJsonUnmarshal)
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
