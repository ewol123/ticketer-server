helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ticketer-ingress ingress-nginx/ingress-nginx
kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s
kubectl create -f ./app