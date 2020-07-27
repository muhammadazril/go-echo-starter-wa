package bootstrap

import (
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
)

// InitMinifier to add helper to bootstrap
func InitMinifier() *minify.M {
	M := minify.New()
	// currently activated only for json minifier
	M.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	return M
}
