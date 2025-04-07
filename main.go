package handler

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
)

var (
	white = color.White
	black = color.Gray{Y: 34}
)

func generateDataJson(shields []Shield) error {
	data, err := os.OpenFile("docs/data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer data.Close()

	return json.NewEncoder(data).Encode(shields)
}

// Handler is the exported function that Vercel uses as the entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	title := query.Get("title")
	hex := query.Get("hex")

	if title == "" || hex == "" {
		http.Error(w, "Missing 'title' or 'hex' query parameter", http.StatusBadRequest)
		return
	}

	icon := Icon{Title: title, Hex: hex}
	shield, err := IconToShield(icon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shield)
}

func main() {
	icons, err := getIcons()
	if (err != nil) {
		log.Panicln(err)
	}

	shields := make([]Shield, len(icons))
	for i, icons := range icons {
		shield, err := IconToShield(icons)
		if err != nil {
			log.Panicln(err)
		}
		shields[i] = *shield
	}

	generateDataJson(shields)

	readme, err := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer readme.Close()

	snippets, err := os.OpenFile("Snippets.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer snippets.Close()

	for _, shield := range shields {
		fmt.Fprintln(readme, shield)
		fmt.Fprintf(snippets, "## %[1]s\n```markdown\n%[1]s\n```\n", shield)
	}

	http.HandleFunc("/api/shield", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		title := query.Get("title")
		hex := query.Get("hex")

		if title == "" || hex == "" {
			http.Error(w, "Missing 'title' or 'hex' query parameter", http.StatusBadRequest)
			return
		}

		icon := Icon{Title: title, Hex: hex}
		shield, err := IconToShield(icon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shield)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
