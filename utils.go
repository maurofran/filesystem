package filesystem

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// normalizeDirname will normalize directory name value.
func normalizeDirname(dirname string) string {
	if dirname == "." {
		return ""
	}
	return dirname
}

// normalizePath will normalize a path.
func normalizePath(path Path) (Path, error) {
	return normalizeRelativePath(path)
}

// Normalizes relative directories in a path.
func normalizeRelativePath(path Path) (Path, error) {
	path = Path(strings.Replace(string(path), "\\", "/", -1))
	path = removeFunkyWhiteSpace(path)

	pathParts := strings.Split(string(path), "/")
	parts := make([]string, 0, len(pathParts))

	for _, part := range pathParts {
		switch part {
		case "":
		case ".":
		case "..":
			if len(parts) == 0 {
				return emptyPath, fmt.Errorf("Path is outside of defined root, path: %s", path)
			}
			parts = parts[:len(parts)-1]
		default:
			parts = append(parts, part)
		}
	}
	return Path(strings.Join(parts, "/")), nil
}

var whiteSpaceRegexp = regexp.MustCompile(`#\p{C}+|^\./#u`)

// Remove unprintable characters and invalid unicode characters.
func removeFunkyWhiteSpace(path Path) Path {
	p := string(path)
	for whiteSpaceRegexp.MatchString(p) {
		p = whiteSpaceRegexp.ReplaceAllString(p, "")
	}
	return Path(p)
}

// Normalize a prefix.
func normalizePrefix(prefix, separator string) string {
	return strings.TrimRightFunc(prefix+separator, unicode.IsSpace) + separator
}
