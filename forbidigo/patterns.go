package forbidigo

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
)

type pattern struct {
	pattern      *regexp.Regexp
	msg          string
	matchPackage bool
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
	matchPackage := false
	for _, groupName := range ptrnRe.SubexpNames() {
		switch groupName {
		case "pkg":
			matchPackage = true
		}
	}
	return &pattern{pattern: ptrnRe, msg: msg, matchPackage: matchPackage}, nil
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
