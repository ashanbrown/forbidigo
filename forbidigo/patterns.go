package forbidigo

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
)

type pattern struct {
	pattern *regexp.Regexp
	msg     string

	// matchWithPackage is set for rules against selector expressions where
	// the selector name gets replaced by the full package name (for
	// imports) or the type including the package (for variables) before
	// checking for a match.
	matchWithPackage bool
}

func parse(ptrn string) (*pattern, error) {
	ptrnRe, err := regexp.Compile(ptrn)
	if err != nil {
		return nil, fmt.Errorf("unable to compile pattern `%s`: %s", ptrn, err)
	}
	re, err := syntax.Parse(ptrn, syntax.Perl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pattern `%s`: %s", ptrn, err)
	}
	msg := extractComment(re)
	matchWithPackage := false
	for _, groupName := range ptrnRe.SubexpNames() {
		switch groupName {
		case "pkg":
			matchWithPackage = true
		}
	}
	return &pattern{pattern: ptrnRe, msg: msg, matchWithPackage: matchWithPackage}, nil
}

// Traverse the leaf submatches in the regex tree and extract a comment, if any
// is present.
func extractComment(re *syntax.Regexp) string {
	for _, sub := range re.Sub {
		subStr := sub.String()
		if strings.HasPrefix(subStr, "#") {
			return strings.TrimSpace(strings.TrimPrefix(sub.String(), "#"))
		}
		if len(sub.Sub) > 0 {
			if comment := extractComment(sub); comment != "" {
				return comment
			}
		}
	}
	return ""
}
