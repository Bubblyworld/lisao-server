package git

import "clanpj/lisao/cmd"

type Client struct {
	repoPath string
}

func NewClient(repoPath string) Client {
	return Client{
		repoPath: repoPath,
	}
}

func (c Client) CheckoutCommit(hash string) error {
	logWriter := cmd.NewLogWriter()
	defer logWriter.Close()

	command := cmd.NewCommand("git checkout " + hash)
	command.Dir = c.repoPath
	command.Stdout = logWriter
	command.Stderr = logWriter

	return command.Run()
}
