echo "run with . or source"
export KOPS_STATE_STORE=gs://kubernetes-cluster22222/
export KOPS_FEATURE_FLAG=AlphaAllowGCE
export PROJECT=p2p-evaluation
echo $KOPS_STATE_STORE
echo $KOPS_FEATURE_FLAG
echo $PROJECT 
