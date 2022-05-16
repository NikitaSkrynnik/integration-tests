// Code generated by gotestmd DO NOT EDIT.
package floatingvl3

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/interdomain"
)

type Suite struct {
	base.Suite
	interdomainSuite interdomain.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.interdomainSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/interdomain/usecases/FloatingVl3")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG3 kubectl delete -k ./cluster3`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2 kubectl delete -k ./cluster2`)
		r.Run(`export KUBECONFIG=$KUBECONFIG1 kubectl delete -k ./cluster1`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`kubectl apply -k ./cluster3`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl apply -k ./cluster1`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl apply -k ./cluster2`)
	r.Run(`nsc2=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `ipAddr2=$(kubectl exec -n ns-vl3-interdomain  $nsc2 -- ifconfig nsm-1)` + "\n" + `ipAddr2=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`nsc1=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `ipAddr1=$(kubectl exec -n ns-vl3-interdomain $nsc1 -- ifconfig nsm-1)` + "\n" + `ipAddr1=$(echo $ipAddr | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)`)
	r.Run(`kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 $ipAddr2`)
	r.Run(`kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 $ipAddr1`)
	r.Run(`kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
}
func (s *Suite) Test() {}
