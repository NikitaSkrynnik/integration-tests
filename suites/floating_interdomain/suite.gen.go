// Code generated by gotestmd DO NOT EDIT.
package floating_interdomain

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/floating_interdomain/dns"
	"github.com/networkservicemesh/integration-tests/suites/floating_interdomain/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/floating_interdomain/spire"
)

type Suite struct {
	base.Suite
	loadbalancerSuite loadbalancer.Suite
	dnsSuite          dns.Suite
	spireSuite        spire.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dnsSuite, &s.spireSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/floating_interdomain")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG1 && kubectl delete ns nsm-system`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2 && kubectl delete ns nsm-system`)
		r.Run(`export KUBECONFIG=$KUBECONFIG3 && kubectl delete ns nsm-system`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k ./clusters-configuration/cluster1`)
	r.Run(`kubectl get services nsmgr-proxy -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k ./clusters-configuration/cluster2`)
	r.Run(`kubectl get services nsmgr-proxy -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl apply -k ./clusters-configuration/cluster3`)
	r.Run(`kubectl get services registry -n nsm-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
}
func (s *Suite) TestFloatingKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/FloatingKernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG2`)
		r.Run(`kubectl delete ns ${NAMESPACE1}`)
		r.Run(`export KUBECONFIG=$KUBECONFIG1`)
		r.Run(`kubectl delete ns ${NAMESPACE2}`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`NAMESPACE1=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE1=${NAMESPACE1:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE1}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    metadata:` + "\n" + `      annotations:` + "\n" + `        registration-name: icmp-server@my.cluster3` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `          - name: NSM_NAME` + "\n" + `            valueFrom:` + "\n" + `              fieldRef:` + "\n" + `                fieldPath: metadata.annotations['registration-name']` + "\n" + `          - name: NSM_CIDR_PREFIX` + "\n" + `            value: 172.16.1.2/31` + "\n" + `          - name: NSM_SERVICE_NAMES` + "\n" + `            value: icmp-responder@my.cluster3` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE1} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `[[ ! -z $NSE ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`NAMESPACE2=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE2=${NAMESPACE2:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE2}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder@my.cluster3/nsm-1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nsc-kernel -n ${NAMESPACE2}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-kernel -n ${NAMESPACE2} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE2} -- ping -c 4 172.16.1.2`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE1} -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloatingKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/FloatingKernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG2`)
		r.Run(`kubectl delete ns ${NAMESPACE1}`)
		r.Run(`export KUBECONFIG=$KUBECONFIG1`)
		r.Run(`kubectl delete ns ${NAMESPACE2}`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`NAMESPACE1=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE1=${NAMESPACE1:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE1}` + "\n" + `` + "\n" + `resources:` + "\n" + `- nse.yaml` + "\n" + `` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE1} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `[[ ! -z $NSE ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`NAMESPACE2=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE2=${NAMESPACE2:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE2}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://my-networkservice-ip@my.cluster3/nsm-1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nsc-kernel -n ${NAMESPACE2}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-kernel -n ${NAMESPACE2} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE2} -- ping -c 4 172.16.1.2`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE1} -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestFloatingVl3() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG3 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster3?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster2?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
		r.Run(`export KUBECONFIG=$KUBECONFIG1 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster1?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster3?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `resources:` + "\n" + `-  https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster1?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-vl3-nse.yaml` + "\n" + `EOF`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster1?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `resources:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster2?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-vl3-nse.yaml` + "\n" + `EOF`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3/cluster2?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`nsc2=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`nsc1=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`ipAddr2=$(kubectl --kubeconfig=$KUBECONFIG2 exec -n ns-vl3-interdomain $nsc2 -- ifconfig nsm-1)` + "\n" + `ipAddr2=$(echo $ipAddr2 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 $ipAddr2`)
	r.Run(`kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`ipAddr1=$(kubectl --kubeconfig=$KUBECONFIG1 exec -n ns-vl3-interdomain $nsc1 -- ifconfig nsm-1)` + "\n" + `ipAddr1=$(echo $ipAddr1 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 $ipAddr1`)
	r.Run(`kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
}
func (s *Suite) TestFloatingVl3_scale_from_zero() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG3 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster3?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster2?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
		r.Run(`export KUBECONFIG=$KUBECONFIG1 && kubectl delete -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster1?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster3?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster1?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl apply -k https://github.com/networkservicemesh/deployments-k8s/examples/floating_interdomain/usecases/FloatingVl3-scale-from-zero/cluster2?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18`)
	r.Run(`nsc2=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`nsc1=$(kubectl get pods -l app=nsc-kernel -n ns-vl3-interdomain --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`ipAddr2=$(kubectl --kubeconfig=$KUBECONFIG2 exec -n ns-vl3-interdomain $nsc2 -- ifconfig nsm-1)` + "\n" + `ipAddr2=$(echo $ipAddr2 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 $ipAddr2`)
	r.Run(`kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc1 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`ipAddr1=$(kubectl --kubeconfig=$KUBECONFIG1 exec -n ns-vl3-interdomain $nsc1 -- ifconfig nsm-1)` + "\n" + `ipAddr1=$(echo $ipAddr1 | grep -Eo 'inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'| cut -c 11-)` + "\n" + `kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 $ipAddr1`)
	r.Run(`kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.0.0` + "\n" + `kubectl exec $nsc2 -n ns-vl3-interdomain -- ping -c 4 169.254.1.0`)
}
func (s *Suite) TestKernel2Vxlan2Kernel() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/Kernel2Vxlan2Kernel")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG1`)
		r.Run(`kubectl delete ns ${NAMESPACE2}`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2`)
		r.Run(`kubectl delete ns ${NAMESPACE1}`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`NAMESPACE1=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE1=${NAMESPACE1:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE1}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.2/31` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE1} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `[[ ! -z $NSE ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`NAMESPACE2=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE2=${NAMESPACE2:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE2}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder@my.cluster2/nsm-1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nsc-kernel -n ${NAMESPACE2}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-kernel -n ${NAMESPACE2} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE2} -- ping -c 4 172.16.1.2`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE1} -- ping -c 4 172.16.1.3`)
}
func (s *Suite) TestKernel2Wireguard2Kernel() {
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/usecases/Kernel2Wireguard2Kernel")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG1`)
		r.Run(`kubectl delete ns ${NAMESPACE2}`)
		r.Run(`export KUBECONFIG=$KUBECONFIG2`)
		r.Run(`kubectl delete ns ${NAMESPACE1}`)
	})
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`NAMESPACE1=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE1=${NAMESPACE1:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE1}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.2/31` + "\n" + `            - name: NSM_PAYLOAD` + "\n" + `              value: IP` + "\n" + `            - name: NSM_SERVICE_NAMES` + "\n" + `              value: my-networkservice-ip` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE1} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')` + "\n" + `[[ ! -z $NSE ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`NAMESPACE2=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/8c7217843c12ffeb9eee8e6ba5ac916c7e475f18/examples/floating_interdomain/usecases/namespace.yaml)[0])` + "\n" + `NAMESPACE2=${NAMESPACE2:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE2}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=8c7217843c12ffeb9eee8e6ba5ac916c7e475f18` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://my-networkservice-ip@my.cluster2/nsm-1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=5m pod -l app=nsc-kernel -n ${NAMESPACE2}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-kernel -n ${NAMESPACE2} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE2} -- ping -c 4 172.16.1.2`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE1} -- ping -c 4 172.16.1.3`)
}
