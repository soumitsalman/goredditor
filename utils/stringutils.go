package utils

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Incase if we have to do some special formatting to remove whitespace charaters
func RemoveWhiteSpace(str string) string {
	white_space := regexp.MustCompile(`\s`)
	return white_space.ReplaceAllString(str, "")
}

func HasWhiteSpace(str string) bool {
	white_space := regexp.MustCompile(`\s`)
	return white_space.MatchString(str)
}

func ExtractTextFromHtml(content string) string {
	//there needs to be multiple runs on the NewDocumentFromReader when '<' and '>' are represented as "&lt;' and '&gt;'
	for count := 3; count > 0; count-- {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
		content = doc.Text()
	}
	return removeMultipleNewLines(content)
}

func removeMultipleNewLines(str string) string {
	white_space := regexp.MustCompile(`[\n]+`)
	return white_space.ReplaceAllString(str, "\n")
}
