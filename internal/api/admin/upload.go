package admin

import (
	"encoding/json"
	"net/http"

	"github.com/henryarin/portfolio-backend-go/internal/storage"
)

// UploadImage handles standalone image uploads (for markdown editor)
func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "no image provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := storage.SaveImage(file, header)
	if err != nil {
		http.Error(w, "failed to save image: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": imageURL,
	})
}
