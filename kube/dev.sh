helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install app-ingress ingress-nginx/ingress-nginx --namespace=prod
kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s -n prod
