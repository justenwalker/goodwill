package mage

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Artifact struct {
	Version string
	OS      string
	Arch    string
	Ext     string
}

func (b Artifact) String() string {
	var sb strings.Builder
	sb.WriteString("goodwill")
	switch b.Ext {
	case "jar":
		if b.Version != "" {
			sb.WriteRune('-')
			sb.WriteString(b.Version)
		}
	default:
		if b.Version != "" {
			sb.WriteRune('_')
			sb.WriteString(b.Version)
		}
		if b.OS != "" {
			sb.WriteRune('_')
			sb.WriteString(b.OS)
		}
		if b.Arch != "" {
			sb.WriteRune('_')
			sb.WriteString(b.Arch)
		}
	}
	if b.Ext != "" {
		sb.WriteRune('.')
		sb.WriteString(b.Ext)
	}
	return sb.String()
}

func JarArtifact(distDir, version, filename string) (*Artifact, error) {
	jar := Artifact{
		Version: version,
		Ext:     "jar",
	}
	if err := sh.Copy(filepath.Join(distDir, jar.String()), filename); err != nil {
		return nil, err
	}
	return &jar, nil
}

func WriteSums(sumFile string, artifacts []string) (err error) {
	var sf *os.File
	sf, err = os.OpenFile(sumFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		cerr := sf.Close()
		if err == nil {
			err = cerr
		}
	}()
	for _, file := range artifacts {
		hash, err := HashFile(file)
		if err != nil {
			return err
		}
		if _, err := sf.WriteString(fmt.Sprintf("%s  %s\n", hash, filepath.Base(file))); err != nil {
			return err
		}
	}
	return nil
}

func VerifySums(distDir string, sumFile string) error {
	data, err := ioutil.ReadFile(sumFile)
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(bytes.NewBuffer(data))
	for sc.Scan() {
		fields := strings.Split(strings.TrimSpace(sc.Text()), "  ")
		if len(fields) != 2 {
			continue
		}
		filename, expected := fields[1], fields[0]
		hash, err := HashFile(filepath.Join(distDir, filename))
		if err != nil {
			return fmt.Errorf("%q: error hashing file: %w", filename, err)
		}
		if hash != expected {
			return fmt.Errorf("%q: hash match. got=%q, expected=%q", filename, hash, expected)
		}
		debug.Printf("%s\tOK\n", filename)
	}
	return nil

}

type Build struct {
	Version   string
	GitCommit string
	BuildTime string
}

type Target struct {
	OS       string
	Arch     string
	Filename string
}

func BuildTarget(distDir string, mod time.Time, build Build, target Target) (*Artifact, error) {
	bin := Artifact{
		OS:      target.OS,
		Arch:    target.Arch,
		Version: build.Version,
	}
	env := make(map[string]string)
	if bin.OS != "" {
		env["GOOS"] = bin.OS
	} else {
		env["GOOS"] = runtime.GOOS
	}
	if bin.Arch != "" {
		env["GOARCH"] = bin.Arch
	} else {
		env["GOARCH"] = runtime.GOARCH
	}
	if target.OS == "windows" {
		bin.Ext = "exe"
	}
	outfile := target.Filename
	if outfile == "" {
		outfile = filepath.Join(distDir, bin.String())
	}
	ldflags := fmt.Sprintf(`-X main.Version=%s -X main.GitCommit=%s -X main.BuildTime=%s`, build.Version, build.GitCommit, build.BuildTime)
	var update bool
	if stat, err := os.Stat(outfile); err != nil {
		if os.IsNotExist(err) {
			update = true
		}
	} else {
		update = stat.ModTime().Before(mod)
	}
	if update {
		err := sh.RunWithV(env, mg.GoCmd(), "build", "-ldflags", ldflags, "-o", outfile)
		if err != nil {
			return nil, err
		}
	}
	return &bin, nil
}

func HashFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
