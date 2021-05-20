package predict

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Dirs returns a predictor that predict directory paths. If a non-empty pattern is given, the
// predicted paths will match that pattern.
func Dirs(pattern string) FilesPredictor {
	return FilesPredictor{pattern: pattern, includeFiles: false}
}

// Dirs returns a predictor that predict file or directory paths. If a non-empty pattern is given,
// the predicted paths will match that pattern.
func Files(pattern string) FilesPredictor {
	return FilesPredictor{pattern: pattern, includeFiles: true}
}

type FilesPredictor struct {
	pattern      string
	includeFiles bool
}

// Predict searches for files according to the given prefix.
// If the only predicted path is a single directory, the search will continue another recursive
// layer into that directory.
func (f FilesPredictor) Predict(prefix string) (options []string) {
	options = f.predictFiles(prefix)

	// If the number of prediction is not 1, we either have many results or have no results, so we
	// return it.
	if len(options) != 1 {
		return
	}

	// Only try deeper, if the one item is a directory.
	if stat, err := os.Stat(options[0]); err != nil || !stat.IsDir() {
		return
	}

	return f.predictFiles(options[0])
}

func (f FilesPredictor) predictFiles(prefix string) []string {
	if strings.HasSuffix(prefix, "/..") {
		return nil
	}

	dir := directory(prefix)
	files := f.listFiles(dir)

	// Add dir if match.
	files = append(files, dir)

	return FilesSet(files).Predict(prefix)
}

func (f FilesPredictor) listFiles(dir string) []string {
	// Set of all file names.
	m := map[string]bool{}

	// List files.
	if files, err := filepath.Glob(filepath.Join(dir, f.pattern)); err == nil {
		for _, file := range files {
			if stat, err := os.Stat(file); err != nil || stat.IsDir() || f.includeFiles {
				m[file] = true
			}
		}
	}

	// List directories.
	if dirs, err := ioutil.ReadDir(dir); err == nil {
		for _, d := range dirs {
			if d.IsDir() {
				m[filepath.Join(dir, d.Name())] = true
			}
		}
	}

	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	return list
}

// directory gives the directory of the given partial path in case that it is not, we fall back to
// the current directory.
func directory(path string) string {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return fixPathForm(path, path)
	}
	dir := filepath.Dir(path)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return fixPathForm(path, dir)
	}
	return "./"
}

// FilesSet predict according to file rules to a given fixed set of file names.
type FilesSet []string

func (s FilesSet) Predict(prefix string) (prediction []string) {
	// add all matching files to prediction
	for _, f := range s {
		f = fixPathForm(prefix, f)

		// test matching of file to the argument
		if matchFile(f, prefix) {
			prediction = append(prediction, f)
		}
	}
	if len(prediction) == 0 {
		return s
	}
	return
}

// MatchFile returns true if prefix can match the file
func matchFile(file, prefix string) bool {
	// special case for current directory completion
	if file == "./" && (prefix == "." || prefix == "") {
		return true
	}
	if prefix == "." && strings.HasPrefix(file, ".") {
		return true
	}

	file = strings.TrimPrefix(file, "./")
	prefix = strings.TrimPrefix(prefix, "./")

	return strings.HasPrefix(file, prefix)
}

// fixPathForm changes a file name to a relative name
func fixPathForm(last string, file string) string {
	// Get wording directory for relative name.
	workDir, err := os.Getwd()
	if err != nil {
		return file
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return file
	}

	// If last is absolute, return path as absolute.
	if filepath.IsAbs(last) {
		return fixDirPath(abs)
	}

	rel, err := filepath.Rel(workDir, abs)
	if err != nil {
		return file
	}

	// Fix ./ prefix of path.
	if rel != "." && strings.HasPrefix(last, ".") {
		rel = "./" + rel
	}

	return fixDirPath(rel)
}

func fixDirPath(path string) string {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}
