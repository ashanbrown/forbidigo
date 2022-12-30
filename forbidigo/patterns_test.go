package forbidigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestParseValidPatterns(t *testing.T) {
	for _, tc := range []struct {
		name            string
		ptrn            string
		expectedComment string
		expectedMatch   string
		expectedPattern string
	}{
		{
			name: "simple expression, no comment",
			ptrn: `fmt\.Errorf`,
		},
		{
			name: "anchored expression, no comment",
			ptrn: `^fmt\.Errorf$`,
		},
		{
			name:            "contains multiple subexpression, with comment",
			ptrn:            `(f)mt\.Errorf(# a comment)?`,
			expectedComment: "a comment",
		},
		{
			name:            "simple expression with comment",
			ptrn:            `fmt\.Println(# Please don't use this!)?`,
			expectedComment: "Please don't use this!",
		},
		{
			name:            "deeply nested expression with comment",
			ptrn:            `fmt\.Println((((# Please don't use this!))))?`,
			expectedComment: "Please don't use this!",
		},
		{
			name:            "anchored expression with comment",
			ptrn:            `^fmt\.Println(# Please don't use this!)?$`,
			expectedComment: "Please don't use this!",
		},
		{
			name:            "match import",
			ptrn:            `{Match: "type", Pattern: "^fmt\\.Println$"}`,
			expectedMatch:   MatchType,
			expectedPattern: `^fmt\.Println$`,
		},
		{
			name: "match import with YAML",
			ptrn: `{Match: type,
Pattern: ^fmt\.Println$
}`,
			expectedMatch:   MatchType,
			expectedPattern: `^fmt\.Println$`,
		},
		{
			name:            "match import with YAML, no line breaks",
			ptrn:            `{Match: type, Pattern: ^fmt\.Println$}`,
			expectedMatch:   MatchType,
			expectedPattern: `^fmt\.Println$`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ptrn, err := parse(tc.ptrn)
			require.Nil(t, err)
			expectedPattern := tc.expectedPattern
			if expectedPattern == "" {
				expectedPattern = tc.ptrn
			}
			assert.Equal(t, expectedPattern, ptrn.re.String(), "pattern")
			assert.Equal(t, tc.expectedComment, ptrn.Msg, "comment")
			expectedMatch := tc.expectedMatch
			if expectedMatch == "" {
				expectedMatch = MatchText
			}
			assert.Equal(t, expectedMatch, ptrn.Match, "match ")
		})
	}
}

func TestParseInvalidPattern_ReturnsError(t *testing.T) {
	_, err := parse(`fmt\`)
	assert.NotNil(t, err)
}

func TestUnmarshalYAML(t *testing.T) {
	for _, tc := range []struct {
		name            string
		yaml            string
		expectedErr     string
		expectedComment string
		expectedMatch   string
		expectedPattern string
	}{
		{
			name: "string: simple expression, no comment",
			yaml: `fmt\.Errorf`,
		},
		{
			name:            "string: contains multiple subexpression, with comment",
			yaml:            `(f)mt\.Errorf(# a comment)?`,
			expectedComment: "a comment",
		},
		{
			name:            "struct: simple expression, no comment",
			yaml:            `Pattern: fmt\.Errorf`,
			expectedPattern: `fmt\.Errorf`,
		},
		{
			name: "match import with YAML",
			yaml: `Match: type
Pattern: ^fmt\.Println$
`,
			expectedMatch:   MatchType,
			expectedPattern: `^fmt\.Println$`,
		},
		{
			name:        "string: invalid regexp",
			yaml:        `fmt\`,
			expectedErr: "unable to compile pattern `fmt\\`: error parsing regexp: trailing backslash at end of expression: ``",
		},
		{
			name:        "stuct: invalid regexp",
			yaml:        `Pattern: fmt\`,
			expectedErr: "unable to compile pattern `fmt\\`: error parsing regexp: trailing backslash at end of expression: ``",
		},
		{
			name:        "invalid match",
			yaml:        `Match: true`,
			expectedErr: `unsupported match string: "true"`,
		},
		{
			name: "invalid struct",
			yaml: `Foo: bar`,
			expectedErr: `pattern is neither a regular expression string (yaml: unmarshal errors:
  line 1: cannot unmarshal !!map into string) nor a Pattern struct (yaml: unmarshal errors:
  line 1: field Foo not found in type forbidigo.Pattern)`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var p YAMLPattern
			err := yaml.UnmarshalStrict([]byte(tc.yaml), &p)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err.Error())
				return
			}
			expectedPattern := tc.expectedPattern
			if expectedPattern == "" {
				expectedPattern = tc.yaml
			}
			assert.Equal(t, expectedPattern, p.re.String(), "pattern")
			assert.Equal(t, tc.expectedComment, p.Msg, "comment")
			expectedMatch := tc.expectedMatch
			if expectedMatch == "" {
				expectedMatch = MatchText
			}
			assert.Equal(t, expectedMatch, p.Match, "match ")
		})
	}
}
