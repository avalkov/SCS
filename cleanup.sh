#!/bin/bash

# Function to delete a Docker container
delete_container() {
  container_name=$1
  echo "Deleting container: $container_name"
  docker rm -f $container_name
}

# Function to delete a Docker image
delete_image() {
  image_name=$1
  echo "Deleting image: $image_name"
  docker rmi $image_name
}

# Delete containers
delete_container "scs_server"
delete_container "scs_client"
delete_container "rabbitmq"

# Delete images
delete_image "scs_server"
delete_image "scs_client"

# Delete volumes
#docker volume rm scs_rabbitmq_data

echo "Cleanup complete."
