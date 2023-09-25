package ovnwebhook

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/maps"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	nodeutil "k8s.io/component-helpers/node/util"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	hotypes "github.com/ovn-org/ovn-kubernetes/go-controller/hybrid-overlay/pkg/types"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/util"
)

// checkNodeAnnot defines additional checks for the allowed annotations
type checkNodeAnnot func(v annotationChange, nodeName string) error

// commonNodeAnnotationChecks holds annotations allowed for ovnkube-node:<nodeName> users in non-IC and IC environments
var commonNodeAnnotationChecks = map[string]checkNodeAnnot{
	util.OvnNodeHostAddresses:            nil,
	util.OvnNodeL3GatewayConfig:          nil,
	util.OvnNodeManagementPortMacAddress: nil,
	util.OvnNodeIfAddr:                   nil,
	util.OvnNodeGatewayMtuSupport:        nil,
	util.OvnNodeChassisID: func(v annotationChange, nodeName string) error {
		if v.action == removed {
			return fmt.Errorf("%s cannot be removed", util.OvnNodeChassisID)
		}
		if v.action == changed {
			return fmt.Errorf("%s cannot be changed once set", util.OvnNodeChassisID)
		}
		return nil
	},
}

// hybridOverlayNodeAnnotationChecks holds annotations allowed for ovnkube-node:<nodeName> users hybrid overlay environments
var hybridOverlayNodeAnnotationChecks = map[string]checkNodeAnnot{
	hotypes.HybridOverlayDRMAC: nil,
	hotypes.HybridOverlayDRIP:  nil,
}

type NodeAdmission struct {
	annotationChecks  map[string]checkNodeAnnot
	annotationKeys    sets.Set[string]
	extraAllowedUsers sets.Set[string]
}

func NewNodeAdmissionWebhook(enableHybridOverlay bool, extraAllowedUsers ...string) *NodeAdmission {
	checks := make(map[string]checkNodeAnnot)
	maps.Copy(checks, commonNodeAnnotationChecks)
	if enableHybridOverlay {
		maps.Copy(checks, hybridOverlayNodeAnnotationChecks)
	}
	return &NodeAdmission{
		annotationChecks:  checks,
		annotationKeys:    sets.New[string](maps.Keys(checks)...),
		extraAllowedUsers: sets.New[string](extraAllowedUsers...),
	}
}

var _ admission.CustomValidator = &NodeAdmission{}

func (p NodeAdmission) ValidateCreate(_ context.Context, _ runtime.Object) (err error) {
	// Ignore creation, the webhook is configured to only handle nodes/status updates
	return nil
}

func (p NodeAdmission) ValidateDelete(_ context.Context, _ runtime.Object) (err error) {
	// Ignore deletion, the webhook is configured to only handle nodes/status updates
	return nil
}

func (p NodeAdmission) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (err error) {
	oldNode := oldObj.(*corev1.Node)
	newNode := newObj.(*corev1.Node)

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return err
	}
	nodeName, isOVNKubeNode := ovnkubeNodeIdentity(req.UserInfo)

	changes := mapDiff(oldNode.Annotations, newNode.Annotations)
	changedKeys := maps.Keys(changes)

	if !isOVNKubeNode {
		if !p.annotationKeys.HasAny(changedKeys...) {
			// the user is not an ovnkube-node and hasn't changed any ovnkube-node annotations
			return nil
		}

		if !p.extraAllowedUsers.Has(req.UserInfo.Username) {
			// the user is not in extraAllowedUsers
			return fmt.Errorf("user %q is not allowed to set the following annotations on node: %q: %v",
				req.UserInfo.Username,
				newNode.Name,
				p.annotationKeys.Intersection(sets.New[string](changedKeys...)).UnsortedList())
		}

		// The user is not ovnkube-node, in this case the nodeName comes from the object
		nodeName = newNode.Name
	}

	for _, key := range changedKeys {
		if check := p.annotationChecks[key]; check != nil {
			if err := check(changes[key], nodeName); err != nil {
				return fmt.Errorf("user: %q is not allowed to set %s on node %q: %v", req.UserInfo.Username, key, newNode.Name, err)
			}
		}
	}

	// All the checks beyond this point are ovnkube-node specific
	// If the user is not ovnkube-node exit here
	if !isOVNKubeNode {
		return nil
	}

	if newNode.Name != nodeName {
		return fmt.Errorf("ovnkube-node on node: %q is not allowed to modify nodes %q annotations", nodeName, newNode.Name)
	}

	// ovnkube-node is not allowed to change annotations outside of it's scope
	if !p.annotationKeys.HasAll(changedKeys...) {
		return fmt.Errorf("ovnkube-node on node: %q is not allowed to set the following annotations: %v",
			nodeName,
			sets.New[string](changedKeys...).Difference(p.annotationKeys).UnsortedList())
	}

	// Verify that nothing but the annotations changed.
	// Since ovnkube-node only has the node/status permissions, it is enough to check .Status and .ObjectMeta only.
	// Ignore .ManagedFields fields which are modified on every update.
	oldNodeShallowCopy := oldNode
	newNodeShallowCopy := newNode
	oldNodeShallowCopy.Annotations = nil
	newNodeShallowCopy.Annotations = nil
	oldNodeShallowCopy.ManagedFields = nil
	newNodeShallowCopy.ManagedFields = nil

	if strings.HasPrefix(newNode.Spec.ProviderID, "gce") {
		oldId, oldCondition := nodeutil.GetNodeCondition(&(oldNodeShallowCopy.Status), corev1.NodeNetworkUnavailable)
		_, newCondition := nodeutil.GetNodeCondition(&(newNodeShallowCopy.Status), corev1.NodeNetworkUnavailable)

		// In GCP OVN-Kubernetes modifies NodeNetworkUnavailable condition on the nodes, allow an update of the condition
		// https://github.com/ovn-org/ovn-kubernetes/blob/e2e442133f16699671bb6564c4b8863229841fd9/go-controller/pkg/ovn/master.go#L507
		if oldId >= 0 && newCondition != nil && !reflect.DeepEqual(oldCondition, newCondition) {
			// Replace the old NodeNetworkUnavailable condition with the new one to make DeepEqual happy
			conditionsDeepCopy := make([]corev1.NodeCondition, len(oldNodeShallowCopy.Status.Conditions))
			copy(conditionsDeepCopy, oldNodeShallowCopy.Status.Conditions)
			conditionsDeepCopy[oldId] = *newCondition

			oldNodeShallowCopy.Status.Conditions = conditionsDeepCopy
		}
	}
	if !apiequality.Semantic.DeepEqual(oldNodeShallowCopy.ObjectMeta, newNodeShallowCopy.ObjectMeta) ||
		!apiequality.Semantic.DeepEqual(oldNodeShallowCopy.Status, newNodeShallowCopy.Status) {
		return fmt.Errorf("ovnkube-node on node: %q is not allowed to modify anything other than annotations", nodeName)
	}

	return nil
}
