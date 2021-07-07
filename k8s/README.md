# Start Nodes
kubectl create -f deployment.yml 

## Check Nodes
kubectl get pods

# Enter Pod/Container
kubectl exec --stdin --tty bcpcp-84b969864f-bcpqx -- /bin/bash

# Delete Deployment
kubectl delete -f deployment.yml 


# Notes
Multi-regional networking can be a minefield.

for Pod-to-Pod communication across clusters:
- on AKS the "advanced networking mode must be set
- on Cilium a "cluster mesh" is needed
- on GKE enable "intranode visibility"
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


### GKE Problem
- Please be aware that free trial accounts for Google Cloud Platform have limited quota during their trial period. In order to increase your quota, please upgrade to a paid account by clicking "Upgrade my account" from the top of any page once logged in to Google Cloud Console.
No more than 8 IPs per region.

## Multi-Regional Cluster
- Set up Kubernetes clusters in multiple regions with different CIDR blocks
- Set up cross-region communication

