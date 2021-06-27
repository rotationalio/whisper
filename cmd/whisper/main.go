package main

import (
	"os"

	"github.com/joho/godotenv"
	whisper "github.com/rotationalio/whisper/pkg"
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
		&cli.BoolFlag{
			Name:    "no-secure",
			Aliases: []string{"S"},
			Usage:   "don't connect with TLS (e.g. for development)",
			EnvVars: []string{"WHISPER_NOTLS", "WHISPER_CLIENT_INSECURE"},
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e", "url", "u"},
			Usage:   "endpoint to connect to the whisper service on",
			EnvVars: []string{"WHISPER_ENDPOINT", "WHISPER_URL"},
			Value:   "localhost:8318",
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
	}

	app.Run(os.Args)
}

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
