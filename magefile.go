// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

//go:build mage
// +build mage

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/sync/errgroup"

	"go.justen.tech/goodwill/internal/mage"
)

var Default = Build

var (
	curDir, _    = filepath.Abs(filepath.Join("."))
	pubFile      = filepath.Join(curDir, "goodwill.pub")
	targetDir    = filepath.Join(curDir, "target")
	testDir      = filepath.Join(curDir, "test")
	terraformDir = filepath.Join(testDir, "terraform")
)

var (
	goClasspathDir = mage.DirOnce(filepath.Join(targetDir, "classes", "go"), 0o755)
	binDir         = mage.DirOnce(filepath.Join(curDir, "bin"), 0o755)
	distDir        = mage.DirOnce(filepath.Join(curDir, "dist"), 0o755)
	cleanDirs      = []string{
		targetDir,
		filepath.Join(curDir, "bin"),
		filepath.Join(curDir, "dist"),
	}
)

var (
	pomVersion = mage.StringOnce(func() (string, error) {
		return sh.Output("mvn", "org.apache.maven.plugins:maven-help-plugin:3.1.0:evaluate", "-Dexpression=project.version", "-q", "-DforceStdout")
	})
	version = mage.StringOnce(func() (string, error) {
		if v := os.Getenv("VERSION"); v != "" {
			return v, nil
		}
		v := pomVersion()
		if sv := strings.TrimSuffix(v, "-SNAPSHOT"); sv != v {
			ts := time.Now().UTC().Format("20060102150405")
			c, err := sh.Output("git", "rev-parse", "--short=9", "HEAD")
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s-%s-%s", sv, ts, c), nil
		}
		return v, nil
	})
	gitCommit = mage.StringOnce(func() (string, error) {
		return sh.Output("git", "rev-parse", "HEAD")
	})
	buildTime = mage.StringOnce(func() (string, error) {
		return time.Now().UTC().Format(time.RFC3339), nil
	})
	sumFile = mage.StringOnce(func() (string, error) {
		version := version()
		return filepath.Join(distDir(), fmt.Sprintf("goodwill_%s_SHA256SUMS", version)), nil
	})
	localBin = mage.StringOnce(func() (string, error) {
		str := filepath.Join(curDir, "bin", "goodwill")
		if runtime.GOOS == "windows" {
			str += ".exe"
		}
		return str, nil
	})
	concordData = mage.LoadOnce(filepath.Join(terraformDir, "concord-env.json"))
)

func mvn(args ...string) error {
	args = append(args,
		fmt.Sprintf("-Dbuild.version=%s", version()),
		fmt.Sprintf("-Dbuild.gitCommit=%s", gitCommit()),
		fmt.Sprintf("-Dbuild.timestamp=%s", buildTime()),
	)
	return sh.RunV("mvn", args...)
}

var debug = log.New(os.Stderr, "", 0)

// Clean removes build artifacts and temporary files
func Clean() error {
	debug.Println("==> cleanup")
	for _, dir := range cleanDirs {
		if err := sh.Rm(dir); err != nil {
			return err
		}
	}
	return nil
}

// License adds licence headers to all files
func License() error {
	debug.Println("===> add license header")
	return sh.RunV("addlicense", "-v",
		"-f", filepath.Join(".", "HEADER"),
		"src", "gw", "internal", "test", "main.go", "magefile.go", "pom.xml")
}

// Bin builds the Go binary for the current os/arch
func Bin() error {
	debug.Println("==> build go binary:", localBin())
	_, err := mage.BuildTarget(binDir(), time.Now(), mage.Build{
		Version:   version(),
		GitCommit: gitCommit(),
		BuildTime: time.Now().UTC().Format(time.RFC3339),
	}, mage.Target{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Filename: localBin(),
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	uberJar *mage.Artifact
)

// UberJAR builds the jar with all its dependencies
func UberJAR() error {
	if uberJar != nil {
		return nil
	}
	mg.SerialDeps(buildAllGoBinaries, copyGoBinaries)
	debug.Println("==> package task")
	err := mvn("package", "-P", "package")
	if err != nil {
		return err
	}
	out := filepath.Join(targetDir, fmt.Sprintf("goodwill-%s-jar-with-dependencies.jar", pomVersion()))
	artifact, err := mage.JarArtifact(distDir(), version(), out)
	if err != nil {
		return err
	}
	uberJar = artifact
	return nil
}

// Build builds the project
func Build() error {
	var g Generate
	mg.SerialDeps(Dependencies, g.All, UberJAR)
	return nil
}

// Package distribution files
func Package() error {
	if err := mvn("clean"); err != nil {
		return err
	}
	mg.SerialDeps(Clean, Build, writeSums, sign, verify)
	return nil
}

// Sign signs the SHA sum file using SIGNIFY_KEY
func Sign() error {
	_, err := os.Stat(sumFile())
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		mg.SerialDeps(Package)
		return nil
	}
	mg.SerialDeps(sign, verify)
	return nil
}

type Nexus mg.Namespace

// Deploy builds, packages, and uploads to maven central
func (m Nexus) Deploy() error {
	mg.Deps(Package)
	debug.Println("==> deploy to maven central")
	return mvn("deploy", "-P", "release")
}

func (m Nexus) Release(stagingID string) error {
	debug.Println("==> release to maven central")
	return mvn("nexus-staging:release", "-P", "release", "-DstagingRepositoryId="+stagingID)
}

func (m Nexus) Drop(stagingID string) error {
	debug.Println("==> drop staging release")
	return mvn("nexus-staging:release", "-P", "release", "-DstagingRepositoryId="+stagingID)
}

func sign() error {
	if os.Getenv("SIGNIFY_KEY") == "" {
		return fmt.Errorf("SIGNIFY_KEY not set")
	}
	file := sumFile()
	key := os.Getenv("SIGNIFY_KEY")
	if key == "" {
		return fmt.Errorf("SIGNIFY_KEY not set")
	}
	return sh.RunV("signify", "-S", "-s", key, "-m", file, "-x", file+".sig")
}

func verify() error {
	debug.Println("==> verify sum signatures")
	if err := sh.RunV("signify", "-V", "-p", pubFile, "-m", sumFile()); err != nil {
		return err
	}
	debug.Println("==> verify sha sums")
	return mage.VerifySums(distDir(), sumFile())
}

// Dependencies ensures Go module dependencies are downloaded
func Dependencies() error {
	debug.Println("==> download go dependencies")
	if err := sh.RunV(mg.GoCmd(), "mod", "tidy"); err != nil {
		return err
	}
	return sh.RunV(mg.GoCmd(), "get")
}

type Generate mg.Namespace

// Generate generates protobuf code for Go and Java
func (g Generate) All() error {
	mg.Deps(g.Go, g.Java)
	return nil
}

// GenerateGo run go:generate
func (g Generate) Go() error {
	debug.Println("==> generate go code")
	return sh.RunV(mg.GoCmd(), "generate")
}

// GenerateJava generates java sources
func (g Generate) Java() error {
	debug.Println("==> generate java code")
	return mvn("generate-sources")
}

var lockGoBinaries sync.Mutex
var goBinaries []mage.Artifact

func addGoBinary(artifact mage.Artifact) {
	lockGoBinaries.Lock()
	defer lockGoBinaries.Unlock()
	goBinaries = append(goBinaries, artifact)
}

// buildGoBinary builds a Go binary for the target os/arch
func buildGoBinary(distDir string, os string, arch string) error {
	t, err := time.Parse(time.RFC3339, buildTime())
	if err != nil {
		return err
	}
	artifact, err := mage.BuildTarget(distDir, t, mage.Build{
		Version:   version(),
		GitCommit: gitCommit(),
		BuildTime: buildTime(),
	}, mage.Target{
		OS:   os,
		Arch: arch,
	})
	if err != nil {
		return err
	}
	addGoBinary(*artifact)
	return nil
}

// buildAllGoBinaries builds Go binaries for all supported os/arch targets
func buildAllGoBinaries() error {

	var deps []interface{}
	for _, target := range []struct {
		OS   string
		Arch string
	}{
		{"linux", "amd64"},
		{"linux", "386"},
		{"darwin", "amd64"},
		{"windows", "amd64"},
		{"windows", "386"},
	} {
		deps = append(deps, mg.F(buildGoBinary, distDir(), target.OS, target.Arch))
	}
	mg.Deps(deps...)
	return nil
}

// writeSums writes the SHA256 sum files for each artifact
func writeSums() error {
	var artifacts []string
	for _, a := range goBinaries {
		artifacts = append(artifacts, filepath.Join(distDir(), a.String()))
	}
	artifacts = append(artifacts, filepath.Join(distDir(), uberJar.String()))
	return mage.WriteSums(sumFile(), artifacts)
}

// Release cuts a new release version and tags the repository
func Release() error {
	ver := os.Getenv("VERSION")
	if ver == "" {
		return fmt.Errorf("VERSION not set")
	}
	debug.Println("==> set project version")
	ver = strings.TrimPrefix(ver, "v")
	if err := mvn("versions:set", "-DnewVersion="+ver); err != nil {
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
	return nil
}

// Snapshot sets the project version to a snapshot
func Snapshot() error {
	debug.Println("==> set project version snapshot")
	ver := os.Getenv("VERSION")
	if ver == "" {
		return fmt.Errorf("VERSION not set")
	}
	if err := mvn("versions:set", "-DnewVersion="+ver+"-SNAPSHOT"); err != nil {
		return fmt.Errorf("error setting project version: %w", err)
	}
	return nil
}

func copyGoBinaries() error {
	debug.Println("==> copy go binaries to classpath")
	for _, b := range goBinaries {
		source := filepath.Join(distDir(), b.String())
		b.Version = ""
		target := filepath.Join(goClasspathDir(), b.String())
		debug.Println("copy", target, "<=", source)
		if err := sh.Copy(target, source); err != nil {
			return err
		}
	}
	return nil
}

type E2E mg.Namespace

// E2EUp starts the end to end test environment
func (e E2E) Up() error {
	cenv := concordEnv()
	chdir := "-chdir=" + terraformDir
	if _, err := os.Stat(filepath.Join(terraformDir, ".terraform")); err != nil {
		if err := sh.RunV("terraform", chdir, "init"); err != nil {
			return err
		}
	}
	if err := sh.RunV("terraform", chdir, "apply", "-var", "concord_api_key="+cenv.AgentKey, "-auto-approve"); err != nil {
		return err
	}
	debug.Println("===> Waiting for Concord")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	return mage.WaitConcordRunning(ctx, cenv)
}

// E2EDown tears down the end to end test environment
func (e E2E) Down() error {
	debug.Println("==> destroy terraform environment")
	if err := sh.RunV("terraform", "-chdir="+terraformDir, "destroy", "-auto-approve"); err != nil {
		return err
	}
	debug.Println("==> cleanup terraform files")
	cleanFiles := []string{
		filepath.Join(terraformDir, ".terraform"),
		filepath.Join(terraformDir, ".terraform.lock.hcl"),
		filepath.Join(terraformDir, "terraform.tfstate"),
		filepath.Join(terraformDir, "terraform.tfstate.backup"),
		filepath.Join(terraformDir, "files", "maven.json"),
		filepath.Join(terraformDir, "files", "bootstrap.ldif"),
		filepath.Join(terraformDir, "files", "concord-agent.conf"),
		filepath.Join(terraformDir, "files", "concord-server.conf"),
	}
	for _, dir := range cleanFiles {
		if err := sh.Rm(dir); err != nil {
			return err
		}
	}
	return nil
}

func vendorTestFlow() error {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = filepath.Join(testDir, "flow")
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func precompileTestFlow() error {
	mg.Deps(Bin)
	return sh.RunV(localBin(), "-debug", "-os", "linux", "-arch", "amd64", "-dir", filepath.Join(testDir, "flow"), "-out", filepath.Join(testDir, "flow", "goodwill.tasks"))
}

func e2eDeps() {
	mg.SerialDeps(vendorTestFlow, precompileTestFlow)
}

// E2ETestPublished runts end to end tests targeting published maven artifacts
func (e E2E) TestPublished() error {
	mg.Deps(e.Up, e2eDeps)
	tests := []e2eTest{
		{
			"published",
			mage.ConcordParams{
				Runtime:      mage.ConcordRuntimeV1,
				Dependencies: true,
				Version:      pomVersion(),
			},
			[]mage.ZipFile{
				{Source: filepath.Join(testDir, "flow", "goodwill.tasks"), Dest: "goodwill.tasks"},
			},
		},
		{
			"published-v2",
			mage.ConcordParams{
				Runtime:      mage.ConcordRuntimeV2,
				Dependencies: true,
				Version:      pomVersion(),
			},
			[]mage.ZipFile{
				{Source: filepath.Join(testDir, "flow", "goodwill.tasks"), Dest: "goodwill.tasks"},
			},
		},
	}
	if err := runE2ETests(tests); err != nil {
		return err
	}
	return waitE2ETests()
}

// E2ETest runs end to end tests
func (e E2E) Test() error {
	mg.Deps(Build, e.Up, e2eDeps)
	jar := filepath.Join(distDir(), uberJar.String())
	tests := []e2eTest{
		{
			"compiled",
			mage.ConcordParams{
				Runtime:   mage.ConcordRuntimeV1,
				GoVersion: "1.20.3",
				UseDocker: true,
			},
			[]mage.ZipFile{
				{Source: jar, Dest: "lib/goodwill.jar"},
				{Source: filepath.Join(testDir, "flow", "goodwill.go"), Dest: "goodwill.go"},
				{Source: filepath.Join(testDir, "flow", "go.mod"), Dest: "go.mod"},
				{Source: filepath.Join(testDir, "flow", "go.sum"), Dest: "go.sum"},
				{Source: filepath.Join(testDir, "flow", "vendor"), Dest: "vendor"},
			},
		},
		{
			"compiled-v2",
			mage.ConcordParams{
				Runtime:   mage.ConcordRuntimeV2,
				GoVersion: "1.20.3",
				UseDocker: true,
			},
			[]mage.ZipFile{
				{Source: jar, Dest: "lib/goodwill.jar"},
				{Source: filepath.Join(testDir, "flow", "goodwill.go"), Dest: "goodwill.go"},
				{Source: filepath.Join(testDir, "flow", "go.mod"), Dest: "go.mod"},
				{Source: filepath.Join(testDir, "flow", "go.sum"), Dest: "go.sum"},
				{Source: filepath.Join(testDir, "flow", "vendor"), Dest: "vendor"},
			},
		},
		{
			"precompiled",
			mage.ConcordParams{
				Runtime:   mage.ConcordRuntimeV1,
				GoVersion: "1.20.3",
				UseDocker: true,
			},
			[]mage.ZipFile{
				{Source: jar, Dest: "lib/goodwill.jar"},
				{Source: filepath.Join(testDir, "flow", "goodwill.tasks"), Dest: "goodwill.tasks"},
			},
		},
		{
			"precompiled-v2",
			mage.ConcordParams{
				Runtime:   mage.ConcordRuntimeV2,
				GoVersion: "1.20.3",
				UseDocker: true,
			},
			[]mage.ZipFile{
				{Source: jar, Dest: "lib/goodwill.jar"},
				{Source: filepath.Join(testDir, "flow", "goodwill.tasks"), Dest: "goodwill.tasks"},
			},
		},
	}
	if err := runE2ETests(tests); err != nil {
		return err
	}
	return waitE2ETests()
}

type e2eTest struct {
	name   string
	params mage.ConcordParams
	files  []mage.ZipFile
}

var e2eTests []string

func waitE2ETests() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	for _, p := range e2eTests {
		processID := p
		group.Go(func() error {
			return waitE2ETest(ctx, processID)
		})
	}
	return group.Wait()
}

func waitE2ETest(ctx context.Context, processID string) error {
	if err := mage.WaitConcordProcess(ctx, concordEnv(), processID); err != nil {
		debug.Println("- FAIL:", processID)
		return fmt.Errorf("%s: %w", processID, err)
	}
	debug.Println("- PASS:", processID)
	return nil
}

func runE2ETests(tests []e2eTest) error {
	var testerrors bool
	for _, test := range tests {
		if _, err := runE2ETest(test.name, test.params, test.files); err != nil {
			debug.Println("[ERROR]", test.name, err)
		}
	}
	debug.Println("===> API Key:", concordEnv().APIKey)
	if testerrors {
		return fmt.Errorf("Error running tests")
	}
	return nil
}

func runE2ETest(name string, params mage.ConcordParams, files []mage.ZipFile) (string, error) {
	cenv := concordEnv()
	debug.Println("===> run e2e test", name)
	var concordYAML bytes.Buffer
	if err := mage.GenerateConcordYaml(&concordYAML, params); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	instanceID, err := mage.NewConcordProcess(ctx, cenv, &concordYAML, files)
	if err != nil {
		return "", err
	}
	e2eTests = append(e2eTests, instanceID)
	debug.Println("Concord Job Submitted:")
	debug.Printf("http://localhost:8001/#process/%s/status", instanceID)
	return instanceID, nil
}

func concordEnv() mage.ConcordEnv {
	var env mage.ConcordEnv
	if err := json.Unmarshal(concordData(), &env); err != nil {
		debug.Fatalln("ERROR", err)
	}
	return env
}
