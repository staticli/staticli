package cali

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// Event holds the json structure for Docker API events
type Event struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// ProgressDetail records the progress achieved downloading an image
type ProgressDetail struct {
	Current int `json:"current,omitempty"`
	Total   int `json:"total,omitempty"`
}

// CreateResponse is the response from Docker API when pulling an image
type CreateResponse struct {
	Id             string         `json:"id"`
	Status         string         `json:"status"`
	ProgressDetail ProgressDetail `json:"progressDetail"`
	Progress       string         `json:"progress,omitempty"`
}

// DockerClient is a slimmed down implementation of the docker cli
type DockerClient struct {
	Cli      *client.Client
	HostConf *container.HostConfig
	NetConf  *network.NetworkingConfig
	Conf     *container.Config
	running  []string
}

// Init initialises the client
func (c *DockerClient) InitDocker() error {
	var cli *client.Client

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(dockerHost, "v1.22", nil, defaultHeaders)

	if err != nil {
		return fmt.Errorf("Could not connect to Docker daemon on %s: %s", dockerHost, err)
	}
	c.Cli = cli
	return nil
}

// NewDockerClient returns a new DockerClient initialised with the API object
func NewDockerClient() *DockerClient {
	c := new(DockerClient)
	c.SetDefaults()
	return c
}

// SetDefaults sets container, host and net configs to defaults. Called when instantiating a new client or can be called
// manually at any time to reset API configs back to empty defaults
func (c *DockerClient) SetDefaults() {
	c.HostConf = &container.HostConfig{Binds: []string{}}
	c.NetConf = &network.NetworkingConfig{}
	c.Conf = &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Env:          []string{},
	}
}

// SetHostConf sets the container.HostConfig struct for the new container
func (c *DockerClient) SetHostConf(h *container.HostConfig) {
	c.HostConf = h
}

// SetNetConf sets the network.NetworkingConfig struct for the new container
func (c *DockerClient) SetNetConf(n *network.NetworkingConfig) {
	c.NetConf = n
}

// SetConf sets the container.Config struct for the new container
func (c *DockerClient) SetConf(co *container.Config) {
	c.Conf = co
}

// AddBind adds a bind mount to the HostConfig
func (c *DockerClient) AddBind(bnd string) {
	c.HostConf.Binds = append(c.HostConf.Binds, bnd)
}

// AddEnvs adds an environment variable to the HostConfig
func (c *DockerClient) AddEnv(key, value string) {
	c.Conf.Env = append(c.Conf.Env, fmt.Sprintf("%s=%s", key, value))
}

// AddBinds adds multiple bind mounts to the HostConfig
func (c *DockerClient) AddBinds(bnds []string) {
	c.HostConf.Binds = append(c.HostConf.Binds, bnds...)
}

// AddEnvs adds multiple envs to the HostConfig
func (c *DockerClient) AddEnvs(envs []string) {
	c.Conf.Env = append(c.Conf.Env, envs...)
}

// SetBinds sets the bind mounts in the HostConfig
func (c *DockerClient) SetBinds(bnds []string) {
	c.HostConf.Binds = bnds
}

// SetEnvs sets the environment variables in the Conf
func (c *DockerClient) SetEnvs(envs []string) {
	c.Conf.Env = envs
}

// SetImage sets the image in Conf
func (c *DockerClient) SetImage(img string) {
	c.Conf.Image = img
}

// Privileged sets whether the container should run as privileged
func (c *DockerClient) Privileged(p bool) {
	c.HostConf.Privileged = p
}

// SetCmd sets the command to run in the container
func (c *DockerClient) SetCmd(cmd []string) {
	c.Conf.Cmd = cmd
}

// SetWorkDir sets the working directory of the container
func (c *DockerClient) SetWorkDir(wd string) {
	c.Conf.WorkingDir = wd
}

// BindFromGit creates a data container with a git clone inside and mounts its volumes inside your app container
// If there is no valid Git repo set in config, the noGit callback function will be executed instead
func (c *DockerClient) BindFromGit(cfg *GitCheckoutConfig, noGit func() error) error {
	cli := NewDockerClient()

	if err := cli.InitDocker(); err != nil {
		return err
	}

	if cfg.Repo != "" {
		// Build code from data volume
		git := cli.Git()

		if cfg.Image != "" {
			git.Image = cfg.Image
		}
		id, err := git.Checkout(cfg)

		if err != nil {
			return err
		}
		c.HostConf.VolumesFrom = []string{id}

		if cfg.RelPath != "" {
			c.SetWorkDir(path.Join(workdir, cfg.RelPath))
		}
	} else {
		// Execute callback
		noGit()
	}
	return nil
}

// StartContainer will create and start a container with logs and optional cleanup
func (c *DockerClient) StartContainer(rm bool, name string) (string, error) {
	log.WithFields(log.Fields{
		"image": c.Conf.Image,
		"envs":  fmt.Sprintf("%v", c.Conf.Env),
		"cmd":   fmt.Sprintf("%v", c.Conf.Cmd),
	}).Debug("Creating new container")

	if err := c.PullImage(c.Conf.Image); err != nil {
		return "", fmt.Errorf("Failed to fetch image: %s", err)
	}
	resp, err := c.Cli.ContainerCreate(context.Background(), c.Conf, c.HostConf, c.NetConf, name)

	if err != nil {
		return "", fmt.Errorf("Failed to create container: %s", err)
	}

	// Clean up on ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	go func() {
		<-ch
		log.Debug("Trapped ctrl+c")

		if err = c.DeleteContainer(resp.ID); err != nil {
			log.Errorf("Failed to remove container: %s", err)
		}
		os.Exit(1)
	}()
	log.WithFields(log.Fields{
		"image": c.Conf.Image,
		"id":    resp.ID[0:12],
	}).Debug("Starting new container")

	// Set the TTY size to match the host terminal
	fd := int(os.Stdin.Fd())

	if !nonInteractive && terminal.IsTerminal(fd) {
		// While we have a container running, create a buffer for the pscli logs
		logBuffer := bufio.NewWriter(os.Stdout)
		log.SetOutput(logBuffer)
		// Write buffer to stdout once detatched from container
		defer logBuffer.Flush()
		// Reset logs to stdout after conection is closed
		defer log.SetOutput(os.Stdout)

		// If we have an interactive terminal then use it!
		ca := types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		}
		hijack, err := c.Cli.ContainerAttach(context.Background(), resp.ID, ca)
		defer hijack.Conn.Close()

		if err != nil {
			return resp.ID, fmt.Errorf("Failed to start container: %s", err)
		}
		oldState, err := terminal.MakeRaw(fd)
		defer terminal.Restore(fd, oldState)

		if err != nil {
			panic(err)
		}

		if err := c.Cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			return resp.ID, fmt.Errorf("Failed to start container: %s", err)
		}

		// Start stdin reader
		go func() {
			defer terminal.Restore(fd, oldState)
			defer hijack.Conn.Close()

			if _, err := io.Copy(hijack.Conn, os.Stdin); err != nil {
				log.Errorf("Write error: %s", err)
			}
		}()

		tw, th, _ := terminal.GetSize(fd)

		if err := c.Cli.ContainerResize(context.Background(), resp.ID, types.ResizeOptions{Height: uint(th), Width: uint(tw)}); err != nil {
			return resp.ID, fmt.Errorf("Failed to start container: %s", err)
		}

		// Start stdout writer
		if _, err := io.Copy(os.Stdout, hijack.Conn); err != nil {
			log.Errorf("Read error: %s", err)
		}
	} else {
		// No terminal, then just pump out the log output
		if err := c.Cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			return resp.ID, fmt.Errorf("Failed to start container: %s", err)
		}
		log.WithFields(log.Fields{
			"image": c.Conf.Image,
			"id":    resp.ID[0:12],
		}).Debug("Fetching log stream")
		logOptions := types.ContainerLogsOptions{Follow: true, ShowStdout: true, ShowStderr: true}
		ls, err := c.Cli.ContainerLogs(context.Background(), resp.ID, logOptions)

		if err != nil {
			return resp.ID, fmt.Errorf("Failed to get container logs: %s", err)
		}

		_, err = io.Copy(os.Stdout, ls)
		if err != nil {
			return resp.ID, fmt.Errorf("Failed to get container logs: %s", err)
		}
	}
	// Container has finished running. Get its exit code
	inspect, err := c.Cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return resp.ID, fmt.Errorf("Failed to inspect Docker container: %s", err)
	}

	if rm {

		if err = c.DeleteContainer(resp.ID); err != nil {
			return resp.ID, fmt.Errorf("Failed to remove container: %s", err)
		}
	}

	if inspect.State.ExitCode != 0 {
		return resp.ID, fmt.Errorf("Non-zero exit status from Docker container")
	}
	return resp.ID, nil
}

// ContainerExists determines if the container with this name exist
func (c *DockerClient) ContainerExists(name string) bool {
	_, err := c.Cli.ContainerInspect(context.Background(), name)

	// Fairly safe assumption
	if err != nil {
		return false
	}
	return true
}

// DeleteContainer - Delete a container
func (c *DockerClient) DeleteContainer(id string) error {
	log.WithFields(log.Fields{
		"id": id[0:12],
	}).Debug("Removing container")

	if err := c.Cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("Failed to remove container: %s", err)
	}
	return nil
}

// ImageExists determines if an image exist locally
func (c *DockerClient) ImageExists(image string) bool {
	log.WithFields(log.Fields{
		"image": image,
	}).Debug("Checking if image exists locally")

	_, _, err := c.Cli.ImageInspectWithRaw(context.Background(), image)

	// Safe assumption?
	if err != nil {
		log.WithFields(log.Fields{
			"image": image,
		}).Debugf("Error inspecting image: %s", err)
		return false
	}
	return true
}

// PullImage - Pull an image locally
func (c *DockerClient) PullImage(image string) error {

	if !c.ImageExists(image) {
		log.WithFields(log.Fields{
			"image": image,
		}).Info("Pulling image layers... please wait")

		resp, err := c.Cli.ImagePull(context.Background(), image, types.ImagePullOptions{})

		if err != nil {
			return fmt.Errorf("API could not fetch \"%s\": %s", image, err)
		}
		scanner := bufio.NewScanner(resp)
		var cr CreateResponse
		bar := pb.New(1)
		// Send progress bar to stderr to keep stdout clean when piping
		bar.Output = os.Stderr
		bar.ShowCounters = true
		bar.ShowTimeLeft = false
		bar.ShowSpeed = false
		bar.Prefix("          ")
		bar.Postfix("          ")
		started := false

		for scanner.Scan() {
			txt := scanner.Text()
			byt := []byte(txt)

			if err := json.Unmarshal(byt, &cr); err != nil {
				return fmt.Errorf("Error decoding json from create image API: %s", err)
			}

			if cr.Status == "Downloading" {

				if !started {
					fmt.Print("\n")
					bar.Total = int64(cr.ProgressDetail.Total)
					bar.Start()
					started = true
				}
				bar.Total = int64(cr.ProgressDetail.Total)
				bar.Set(cr.ProgressDetail.Current)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Failed to get logs: %s", err)
		}
		bar.Finish()
		fmt.Print("\n")
	}
	return nil
}
