package storage

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var storage string
var static string

func Init(staticPath, storagePath string) error {
	storage = storagePath
	static = staticPath
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

	if m := mime.TypeByExtension(realName); m != "" {
		w.Header().Set("Content-Type", m)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+realName)
	http.ServeContent(w, r, realName, stat.ModTime(), file)
	return nil
}

func Store(filename string, r io.Reader) error {
	file, err := os.Create(filepath.Join(storage, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func Delete(filename string) {
	os.Remove(filepath.Join(storage, filename))
}
