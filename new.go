package main

import (
	"io"
	"os"
	"path"
)

func newFile(config Config) error {
	target := path.Join(config.Name, config.Filename)
	if _, err := os.Stat(target); !config.ForceDownload && err == nil {
		return nil
	} else if config.ForceDownload || os.IsNotExist(err) {
		return Copy(path.Join("templates", config.Filename), target)
	} else {
		return err
	}
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
