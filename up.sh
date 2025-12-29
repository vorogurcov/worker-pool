# Build image with new version

docker build -t worker-pool:0.0.2 ./backend

DB_FOLDER="./infra/db"
WORKER_POOL_FOLDER="./infra/worker-pool"
PROMETHEUS_FOLDER="./infra/monitoring/prometheus"


kubectl apply -f ${DB_FOLDER}/k8s-pvc.yaml

kubectl apply -f ${DB_FOLDER}/k8s-secret.yaml

kubectl apply -f ${DB_FOLDER}/k8s-deployment.yaml

kubectl apply -f ${DB_FOLDER}/k8s-service.yaml

kubectl apply -f ${WORKER_POOL_FOLDER}/k8s-deployment.yaml

kubectl apply -f ${WORKER_POOL_FOLDER}/k8s-service.yaml

kubectl apply -f ${PROMETHEUS_FOLDER}/rbac.yaml

kubectl apply -f ${PROMETHEUS_FOLDER}/k8s-config-map.yaml

kubectl apply -f ${PROMETHEUS_FOLDER}/k8s-deployment.yaml

kubectl apply -f ${PROMETHEUS_FOLDER}/k8s-service.yaml

