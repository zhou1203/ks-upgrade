package globalrolebindings

import (
	"encoding/json"
	"fmt"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"kubesphere.io/ks-upgrade/pkg/task"
)

type globalRoleBindingMigrateTask struct {
	k8sClient kubernetes.Interface
}

func NewGlobalRoleBindingMigrateTask(k8sClient kubernetes.Interface) task.UpgradeTask {
	return &globalRoleBindingMigrateTask{k8sClient: k8sClient}
}

func (t *globalRoleBindingMigrateTask) Run() error {

	globalRoleBindings := make([]GlobalRoleBinding, 0)

	clusterRoleBindings, err := t.k8sClient.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})

	if err != nil {
		klog.Error(err)
		return err
	}

	migrateMapping := map[string]string{
		"cluster-admin":      "platform-admin",
		"cluster-regular":    "platform-regular",
		"workspaces-manager": "workspaces-manager",
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		globalRole := migrateMapping[clusterRoleBinding.RoleRef.Name]
		if globalRole == "" {
			continue
		}
		for _, subject := range clusterRoleBinding.Subjects {
			if subject.Kind == "User" {
				globalRoleBinding := GlobalRoleBinding{
					TypeMeta: metav1.TypeMeta{
						Kind:       "GlobalRoleBinding",
						APIVersion: "iam.kubesphere.io/v1alpha2",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: fmt.Sprintf("%s-%s", subject.Name, globalRole),
						Labels: map[string]string{
							"iam.kubesphere.io/user-ref": subject.Name,
						},
					},
					Subjects: []rbacv1.Subject{{Name: subject.Name, Kind: "User", APIGroup: "rbac.authorization.k8s.io"}},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "iam.kubesphere.io",
						Kind:     "GlobalRole",
						Name:     globalRole,
					},
				}
				globalRoleBindings = append(globalRoleBindings, globalRoleBinding)
			}
		}
	}

	cli := t.k8sClient.(*kubernetes.Clientset)
	for _, globalRoleBinding := range globalRoleBindings {
		outputData, _ := json.Marshal(globalRoleBinding)
		klog.Infof("migrate globalRoleBinding: %s: %s", globalRoleBinding.Name, string(outputData))
		err := cli.RESTClient().
			Post().
			AbsPath(fmt.Sprintf("/apis/iam.kubesphere.io/v1alpha2/globalrolebindings")).
			Body(outputData).
			Do().Error()

		if err != nil && !errors.IsAlreadyExists(err) {
			klog.Error(err)
			return err
		}
	}

	return nil
}
