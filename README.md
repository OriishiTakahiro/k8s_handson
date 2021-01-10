# Kubernetesハンズオン

Kubernetesの概念を勉強する上で使うときの感覚を掴んでもらうためのハンズオンとなります。

## 対象範囲

本ハンズオンでは以下の概念を触ってもらいます。

- Pod
- Job
- CronJob
- ReplicaSet
- Service
- Ingress

## 想定環境

本ハンズオンではkindを使うので、以下の手順を行なっている必要があります。

- [kubectlのインストール](https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/)
- [Dockerのインストール](https://docs.docker.com/get-docker/)
- [kindのインストール](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

## 実施手順

### Cluster

最初にclusterを作成します。

clusterとはKubernetesが動作するサーバ群のことで、kindではDockerコンテナを使ってclusterをエミュレートします。

```sh
 >> kind create cluster --config cluster.yaml                                                                                                                                                                                                                                                        ?[main]
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.19.1) 🖼
 ✓ Preparing nodes 📦 📦 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

### Pod

最初にPodを立ててみましょう。

nodeが物理サーバのようなものであれば、Podは1台のVMのようなものです。

Podでは1台以上のコンテナが稼働しており、同一Pod上のコンテナはIPアドレスやポート、ローカルストレージを共有します。

```sh
# Podの作成
>> kubectl apply -f pod.yaml
pod/bastion configured

# Podの一覧
>> k get pods
NAME      READY   STATUS    RESTARTS   AGE
bastion   1/1     Running   0          7m42s

# Podの詳細を取得
>> k describe pods bastion
Name:         bastion
Namespace:    default
Priority:     0
Node:         kind-worker/172.20.0.2
Start Time:   Mon, 11 Jan 2021 02:20:30 +0900
Labels:       run=bastion
Annotations:  <none>
Status:       Running
IP:           10.244.3.4
IPs:
  IP:  10.244.3.4
Containers:
  bastion:
    Container ID:  containerd://84d28de7e9d31b0c733d0b9a246ff91495ce1fb4982762900df1bfc3ef1e4092
    Image:         debian:8.11
    Image ID:      docker.io/library/debian@sha256:b6355f8f0cae1ad0f8f70f1fa276e3b86ed9d3081c44cfe436ef7133bf10ff19
    Port:          <none>
    Host Port:     <none>
    Command:
      bash
    Args:
      -c
      while true; do sleep 1; done
    State:          Running
      Started:      Mon, 11 Jan 2021 02:20:31 +0900
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-s7m66 (ro)
Conditions:
  Type              Status
  Initialized       True
  Ready             True
  ContainersReady   True
  PodScheduled      True
Volumes:
  default-token-s7m66:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-s7m66
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                 node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age    From               Message
  ----    ------     ----   ----               -------
  Normal  Scheduled  7m52s  default-scheduler  Successfully assigned default/bastion to kind-worker
  Normal  Pulled     7m51s  kubelet            Container image "debian:8.11" already present on machine
  Normal  Created    7m51s  kubelet            Created container bastion
  Normal  Started    7m51s  kubelet            Started container bastion

# またPodの中でシェルを立ち上げて操作することも可能
>> kubectl exec -ti bastion -- bash
root@bastion:/# exit
exit
```

### ReplicaSet

次にReplicaSetを作ってみます。

ReplicaSetはPodの集合で、望んだレプリカ数だけPodが起動しているように維持してくれます。

レプリカ数の指定はPodに付与するlabelで判定します。

```sh
# ReplicaSetの作成
>> kubectl apply -f rs.yaml
replicaset.apps/rs-sample created

# ReplicaSetの一覧
>> kubectl get replicaset
NAME        DESIRED   CURRENT   READY   AGE
rs-sample   3         3         3       36s

# replicasの数だけPodが起動している
>> kubectl get pods -l app=rs-match-label
NAME              READY   STATUS    RESTARTS   AGE
rs-sample-bnzpq   1/1     Running   0          43s
rs-sample-g5g6q   1/1     Running   0          43s
rs-sample-n58pq   1/1     Running   0          43s

# 試しにPodを削除
>> kubectl delete pods rs-sample-bnzpq
pod "rs-sample-bnzpq" deleted

# 新しいPodが立ち上がってPod数を維持
>> kubectl get pods -l app=rs-match-label
NAME              READY   STATUS    RESTARTS   AGE
rs-sample-g5g6q   1/1     Running   0          2m15s
rs-sample-n58pq   1/1     Running   0          2m15s
rs-sample-xt62b   1/1     Running   0          39s

# ReplicaSetの削除
>> kubectl delete replicasets rs-sample
replicaset.apps "rs-sample" deleted
```

### Deployment

次にDeploymentを作ってみます。

DeploymentはReplicaSetの世代管理をするものです。

ReplicaSetにアップデートを仕掛ける際に一台ずつ立ち上げては落とすローリングアップデートを行なってくれます。

(ローリングアップデートを行うことによりサービスは瞬断を免れることができます)

```sh
# Deploymentの作成 (Deploymentは内部的にReplicaSetを作る)
>> kubectl apply -f deploy.yaml
deployment.apps/deploy-sample created

# Deploymentの一覧
>> kubectl get deployments
NAME            READY   UP-TO-DATE   AVAILABLE   AGE
deploy-sample   3/3     3            3           30s

# 別のターミナルを開き、リアルタイムでReplicaSetのレプリカ数の変化を観察
# (表示は元のターミナルで次のコマンドを実行するまで表示されない)
>> kubectl get replicasets -w
NAME                       DESIRED   CURRENT   READY   AGE
deploy-sample-785c79f7f8   3         3         3       16s
deploy-sample-7b7fc7c5d    1         0         0       0s
deploy-sample-7b7fc7c5d    1         0         0       0s
deploy-sample-7b7fc7c5d    1         1         0       0s
deploy-sample-7b7fc7c5d    1         1         1       3s
deploy-sample-785c79f7f8   2         3         3       33s
deploy-sample-7b7fc7c5d    2         1         1       3s
deploy-sample-785c79f7f8   2         3         3       33s
deploy-sample-7b7fc7c5d    2         1         1       3s
deploy-sample-785c79f7f8   2         2         2       33s
deploy-sample-7b7fc7c5d    2         2         1       3s
deploy-sample-7b7fc7c5d    2         2         2       7s
deploy-sample-785c79f7f8   1         2         2       37s
deploy-sample-785c79f7f8   1         2         2       37s
deploy-sample-7b7fc7c5d    3         2         2       7s
deploy-sample-7b7fc7c5d    3         2         2       7s
deploy-sample-785c79f7f8   1         1         1       37s
deploy-sample-7b7fc7c5d    3         3         2       7s
deploy-sample-7b7fc7c5d    3         3         3       10s
deploy-sample-785c79f7f8   0         1         1       40s
deploy-sample-785c79f7f8   0         1         1       40s
deploy-sample-785c79f7f8   0         0         0       40s

# アップデートを仕掛ける
>> kubectl apply -f deploy-2.yaml
deployment.apps/deploy-sample configured
```

### Service

次にServiceを作ってみます。

ServiceはFQDN経由でPodへのアクセスを可能にし、複数のPodへのロードバランシングやヘルスチェックも行います。

Service経由でのアクセス先はPodに付与されるlabelで判定します。

```sh
# Serviceの作成
>> kubectl apply -f svc.yaml
service/svc-sample created

# Serviceの一覧
>> kubectl get services
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP          49m
svc-sample   NodePort    10.96.222.105   <none>        8900:30262/TCP   7s

# bastionから何度かcurlを投げると異なるIPが返ってくる (異なるPodにアクセス)
>> kubectl exec -ti bastion -- bash
root@bastion:/# apt update
root@bastion:/# apt install -y curl
root@bastion:/# curl http://svc-sample.default.svc.cluster.local:8900/ip
10.244.1.7
root@bastion:/# curl http://svc-sample.default.svc.cluster.local:8900/ip
10.244.2.7
root@bastion:/# curl http://svc-sample.default.svc.cluster.local:8900/ip
10.244.3.10
```

### Ingress

Serviceだとcluster内でのPod間通信しか行えませんでした。

次にcluster外部からのアクセスを可能にするIngressを作ってみましょう。

ここではNginxが起動するPodをIngressとして使っていますが、GCPではCLB、EKSではELBなどのロードバランサがIngressの実体として使われます。

```sh
# nginx-ingress-controllerを適用
>> kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
namespace/ingress-nginx created
serviceaccount/ingress-nginx created
configmap/ingress-nginx-controller created
clusterrole.rbac.authorization.k8s.io/ingress-nginx created
clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx created
role.rbac.authorization.k8s.io/ingress-nginx created
rolebinding.rbac.authorization.k8s.io/ingress-nginx created
service/ingress-nginx-controller-admission created
service/ingress-nginx-controller created
deployment.apps/ingress-nginx-controller created
validatingwebhookconfiguration.admissionregistration.k8s.io/ingress-nginx-admission created
serviceaccount/ingress-nginx-admission created
clusterrole.rbac.authorization.k8s.io/ingress-nginx-admission created
clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
role.rbac.authorization.k8s.io/ingress-nginx-admission created
rolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
job.batch/ingress-nginx-admission-create created
job.batch/ingress-nginx-admission-patch created

# nginx-ingress-controllerが適用されるまで待機
>> kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
pod/ingress-nginx-controller-55bc59c885-7jcmn condition met

# Ingressの作成
>> kubectl apply -f ingress.yaml
ingress.networking.k8s.io/ingress-sample created

# Ingressの一覧
>> kubectl get ingress
Warning: extensions/v1beta1 Ingress is deprecated in v1.14+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress
NAME             CLASS    HOSTS   ADDRESS   PORTS   AGE
ingress-sample   <none>   *                 80      7s

# ホストからリクエストが到達するようになっている
>> curl http://localhost/ip
10.244.1.7
```

### CronJob

Kubernetesにはバッチ処理を定期実行するためのCronJobがあります。

CronJobは実行するごとに新しいPodを作成し、所定の処理を実行します。

```sh
# Deploymentのレプリカ数を1にする
>> kubectl scale deployments deploy-sample --replicas=1
deployment.apps/deploy-sample scaled

# Deploymentのレプリカ数を1にする
>> kubectl scale deployments deploy-sample --replicas=1
deployment.apps/deploy-sample scaled

# CronJob作成
>> kubectl apply -f cj.yaml
cronjob.batch/cj-sample created

# CronJob一覧
>> kubectl get cronjobs
NAME        SCHEDULE    SUSPEND   ACTIVE   LAST SCHEDULE   AGE
cj-sample   * * * * *   False     0        50s             5m6s

# CronJob実行用のpodが作られている
>> k get pods
NAME                            READY   STATUS      RESTARTS   AGE
cj-sample-1610308920-hf78f      0/1     Completed   0          23s
deploy-sample-7b7fc7c5d-2fkn2   1/1     Running     0          19m

# CronJobが実行された時刻が返ってくる
>> curl  http://localhost/date
Sun Jan 10 20:04:02 UTC 2021
```

### おかたづけ

```sh
>> kind delete cluster
```
