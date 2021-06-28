package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	whisper "github.com/rotationalio/whisper/pkg"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/urfave/cli/v2"
)

func main() {
	// Load the .env file if it exits
	godotenv.Load()

	// Instantiate the CLI application
	app := cli.NewApp()
	app.Name = "whisper"
	app.Usage = "interactions with the whisper secret manager"
	app.Version = whisper.Version()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e", "url", "u"},
			Usage:   "endpoint to connect to the whisper service on",
			EnvVars: []string{"WHISPER_ENDPOINT", "WHISPER_URL"},
			Value:   "http://localhost:8318",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run the whisper server",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"a"},
					Usage:   "address to bind the whisper server on",
					EnvVars: []string{"WHISPER_BIND_ADDR"},
				},
			},
		},
		{
			Name:     "status",
			Usage:    "get the whisper server status",
			Category: "client",
			Before:   initClient,
			Action:   status,
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// Server Actions
//===========================================================================

func serve(c *cli.Context) (err error) {
	// Create server configuration
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Update from CLI flags
	conf.BindAddr = c.String("addr")
	if c.Bool("no-secure") {
		conf.UseTLS = false
	}

	// Create and run the whisper server
	var server *whisper.Server
	if server, err = whisper.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = server.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

//===========================================================================
// Client Actions
//===========================================================================

var client v1.Service

func status(c *cli.Context) (err error) {
	var rep *v1.StatusReply
	if rep, err = client.Status(); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(rep)
}

//===========================================================================
// Helper Functions
//===========================================================================

func initClient(c *cli.Context) (err error) {
	if client, err = v1.New(c.String("endpoint")); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func printJSON(v interface{}) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(v, "", "  "); err != nil {
		return cli.Exit("could not marshal json response", 2)
	}
	fmt.Println(string(data))
	return nil
}
