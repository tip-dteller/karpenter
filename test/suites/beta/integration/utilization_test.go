/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/aws/karpenter-core/pkg/test"
	"github.com/aws/karpenter/test/pkg/debug"
)

var _ = Describe("Utilization", Label(debug.NoWatch), Label(debug.NoEvents), func() {
	It("should provision one pod per node", func() {
		nodePool.Spec.Template.Spec.Requirements = append(nodePool.Spec.Template.Spec.Requirements, v1.NodeSelectorRequirement{
			Key:      v1.LabelInstanceTypeStable,
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"t3a.small"},
		})

		deployment := test.Deployment(test.DeploymentOptions{
			Replicas:   100,
			PodOptions: test.PodOptions{ResourceRequirements: v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1.5")}}}})

		env.ExpectCreated(nodeClass, nodePool, deployment)
		env.EventuallyExpectHealthyPodCount(labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels), int(*deployment.Spec.Replicas))
		env.ExpectCreatedNodeCount("==", int(*deployment.Spec.Replicas)) // One pod per node enforced by instance size
	})
})
