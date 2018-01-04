package collector

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/prometheus/common/log"
	"strings"
)


func getGlusterBinary(glusterPath string) (string, error) {

	switch glusterPath {
	// NoDefine
	case "":
		out, err := exec.Command("which","gluster").Output()

		// Trim `out` with '\n'
		rout := strings.TrimSuffix(string(out), "\n")

		if err != nil {
			log.Fatal("Please Make sure Gluster installed correctly. Cannot find gluster binary.")
			return rout, err
		}
		return rout, err
		// Has Define
	default:
		// Check Exists
		_, err := PathExists(glusterPath)
		if err != nil {
			return "", err
		}
		return glusterPath, nil
	}
}

func execGlusterCommand(arg ...string) (*bytes.Buffer, error) {
	glusterCmd, getErr := getGlusterBinary("")
	if getErr != nil {
		log.Error(getErr)
	}

	stdoutBuffer := &bytes.Buffer{}
	argXML := append(arg, "--xml")
	glusterExec := exec.Command(glusterCmd, argXML...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		log.Errorf("tried to execute %v and got error: %v", arg, err)
		return stdoutBuffer, err
	}
	return stdoutBuffer, nil
}

// ContainsVolume checks a slice if it contains an element
func ContainsVolume(slice []string, element string) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
