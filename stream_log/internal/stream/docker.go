package stream

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
)

func DockerStreamLogs(ctx context.Context, containerName string, streamChan chan<- string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
		Timestamps: false,
	}

	logReader, err := cli.ContainerLogs(ctx, containerName, options)
	if err != nil {
		return err
	}
	defer logReader.Close()

	log.Println("logReader connected successfuly")
	log.Println("Reading logs from Docker container")
	log.Println(containerName, options)

	scanner := bufio.NewScanner(logReader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		case streamChan <- scanner.Text():
			log.Println("receive log: ", scanner.Text())
		}
	}
	return scanner.Err()
}
