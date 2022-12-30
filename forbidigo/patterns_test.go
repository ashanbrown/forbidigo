package forbidigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
