// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "os"

// PathError reports an error and the file path where it occurred.
type PathError struct {
	Path string;
	Error Error;
}

func (p *PathError) String() string {
	return p.Path + ": " + p.Error.String();
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func MkdirAll(path string, perm int) Error {
	// If path exists, stop with success or error.
	dir, err := os.Lstat(path);
	if err == nil {
		if dir.IsDirectory() {
			return nil;
		}
		return &PathError{path, ENOTDIR};
	}

	// Doesn't already exist; make sure parent does.
	i := len(path);
	for i > 0 && path[i-1] == '/' {	// Skip trailing slashes.
		i--;
	}

	j := i;
	for j > 0 && path[j-1] != '/' {	// Scan backward over element.
		j--;
	}

	if j > 0 {
		// Create parent
		err = MkdirAll(path[0:j-1], perm);
		if err != nil {
			return err;
		}
	}

	// Now parent exists, try to create.
	err = Mkdir(path, perm);
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path);
		if err1 == nil && dir.IsDirectory() {
			return nil;
		}
		return &PathError{path, err};
	}
	return nil;
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters.
func RemoveAll(path string) Error {
	// Simple case: if Remove works, we're done.
	err := Remove(path);
	if err == nil {
		return nil;
	}

	// Otherwise, is this a directory we need to recurse into?
	dir, err1 := os.Lstat(path);
	if err1 != nil {
		return &PathError{path, err1};
	}
	if !dir.IsDirectory() {
		// Not a directory; return the error from Remove.
		return &PathError{path, err};
	}

	// Directory.
	fd, err := Open(path, os.O_RDONLY, 0);
	if err != nil {
		return &PathError{path, err};
	}
	defer fd.Close();

	// Remove contents & return first error.
	err = nil;
	for {
		names, err1 := fd.Readdirnames(100);
		for i, name := range names {
			err1 := RemoveAll(path + "/" + name);
			if err1 != nil && err == nil {
				err = err1;
			}
		}
		// If Readdirnames returned an error, use it.
		if err1 != nil && err == nil {
			err = &PathError{path, err1};
		}
		if len(names) == 0 {
			break;
		}
	}

	// Remove directory.
	err1 = Remove(path);
	if err1 != nil && err == nil {
		err = &PathError{path, err1};
	}
	return err;
}
