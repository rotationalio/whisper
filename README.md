# Whisper

[![Go Reference](https://pkg.go.dev/badge/github.com/rotationalio/whisper.svg)](https://pkg.go.dev/github.com/rotationalio/whisper)
[![Go Report Card](https://goreportcard.com/badge/github.com/rotationalio/whisper)](https://goreportcard.com/report/github.com/rotationalio/whisper)
![GitHub Actions CI](https://github.com/rotationalio/whisper/actions/workflows/build.yaml/badge.svg?branch=main)
![GitHub Actions CD](https://github.com/rotationalio/whisper/actions/workflows/release.yaml/badge.svg)
[![codecov](https://codecov.io/gh/rotationalio/whisper/branch/main/graph/badge.svg?token=64KYN8JYL4)](https://codecov.io/gh/rotationalio/whisper)


**There are many secrets management utilities, this one is ours … shhh**

## Command Line Application

The `whisper` CLI program makes it easy to interact with the Whisper API from the command line and securely share small files such as certificates or configurations (up to 48MiB for now).

To install the CLI application, download the appropriate binary from the [releases page](https://github.com/rotationalio/whisper/releases) and extract to a directory in your `$PATH` such as `~/bin`. Alternatively, if you have `$GOPATH/bin` in your `$PATH` you can fetch and install the binary with:

```
$ go get github.com/rotationalio/whisper/...
```

Once installed, make sure that you can execute the `whisper` command as follows:

```
$ whisper --version
```

You should occassionally check that the server version matches your binary version:

```
$ whisper status
{
  "status": "ok",
  "timestamp": "2021-07-15T18:24:50.343612976Z",
  "version": "1.0"
}
```

If the version doesn't match, please download the newest release of Whisper!

### Creating and Fetching a Secret

The simplest way to create a secret is as follows:

```
$ whisper create -s "the eagle flies at midnight"
{
  "token": "2nmwJzFnZ_iaa71wtjHKJbf-w_-P-g_qSi9qox3BfsY",
  "expires": "2021-07-22T18:15:33.459874936Z"
}
```

Then to fetch the secret:

```
$ whisper fetch 2nmwJzFnZ_iaa71wtjHKJbf-w_-P-g_qSi9qox3BfsY
{
  "secret": "the eagle flies at midnight",
  "is_base64": false,
  "created": "2021-07-15T18:15:33.459874936Z",
  "accesses": 1
}
```

This secret is a one time secret, meaning that once you fetch the secret with the token it will be deleted. The expiration date says when the secret will be automatically destroyed, and it must be fetched by this date. If you try to fetch the secret again, you'll get a `404 Not Found` error.

To change the number of access attempts to 7 and the time to live to 24 hours, use the following arguments:

```
$ whisper create -s "does the red robin crow at dusk?" -a 7 -l 24h
{
  "token": "Fv9dbbZXimubAx-snLz4lDoaRP779fqnTQyIGFAgBKk",
  "expires": "2021-07-16T19:29:21.77850402Z"
}
```

You can now access the secret 7 times before it's deleted or for 24 hours, whichever comes first. Specify the lifetime argument `-l` as a Golang duration `#h#m#s#ms#us#ns`. Note that the expiration must be at least 1 minute in the future. If you specify a negative number of attempts, e.g. `-a -1` then the secret can be accessed an unlimited number of times before the expiration time.

```
$ whisper fetch Fv9dbbZXimubAx-snLz4lDoaRP779fqnTQyIGFAgBKk
{
  "secret": "does the red robin crow at dusk?",
  "is_base64": false,
  "created": "2021-07-15T19:29:21.77850402Z",
  "accesses": 3
}
```

You can also create a random secret, e.g. if you're generating a random password as follows:

```
$ whisper create -G 14
```

The number following the `-G` argument specifies how long the generated secret should be.

### Password protecting secrets

If you'd like to add a passphrase to further control access to a secret, you can do so using the `-p` or `-g` commands, where `-p` allows you to specify the password and `-g N` generates a random password of length `N`:

```
$ whisper create -g 16 -s "The chimera has three sets of teeth"
Password for retrieval: +vDCk7T1wx%b5RPu
{
  "token": "5PeRU1spKSs2VjuUXv6gjo074QFR1Mcco9VGpweULEE",
  "expires": "2021-07-22T19:38:34.877795987Z"
}
```

To fetch the secret you must specify the correct password, otherwise the command will error with a `401 Unauthorized`.

```
$ whisper fetch -p +vDCk7T1wx%b5RPu 5PeRU1spKSs2VjuUXv6gjo074QFR1Mcco9VGpweULEE
{
  "secret": "The chimera has three sets of teeth",
  "is_base64": false,
  "created": "2021-07-15T19:38:34.877795987Z",
  "accesses": 1
}
```

Note that if you create a password with a secret, you'll also need to provide the password to destroy it.

### Creating and Fetching a Secret File

One of the best reasons to use the CLI application is to share configurations, certificates, and secrets for development projects. You can share files as secrets as follows:

```
$ whisper create -i secret.txt
{
  "token": "LuxV1DAmmdQH2iLez-76m7tvF-3_qNipfZECJx6ALDI",
  "expires": "2021-07-22T22:37:06.619848694Z"
}
```

The `-i` flag specifies the path to a secret file which will be base64 encoded and stored in the secret server. When you fetch the secret that has a file name, it will automatically be saved with the original file name in the current working directory:

```
$ whisper fetch LuxV1DAmmdQH2iLez-76m7tvF-3_qNipfZECJx6ALDI
secret written to secret.txt
```

You can save it to a different directory or to a different file name with the `-o` flag:

```
$ whisper fetch -o fixtures/ LuxV1DAmmdQH2iLez-76m7tvF-3_qNipfZECJx6ALDI
secret written to fixtures/secret.txt
```

Note that you can also save non-file secrets to disk using the `-o` flag as well!

### Destroying Secrets

If you'd like to destroy a secret before it expires without fetching it, use the following command:

```
$ whisper destroy y-CP64rt-tNuy3zeOb2Au52980ALquBg4J6JtSR8fKw
{
  "destroyed": true
}
```

Note that if a secret is password protected, the password is required to destroy it. If the secret is not found then a `404 Not Found` will be returned.

## API Details

To develop against the Whisper REST API load the [Postman](https://www.postman.com/) collection found here: [fixtures/postman_collection.json](fixtures/postman_collection.json).

## Docker

Docker images are used for deployment to Google Cloud Run and Kubernetes clusters and can also be used for development.

### Development

The `docker-compose.yml` configuration is intended to help run local services when developing either the API backend or the React frontend. Profiles assist in running either one set of services or all of the services. To build the images run:

```
$ docker compose --profile=all build
```

The `all` profile will build all images defined by the `docker-compose.yml` function. It should create two images `rotationalio/whisper-api:local` and `rotationalio/whisper-ui:local`. To run just the API:

```
$ docker compose --profile=backend up
```

Note, you may have to create a `fixtures/whisper-sa.json` with a Google service account file in order for the backend to run. Alternatively to run just the front-end:

```
$ docker compose --profile=frontend up
```

You can also run both with the `all` profile as in the build example.

Docker Compose runs the services in a development/debug mode. This means verbose logging from the API server as well as the use of the Mock in-memory secrets database rather than using Google Secret Manager directly; this makes it easier to develop against locally.

## Build and Deploy

This should be handled by GitHub actions. If you want to manually build the images and push them to Dockerhub and gcr.io you can run the following script:

```
$ ./containers/build.sh
```

By default this will tag the build with the Git revision hash, but you can specify a different tag as the first argument as follows:

```
$ ./containers/build.sh v1.0
```