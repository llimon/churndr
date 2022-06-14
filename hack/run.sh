if [ "$1" != "build" -a "$1" != "run" ]; then
    echo "Syntax: $0 [ build | run ]"
fi

if [ "$1" == "build" -o "$1" == "run" ]; then
    export GO111MODULE=on
    go build -mod=vendor -o churndrcontroller
fi
if [ "$1" == "run" ]; then
    ./churndrcontroller -kubeconfig=$HOME/.kube/config
fi
