// Code generated by gotestmd DO NOT EDIT.
package basic

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire/single_cluster"
)

type Suite struct {
	base.Suite
	single_clusterSuite single_cluster.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.single_clusterSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/basic")
	s.T().Cleanup(func() {
		r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl delete mutatingwebhookconfiguration ${WH}` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/basic?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`WH=$(kubectl get pods -l app=admission-webhook-k8s -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `kubectl wait --for=condition=ready --timeout=1m pod ${WH} -n nsm-system`)
}
func (s *Suite) TestKernel2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2kernel`)
	})
	r.Run(`kubectl create ns ns-kernel2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2memif`)
	})
	r.Run(`kubectl create ns ns-kernel2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-kernel2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-kernel2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2memif -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-kernel2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2vxlan2kernel`)
	})
	r.Run(`kubectl create ns ns-kernel2vxlan2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Vxlan2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2vxlan2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2vxlan2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2vxlan2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Vxlan2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Vxlan2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2vxlan2memif`)
	})
	r.Run(`kubectl create ns ns-kernel2vxlan2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Vxlan2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2vxlan2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-kernel2vxlan2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2vxlan2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-kernel2vxlan2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2vxlan2memif -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-kernel2vxlan2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2wireguard2kernel`)
	})
	r.Run(`kubectl create ns ns-kernel2wireguard2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Wireguard2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2wireguard2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-kernel2wireguard2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2wireguard2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-kernel2wireguard2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2wireguard2kernel -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ns-kernel2wireguard2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestKernel2Wireguard2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Wireguard2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-kernel2wireguard2memif`)
	})
	r.Run(`kubectl create ns ns-kernel2wireguard2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Kernel2Wireguard2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=alpine -n ns-kernel2wireguard2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-kernel2wireguard2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=alpine -n ns-kernel2wireguard2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-kernel2wireguard2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ns-kernel2wireguard2memif -- ping -c 4 172.16.1.100`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-kernel2wireguard2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2kernel`)
	})
	r.Run(`kubectl create ns ns-memif2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-memif2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-memif2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2kernel" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec ${NSE} -n ns-memif2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2memif`)
	})
	r.Run(`kubectl create ns ns-memif2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-memif2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-memif2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2vxlan2kernel`)
	})
	r.Run(`kubectl create ns ns-memif2vxlan2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Vxlan2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2vxlan2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-memif2vxlan2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-memif2vxlan2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2vxlan2kernel" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec ${NSE} -n ns-memif2vxlan2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Vxlan2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Vxlan2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2vxlan2memif`)
	})
	r.Run(`kubectl create ns ns-memif2vxlan2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Vxlan2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2vxlan2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2vxlan2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2vxlan2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-memif2vxlan2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2vxlan2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-memif2vxlan2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
func (s *Suite) TestMemif2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2wireguard2kernel`)
	})
	r.Run(`kubectl create ns ns-memif2wireguard2kernel`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Wireguard2Kernel?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2wireguard2kernel`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ns-memif2wireguard2kernel`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2wireguard2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ns-memif2wireguard2kernel --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2wireguard2kernel" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`kubectl exec ${NSE} -n ns-memif2wireguard2kernel -- ping -c 4 172.16.1.101`)
}
func (s *Suite) TestMemif2Wireguard2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Wireguard2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-memif2wireguard2memif`)
	})
	r.Run(`kubectl create ns ns-memif2wireguard2memif`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Memif2Wireguard2Memif?ref=7be3b3ec60069b2a91e87215780814ae088c7456`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ns-memif2wireguard2memif`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ns-memif2wireguard2memif`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ns-memif2wireguard2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ns-memif2wireguard2memif --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "ns-memif2wireguard2memif" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "ns-memif2wireguard2memif" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}
