### update volcano
```
# volcano-scheduler-gpu.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: volcano-scheduler-gpu-configmap
  namespace: kube-stack
data:
  volcano-scheduler.conf: |
    actions: "enqueue, allocate, backfill"
    tiers:
    - plugins:
      - name: priority
      - name: gang
        enablePreemptable: false
      - name: conformance
    - plugins:
      - name: overcommit
      - name: drf
        enablePreemptable: false
      - name: predicates
      - name: proportion
      - name: nodeorder
      - name: binpack
      - name: gpupredicates
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: volcano-scheduler-gpu
  namespace: kube-stack
  labels:
    app: volcano-scheduler-gpu
spec:
  replicas: 1
  selector:
    matchLabels:
      app: volcano-scheduler-gpu
  template:
    metadata:
      labels:
        app: volcano-scheduler-gpu
    spec:
      serviceAccount: volcano-scheduler
      priorityClassName: system-cluster-critical
      containers:
        - name: volcano-scheduler-gpu
          image: g-ubjg5602-docker.pkg.coding.net/iscas-system/containers/volcano-scheduler-gpu:v1.0.0
          args:
            - --logtostderr
            - --scheduler-conf=/volcano.scheduler/volcano-scheduler.conf
            - --enable-healthz=true
            - --enable-metrics=true
            - --leader-elect=false
            - --scheduler-name=volcano
            - -v=5
            - 2>&1
          env:
            - name: DEBUG_SOCKET_DIR
              value: /tmp/klog-socks
          imagePullPolicy: Always
          volumeMounts:
            - name: scheduler-config-gpu
              mountPath: /volcano.scheduler
            - name: klog-sock
              mountPath: /tmp/klog-socks
      volumes:
        - name: scheduler-config-gpu
          configMap:
            name: volcano-scheduler-gpu-configmap
        - name: klog-sock
          hostPath:
            path: /tmp/klog-socks

```