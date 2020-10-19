# bc_p2p
A peer to peer implementation of confidential bibliographic coupling detection.

---
#### Run multiple instances from your terminal, then:
PUT: Type `PUT my-key my-value`

GET: Type `GET my-key`

BATCH: Type `BATCH my-features-file`

CHECK: TYPE `CHECK my-features-file`

Close: Type `Ctrl-c`

#### Run automated sequential tests

/Data
xstarttmux.sh: starts the desired amount of instances on a host
batchGET.sh: queries 1000 hashes
batchPUT.sh: posts 1000 hashes
putget.sh: posts and queries a single hash
---

*if you use Intellij - a library requires the following setting (https://github.com/intellij-rust/intellij-rust/issues/3628)
