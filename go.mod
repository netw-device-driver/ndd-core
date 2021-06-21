module github.com/netw-device-driver/ndd-core

go 1.16

require (
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/Masterminds/semver v1.5.0
	github.com/aws/aws-sdk-go v1.31.6 // indirect
	github.com/docker/cli v0.0.0-20200915230204-cd8016b6bcc5 // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200926000217-2617742802f6+incompatible // indirect
	github.com/google/go-containerregistry v0.4.1
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20210330174036-3259211c1f24
	github.com/netw-device-driver/ndd-runtime v0.3.27
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.4.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/yaml v1.2.0
)
