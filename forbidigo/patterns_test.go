package forbidigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseValidPatterns(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		ptrn                 string
		expectedComment      string
		expectedMatchPackage bool
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
			name:                 "match package with non-empty group",
			ptrn:                 `^(?P<pkg>fmt).Println$`,
			expectedMatchPackage: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ptrn, err := parse(tc.ptrn)
			require.Nil(t, err)
			assert.Equal(t, tc.ptrn, ptrn.pattern.String())
			assert.Equal(t, tc.expectedComment, ptrn.msg)
			assert.Equal(t, tc.expectedMatchPackage, ptrn.matchPackage, "match pattern")
		})
	}
}

func TestParseInvalidPattern_ReturnsError(t *testing.T) {
	_, err := parse(`fmt\`)
	assert.NotNil(t, err)
}
