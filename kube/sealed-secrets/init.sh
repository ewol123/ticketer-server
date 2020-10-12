kubectl create namespace prod
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/controller.yaml -n prod
kubectl wait --for condition=ready pod --all --timeout=120s -n kube-system
kubeseal --fetch-cert \--controller-namespace=kube-system \--controller-name=sealed-secrets-controller \ > pub-cert.pem