package ssh

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"suse.com/caaspctl/pkg/caaspctl"
	"suse.com/caaspctl/internal/pkg/caaspctl/deployments"
	"suse.com/caaspctl/internal/pkg/caaspctl/deployments/ssh/assets"
)

var (
	stateMap = map[string]deployments.Runner{
		"cni.deploy": cniDeploy(),
		"kubeadm.init": kubeadmInit(),
		"kubeadm.join": kubeadmJoin(),
		"kubelet.configure": kubeletConfigure(),
		"kubelet.enable": kubeletEnable(),
		"kubernetes.upload-secrets": kubernetesUploadSecrets(),
	}
)

type Target struct {
	Node   string
	User   string
	Sudo   bool
	Client *ssh.Client
}

func NewTarget(target, user string, sudo bool) deployments.Target {
	return &Target{
		Node: target,
		User: user,
		Sudo: sudo,
	}
}

func (t *Target) Apply(states ...string) error {
	for _, state := range states {
		if state, stateExists := stateMap[state]; stateExists {
			state.Run(t)
		} else {
			log.Fatalf("State does not exist: %s", state)
		}
	}
	return nil
}

func (t *Target) Target() string {
	return t.Node
}

func (t *Target) UploadFile(sourcePath, targetPath string) error {
	if contents, err := ioutil.ReadFile(sourcePath); err == nil {
		return t.UploadFileContents(targetPath, string(contents))
	}
	return nil
}

func (t *Target) UploadFileContents(targetPath, contents string) error {
	if target := sshTarget(t); target != nil {
		dir, _ := path.Split(targetPath)
		encodedContents := base64.StdEncoding.EncodeToString([]byte(contents))
		target.ssh("mkdir", "-p", dir)
		target.sshWithStdin(encodedContents, "base64", "-d", "-w0", fmt.Sprintf("> %s", targetPath))
	}
	return errors.New("cannot access SSH target")
}

func (t *Target) DownloadFileContents(sourcePath string) (string, error) {
	if target := sshTarget(t); target != nil {
		if stdout, _, err := target.ssh("base64", "-w0", sourcePath); err == nil {
			decodedStdout, err := base64.StdEncoding.DecodeString(stdout)
			if err != nil {
				return "", err
			}
			return string(decodedStdout), nil
		} else {
			return "", err
		}
	}
	return "", errors.New("cannot access SSH target")
}

func (t *Target) ssh(command string, args ...string) (stdout string, stderr string, error error) {
	return t.sshWithStdin("", command, args...)
}

func (t *Target) sshWithStdin(stdin string, command string, args ...string) (stdout string, stderr string, error error) {
	if t.Client == nil {
		t.initClient()
	}
	session, err := t.Client.NewSession()
	if err != nil {
		return "", "", err
	}
	if len(stdin) > 0 {
		session.Stdin = bytes.NewBufferString(stdin)
	}
	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderrReader, err := session.StderrPipe()
	if err != nil {
		return "", "", err
	}
	finalCommand := strings.Join(append([]string{command}, args...), " ")
	if t.Sudo {
		finalCommand = fmt.Sprintf("sudo sh -c '%s'", finalCommand)
	}
	log.Printf("running command: %s", finalCommand)
	if err := session.Start(finalCommand); err != nil {
		return "", "", err
	}
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	go readerStreamer(stdoutReader, stdoutChan, "stdout")
	go readerStreamer(stderrReader, stderrChan, "stderr")
	if err := session.Wait(); err != nil {
		return "", "", err
	}
	stdout = <-stdoutChan
	stderr = <-stderrChan
	return
}

func readerStreamer(reader io.Reader, outputChan chan<- string, description string) {
	result := bytes.Buffer{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		result.Write([]byte(scanner.Text()))
		fmt.Printf("%s: %s\n", description, scanner.Text())
	}
	outputChan <- result.String()
}

func (t *Target) initClient() {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatalf("net.Dial: %v", err)
	}
	agentClient := agent.NewClient(conn)
	config := &ssh.ClientConfig{
		User: t.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	t.Client, err = ssh.Dial("tcp", t.Node, config)
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}
}

func sshTarget(target deployments.Target) *Target {
	if target, ok := target.(*Target); ok {
		return target
	}
	log.Fatal("Target is of the wrong type")
	return nil
}

func cniDeploy() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		// Deploy locally
		return nil
	}
	return runner
}

func kubeadmInit() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		t.UploadFile(caaspctl.KubeadmInitConfFile(), "/tmp/kubeadm.conf")
		if target := sshTarget(t); target != nil {
			target.ssh("systemctl", "enable", "--now", "docker")
			target.ssh("systemctl", "stop", "kubelet")
			target.ssh("kubeadm", "init", "--config", "/tmp/kubeadm.conf", "--skip-token-print")
			target.ssh("rm", "/tmp/kubeadm.conf")
		}
		return nil
	}
	return runner
}

func kubeadmJoin() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		// FIXME: ereslibre
		t.UploadFile(ConfigPath(caaspctl.MasterRole, t.Target()), "/tmp/kubeadm.conf")
		if target := sshTarget(t); target != nil {
			target.ssh("systemctl", "enable", "--now", "docker")
			target.ssh("systemctl", "stop", "kubelet")
			target.ssh("kubeadm", "join", "--config", "/tmp/kubeadm.conf")
			target.ssh("rm", "/tmp/kubeadm.conf")
		}
		return nil
	}
	return runner
}

func kubeletConfigure() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		t.UploadFileContents("/lib/systemd/system/kubelet.service", assets.KubeletService)
		t.UploadFileContents("/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", assets.KubeadmService)
		t.UploadFileContents("/etc/sysconfig/kubelet", assets.KubeletSysconfig)
		if target := sshTarget(t); target != nil {
			target.ssh("systemctl", "daemon-reload")
		}
		return nil
	}
	return runner
}

func kubeletEnable() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		if target := sshTarget(t); target != nil {
			target.ssh("systemctl", "enable", "kubelet")
		}
		return nil
	}
	return runner
}

func kubernetesUploadSecrets() deployments.Runner {
	runner := struct{ deployments.State }{}
	runner.DoRun = func(t deployments.Target) error {
		for _, file := range deployments.Secrets {
			t.UploadFile(file, path.Join("/etc/kubernetes", file))
		}
		return nil
	}
	return runner
}
