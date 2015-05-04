package compiler

import (
	"io/ioutil"
	"path/filepath"
)

// Gets files, storing the contents in-memory for faster future lookups
type FileCache struct {
	cache map[string][]byte
}

func NewFileCache() *FileCache {
	return &FileCache{
		cache: make(map[string][]byte, 100),
	}
}

// Gets the given file
func (self *FileCache) Get(path string) ([]byte, error) {
	abs, err := filepath.Abs(path)

	if err != nil {
		return []byte{}, err
	}

	contents, ok := self.cache[abs]

	if !ok {
		contents, err := ioutil.ReadFile(abs)

		if err == nil {
			self.cache[abs] = contents
		}

		return contents, err
	} else {
		return contents, nil
	}
}

func (self *FileCache) Invalidate(path string) error {
	abs, err := filepath.Abs(path)

	if err != nil {
		return err
	} else {
		delete(self.cache, abs)
		return nil
	}
}
