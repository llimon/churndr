kubectl apply -f artifacts/examples/crd.yaml
for a in `ls artifacts/examples/example-*.yaml`; do
    kubectl apply -f $a
done
