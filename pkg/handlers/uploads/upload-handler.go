package uploads

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	Types "forum/pkg/types"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) Types.UploadedImage {

	const maxUploadSize = 20 * 1024 * 1024 // 20MB

	err := r.ParseMultipartForm(maxUploadSize + 1)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		log.Println(err)
		return Types.UploadedImage{}
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)

	var b bytes.Buffer
	formData := multipart.NewWriter(&b)

	fw, err := formData.CreateFormField("source")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}

	if _, err = fw.Write([]byte(encoded)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}
	formData.Close()

	apiKey, exists := os.LookupEnv("FREEIMAGEHOST_API_KEY")

	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("freeimage.host API key not found")
		return Types.UploadedImage{}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://freeimage.host/api/1/upload?key=%s", apiKey), &b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}

	req.Header.Set("Content-Type", formData.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return Types.UploadedImage{}
	}

	var uploadedImage Types.UploadedImage

	jsonUnmarshalErr := json.Unmarshal(respBody, &uploadedImage)

	if jsonUnmarshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(jsonUnmarshalErr)
		return Types.UploadedImage{}
	}

	return uploadedImage
}
