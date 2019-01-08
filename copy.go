package syugo

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return rcopy(src, dest, info)
}

func rcopy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

func fcopy(src, dest string, info os.FileInfo) (err error) {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err = s.Close()
	}()

	_, err = io.Copy(f, s)
	return err
}

func dcopy(srcdir, destdir string, info os.FileInfo) error {
	if err := os.MkdirAll(destdir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := rcopy(cs, cd, content); err != nil {
			return err
		}
	}
	return nil
}

func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(src, dest)
}
