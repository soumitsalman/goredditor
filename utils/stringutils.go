package utils

import "regexp"

//Incase if we have to do some special formatting to remove whitespace charaters
func RemoveWhiteSpace(str string) string {
	white_space := regexp.MustCompile(`\s`)
	return white_space.ReplaceAllString(str, "")
}

func HasWhiteSpace(str string) bool {
	white_space := regexp.MustCompile(`\s`)
	return white_space.MatchString(str)
}
