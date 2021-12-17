#testground run single --plan=coop_bc --testcase=bc --runner=cluster:k8s --builder=docker:go --instances=1
testground run single --plan=coopbc --testcase=bc --runner=cluster:k8s --builder=docker:go --instances=10 --collect --build-cfg go_proxy_mode=direct
