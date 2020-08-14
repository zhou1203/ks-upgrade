package workspaces

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type WorkspaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederatedWorkspaceSpec `json:"spec,omitempty"`
}

type FederatedWorkspace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederatedWorkspaceSpec `json:"spec"`
}

type FederatedWorkspaceSpec struct {
	Template  Workspace              `json:"template"`
	Placement GenericPlacementFields `json:"placement"`
	Overrides []GenericOverrideItem  `json:"overrides,omitempty"`
}

type Workspace struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkspaceSpec `json:"spec,omitempty"`
}

type WorkspaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workspace `json:"items"`
}

type WorkspaceSpec struct {
	Manager          string `json:"manager,omitempty"`
	NetworkIsolation *bool  `json:"networkIsolation,omitempty"`
}

type GenericClusterReference struct {
	Name string `json:"name"`
}

type GenericPlacementFields struct {
	Clusters        []GenericClusterReference `json:"clusters,omitempty"`
	ClusterSelector *metav1.LabelSelector     `json:"clusterSelector,omitempty"`
}

type ClusterOverride struct {
	Op    string               `json:"op,omitempty"`
	Path  string               `json:"path"`
	Value runtime.RawExtension `json:"value,omitempty"`
}

type GenericOverrideItem struct {
	ClusterName      string            `json:"clusterName"`
	ClusterOverrides []ClusterOverride `json:"clusterOverrides,omitempty"`
}
