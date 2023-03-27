package mySSH

import (
	"bytes"
	"fmt"

	"io"
	"net"
	"os"
	"strings"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

type NodeItem struct {
	NodeIndex    string
	IPaddress    string
	UserName     string
	Password     string
	AbsolutePath string
}

type NodeOperationItem struct {
	OperationName           string
	OperationContent        string
	OperationRefinedContent string
	OperationPrefix         string
	OperationPostfix        string
	IsSubstitute            bool
}

var LegalOperationName = []string{"copy", "command", "copyN", "getN", "commandD"}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetRealNameFromPattern(oldString string, index string) string {
	indexPos := strings.Index(oldString, "%")

	if indexPos > -1 {
		preName := oldString[:indexPos]
		afterName := oldString[indexPos+1:]
		refinedName := preName + index + afterName
		return refinedName
	}
	return oldString
}
func DirectImplementWithoutPath(currNode NodeItem, command string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()
	/* excute the command */
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		fmt.Printf("Failed to excute command: " + err.Error() + "\n")
		fmt.Fprintf(printOutput, "Failed to excute command: "+err.Error()+"\n")
		return err
	}

	fmt.Printf("Output: " + stdoutBuf.String() + "\n")
	fmt.Fprintf(printOutput, "Output: "+stdoutBuf.String()+"\n")

	return nil
}

func DirectImplement(currNode NodeItem, command string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()
	/* excute the command */
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run("cd " + currNode.AbsolutePath + "&&" + command)
	if err != nil {
		fmt.Printf("Failed to excute command: " + err.Error() + "\n")
		fmt.Fprintf(printOutput, "Failed to excute command: "+err.Error()+"\n")
		return err
	}

	fmt.Printf("Output: " + stdoutBuf.String() + "\n")
	fmt.Fprintf(printOutput, "Output: "+stdoutBuf.String()+"\n")

	return nil
}
func TransferFile(currNode NodeItem, filePath string, destName string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()

	dest := currNode.AbsolutePath + destName
	err = scp.CopyPath(filePath, dest, session)
	if err != nil {
		fmt.Printf("Transfering file error with the destination: " + err.Error())
		fmt.Fprintf(printOutput, "Transfering file error with the destination: "+err.Error())
		return err
	}
	return nil
}

func DownloadFile(currNode NodeItem, remotePath string, destLocalPath string, printOutput io.Writer) error {
	sshConfig := &ssh.ClientConfig{
		User: currNode.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(currNode.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", currNode.IPaddress+":22", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to dial: "+err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: " + err.Error())
		fmt.Fprintf(printOutput, "Failed to create session: "+err.Error())
		return err
	}
	defer session.Close()
	r, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(destLocalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := "cat " + remotePath
	if err := session.Start(cmd); err != nil {
		return err
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}

func TestOutput() {
	fmt.Println("hello world")
}
