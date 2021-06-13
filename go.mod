module github.com/netw-device-driver/ndd-core

go 1.16

require (
	github.com/Masterminds/semver v1.5.0
	github.com/crossplane/crossplane v1.2.2
	github.com/crossplane/crossplane-runtime v0.13.0
	github.com/google/go-containerregistry v0.4.1
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20210330174036-3259211c1f24
	github.com/netw-device-driver/ndd-runtime v0.3.4
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.4.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/yaml v1.2.0
)
