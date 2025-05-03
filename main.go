package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ethan/seo-score-analyzer/analyzer"
	"github.com/gorilla/mux"
)

type AnalysisRequest struct {
	URL         string   `json:"url"`
	Content     string   `json:"content"`
	Keywords    []string `json:"keywords"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

type SEOScore struct {
	URL               string  `json:"url"`
	KeywordScore      float64 `json:"keyword_score"`
	MetaDescription   float64 `json:"meta_description_score"`
	TitleTag          float64 `json:"title_tag_score"`
	Readability       float64 `json:"readability_score"`
	InternalLinks     float64 `json:"internal_links_score"`
	ExternalLinks     float64 `json:"external_links_score"`
	ImageAltText      float64 `json:"image_alt_text_score"`
	SocialTags        float64 `json:"social_tags_score"`
	OverallScore      float64 `json:"overall_score"`
}

func analyzeSEO(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	seoAnalyzer := analyzer.NewAnalyzer(req.Content, req.URL)
	seoAnalyzer.SetKeywords(req.Keywords)
	seoAnalyzer.SetTitle(req.Title)
	seoAnalyzer.SetDescription(req.Description)

	internalScore, externalScore := seoAnalyzer.CalculateLinkScores()

	score := SEOScore{
		URL:               req.URL,
		KeywordScore:      seoAnalyzer.CalculateKeywordScore(),
		MetaDescription:   seoAnalyzer.CalculateMetaDescriptionScore(),
		TitleTag:          seoAnalyzer.CalculateTitleTagScore(),
		Readability:       seoAnalyzer.CalculateReadabilityScore(),
		InternalLinks:     internalScore,
		ExternalLinks:     externalScore,
		ImageAltText:      seoAnalyzer.CalculateImageAltTextScore(),
		SocialTags:        seoAnalyzer.CalculateSocialTagsScore(),
		OverallScore:      seoAnalyzer.CalculateOverallScore(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(score)
}

func main() {
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/analyze", analyzeSEO).Methods("POST")

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
} 