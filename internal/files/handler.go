package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileHandler struct {
	rootDir string
	mut     sync.RWMutex
}

func NewFileHandler(rootDir string) *FileHandler {
	return &FileHandler{rootDir: rootDir}
}

func (fh *FileHandler) Create(fileName string) error {
	fh.mut.Lock()
	defer fh.mut.Unlock()

	fullPath, err := fh.resolveFullPath(fileName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	return f.Close()
}

func (fh *FileHandler) Delete(fileName string) error {
	fh.mut.Lock()
	defer fh.mut.Unlock()

	fullPath, err := fh.resolveFullPath(fileName)
	if err != nil {
		return err
	}

	return os.Remove(fullPath)
}

func (fh *FileHandler) Write(fileName string, data io.Reader) (int, error) {
	fh.mut.Lock()
	defer fh.mut.Unlock()

	fullPath, err := fh.resolveFullPath(fileName)
	if err != nil {
		return 0, err
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n64, err := io.Copy(file, data)
	if err != nil {
		return 0, err
	}
	return int(n64), nil
}

func (fh *FileHandler) Read(fileName string, size uint64) ([]byte, error) {
	fh.mut.RLock()
	defer fh.mut.RUnlock()

	fullPath, err := fh.resolveFullPath(fileName)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fullPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if size == 0 {
		return io.ReadAll(file)
	}

	data := make([]byte, size)
	n, err := io.ReadAtLeast(file, data, 1)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			if n == 0 {
				return []byte{}, nil
			}
			return data[:n], nil
		}
		return nil, err
	}

	return data[:n], nil
}

func (fh *FileHandler) Clear() error {
	fh.mut.Lock()
	defer fh.mut.Unlock()

	return os.RemoveAll(fh.rootDir)
}

func (fh *FileHandler) resolveFullPath(path string) (string, error) {
	rel, err := TransformPath(path)
	if err != nil {
		return "", err
	}
	return filepath.Join(fh.rootDir, rel), nil
}

func (fh *FileHandler) RootDir() string {
	return fh.rootDir
}
