package cali

import (
	"crypto/md5"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// GitCheckoutConfig is input for Git.Checkout
type GitCheckoutConfig struct {
	Repo, Branch, RelPath, Image string
}

const gitImage = "indiehosters/git:latest"

// Git returns a new instance
func (c *DockerClient) Git() *Git {
	return &Git{c: c, Image: gitImage}
}

// Git is used to interact with containerised git
type Git struct {
	c     *DockerClient
	Image string
}

// GitCheckout will create and start a container, checkout repo and leave container stopped
// so volume can be imported
func (g *Git) Checkout(cfg *GitCheckoutConfig) (string, error) {
	containerName, err := cfg.GetContainerName()
	if err != nil {
		return "", fmt.Errorf("Failed to create data container for %s: %s", cfg.Repo, err)
	}

	if g.c.ContainerExists(containerName) {
		log.Infof("Existing data container found: %s", containerName)

		if _, err := g.Pull(containerName); err != nil {
			log.Warnf("Git pull error: %s", err)
			return containerName, err
		}
		return containerName, nil
	} else {
		log.WithFields(log.Fields{
			"git_url": cfg.Repo,
			"image":   g.Image,
		}).Info("Creating data containers")

		co := container.Config{
			Cmd:          []string{"clone", cfg.Repo, "-b", cfg.Branch, "--depth", "1", "."},
			Image:        gitImage,
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			WorkingDir:   "/tmp/workspace",
			Entrypoint:   []string{"git"},
		}
		hc := container.HostConfig{
			Binds: []string{
				"/tmp/workspace",
				fmt.Sprintf("%s/.ssh:/root/.ssh", os.Getenv("HOME")),
			},
		}
		nc := network.NetworkingConfig{}

		g.c.SetConf(&co)
		g.c.SetHostConf(&hc)
		g.c.SetNetConf(&nc)

		id, err := g.c.StartContainer(false, containerName)

		if err != nil {
			return "", fmt.Errorf("Failed to create data container for %s: %s", cfg.Repo, err)
		}
		return id, nil
	}
}

func (g *Git) Pull(name string) (string, error) {
	co := container.Config{
		Cmd:          []string{"pull"},
		Image:        g.Image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/tmp/workspace",
		Entrypoint:   []string{"git"},
	}
	hc := container.HostConfig{
		VolumesFrom: []string{name},
		Binds: []string{
			fmt.Sprintf("%s/.ssh:/root/.ssh", os.Getenv("HOME")),
		},
	}
	nc := network.NetworkingConfig{}

	g.c.SetConf(&co)
	g.c.SetHostConf(&hc)
	g.c.SetNetConf(&nc)

	return g.c.StartContainer(true, "")
}

// GetContainerName returns a container name for provided Git config
func (cfg GitCheckoutConfig) GetContainerName() (string, error) {
	repoName, err := repoNameFromUrl(cfg.Repo)
	if err != nil {
		return "", fmt.Errorf("Failed to get container name for %s: %s", cfg.Repo, err)
	}
	containerName := repoName

	if cfg.RelPath == "." || cfg.RelPath == "" {
		containerName = fmt.Sprintf("data_%s_%s_%x",
			repoName,
			strings.Replace(cfg.Branch, "/", "-", -1),
			md5.Sum([]byte(cfg.Repo)),
		)
	} else {
		containerName = fmt.Sprintf("data_%s_%s_%s_%x",
			repoName,
			strings.Replace(cfg.RelPath, "/", "-", -1),
			strings.Replace(cfg.Branch, "/", "-", -1),
			md5.Sum([]byte(cfg.Repo)),
		)
	}

	return containerName, nil
}

// repoNameFromUrl takes a git repo URL and returns a string
// representing the repository name
func repoNameFromUrl(url string) (string, error) {

	// Strip out the https:// or git:// protocol
	protocolRe := regexp.MustCompile("^.*//")
	url = protocolRe.ReplaceAllString(url, "")

	// Remove trailing .git
	dotGitRe := regexp.MustCompile(".git$")
	url = dotGitRe.ReplaceAllString(url, "")

	// Remove user@
	userAtRe := regexp.MustCompile(".*@")
	url = userAtRe.ReplaceAllString(url, "")

	// Actual regex for container names: [a-zA-Z0-9][a-zA-Z0-9_.-]
	// https://github.com/moby/moby/blob/master/daemon/names/names.go
	// but to simplify, as we're doing an inverse match, just use [a-zA-Z0-9]
	nonContainerRe := regexp.MustCompile("[^a-zA-Z0-9]")
	repoName := nonContainerRe.ReplaceAllString(url, "-")

	return repoName, nil
}
