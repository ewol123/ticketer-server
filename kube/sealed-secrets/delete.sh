kubectl delete secret creds -n prod
kubectl delete sealedsecret creds -n prod
kubectl delete secret -n prod -l sealedsecrets.bitnami.com/sealed-secrets-key
kubectl delete -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/controller.yaml
