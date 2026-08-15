package bibliography

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/lsdch/biome/types"
)

var (
	authorYearRegex = regexp.MustCompile(`(?s)^(.*?)(?:\(?\d{4}\)?.*)$`)
	yearRegex       = regexp.MustCompile(`\(?(\d{4})\)?`)
	doiRegex        = regexp.MustCompile(`10\.\d{4,9}/[-._;()/:A-Za-z0-9]+`)
)

func ExtractAuthorString(verbatim string) string {
	m := authorYearRegex.FindStringSubmatch(verbatim)
	if len(m) < 2 {
		return ""
	}

	return strings.Trim(strings.TrimSpace(m[1]), ".")
}

func ExtractYear(verbatim string) *int32 {
	m := yearRegex.FindStringSubmatch(verbatim)
	if len(m) < 2 {
		return nil
	}

	year, err := strconv.Atoi(m[1])
	if err != nil {
		return nil
	}

	y := int32(year)
	return &y
}

func ExtractDOI(verbatim string) *types.DOI {
	doiStr := doiRegex.FindString(verbatim)
	if doiStr == "" {
		return nil
	}

	doi := types.ParseDOI(strings.TrimRight(doiStr, ".,;:"))
	return &doi
}

var (
	initialsPattern = regexp.MustCompile(`^[A-Z](?:\.[A-Z])?\.?$`)
)

// ParseAuthors transforms a raw author string into a slice of individual author names.
// It handles various formats and separators, including commas, semicolons, and conjunctions like "and" or "&".
//
// The function also normalizes spaces and trims unnecessary punctuation from each author's name.
//
// Examples :
//
//	"Smith J., Doe A."              -> ["Smith J.", "Doe A."]
//	"Smith, J., Doe, A."            -> ["Smith, J.", "Doe, A."]
//	"Smith J.; Doe A."              -> ["Smith J.", "Doe A."]
//	"Smith J. and Doe A."           -> ["Smith J.", "Doe A."]
func ParseAuthors(raw string) []string {
	raw = normalizeSpaces(raw)

	if raw == "" {
		return nil
	}

	// Non ambiguous separators: semicolon, "and", "&"
	for _, sep := range []string{";", " and ", " & "} {
		if strings.Contains(strings.ToLower(raw), sep) {
			return cleanAuthors(strings.Split(raw, sep))
		}
	}

	// Case "Name, I., Name, I."
	if authors := parseCommaSeparatedAuthors(raw); len(authors) > 0 {
		return authors
	}

	// Case "Name I., Name I."
	return cleanAuthors(strings.Split(raw, ","))
}

func parseCommaSeparatedAuthors(raw string) []string {
	parts := strings.Split(raw, ",")

	if len(parts) < 2 {
		return nil
	}

	var authors []string
	var current strings.Builder

	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])

		if current.Len() == 0 {
			current.WriteString(part)
			continue
		}

		// If the next part is an initial,
		// it belongs to the previous name.
		if initialsPattern.MatchString(part) {
			current.WriteString(", ")
			current.WriteString(part)
			continue
		}

		authors = append(authors, current.String())
		current.Reset()
		current.WriteString(part)
	}

	if current.Len() > 0 {
		authors = append(authors, current.String())
	}

	if len(authors) <= 1 {
		return nil
	}

	return cleanAuthors(authors)
}

func cleanAuthors(authors []string) []string {
	result := make([]string, 0, len(authors))

	for _, author := range authors {
		author = strings.Trim(author, " ,;")
		author = normalizeSpaces(author)

		if author != "" {
			result = append(result, author)
		}
	}

	return result
}

func normalizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
