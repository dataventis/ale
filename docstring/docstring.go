package docstring

import (
	"fmt"

	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/internal/maps"
	"github.com/kode4food/ale/internal/markdown"
)

// Error messages
const (
	ErrSymbolNotDocumented = "symbol not documented: %s"
)

const extension = ".md"

var (
	docStringCache     = map[string][]byte{}
	docStringCacheOnce = do.Once()
)

// Get resolves a registered docstring entry by name
func Get(n string) (string, error) {
	ensureDocStringCache()
	res, ok := docStringCache[n]
	if ok {
		return string(res), nil
	}
	return "", fmt.Errorf(ErrSymbolNotDocumented, n)
}

// MustGet resolves a registered docstring entry by name or explodes
func MustGet(n string) string {
	res, err := Get(n)
	if err != nil {
		panic(err)
	}
	return res
}

// Names returns the registered names of available docstring entries
func Names() []string {
	ensureDocStringCache()
	return maps.SortedKeys(docStringCache)
}

func ensureDocStringCache() {
	docStringCacheOnce(func() {
		for _, filename := range assetNames() {
			doc, _ := getAsset(filename)
			meta := markdown.ParseHeader(string(doc))
			if names := meta.Names; len(names) > 0 {
				for _, name := range names {
					docStringCache[name] = doc
				}
			} else {
				n := filename[0 : len(filename)-len(extension)]
				docStringCache[n] = doc
			}
		}
	})
}
