export operator=ndd-core
export github_org=netw-device-driver
mkdir $operator
cd $operator
go mod init $github_org/$operator

## initialize opera

operator-sdk init --domain ndd.henderiw.be --license apache2 --owner "Wim Henderickx" --repo=github.com/$github_org/$operator

## Create a frigates API with Group: testgroup, Version: v1alpha1 and Kind: Test1

kubebuilder edit --multigroup=true

### group driver
operator-sdk create api --group driver --version v1 --kind NetworkNode
operator-sdk create api --group driver --version v1 --kind DeviceDriver

### group pkg
operator-sdk create api --group pkg --version v1 --kind Provider
operator-sdk create api --group pkg --version v1 --kind Package



export operator=ndd-runtime
export github_org=netw-device-driver
mkdir $operator
cd $operator
go mod init $github_org/$operator


## makefile changes

IMG ?= henderiw/nddcore:latest

## dockerfile changes

COPY cmd/ cmd/

## kustomize

update file config/manager/manager.yaml
    replace -> control-plane: controller-manager -> app: ndd
    namespace: system -> ndd-system
    image: controller:latest -> henderiw/nddcore:latest
    imagePullPolicy: Always
    add args and env
        - command:
          - ndd
            args:
            - start
            - --leader-elect
            - --debug
            env:
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                fieldPath: metadata.namespace

update file config/default/kustomization.yaml
    -> namespace: ndd-system
    -> namePrefix: ndd-


kustomize build config/default > manifest.yaml