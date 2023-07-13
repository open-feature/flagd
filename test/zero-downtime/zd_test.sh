#!/bin/sh

set -eu

# Store the flagD image to a helper variable
IMG_ORIGINAL=$IMG

# Create pod requesting the values from flagD
envsubst < test/zero-downtime/test-pod.yaml | kubectl apply -f - -n $ZD_TEST_NAMESPACE

for count in 1 2 3;
do
    # Update the flagD deployment with the second image
    IMG=$IMG_ZD
    envsubst < config/deployments/flagd/deployment.yaml | kubectl apply -f - -n $FLAGD_DEV_NAMESPACE
    kubectl wait --for=condition=available deployment/flagd -n $FLAGD_DEV_NAMESPACE --timeout=30s

    # Wait until the client pod executes curl requests agains flagD
    sleep 20

    # Update the flagDT deployment back to original image
    IMG=$IMG_ORIGINAL
    envsubst < config/deployments/flagd/deployment.yaml | kubectl apply -f - -n $FLAGD_DEV_NAMESPACE
    kubectl wait --for=condition=available deployment/flagd -n $FLAGD_DEV_NAMESPACE --timeout=30s

    # Wait until the client pod executes curl requests agains flagD
    sleep 20
done

# Pod will fail only when it fails to get a proper response from curl (that means we do not have zero downtime)
# If it is still running, the last curl request was successfull.
kubectl wait --for=condition=ready pod/test-zd -n $ZD_TEST_NAMESPACE --timeout=30s

# If curl request once not successful and another curl request was, pod might be in a ready state again.
# Therefore we need to check that the restart count is equal to zero -> this means every request provided valid data.
restart_count=$(kubectl get pods test-zd -o=jsonpath='{.status.containerStatuses[0].restartCount}' -n $ZD_TEST_NAMESPACE)
if [ "$restart_count" -ne 0 ]; then
    echo "Restart count of the test-zd pod is not equal to zero."
    exit 1
fi

# Cleanup only when the test passed
kubectl delete ns $ZD_TEST_NAMESPACE --ignore-not-found=true

