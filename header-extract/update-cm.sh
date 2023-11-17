#!/bin/bash

# get original configmap
kubectl get cm -n deepflow deepflow-agent -o jsonpath='{.data}'

# patch to configmap
kubectl patch  -n deepflow configmap/deepflow-agent \
--type merge \
-p '{"data":{"deepflow-agent.yaml":"controller-ips:\n- deepflow-server\n","header.yaml":"port-white-list:\n- \"8080\""}}'
