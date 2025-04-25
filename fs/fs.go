package fs

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path"
)

func GetFilenameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return string(rand.Int())
	}

	f := path.Base(u.Path)
	if f == "" || f == "/" {
		// recurso raíz, usar host + timestamp
		f = fmt.Sprintf("%s.html", u.Hostname())
	}
	return f
}

func EnsureDir(dirPath string) string {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		println("Cannot write in path: ", dirPath)
		os.Exit(1)
	}
	return dirPath
}

func WriteFile(path string, data []byte) bool {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return false
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
