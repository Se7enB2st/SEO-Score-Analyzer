# SEO Score Analyzer

A Go-based web application that analyzes blog posts, articles, or web pages for search engine optimization (SEO) performance. It helps content creators identify areas for improvement to boost their rankings and visibility.

## Features

- **Keyword Optimization**: Checks for keyword presence, frequency, and density
- **Meta Description Analysis**: Validates the length and keyword usage in meta descriptions
- **Title Tag Evaluation**: Ensures the title is engaging and keyword-optimized
- **Readability Check**: Scores the content for readability using the Flesch Reading Ease algorithm
- **Internal and External Linking**: Verifies the presence and quality of internal and external links
- **Image Alt Text Verification**: Confirms all images have descriptive alt texts
- **Social Sharing Tags**: Checks for Open Graph and Twitter Card tags

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```
2. Open your browser and navigate to `http://localhost:8080`
3. Enter the URL, title, meta description, keywords, and content to analyze
4. Click "Analyze SEO" to get your SEO score and recommendations

## API Endpoints

### POST /api/analyze

Analyzes the provided content and returns SEO scores.

Request body:
```json
{
    "url": "https://example.com",
    "title": "Page Title",
    "description": "Meta description",
    "keywords": ["keyword1", "keyword2"],
    "content": "HTML content"
}
```

Response:
```json
{
    "url": "https://example.com",
    "keyword_score": 85.5,
    "meta_description_score": 90.0,
    "title_tag_score": 95.0,
    "readability_score": 75.5,
    "internal_links_score": 80.0,
    "external_links_score": 70.0,
    "image_alt_text_score": 100.0,
    "social_tags_score": 60.0,
    "overall_score": 82.0
}
```

## Scoring System

Each component is scored on a scale of 0-100, with the overall score being a weighted average:

- Keyword Score (20%)
- Meta Description Score (15%)
- Title Tag Score (15%)
- Readability Score (15%)
- Link Scores (15%)
- Image Alt Text Score (10%)
- Social Tags Score (10%)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.