package tasks

import (
	"errors"
	"flag"
	"log"
	"path/filepath"

	"clanpj/lisao/cmd"
	"clanpj/lisao/cmd/git"
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
// TODO(guy) workers should probably pass context down
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

	return git.
		NewClient(absolutePath).
		CheckoutCommit(bi.commitHash)
}

func (bi BuildInfo) buildMain() error {
	log.Printf("builds: installing package %s", bi.mainPath)

	absolutePath, err := absolutePathToRepo()
	if err != nil {
		return err
	}

	return cmd.
		NewCommand("go build").
		WithFlag("-i").
		SetParam("-o", bi.outputPath).
		WithArg(bi.mainPath).
		SetEnv("GOPATH", absolutePath).
		CD(absolutePath).
		Do()
}
