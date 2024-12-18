#!/bin/sh


if [ -z "$1" ]; then
    echo "image not provided"
    exit 1
fi
IMAGE="$1"

if [ -z "$2" ]; then
    echo "source not provided"
    exit 1
fi
SRC="$2"

if [ -z "$3" ]; then
    echo "cp destination not provided"
    exit 1
fi
DEST="$3"

echo "copying from: $IMAGE:$SRC to $DEST"

CONTAINER_ID=$(docker create "${IMAGE}")
echo "container created with id: ${CONTAINER_ID}"

docker cp "${CONTAINER_ID}":"${SRC}" "${DEST}"
