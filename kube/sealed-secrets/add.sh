kubectl kustomize ./prod > ./prod/secret.yaml
kubeseal --format=yaml --cert=pub-cert.pem < ./prod/secret.yaml > ./prod/sealed-secret.yaml
kubectl apply -f ./prod/sealed-secret.yaml -n prod