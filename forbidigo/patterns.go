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

// A YAMLPattern pattern in a YAML string may be represented either by a string
// (the traditional regular expression syntax) or a struct (for more complex
// patterns).
type YAMLPattern Pattern

func (p *YAMLPattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try struct first. It's unlikely that a regular expression string
	// is valid YAML for a struct.
	var pattern Pattern
	if err := unmarshal(&pattern); err != nil {
		errStr := err.Error()
		// Didn't work, try plain string.
		var ptrn string
		if err := unmarshal(&ptrn); err != nil {
			return fmt.Errorf("pattern is neither a regular expression string (%s) nor a Pattern struct (%s)", err.Error(), errStr)
		}
		p.Pattern = ptrn
	} else {
		*p = YAMLPattern(pattern)
	}
	return ((*Pattern)(p)).validate(p.Pattern)
}

var _ yaml.Unmarshaler = &YAMLPattern{}

// parse accepts a regular expression or, if the string starts with { or contains a line break, a
// JSON or YAML representation of a Pattern.
func parse(ptrn string) (*Pattern, error) {
	pattern := &Pattern{}

	if strings.HasPrefix(strings.TrimSpace(ptrn), "{") ||
		strings.Contains(ptrn, "\n") {
		// Embedded JSON or YAML. We can decode both with the YAML decoder.
		if err := yaml.UnmarshalStrict([]byte(ptrn), pattern); err != nil {
			return nil, fmt.Errorf("parsing as JSON or YAML failed: %v", err)
		}
		ptrn = pattern.Pattern
	}

	if err := pattern.validate(ptrn); err != nil {
		return nil, err
	}
	return pattern, nil
}

func (p *Pattern) validate(ptrn string) error {
	ptrnRe, err := regexp.Compile(ptrn)
	if err != nil {
		return fmt.Errorf("unable to compile pattern `%s`: %s", ptrn, err)
	}
	re, err := syntax.Parse(ptrn, syntax.Perl)
	if err != nil {
		return fmt.Errorf("unable to parse pattern `%s`: %s", ptrn, err)
	}
	msg := extractComment(re)
	if msg != "" {
		p.Msg = msg
	}
	p.re = ptrnRe

	if p.Match == "" {
		p.Match = MatchText
	}
	switch p.Match {
	case MatchText, MatchType:
		// okay
	default:
		return fmt.Errorf("unsupported match string: %q", p.Match)
	}

	return nil
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
