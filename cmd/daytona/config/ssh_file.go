// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var sshHomeDir string

func ensureSshFilesLinked() error {
	// Make sure ~/.ssh/config file exists if not create it
	sshDir := filepath.Join(sshHomeDir, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(sshDir, 0700)
		if err != nil {
			return err
		}
		err = os.WriteFile(configPath, []byte{}, 0600)
		if err != nil {
			return err
		}
	}

	// Make sure daytona_config file exists
	daytonaConfigPath := filepath.Join(sshDir, "daytona_config")

	_, err = os.Stat(daytonaConfigPath)
	if os.IsNotExist(err) {
		err := os.WriteFile(daytonaConfigPath, []byte{}, 0600)
		if err != nil {
			return err
		}
	}

	// Make sure daytona_config is included
	configFile := filepath.Join(sshDir, "config")
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		err := os.WriteFile(configFile, []byte("Include daytona_config\n\n"), 0600)
		if err != nil {
			return err
		}
	} else {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), "Include daytona_config\n\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config", "")
		newContent = "Include daytona_config\n\n" + newContent
		err = os.WriteFile(configFile, []byte(newContent), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnlinkSshFiles() error {
	sshDirPath := filepath.Join(sshHomeDir, ".ssh")
	sshConfigPath := filepath.Join(sshDirPath, "config")
	daytonaConfigPath := filepath.Join(sshDirPath, "daytona_config")

	// Remove the include line from the config file
	_, err := os.Stat(sshConfigPath)
	if os.IsExist(err) {
		content, err := os.ReadFile(sshConfigPath)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), "Include daytona_config\n\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config", "")
		err = os.WriteFile(sshConfigPath, []byte(newContent), 0600)
		if err != nil {
			return err
		}
	}

	// Remove the daytona_config file
	_, err = os.Stat(daytonaConfigPath)
	if os.IsExist(err) {
		err = os.Remove(daytonaConfigPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add ssh entry

func generateSshConfigEntry(profileId, workspaceId, projectName, knownHostsPath string, gpgForward bool) (string, error) {
	daytonaPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	tab := "\t"
	projectHostname := GetProjectHostname(profileId, workspaceId, projectName)

	config := fmt.Sprintf("Host %s\n"+
		tab+"User daytona\n"+
		tab+"StrictHostKeyChecking no\n"+
		tab+"UserKnownHostsFile %s\n"+
		tab+"ProxyCommand \"%s\" ssh-proxy %s %s %s\n"+
		tab+"ForwardAgent yes\n", projectHostname, knownHostsPath, daytonaPath, profileId, workspaceId, projectName)
	if gpgForward {
		localSocket, _ := getLocalGPGSocket()
		remoteSocket, _ := getRemoteGPGSocket(projectHostname)
		config += fmt.Sprintf(
			tab+"StreamLocalBindUnlink yes\n"+
				tab+"RemoteForward %s:%s\n\n", remoteSocket, localSocket)
	} else {
		config += "\n"
	}

	return config, nil
}

func EnsureSshConfigEntryAdded(profileId, workspaceName, projectName string, gpgForward bool) error {
	err := ensureSshFilesLinked()
	if err != nil {
		return err
	}

	sshDir := filepath.Join(sshHomeDir, ".ssh")
	configPath := filepath.Join(sshDir, "daytona_config")

	knownHostsFile := getKnownHostsFile()

	// Read existing content from the file
	existingContent, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Generate SSH config entry without GPG forwarding
	err = appendSshConfigEntry(configPath, profileId, workspaceName, projectName, knownHostsFile, false, existingContent)
	if err != nil {
		return err
	}

	if gpgForward {
		// Generate SSH config entry with GPG forwarding and override previous config
		err = appendSshConfigEntry(configPath, profileId, workspaceName, projectName, knownHostsFile, true, existingContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func getKnownHostsFile() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}

func appendSshConfigEntry(configPath, profileId, workspaceName, projectName, knownHostsFile string, gpgForward bool, existingContent []byte) error {
	data, err := generateSshConfigEntry(profileId, workspaceName, projectName, knownHostsFile, gpgForward)
	if err != nil {
		return err
	}

	// Check if the entry already exists
	if strings.Contains(string(existingContent), data) {
		return nil
	}

	// Combine the new data with existing content
	newData := data + string(existingContent)

	// Open the file for writing
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(newData)
	return err
}

func getLocalGPGSocket() (string, error) {
	cmd := exec.Command("gpgconf", "--list-dir", "agent-extra-socket")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get local GPG socket: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func getRemoteGPGSocket(projectHostname string) (string, error) {
	cmd := exec.Command("ssh", projectHostname, "gpgconf --list-dir agent-socket")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote GPG socket: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func ExportGPGKey(keyID, projectHostname string) error {
	exportCmd := exec.Command("gpg", "--export", keyID)
	var output bytes.Buffer
	exportCmd.Stdout = &output

	if err := exportCmd.Run(); err != nil {
		return err
	}

	importCmd := exec.Command("ssh", projectHostname, "gpg --import")
	importCmd.Stdin = &output

	return importCmd.Run()
}

func RemoveWorkspaceSshEntries(profileId, workspaceId string) error {
	sshDir := filepath.Join(sshHomeDir, ".ssh")
	configPath := filepath.Join(sshDir, "daytona_config")

	// Read existing content from the file
	existingContent, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}

	regex := regexp.MustCompile(fmt.Sprintf(`Host %s-%s-\w+\n(?:\t.*\n?)*`, profileId, workspaceId))
	newContent := regex.ReplaceAllString(string(existingContent), "")

	newContent = strings.Trim(newContent, "\n")

	// Open the file for writing
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(newContent)
	if err != nil {
		return err
	}

	return nil
}

func GetProjectHostname(profileId, workspaceId, projectName string) string {
	return fmt.Sprintf("%s-%s-%s", profileId, workspaceId, projectName)
}

func init() {
	if runtime.GOOS == "windows" {
		sshHomeDir = os.Getenv("USERPROFILE")
	} else {
		sshHomeDir = os.Getenv("HOME")
	}
}
