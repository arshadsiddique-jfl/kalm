module github.com/kapp-staging/kapp/operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/jetstack/cert-manager v0.13.1
	github.com/kapp-staging/kapp/controller v0.0.0-20200709120351-d96177bbf37a
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
	istio.io/client-go v0.0.0-20200324231647-289a91f51a8e
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.4.0
)
