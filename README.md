# k8s 集群同步工具

## 使用方法

```shell
# 同步 deployment
go run main.go -s $HOME/.kube/source  -t $HOME/.kube/target -n default deployment


# 同步 ingress
go run main.go -s $HOME/.kube/source  -t $HOME/.kube/target -n default ingress


# 同步 service
go run main.go -s $HOME/.kube/source  -t $HOME/.kube/target -n default service

```