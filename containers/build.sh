#!/bin/bash
# A helper script for building docker images

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-t TAG] [-p PLATFORM]
A helper script for building whisper docker images
Flags are as follows (getopt required):

    -h  display this help and exit
    -t  the tag to tag the images with
    -p  the target platform to build the docker images for

Unless otherwise specified TAG is the git hash and PLATFORM is
linux/amd64 when deploying to ensure the correct images are deployed.

NOTE: realpath is required; you can install it on OS X with

    $ brew install coreutils
EOF
}

ask() {
    local prompt default reply

    if [[ ${2:-} = 'Y' ]]; then
        prompt='Y/n'
        default='Y'
    elif [[ ${2:-} = 'N' ]]; then
        prompt='y/N'
        default='N'
    else
        prompt='y/n'
        default=''
    fi

    while true; do

        # Ask the question (not using "read -p" as it uses stderr not stdout)
        echo -n "$1 [$prompt] "

        # Read the answer (use /dev/tty in case stdin is redirected from somewhere else)
        read -r reply </dev/tty

        # Default?
        if [[ -z $reply ]]; then
            reply=$default
        fi

        # Check if the reply is valid
        case "$reply" in
            Y*|y*) return 0 ;;
            N*|n*) return 1 ;;
        esac

    done
}

# Helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")
DOTENV="$REPO/.env"

# Set environment variables for the build process
export GIT_REVISION=$(git rev-parse --short HEAD)

# Load .env file from project root if it exists
if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi

# Parse command line options with getopt
OPTIND=1
TAG=${GIT_REVISION}
PLATFORM="linux/amd64"

while getopts htp: opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        t)  TAG=$OPTARG
            ;;
        p)  PLATFORM=$OPTARG
            ;;
        *)
            show_help >&2
            exit 2
            ;;
    esac
done
shift "$((OPTIND-1))"


if ! ask "Continue with tag $TAG?" N; then
    exit 1
fi

# Build the images
docker buildx build --platform $PLATFORM -t rotationalio/whisper-api:$TAG -f $DIR/api/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
docker buildx build --platform $PLATFORM -t rotationalio/whisper-ui:$TAG -f $DIR/web/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO

docker tag rotationalio/whisper-api:$TAG gcr.io/rotationalio-habanero/whisper-api:$TAG
docker tag rotationalio/whisper-ui:$TAG gcr.io/rotationalio-habanero/whisper-ui:$TAG

docker push rotationalio/whisper-api:$TAG
docker push rotationalio/whisper-ui:$TAG
docker push gcr.io/rotationalio-habanero/whisper-api:$TAG
docker push gcr.io/rotationalio-habanero/whisper-ui:$TAG
