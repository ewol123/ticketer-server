kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
kubectl apply -f ./serviceaccount.yaml
kubectl get secret -n kube-system

@echo off
echo find your display token "kubectl describe secret <ServiceAccountName-token-xxxxx> -n kube-system"
echo go to http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/. choose token authentication and paste your token

kubectl proxy