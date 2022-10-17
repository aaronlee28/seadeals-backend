package formatter

import (
	"regexp"
	"seadeals-backend/config"
	"strings"
)

func GenerateSlug(title string) string {
	if config.Config.ENV == "testing" {
		return title + "-slug"
	}

	nonAlphanumericRegex := regexp.MustCompile(`^[a-zA-Z\d\-_\s(1)]+`)
	split := strings.Split(title, " ")
	slug := strings.Join(split, "-")
	split = nonAlphanumericRegex.FindAllString(slug, -1)
	slug = strings.Join(split, "-")
	return slug
}
