# Start Nodes
kubectl create -f deployment.yml 

## Check Nodes
kubectl get pods

# Enter Pod/Container
kubectl exec --stdin --tty bcpcp-84b969864f-bcpqx -- /bin/bash

# Delete Deployment
kubectl delete -f deployment.yml 
