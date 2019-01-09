package syugo

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type (
	Syugo struct {
		Collects []*Collect
	}

	Collect struct {
		Repository string
		Requests   []string
		Version    string
		Dir        string
	}
)

func NewSyugo(cs []*Collect) (*Syugo, error) {
	if len(cs) == 0 {
		return nil, errors.New("collect is required")
	}
	for i, v := range cs {
		if len(v.Repository) == 0 {
			return nil, fmt.Errorf("collect[%d] is repository required.", i)
		}
	}

	return &Syugo{
		Collects: cs,
	}, nil
}

func (s *Syugo) Run() (err error) {
	tmpDir := path.Join(os.TempDir(), random())
	defer func() {
		err = os.RemoveAll(tmpDir)
	}()

	for _, v := range s.Collects {
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

func (c *Collect) checkout(dir string) error {
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
