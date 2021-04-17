// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

// +build mage

package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	version    string
	goBinaries []artifact
	jar        artifact
	localBin   string
	sumFile    string
	lock       sync.Mutex
)

var (
	curDir, _      = filepath.Abs(filepath.Join("."))
	binDir         = filepath.Join(curDir, "bin")
	distDir        = filepath.Join(curDir, "dist")
	targetDir      = filepath.Join(curDir, "target")
	goClasspathDir = filepath.Join(targetDir, "classes", "go")
	testDir        = filepath.Join(curDir, "test")
	terraformDir   = filepath.Join(testDir, "terraform")
	cleanDirs      = []string{binDir, distDir, targetDir}
)

var Default = Package

var debug = log.New(os.Stderr, "", 0)

// Clean up build artifacts
func Clean() error {
	debug.Println("==> cleanup")
	for _, dir := range cleanDirs {
		if err := sh.Rm(dir); err != nil {
			return err
		}
	}
	return nil
}

func License() error {
	debug.Println("===> add license header")
	return sh.RunV("addlicense", "-v",
		"-f", filepath.Join(".", "HEADER"),
		"src", "gw", "internal", "test", "main.go", "magefile.go", "pom.xml")
}

// Build the Go binary only
func Build() error {
	localBin = filepath.Join(".", "bin", "goodwill")
	if runtime.GOOS == "windows" {
		localBin += ".exe"
	}
	debug.Println("==> build go binary:", localBin)
	return sh.RunV(mg.GoCmd(), "build", "-o", localBin)
}

// Build the task JAR
func PackageJAR() error {
	if err := os.MkdirAll(distDir, os.FileMode(0o755)); err != nil {
		return err
	}
	mg.SerialDeps(pomVersion, buildAllGoBinaries, copyGoBinaries)
	debug.Println("==> package task")
	err := sh.RunV("mvn", "package")
	if err != nil {
		return err
	}
	out := filepath.Join(targetDir, fmt.Sprintf("goodwill-%s-jar-with-dependencies.jar", version))
	jar = artifact{
		Filename: filepath.Join(distDir, fmt.Sprintf("goodwill-%s.jar", version)),
	}
	jar.Hash, err = hashFile(out)
	if err != nil {
		return err
	}
	return sh.Copy(jar.Filename, out)
}

// Build and package distribution files
func Package() error {
	mg.SerialDeps(Dependencies, Generate, PackageJAR, sha256Sums)
	return nil
}

func Sign() error {
	mg.Deps(pomVersion)
	sumFile = filepath.Join(distDir, fmt.Sprintf("goodwill_%s_SHA256SUMS", version))
	key := os.Getenv("SIGNIFY_KEY")
	if key == "" {
		return fmt.Errorf("SIGNIFY_KEY not set")
	}
	return sh.RunV("signify", "-S", "-s", key, "-m", sumFile, "-x", sumFile+".sig")
}

func sha256Sums() (err error) {
	var sf *os.File
	sumFile = filepath.Join(distDir, fmt.Sprintf("goodwill_%s_SHA256SUMS", version))
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
	for _, b := range append(goBinaries, jar) {
		_, err := sf.WriteString(fmt.Sprintf("%s  %s\n", b.Hash, filepath.Base(b.Filename)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Ensure Go module dependencies are downloaded
func Dependencies() error {
	debug.Println("==> download go dependencies")
	return sh.RunV(mg.GoCmd(), "mod", "download")
}

// Generate Protobuf code for Go and Java
func Generate() error {
	mg.Deps(GenerateGo, GenerateJava)
	return nil
}

// Generate Go code
func GenerateGo() error {
	debug.Println("==> generate go code")
	return sh.RunV(mg.GoCmd(), "generate")
}

// Generate Java Code
func GenerateJava() error {
	debug.Println("==> generate java code")
	return sh.RunV("mvn", "generate-sources")
}

// Build Go binary for all architectures
func buildAllGoBinaries() error {
	if err := os.MkdirAll(distDir, os.FileMode(0o755)); err != nil {
		return err
	}
	var deps []interface{}
	for _, t := range []struct {
		OS   string
		Arch string
	}{
		{"linux", "amd64"},
		{"linux", "386"},
		{"darwin", "amd64"},
		{"windows", "amd64"},
		{"windows", "386"},
	} {
		deps = append(deps, mg.F(buildGoBinary, distDir, t.OS, t.Arch))
	}
	mg.Deps(deps...)
	return nil
}

func Release() error {
	mg.Deps(Clean)
	mg.SerialDeps(Dependencies, Generate)
	debug.Println("==> set project version")
	ver := os.Getenv("VERSION")
	if ver == "" {
		return fmt.Errorf("VERSION not set")
	}
	ver = strings.TrimPrefix(ver, "v")
	if err := sh.Run("mvn", "versions:set", "-DnewVersion="+ver); err != nil {
		return fmt.Errorf("error setting project version: %w", err)
	}
	sh.Rm("pom.xml.versionsBackup")
	debug.Println("==> make release commit")
	if err := sh.Run("git", "add", "pom.xml"); err != nil {
		return err
	}
	if err := sh.Run("git", "commit", "-m", fmt.Sprintf("Prepare Release: %s", ver)); err != nil {
		return err
	}
	debug.Println("==> tag release")
	if err := sh.Run("git", "tag", "--annotate", "-m", "Release v"+ver, "v"+ver); err != nil {
		return err
	}
	mg.SerialDeps(PackageJAR, sha256Sums)
	return nil
}

func Snapshot() error {
	debug.Println("==> set project version snapshot")
	ver := os.Getenv("VERSION")
	if ver == "" {
		return fmt.Errorf("VERSION not set")
	}
	if err := sh.Run("mvn", "versions:set", "-DnewVersion="+ver+"-SNAPSHOT"); err != nil {
		return fmt.Errorf("error setting project version: %w", err)
	}
	return nil
}

func pomVersion() error {
	v, err := sh.Output("mvn", "org.apache.maven.plugins:maven-help-plugin:3.1.0:evaluate", "-Dexpression=project.version", "-q", "-DforceStdout")
	version = v
	return err
}

func copyGoBinaries() error {
	debug.Println("==> copy go binaries to classpath")
	if err := os.MkdirAll(goClasspathDir, 0o755); err != nil {
		return err
	}
	for _, b := range goBinaries {
		target := filepath.Join(goClasspathDir, b.NoVersion())
		if err := sh.Copy(target, b.Filename); err != nil {
			return err
		}
	}
	return nil
}

const (
	authToken   = "LcJodmNtim3i1XfY0Pivsw"
	orgName     = "Default"
	projectName = "Test"
)

func E2EUp() error {
	chdir := "-chdir=" + terraformDir
	if _, err := os.Stat(filepath.Join(terraformDir, ".terraform")); err != nil {
		if err := sh.RunV("terraform", chdir, "init"); err != nil {
			return err
		}
	}
	return sh.RunV("terraform", chdir, "apply", "-auto-approve")

}

func E2EDown() error {
	return sh.RunV("terraform", "-chdir="+terraformDir, "destroy", "-auto-approve")
}

func vendorTestFlow() error {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = filepath.Join(testDir, "flow")
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func precompileTestFlow() error {
	mg.Deps(Build)
	return sh.RunV(localBin, "-debug", "-os", "linux", "-arch", "amd64", "-dir", filepath.Join(testDir, "flow"), "-out", filepath.Join(testDir, "flow", "goodwill.tasks"))
}

func E2E() error {
	mg.Deps(Package, E2EUp)
	mg.SerialDeps(vendorTestFlow, precompileTestFlow)
	debug.Println("===> API Key:", authToken)

	runE2ETest("compiled", []payloadFile{
		{filepath.Join(testDir, "concord.yml"), "concord.yml"},
		{jar.Filename, "lib/goodwill.jar"},
		{filepath.Join(testDir, "flow", "goodwill.go"), "goodwill.go"},
		{filepath.Join(testDir, "flow", "go.mod"), "go.mod"},
		{filepath.Join(testDir, "flow", "go.sum"), "go.sum"},
		{filepath.Join(testDir, "flow", "vendor"), "vendor"},
	})
	runE2ETest("precompiled", []payloadFile{
		{filepath.Join(testDir, "concord.yml"), "concord.yml"},
		{jar.Filename, "lib/goodwill.jar"},
		{filepath.Join(testDir, "flow", "goodwill.tasks"), "goodwill.tasks"},
	})
	return nil
}

type payloadFile struct {
	From string
	To   string
}

func runE2ETest(name string, files []payloadFile, ) error {
	debug.Println("===> run e2e test", name)
	var buf bytes.Buffer
	mpw := multipart.NewWriter(&buf)
	mpw.WriteField("org", orgName)
	mpw.WriteField("project", projectName)
	payload, err := mpw.CreateFormFile("archive", "payload.zip")
	if err != nil {
		return err
	}
	zw := zip.NewWriter(payload)
	for _, file := range files {
		stat, err := os.Stat(file.From)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			err = addZipDir(zw, file.From, file.To)
		} else {
			err = addZipFile(zw, file.From, file.To)
		}
		if err != nil {
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("close zip file: %w", err)
	}
	if err := mpw.Close(); err != nil {
		return fmt.Errorf("close multipart file: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8001/api/v1/process", &buf)
	if err != nil {
		return fmt.Errorf("could not create http request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", mime.FormatMediaType("multipart/form-data", map[string]string{"boundary": mpw.Boundary()}))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send http request: %w", err)
	}
	defer resp.Body.Close()
	rbody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: %s", resp.Status, string(rbody))
	}
	response := struct {
		InstanceID string `json:"instanceId"`
	}{}
	if err := json.Unmarshal(rbody, &response); err != nil {
		return fmt.Errorf("could not parse concord response: %w\nresponse:\n%s", err, string(rbody))
	}
	debug.Println("Concord Job Submitted:")
	debug.Printf("http://localhost:8001/#process/%s/status", response.InstanceID)
	return nil
}

func addZipDir(zw *zip.Writer, dir string, dest string) error {
	return filepath.Walk(dir, func(filename string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return addZipFile(zw, filename, path.Join(dest, strings.TrimPrefix(filename, dir)))
	})
}

func addZipFile(zw *zip.Writer, filename string, dest string) error {
	//debug.Printf("> payload: %s <- %s", dest, filename)
	w, err := zw.Create(dest)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func addBinary(b artifact) {
	lock.Lock()
	goBinaries = append(goBinaries, b)
	lock.Unlock()
}

type artifact struct {
	Filename string
	OS       string
	Arch     string
	Ext      string
	Hash     string
}

func (b artifact) String() string {
	var sb strings.Builder
	sb.WriteString("goodwill_")
	sb.WriteString(version)
	if b.OS != "" {
		sb.WriteRune('_')
		sb.WriteString(b.OS)
	}
	if b.Arch != "" {
		sb.WriteRune('_')
		sb.WriteString(b.Arch)
	}
	if b.Ext != "" {
		sb.WriteRune('.')
		sb.WriteString(b.Ext)
	}
	return sb.String()
}

func (b artifact) NoVersion() string {
	var sb strings.Builder
	sb.WriteString("goodwill")
	if b.OS != "" {
		sb.WriteRune('_')
		sb.WriteString(b.OS)
	}
	if b.Arch != "" {
		sb.WriteRune('_')
		sb.WriteString(b.Arch)
	}
	if b.Ext != "" {
		sb.WriteRune('.')
		sb.WriteString(b.Ext)
	}
	return sb.String()
}

func buildGoBinary(dir, goos, goarch string) error {
	bin := artifact{
		OS:   goos,
		Arch: goarch,
	}
	if goos == "windows" {
		bin.Ext = "exe"
	}
	bin.Filename = filepath.Join(dir, bin.String())
	err := sh.RunWithV(map[string]string{
		"GOOS":   goos,
		"GOARCH": goarch,
	}, mg.GoCmd(), "build", "-o", bin.Filename)
	if err != nil {
		return err
	}
	bin.Hash, err = hashFile(bin.Filename)
	if err != nil {
		return err
	}
	addBinary(bin)
	return nil
}

func hashFile(filename string) (string, error) {
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
