package salt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Target struct {
	Node string
	User string
	Sudo bool
}

func ssh(masterConfig MasterConfig, command string, args ...string) (string, string, error) {
	saltArgs := []string{
		"-v",
		"-c",
		masterConfig.GetTempDir(),
		"-i",
		"--state-verbose=true",
	}

	defer os.RemoveAll(masterConfig.GetTempDir())

	saltArgs = append(
		saltArgs,
		"target",
		command,
	)

	saltArgs = append(saltArgs, args...)

	cmd := exec.Command("salt-ssh", saltArgs...)
	fmt.Printf("Command is %+v\n", cmd)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()
	stdout := stdOut.String()
	stderr := stdErr.String()
	//TODO: print stdout only when in debug mode
	fmt.Printf("stdout is %s\n", stdout)
	fmt.Printf("stderr is %s\n", stderr)

	return stdout, stderr, err
}
