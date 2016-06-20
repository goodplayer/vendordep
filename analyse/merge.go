package analyse

import (
	"strings"
)

func MergeUnimportedUrlPaths(allImports []string, existImportRootPaths []string) []string {
	l := make([]string, len(allImports))
	i := 0
	for _, v := range allImports {
		if !isUrlImport(v) {
			continue
		}
		notExist := true
		for _, exist := range existImportRootPaths {
			if strings.HasPrefix(v, exist) {
				notExist = false
				break
			}
		}
		if notExist {
			l[i] = v
			i++
		}
	}
	return l[:i]
}

func isUrlImport(path string) bool {
	// standard library path does not has '.', may be.
	paths := strings.Split(path, "/")
	return strings.Contains(paths[0], ".")
}
