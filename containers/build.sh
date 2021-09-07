#!/bin/bash

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


if [ -z "$1" ]; then
    TAG=$(git rev-parse --short HEAD)
else
    TAG=$1
fi


if ! ask "Continue with tag $TAG?" N; then
    exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")

docker build -t rotationalio/whisper-api:$TAG -f $DIR/api/Dockerfile $REPO
docker build -t rotationalio/whisper-ui:$TAG -f $DIR/web/Dockerfile $REPO

docker tag rotationalio/whisper-api:$TAG gcr.io/rotationalio-habanero/whisper-api:$TAG
docker tag rotationalio/whisper-ui:$TAG gcr.io/rotationalio-habanero/whisper-ui:$TAG

docker push rotationalio/whisper-api:$TAG
docker push rotationalio/whisper-ui:$TAG
docker push gcr.io/rotationalio-habanero/whisper-api:$TA
docker push gcr.io/rotationalio-habanero/whisper-ui:$TA