package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	// Log ...
	Log = log.Log

	// SetLogger sets a concrete logging implementation for all deferred Loggers.
	SetLogger = log.SetLogger
)

// Config ...
type Config struct {
	IDRsa     string
	ImageName string
}

// Build creates a docker image with an id_rsa
func (c *Config) Build() error {
	//
	logger := Log.WithName("builder")

	//
	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Error(err, "Failed while instantiating docker client")
		return err
	}

	//
	dockerContext, err := c.createDockerContext()
	if err != nil {
		logger.Error(err, "Failed while seting up docker images")
		return err
	}

	_, err = cli.ImageBuild(
		context.Background(),
		dockerContext,
		types.ImageBuildOptions{
			Context: dockerContext,
			Tags:    []string{c.ImageName},
			NoCache: true,
		})
	if err != nil {
		logger.Error(err, "Failed while building docker image")
		return err
	}

	return nil
}

func (c *Config) createDockerContext() (io.Reader, error) {
	//
	logger := Log.WithName("builder")

	//
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	dockerfile := `FROM scratch
				  COPY /id_rsa /id_rsa
				  `
	var files = []struct {
		Name, Body string
	}{
		{"id_rsa", c.IDRsa},
		{"Dockerfile", dockerfile},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			logger.Error(err, "Failed while writing tarball header")
			return nil, err
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			logger.Error(err, "Failed while writing tarball body")
			return nil, err
		}
	}
	// @TODO: Should i defer this?
	if err := tw.Close(); err != nil {
		logger.Error(err, "Failed while closing the tarball")
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
