#!/bin/sh

set -eu

# Store the flagd-proxy image to a helper variable
FLAGD_PROXY_IMG_ORIGINAL=$FLAGD_PROXY_IMG

# Deploy flagd-proxy and the client pod
kubectl delete ns $ZD_TEST_NAMESPACE_FLAGD_PROXY --ignore-not-found=true
kubectl create ns $ZD_TEST_NAMESPACE_FLAGD_PROXY
envsubst < test/zero-downtime-flagd-proxy/manifests/proxy/deployment.yaml | kubectl apply -f - -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
kubectl apply -f test/zero-downtime-flagd-proxy/manifests/proxy/service.yaml -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
kubectl apply -f test/zero-downtime-flagd-proxy/manifests/proxy/flag-config.yaml -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
kubectl wait --for=condition=available deployment/flagd-proxy -n $ZD_TEST_NAMESPACE_FLAGD_PROXY --timeout=60s
envsubst < test/zero-downtime-flagd-proxy/manifests/pod.yaml | kubectl apply -f - -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
kubectl wait --for=condition=ready pod/zd-test -n $ZD_TEST_NAMESPACE_FLAGD_PROXY --timeout=30s

# Wait until connections from client to flagd-proxy are established
sleep 20

for count in 1 2 3;
do
    # Update the flagd-proxy deployment with the second image
    FLAGD_PROXY_IMG=$FLAGD_PROXY_IMG_ZD
    envsubst < test/zero-downtime-flagd-proxy/manifests/proxy/deployment.yaml | kubectl apply -f - -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
    kubectl wait --for=condition=available deployment/flagd-proxy -n $ZD_TEST_NAMESPACE_FLAGD_PROXY --timeout=60s

    # Wait until connections from client to flagd-proxy are re-established
    # 20s should be enough, as terminationGracePeriod is set to 10s
    sleep 20

    # Update the flagd-proxy deployment back to original image
    FLAGD_PROXY_IMG=$FLAGD_PROXY_IMG_ORIGINAL
    envsubst < test/zero-downtime-flagd-proxy/manifests/proxy/deployment.yaml | kubectl apply -f - -n $ZD_TEST_NAMESPACE_FLAGD_PROXY
    kubectl wait --for=condition=available deployment/flagd-proxy -n $ZD_TEST_NAMESPACE_FLAGD_PROXY --timeout=60s

    # Wait until connections from client to flagd-proxy are re-established
    # 20s should be enough, as terminationGracePeriod is set to 10s
    sleep 20
done

# Pod will fail only when it fails to re-connect (that means we do not have zero downtime)
# If it is still running, the last re-connection was successfull.
kubectl wait --for=condition=ready pod/zd-test -n $ZD_TEST_NAMESPACE_FLAGD_PROXY --timeout=30s

# If re-connection was once not successful and another re-connection was, pod might be in a ready state again.
# Therefore we need to check that the restart count is equal to zero -> this means every re-connection was ok.
restart_count=$(kubectl get pods zd-test -o=jsonpath='{.status.containerStatuses[0].restartCount}' -n $ZD_TEST_NAMESPACE_FLAGD_PROXY)
if [ "$restart_count" -ne 0 ]; then
    echo "Restart count of the zd-test pod is not equal to zero."
    exit 1
fi

# Cleanup only when the test passed
kubectl delete ns $ZD_TEST_NAMESPACE_FLAGD_PROXY --ignore-not-found=true

