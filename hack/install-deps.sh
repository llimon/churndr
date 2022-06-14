go install k8s.io/code-generator/cmd/client-gen
go install k8s.io/code-generator/cmd/lister-gen
go install k8s.io/code-generator/cmd/informer-gen
go install k8s.io/code-generator/cmd/deepcopy-gen

go install k8s.io/code-generator/cmd/defaulter-gen
go get k8s.io/code-generator/cmd/defaulter-gen@v0.22.2
go install k8s.io/code-generator/cmd/defaulter-gen
cp ~/go/bin/defaulter-gen hack/cmd/
