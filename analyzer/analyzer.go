package analyzer

import (
	"math"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Analyzer struct {
	Content     string
	URL         string
	Keywords    []string
	Title       string
	Description string
}

func NewAnalyzer(content, url string) *Analyzer {
	return &Analyzer{
		Content: content,
		URL:     url,
	}
}

func (a *Analyzer) SetKeywords(keywords []string) {
	a.Keywords = keywords
}

func (a *Analyzer) SetTitle(title string) {
	a.Title = title
}

func (a *Analyzer) SetDescription(description string) {
	a.Description = description
}

func (a *Analyzer) CalculateKeywordScore() float64 {
	if len(a.Keywords) == 0 || len(a.Content) == 0 {
		return 0.0
	}

	totalScore := 0.0
	contentWords := strings.Fields(strings.ToLower(a.Content))
	totalWords := len(contentWords)

	for _, keyword := range a.Keywords {
		keyword = strings.ToLower(keyword)
		keywordCount := 0
		for _, word := range contentWords {
			if word == keyword {
				keywordCount++
			}
		}
		keywordDensity := float64(keywordCount) / float64(totalWords)
		score := math.Min(keywordDensity*100, 100) // Cap at 100%
		totalScore += score
	}

	return totalScore / float64(len(a.Keywords))
}

func (a *Analyzer) CalculateMetaDescriptionScore() float64 {
	if a.Description == "" {
		return 0.0
	}

	length := len(a.Description)
	score := 0.0

	// Ideal length is between 120-160 characters
	if length >= 120 && length <= 160 {
		score += 50
	} else if length > 160 {
		score += math.Max(0, 50-(float64(length-160)*0.5))
	} else {
		score += math.Max(0, 50-(float64(120-length)*0.5))
	}

	// Check for keyword presence
	if len(a.Keywords) > 0 {
		description := strings.ToLower(a.Description)
		for _, keyword := range a.Keywords {
			if strings.Contains(description, strings.ToLower(keyword)) {
				score += 50 / float64(len(a.Keywords))
			}
		}
	}

	return math.Min(score, 100)
}

func (a *Analyzer) CalculateTitleTagScore() float64 {
	if a.Title == "" {
		return 0.0
	}

	score := 0.0
	length := len(a.Title)

	// Ideal length is between 50-60 characters
	if length >= 50 && length <= 60 {
		score += 50
	} else if length > 60 {
		score += math.Max(0, 50-(float64(length-60)*0.5))
	} else {
		score += math.Max(0, 50-(float64(50-length)*0.5))
	}

	// Check for keyword presence
	if len(a.Keywords) > 0 {
		title := strings.ToLower(a.Title)
		for _, keyword := range a.Keywords {
			if strings.Contains(title, strings.ToLower(keyword)) {
				score += 50 / float64(len(a.Keywords))
			}
		}
	}

	return math.Min(score, 100)
}

func (a *Analyzer) CalculateReadabilityScore() float64 {
	if len(a.Content) == 0 {
		return 0.0
	}

	// Simple implementation of Flesch Reading Ease
	words := strings.Fields(a.Content)
	sentences := strings.Split(a.Content, ".")
	
	// Remove empty sentences
	var validSentences []string
	for _, s := range sentences {
		if strings.TrimSpace(s) != "" {
			validSentences = append(validSentences, s)
		}
	}

	if len(validSentences) == 0 || len(words) == 0 {
		return 0.0
	}

	// Count syllables (simplified)
	syllables := 0
	for _, word := range words {
		syllables += countSyllables(word)
	}

	// Flesch Reading Ease formula
	score := 206.835 - 1.015*(float64(len(words))/float64(len(validSentences))) - 84.6*(float64(syllables)/float64(len(words)))

	// Normalize score to 0-100 range
	normalizedScore := (score + 60) / 1.6
	return math.Max(0, math.Min(100, normalizedScore))
}

func countSyllables(word string) int {
	// Simplified syllable counting
	word = strings.ToLower(word)
	vowels := "aeiouy"
	count := 0
	prevChar := ' '

	for _, char := range word {
		if strings.ContainsRune(vowels, char) && !strings.ContainsRune(vowels, prevChar) {
			count++
		}
		prevChar = char
	}

	if count == 0 {
		count = 1
	}

	return count
}

func (a *Analyzer) CalculateLinkScores() (float64, float64) {
	doc, err := html.Parse(strings.NewReader(a.Content))
	if err != nil {
		return 0.0, 0.0
	}

	internalLinks := 0
	externalLinks := 0
	totalLinks := 0

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					totalLinks++
					linkURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					baseURL, err := url.Parse(a.URL)
					if err != nil {
						continue
					}

					if linkURL.Host == "" || linkURL.Host == baseURL.Host {
						internalLinks++
					} else {
						externalLinks++
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Score based on link presence and distribution
	internalScore := 0.0
	externalScore := 0.0

	if totalLinks > 0 {
		internalScore = math.Min(100, float64(internalLinks)*10)
		externalScore = math.Min(100, float64(externalLinks)*10)
	}

	return internalScore, externalScore
}

func (a *Analyzer) CalculateImageAltTextScore() float64 {
	doc, err := html.Parse(strings.NewReader(a.Content))
	if err != nil {
		return 0.0
	}

	totalImages := 0
	imagesWithAlt := 0

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			totalImages++
			for _, attr := range n.Attr {
				if attr.Key == "alt" && attr.Val != "" {
					imagesWithAlt++
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if totalImages == 0 {
		return 100.0 // No images is considered perfect
	}

	return (float64(imagesWithAlt) / float64(totalImages)) * 100
}

func (a *Analyzer) CalculateSocialTagsScore() float64 {
	doc, err := html.Parse(strings.NewReader(a.Content))
	if err != nil {
		return 0.0
	}

	requiredTags := map[string]bool{
		"og:title":       false,
		"og:description": false,
		"og:image":       false,
		"twitter:title":  false,
		"twitter:description": false,
		"twitter:image":  false,
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var property, content string
			for _, attr := range n.Attr {
				if attr.Key == "property" || attr.Key == "name" {
					property = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if _, exists := requiredTags[property]; exists && content != "" {
				requiredTags[property] = true
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	score := 0.0
	for _, present := range requiredTags {
		if present {
			score += 100.0 / float64(len(requiredTags))
		}
	}

	return score
}

func (a *Analyzer) CalculateOverallScore() float64 {
	keywordScore := a.CalculateKeywordScore()
	metaScore := a.CalculateMetaDescriptionScore()
	titleScore := a.CalculateTitleTagScore()
	readabilityScore := a.CalculateReadabilityScore()
	internalScore, externalScore := a.CalculateLinkScores()
	imageScore := a.CalculateImageAltTextScore()
	socialScore := a.CalculateSocialTagsScore()

	// Weighted average
	overallScore := (keywordScore * 0.2) +
		(metaScore * 0.15) +
		(titleScore * 0.15) +
		(readabilityScore * 0.15) +
		((internalScore+externalScore)/2 * 0.15) +
		(imageScore * 0.1) +
		(socialScore * 0.1)

	return math.Round(overallScore*100) / 100
} 