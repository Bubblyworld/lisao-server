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

	return cmd.
		NewCommand("git checkout").
		WithArg(hash).
		CD(c.repoPath).
		LogTo(logWriter).
		Do()
}
