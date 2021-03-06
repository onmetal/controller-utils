//go:generate go run github.com/golang/mock/mockgen -copyright_file ../../../hack/boilerplate.go.txt -package metautils -destination=funcs.go github.com/onmetal/controller-utils/mock/controller-utils/metautils EachListItemFunc
package metautils

import "sigs.k8s.io/controller-runtime/pkg/client"

type EachListItemFunc interface {
	Call(obj client.Object) error
}
