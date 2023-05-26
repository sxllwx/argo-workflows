package entrypoint

import (
	"context"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/lru"

	"github.com/argoproj/argo-workflows/v3/config"
)

type Interface interface {
	Lookup(ctx context.Context, image string, options Options) (*Image, error)
}

type Options struct {
	Namespace          string
	ServiceAccountName string
	ImagePullSecrets   []apiv1.LocalObjectReference
	Transport          http.RoundTripper
}

type Image struct {
	Entrypoint []string
	Cmd        []string
}

func New(kubernetesClient kubernetes.Interface, config map[string]config.Image, insecureRegistries []string) Interface {
	return &cacheIndex{
		lru.New(1024),
		chainIndex{
			configIndex(config),
			&containerRegistryIndex{
				kubernetesClient:   kubernetesClient,
				insecureRegistries: sets.NewString(insecureRegistries...),
			},
		},
	}
}
