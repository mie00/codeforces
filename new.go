package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func newFile(config Config) error {
	target := path.Join(config.Name, config.Filename)
	if _, err := os.Stat(target); !config.ForceWrite && err == nil {
		return fmt.Errorf("file %s already exists", target)
	} else if config.ForceWrite || os.IsNotExist(err) {
	} else {
		return err
	}
	src := path.Join("templates", config.Filename)
	if _, err := os.Stat(src); err == nil {
	} else if os.IsNotExist(err) {
		return fmt.Errorf("cannot find template %s", src)
	} else {
		return err
	}
	return Copy(src, target)
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
