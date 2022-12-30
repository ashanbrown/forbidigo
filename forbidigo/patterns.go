package forbidigo

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	MatchText = "text"
	MatchType = "type"
)

// Pattern matches code that is not supposed to be used.
type Pattern struct {
	re *regexp.Regexp

	// Pattern is the regular expression string that is used for matching.
	Pattern string `yaml:"Pattern"`

	// Match defines whether the regular expression is matched against the
	// source code literally ("text", the default) or whether type
	// information is used to determine what is being referenced ("type").
	Match string `yaml:"Match"`

	// Msg gets printed in addition to the normal message if a match is
	// found.
	Msg string `yaml:"Msg"`
}

// parse accepts a regular expression or, if the string starts with {, a
// JSON or YAML representation of a Pattern.
func parse(ptrn string) (*Pattern, error) {
	pattern := &Pattern{}

	if strings.HasPrefix(strings.TrimSpace(ptrn), "{") {
		// Embedded JSON or YAML. We can decode both with the YAML decoder.
		if err := yaml.UnmarshalStrict([]byte(ptrn), pattern); err != nil {
			return nil, fmt.Errorf("parsing as JSON or YAML failed: %v", err)
		}
		ptrn = pattern.Pattern
	}

	ptrnRe, err := regexp.Compile(ptrn)
	if err != nil {
		return nil, fmt.Errorf("unable to compile pattern `%s`: %s", ptrn, err)
	}
	re, err := syntax.Parse(ptrn, syntax.Perl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pattern `%s`: %s", ptrn, err)
	}
	msg := extractComment(re)
	if msg != "" {
		pattern.Msg = msg
	}
	pattern.re = ptrnRe
	if pattern.Match == "" {
		pattern.Match = MatchText
	}

	switch pattern.Match {
	case MatchText, MatchType:
		// okay
	default:
		return nil, fmt.Errorf("unsupported match string: %q", pattern.Match)
	}

	return pattern, nil
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
