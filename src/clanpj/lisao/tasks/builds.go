package tasks

import (
	"errors"
	"flag"
	"log"
	"os/exec"
	"path/filepath"
)

var pathToRepo = flag.String("repo_path", "", "Path to lisao-bot git repo.")

// BuildInfo stores info about a commit hash and go path to build.
type BuildInfo struct {
	mainPath   string
	outputPath string

	commitHash string
}

func NewBuildInfo(commitHash, mainPath, outputPath string) BuildInfo {
	return BuildInfo{
		commitHash: commitHash,
		mainPath:   mainPath,
		outputPath: outputPath,
	}
}

// TODO(guy) add logging for stdout/stderr from exec.Run().
func DoBuild(work interface{}) error {
	buildInfo, ok := work.(BuildInfo)
	if !ok {
		return errors.New("builds: received wrong type of work, should be BuildInfo")
	}

	err := buildInfo.checkoutCommit()
	if err != nil {
		return err
	}

	return buildInfo.buildMain()
}

func absolutePathToRepo() (string, error) {
	if *pathToRepo == "" {
		return "", errors.New("builds: no path to repo given")
	}

	return filepath.Abs(*pathToRepo)
}

func (bi BuildInfo) checkoutCommit() error {
	log.Printf("builds: checking out commit %s", bi.commitHash)

	absolutePath, err := absolutePathToRepo()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "checkout", bi.commitHash)
	cmd.Dir = absolutePath

	return cmd.Run()
}

func (bi BuildInfo) buildMain() error {
	log.Printf("builds: installing package %s", bi.mainPath)

	absolutePath, err := absolutePathToRepo()
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-i", "-o", bi.outputPath, bi.mainPath)
	cmd.Env = append(cmd.Env, "GOPATH="+absolutePath)
	cmd.Dir = absolutePath
	return cmd.Run()
}
