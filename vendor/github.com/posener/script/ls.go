package script

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Files is a stream of a list of files. A user can either use the file list directly or the the
// created stream. In the stream, each line contains a path to a file.
type Files struct {
	Stream
	Files []FileInfo
}

// File contains information about a file.
type FileInfo struct {
	// FileInfo contains information about the file.
	os.FileInfo
	// Path is the path of the file. It may be relative or absolute, depending on how the `Ls`
	// command was invoked.
	Path string
}

// Ls returns a stream of a list files. In the returned stream, each line will contain a path to
// a single file.
//
// If the provided paths list is empty, the local directory will be listed.
//
// The provided paths may be relative to the local directory or absolute - this will influence the
// format of the returned paths in the output.
//
// If some provided paths correlate to the arguments correlate to the same file, it will also appear
// multiple times in the output.
//
// If any of the paths fails to be listed, it will result in an error in the output, but the stream
// will still conain all paths that were successfully listed.
//
// Shell command: `ls`.
func Ls(paths ...string) Files {
	// Default to local directory.
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	var (
		files  []FileInfo
		errors *multierror.Error
	)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("stat path: %s", err))
			continue
		}

		// Path is a single file.
		if !info.IsDir() {
			files = append(files, FileInfo{Path: path, FileInfo: info})
			continue
		}

		// Path is a directory.
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("read dir: %s", err))
			continue
		}

		for _, info := range infos {
			files = append(files, FileInfo{Path: filepath.Join(path, info.Name()), FileInfo: info})
		}
	}

	return Files{
		Stream: Stream{
			stage: fmt.Sprintf("ls (%+v)", paths),
			r:     &filesReader{files: files},
			err:   errors.ErrorOrNil(),
		},
		Files: files,
	}
}

// filesReader reads from a file info list.
type filesReader struct {
	files []FileInfo
	// seek indicates which file to write for the next Read function call.
	seek int
}

func (f *filesReader) Read(out []byte) (int, error) {
	if f.seek >= len(f.files) {
		return 0, io.EOF
	}

	line := []byte(f.files[f.seek].Path + "\n")
	f.seek++

	n := copy(out, line)
	return n, nil
}
