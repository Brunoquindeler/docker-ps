package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

var (
	containerListOptions = types.ContainerListOptions{}

	usageMessage = "Usage: (docker ps) or (docker ps -a)"
)

func main() {
	argsLength := len(os.Args)
	if argsLength < 2 || argsLength > 3 {
		fmt.Println(usageMessage)
		os.Exit(1)
	}

	arg1 := os.Args[1]
	if arg1 != "ps" {
		fmt.Println(usageMessage)
		os.Exit(1)
	}

	if argsLength > 2 {
		arg2 := os.Args[2]
		if arg2 != "-a" {
			fmt.Println(usageMessage)
			os.Exit(1)
		}
		containerListOptions.All = true
	}

	ctx := context.Background()
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	containers, err := c.ContainerList(ctx, containerListOptions)
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"})

	for _, container := range containers {
		cmd := fmt.Sprintf(`"%s..."`, container.Command)
		if len(container.Command) > 19 {
			cmd = fmt.Sprintf(`"%s..."`, container.Command[:19])
		}

		timestamp := container.Created
		t := time.Unix(timestamp, 0)
		diffStr := humanize.RelTime(t, time.Now(), "", "") + "ago"

		ports := ""
		if len(container.Ports) > 0 {
			ports = fmt.Sprintf("%s:%d->%d/%s", container.Ports[0].IP, container.Ports[0].PrivatePort, container.Ports[0].PublicPort, container.Ports[0].Type)
		}

		table.Append([]string{container.ID[:12], container.Image, cmd, diffStr, container.Status, ports, container.Names[0][1:]})
	}

	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetCenterSeparator("")

	table.Render()
}
