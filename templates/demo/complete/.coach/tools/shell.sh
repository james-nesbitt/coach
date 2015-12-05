#!/bin/bash

# What shell image to use
IMAGE="${IMAGE:-jamesnesbitt/wunder-developershell}"

# Project name, used for hostname, and to build container names
PROJECT_NAME="${PROJECT_NAME:-project}"
# Used to build container names
PROJECT_INSTANCE="${PROJECT_INSTANCE:-single}"

# Use some standard docker flags
DOCKER_FLAGS="${DOCKER_FLAGS:---rm --tty --interactive=true --hostname=${PROJECT_NAME}}"

# EACH ENTRY IS A CONTAINER NODE TO LINK TO
LINKS="${LINKS:-}"

# EACH ENTRY IS A CONTAINER TO MOUNT ALL VOLUMES FROM
VOLUMESFROM="${VOLUMESFROM:-}"
# These indicate ENVs that hold volume definitions
USE_VOLUMES="${USE_VOLUMES:-}"

# These ENVs will be added to the container
USE_ENVS="${USE_ENVS:-}"

# Process all of the above configs to make a docker run command

for LINK in ${LINKS}; do
  DOCKER_FLAGS="${DOCKER_FLAGS} --link=${PROJECT_NAME}_${LINK}_${PROJECT_INSTANCE}:${LINK}.app"
done
for VOLUME in ${USE_VOLUMES}; do
  DOCKER_FLAGS="${DOCKER_FLAGS} --volume=${!VOLUME}"
done
for VOLUMEFROM in ${VOLUMESFROM}; do
  DOCKER_FLAGS="${DOCKER_FLAGS} --volumes-from=${PROJECT_NAME}_${VOLUMEFROM}_${PROJECT_INSTANCE}"
done
for ENV in ${USE_ENVS}; do
  DOCKER_FLAGS="${DOCKER_FLAGS} --env=${ENV}=${!ENV}"
done

# Run Docker
docker run ${DOCKER_FLAGS} ${IMAGE} $@
