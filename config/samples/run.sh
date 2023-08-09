export POOL_SIZE=
export CLUSTER_NAME=
export PUB_SSH_KEY=
export VM_SIZE=

clusterctl generate cluster rackncluster \
  --kubernetes-version v1.27.3 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  > rackncluster.yaml