package storage

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var storage string

func Init(storagePath string) error {
	storage = storagePath
	if storage == "" {
		storage = filepath.Join(os.TempDir(), "cybersec-file-storage")
	}
	if err := os.MkdirAll(storage, os.ModePerm); err != nil {
		return err
	}
	stat, err := os.Stat(storage)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("storage path is not a directory")
	}
	// We could check if the directory is writable here, but that is beyond the
	// scope of this project. If you try to use a root-owned directory that's on you.
	return nil
}

func Close() {}

func Get(w http.ResponseWriter, r *http.Request, storedName, realName string) error {
	file, err := os.Open(filepath.Join(storage, storedName))
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	http.ServeContent(w, r, realName, stat.ModTime(), file)
	return nil
}

func Store(filename string, r io.Reader) {

}
