# Start Nodes
kubectl create -f deployment.yml 

## Check Nodes
kubectl get pods

# Enter Pod/Container
kubectl exec --stdin --tty bcpcp-84b969864f-bcpqx -- /bin/bash

# Delete Deployment
kubectl delete -f deployment.yml 


# Notes
Pod-to-Pod communication across clusters works out of the box on GKE and EKS
- on AKS the "advanced networking mode must be set
- on Cilium a "cluster mesh" is needed
- others are not tested