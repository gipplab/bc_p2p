# testground run single --plan=coop_bc --testcase=upload --runner=local:docker --builder=docker:go --instances=1
testground run single --plan=coop_bc --testcase=bc --runner=local:docker --builder=docker:go --instances=100 --collect
