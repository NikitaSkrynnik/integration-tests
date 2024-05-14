// Code generated by gotestmd DO NOT EDIT.
package spiffe_federation

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
)

type Suite struct {
	base.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/k8s_monolith/external_nse/spiffe_federation")
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/9197ba33f5493f33830cb1eb192179b83ff674bc/examples/k8s_monolith/external_nse/spiffe_federation/clusterspiffeid-template.yaml`)
	r.Run(`kubectl apply -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/9197ba33f5493f33830cb1eb192179b83ff674bc/examples/spire/base/clusterspiffeid-webhook-template.yaml`)
	r.Run(`bundlek8s=$(kubectl exec spire-server-0 -n spire -- bin/spire-server bundle show -format spiffe)` + "\n" + `bundledock=$(docker exec nse-simple-vl3-docker bin/spire-server bundle show -format spiffe)`)
	r.Run(`echo $bundledock | kubectl exec -i spire-server-0 -n spire -- bin/spire-server bundle set -format spiffe -id "spiffe://docker.nsm/cmd-nse-simple-vl3-docker"` + "\n" + `echo $bundlek8s | docker exec -i nse-simple-vl3-docker bin/spire-server bundle set -format spiffe -id "spiffe://k8s.nsm"`)
}
func (s *Suite) Test() {}
