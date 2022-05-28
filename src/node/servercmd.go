package main

import (
	"github.com/urfave/cli/v2"
)

var (
	nodeIdFlag = &cli.IntFlag{
		Name:     "id",
		Usage:    "id",
		Required: true,
	}
	nodeSubCommand = &cli.Command{
		Name:        "node",
		Usage:       "start node node",
		Description: "start node node",
		ArgsUsage:   "<id>",
		Flags: []cli.Flag{
			nodeIdFlag,
		},
		Action: func(c *cli.Context) error {
			nodeId := c.Int("id")
			powID = nodeId

			for i := 0; i < 10; i++ {
				go func(nodeID int) {
					server := NewServer(nodeID)
					server.Start()
				}(powID + i)
			}
			Work()

			return nil
		},
	}
	clientSubCommand = &cli.Command{
		Name:        "client",
		Usage:       "start node client",
		Description: "start node client",
		ArgsUsage:   "",
		Action: func(c *cli.Context) error {
			client := NewClient()
			client.Start()
			return nil
		},
	}
	PBFTCommand = &cli.Command{
		Name:        "pbft",
		Usage:       "pbft commands",
		ArgsUsage:   "",
		Category:    "pbft Commands",
		Description: "",
		Subcommands: []*cli.Command{
			nodeSubCommand,
			clientSubCommand,
		},
	}
)
