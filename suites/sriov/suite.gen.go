// Code generated by gotestmd DO NOT EDIT.
package sriov

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
	r := s.Runner("../deployments-k8s/examples/sriov")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration nsm-mutating-webhook` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/sriov?ref=28570ae6d97926e17351dacf34b458170a06b22b`)
}
func (s *Suite) TestSriovKernel2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/SriovKernel2Noop")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ns-sriov-kernel2noop`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/SriovKernel2Noop/ponger?ref=28570ae6d97926e17351dacf34b458170a06b22b`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=ponger`)
	r.Run(`kubectl -n ns-sriov-kernel2noop exec deploy/ponger -- ip a | grep "172.16.1.100"`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/SriovKernel2Noop?ref=28570ae6d97926e17351dacf34b458170a06b22b`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel`)
	r.Run(`kubectl -n ns-sriov-kernel2noop wait --for=condition=ready --timeout=1m pod -l app=nse-noop`)
	r.Run(`kubectl -n ns-sriov-kernel2noop exec deployments/nsc-kernel -- ping -c 4 172.16.1.100`)
}
func (s *Suite) TestVfio2Noop() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Vfio2Noop")
	s.T().Cleanup(func() {
		r.Run(`kubectl -n ns-vfio2noop exec deployments/nse-vfio --container ponger -- /bin/bash -c '\` + "\n" + `  (sleep 10 && kill $(pgrep "pingpong")) 1>/dev/null 2>&1 &             \` + "\n" + `'`)
		r.Run(`kubectl delete ns ns-vfio2noop`)
	})
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/use-cases/Vfio2Noop?ref=28570ae6d97926e17351dacf34b458170a06b22b`)
	r.Run(`kubectl -n ns-vfio2noop wait --for=condition=ready --timeout=1m pod -l app=nsc-vfio`)
	r.Run(`kubectl -n ns-vfio2noop wait --for=condition=ready --timeout=1m pod -l app=nse-vfio`)
	r.Run(`function dpdk_ping() {` + "\n" + `  err_file="$(mktemp)"` + "\n" + `  trap 'rm -f "${err_file}"' RETURN` + "\n" + `` + "\n" + `  client_mac="$1"` + "\n" + `  server_mac="$2"` + "\n" + `` + "\n" + `  command="/root/dpdk-pingpong/build/pingpong \` + "\n" + `      --no-huge                               \` + "\n" + `      --                                      \` + "\n" + `      -n 500                                  \` + "\n" + `      -c                                      \` + "\n" + `      -C ${client_mac}                        \` + "\n" + `      -S ${server_mac}` + "\n" + `      "` + "\n" + `  out="$(kubectl -n ns-vfio2noop exec deployments/nsc-vfio --container pinger -- /bin/bash -c "${command}" 2>"${err_file}")"` + "\n" + `` + "\n" + `  if [[ "$?" != 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if ! pong_packets="$(echo "${out}" | grep "rx .* pong packets" | sed -E 's/rx ([0-9]*) pong packets/\1/g')"; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  if [[ "${pong_packets}" == 0 ]]; then` + "\n" + `    echo "${out}"` + "\n" + `    cat "${err_file}" 1>&2` + "\n" + `    return 1` + "\n" + `  fi` + "\n" + `` + "\n" + `  echo "${out}"` + "\n" + `  return 0` + "\n" + `}`)
	r.Run(`dpdk_ping "0a:55:44:33:22:00" "0a:55:44:33:22:11"`)
}
