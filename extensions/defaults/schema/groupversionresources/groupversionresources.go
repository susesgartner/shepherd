package groupversionresources

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Node() schema.GroupVersionResource {
	node := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}

	return node
}

func Pod() schema.GroupVersionResource {
	pod := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	return pod
}

func ConfigMap() schema.GroupVersionResource {
	configMap := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}

	return configMap
}

func CustomResourceDefinition() schema.GroupVersionResource {
	customResourceDefinition := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	return customResourceDefinition
}

func Ingress() schema.GroupVersionResource {
	ingress := schema.GroupVersionResource{
		Group:    "networking.k8s.io",
		Version:  "v1",
		Resource: "ingresses",
	}

	return ingress
}

func Project() schema.GroupVersionResource {
	project := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projects",
	}

	return project
}

func Role() schema.GroupVersionResource {
	role := schema.GroupVersionResource{
		Group:    rbacv1.SchemeGroupVersion.Group,
		Version:  rbacv1.SchemeGroupVersion.Version,
		Resource: "roles",
	}

	return role
}

func ClusterRole() schema.GroupVersionResource {
	clusterRole := schema.GroupVersionResource{
		Group:    rbacv1.SchemeGroupVersion.Group,
		Version:  rbacv1.SchemeGroupVersion.Version,
		Resource: "clusterroles",
	}

	return clusterRole
}

func RoleBinding() schema.GroupVersionResource {
	roleBinding := schema.GroupVersionResource{
		Group:    rbacv1.SchemeGroupVersion.Group,
		Version:  rbacv1.SchemeGroupVersion.Version,
		Resource: "rolebindings",
	}

	return roleBinding
}

func ClusterRoleBinding() schema.GroupVersionResource {
	clusterRoleBinding := schema.GroupVersionResource{
		Group:    rbacv1.SchemeGroupVersion.Group,
		Version:  rbacv1.SchemeGroupVersion.Version,
		Resource: "clusterrolebindings",
	}

	return clusterRoleBinding
}

func GlobalRole() schema.GroupVersionResource {
	globalRole := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "globalroles",
	}

	return globalRole
}

func GlobalRoleBinding() schema.GroupVersionResource {
	globalRoleBinding := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "globalrolebindings",
	}

	return globalRoleBinding
}

func ClusterRoleTemplateBinding() schema.GroupVersionResource {
	clusterRoleTemplateBinding := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "clusterroletemplatebindings",
	}

	return clusterRoleTemplateBinding
}

func RoleTemplate() schema.GroupVersionResource {
	roleTemplate := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "roletemplates",
	}

	return roleTemplate
}

func ProjectRoleTemplateBinding() schema.GroupVersionResource {
	projectRoleTemplateBinding := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projectroletemplatebindings",
	}

	return projectRoleTemplateBinding
}

func ResourceQuota() schema.GroupVersionResource {
	resourceQuota := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "resourcequotas",
	}

	return resourceQuota
}

func Secret() schema.GroupVersionResource {
	secret := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}

	return secret
}

func Service() schema.GroupVersionResource {
	service := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}

	return service
}

func StorageClass() schema.GroupVersionResource {
	storageClass := schema.GroupVersionResource{
		Group:    "storage.k8s.io",
		Version:  "v1",
		Resource: "storageclasses",
	}

	return storageClass
}

func Token() schema.GroupVersionResource {
	token := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "tokens",
	}

	return token
}

func PersistentVolumeClaim() schema.GroupVersionResource {
	persistentVolumeClaim := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumeclaims",
	}

	return persistentVolumeClaim
}

func PersistentVolume() schema.GroupVersionResource {
	persistentVolume := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumes",
	}

	return persistentVolume
}

func Namespace() schema.GroupVersionResource {
	namespace := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}

	return namespace
}

func Daemonset() schema.GroupVersionResource {
	daemonset := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "daemonsets",
	}

	return daemonset
}

func Deployment() schema.GroupVersionResource {
	deployment := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	return deployment
}

func Job() schema.GroupVersionResource {
	job := schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	}

	return job
}

func CronJob() schema.GroupVersionResource {
	cronJob := schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1beta1",
		Resource: "cronjobs",
	}

	return cronJob
}
