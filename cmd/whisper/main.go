package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

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
			Name:     "create",
			Usage:    "create a whisper secret",
			Category: "client",
			Before:   initClient,
			Action:   create,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "secret",
					Aliases: []string{"s"},
					Usage:   "input the secret as a string on the command line",
				},
				&cli.IntFlag{
					Name:    "generate",
					Aliases: []string{"g"},
					Usage:   "generate a random secret of the specified length",
				},
				&cli.StringFlag{
					Name:    "in",
					Aliases: []string{"i", "u", "upload"},
					Usage:   "upload a file as the secret contents",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "specify a password to access the secret",
				},
				&cli.IntFlag{
					Name:    "generate-password",
					Aliases: []string{"G", "gp"},
					Usage:   "generate a random password of the specified length",
				},
				&cli.IntFlag{
					Name:    "accesses",
					Aliases: []string{"a"},
					Usage:   "set number of allowed accesses; default 1, -1 for unlimited until expiration",
				},
				&cli.DurationFlag{
					Name:    "lifetime",
					Aliases: []string{"l", "e", "expires", "expires-after"},
					Usage:   "specify the lifetime of the secret before it is deleted",
				},
				&cli.BoolFlag{
					Name:    "b64encoded",
					Aliases: []string{"b", "b64"},
					Usage:   "specify if the secret is base64 encoded (true if uploading a file, false if generated)",
				},
			},
		},
		{
			Name:      "fetch",
			Usage:     "fetch a whisper secret by its token",
			ArgsUsage: "token",
			Category:  "client",
			Before:    initClient,
			Action:    fetch,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "specify a password to access the secret",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o", "d", "download"},
					Usage:   "download the secret to a file or to a directory",
				},
			},
		},
		{
			Name:      "destroy",
			Usage:     "destroy a whisper secret by its token",
			ArgsUsage: "token",
			Category:  "client",
			Before:    initClient,
			Action:    destroy,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "specify a password to access the secret",
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

func create(c *cli.Context) (err error) {
	// Create the request
	req := &v1.CreateSecretRequest{
		Password: c.String("password"),
		Accesses: c.Int("accesses"),
		Lifetime: v1.Duration(c.Duration("lifetime")),
	}

	// Add the secret to the request via one of the command line options
	switch {
	case c.String("secret") != "":
		if c.Int("generate") != 0 || c.String("in") != "" {
			return cli.Exit("specify only one of secret, generate, or in path", 1)
		}

		// Basic secret provided via the CLI
		req.Secret = c.String("secret")
		req.IsBase64 = c.Bool("b64encoded")

	case c.String("in") != "":
		if c.Int("generate") != 0 {
			// The check for secret has already been done
			return cli.Exit("specify only one of secret, generate, or in path", 1)
		}

		// Load the secret as base64 encoded data from a file
		var data []byte
		if data, err = ioutil.ReadFile(c.String("in")); err != nil {
			return cli.Exit(err, 1)
		}
		req.Filename = filepath.Base(c.String("in"))
		req.Secret = base64.StdEncoding.EncodeToString(data)
		req.IsBase64 = true

	case c.Int("generate") != 0:
		// Generate a random secret of the specified length
		if req.Secret, err = generateRandomSecret(c.Int("generate")); err != nil {
			return cli.Exit(err, 1)
		}
		req.IsBase64 = false

	default:
		// No secret was specified at all?
		return cli.Exit("specify at least one of secret, generate, or in path", 1)
	}

	// Handle password generation if requested
	if gp := c.Int("generate-password"); gp > 0 {
		if req.Password != "" {
			return cli.Exit("specify either password or generate password, not both", 1)
		}
		if req.Password, err = generateRandomSecret(gp); err != nil {
			return cli.Exit(err, 1)
		}
		// Print the password so that it can be used to retrieve the secret later
		fmt.Printf("Password for retrieval: %s\n", req.Password)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *v1.CreateSecretReply
	if rep, err = client.CreateSecret(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func fetch(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		return cli.Exit("specify one token to fetch the secret for", 1)
	}

	token := c.Args().First()
	password := c.String("password")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *v1.FetchSecretReply
	if rep, err = client.FetchSecret(ctx, token, password); err != nil {
		return cli.Exit(err, 1)
	}

	// Figure out where to write the file to; if out is a directory, write the
	var path string

	// If the user has specified an output location, handle it.
	if out := c.String("out"); out != "" {
		var isDir bool
		if isDir, err = isDirectory(out); err == nil && isDir {
			if rep.Filename != "" {
				path = filepath.Join(out, rep.Filename)
			} else {
				path = filepath.Join(out, "secret.dat")
			}
		} else {
			path = out
		}
	} else {
		// If the user didn't specify an out location and the response has a filename,
		// write the file with the specified filename in the current working directory.
		path = rep.Filename
	}

	// If we've discovered a path to write the file to, write it there, decoding the
	// data as necessary from base64. Otherwise print the json to stdout and exit.
	if path != "" {
		fmt.Println(path)
		var data []byte
		if rep.IsBase64 {
			if data, err = base64.StdEncoding.DecodeString(rep.Secret); err != nil {
				return cli.Exit(err, 1)
			}
		} else {
			data = []byte(rep.Secret)
		}

		if err = ioutil.WriteFile(path, data, 0644); err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Printf("secret written to %s\n", path)
		return nil
	}

	// Simply print the JSON response as the last case.
	// TODO: should we provide a flag to just print the secret for copy and paste?
	return printJSON(rep)
}

func destroy(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		return cli.Exit("specify one token to fetch the secret for", 1)
	}

	token := c.Args().First()
	password := c.String("password")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *v1.DestroySecretReply
	if rep, err = client.DestroySecret(ctx, token, password); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(rep)
}

func status(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *v1.StatusReply
	if rep, err = client.Status(ctx); err != nil {
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

func generateRandomSecret(n int) (s string, err error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_%=+"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func isDirectory(path string) (isDir bool, err error) {
	var fi fs.FileInfo
	if fi, err = os.Stat(path); err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func printJSON(v interface{}) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(v, "", "  "); err != nil {
		return cli.Exit("could not marshal json response", 2)
	}
	fmt.Println(string(data))
	return nil
}
