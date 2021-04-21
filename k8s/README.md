# Start Nodes
kubectl create -f deployment.yml 

## Check Nodes
kubectl get pods

# Enter Pod/Container
kubectl exec --stdin --tty bcpcp-84b969864f-bcpqx -- /bin/bash

# Delete Deployment
kubectl delete -f deployment.yml 


# Notes
Multi-cluster networking can be a minefield.

Pod-to-Pod communication across clusters should work out of the box on GKE and EKS
- on AKS the "advanced networking mode must be set
- on Cilium a "cluster mesh" is needed
- others are not tested

## If pods are not able to communicate:

### Fix 1
If Pods cannot find each other directly
- use kubeadm join 
- or join flags
(make sure the pods have static ips)

### Fix 2
- create a public load balancer for each pod to expose the pods
(this causes a big overhead and extra costs)



