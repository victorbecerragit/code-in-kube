#! /usr/bin/env bash
set -eou pipefail

USER=victorbecerra  # non-root user to move kubeconfig file into their home dir
USER_HOME="/home/$USER" # path to non-root home folder
CURRENT_DIR="$(dirname $(readlink -f $0))"

CLUSTER_NAME=kata-local

APP_MOUNT_FOLDERS=(
  tmp
)

NODE_MOUNT_FOLDERS=(
  "/tmp/control-plane-1"
  "/tmp/worker-1"
  "/tmp/worker-2"
  "/tmp/worker-3"
)

create_node_folders() {
  for node_folder in "${NODE_MOUNT_FOLDERS[@]}"; do
    if [ -e "$node_folder" ]; then
      echo "Removing folder $node_folder..."
      rm -rfI "$node_folder"
    fi
    for app_folder in "${APP_MOUNT_FOLDERS[@]}"; do
      echo "Creating folder ${node_folder}/${app_folder} ..."
      mkdir -p "${node_folder}/${app_folder}"
    done
  done
}

if [ -z "$(kind get clusters -q)" ]
then
  echo "Creating cluster..."
  create_node_folders
  kind create cluster --name "$CLUSTER_NAME" --config="${CURRENT_DIR}/cluster-config.yaml"
  kind get kubeconfig --name "$CLUSTER_NAME" -q > kubeconfig
  chown "$USER" kubeconfig
  mv kubeconfig "$USER_HOME/.kube/config"
else
  echo "Cluster already there. Skipping creation."
fi
