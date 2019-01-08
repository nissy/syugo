package syugo

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type (
	Collects []*Collect

	Collect struct {
		Repository string   `toml:"repository" validate:"required"`
		Requests   []string `toml:"requests"`
		Version    string   `toml:"version" validate:"required"`
		Dir        string   `toml:"dir"`
	}
)

func (cs Collects) Run() (err error) {
	tmpDir := path.Join(os.TempDir(), random())
	defer func() {
		err = os.RemoveAll(tmpDir)
	}()

	for _, v := range cs {
		if err := v.checkout(path.Join(tmpDir, v.Dir)); err != nil {
			return err
		}
	}

	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	return copy(tmpDir, curDir)
}

func (c *Collect) checkout(dir string) (err error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := runCmd("git", "-C", dir, "init"); err != nil {
		return err
	}
	if err := runCmd("git", "-C", dir, "remote", "add", "origin", c.Repository); err != nil {
		return err
	}
	if len(c.Requests) > 0 {
		if err := runCmd("git", "-C", dir, "config", "core.sparsecheckout", "true"); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(dir, ".git", "info", "sparse-checkout"), []byte(strings.Join(c.Requests, "\n")), os.ModePerm); err != nil {
			return err
		}
	}
	if err := runCmd("git", "-C", dir, "pull", "origin", c.Version); err != nil {
		return err
	}
	if err := os.RemoveAll(path.Join(dir, ".git")); err != nil {
		return err
	}

	return nil
}

func runCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func random() string {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		panic(err)
	}

	return strconv.FormatUint(n, 36)
}
