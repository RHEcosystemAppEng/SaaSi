package installer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	authv1T "github.com/openshift/api/authorization/v1"
	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterRolesInspector struct {
	config              *rest.Config
	authClient          *authv1.AuthorizationV1Client
	clusterRoleBindings *authv1T.ClusterRoleBindingList
}

func NewClusterRolesInspector() *ClusterRolesInspector {
	return &ClusterRolesInspector{}
}

func (c *ClusterRolesInspector) LoadClusterRoles() {
	err := c.connectCluster()
	if err != nil {
		log.Fatalf("Cannot connect to default cluster: %s", err)
	}
	log.Print("Cluster connected")

	c.authClient, err = authv1.NewForConfig(c.config)
	if err != nil {
		log.Fatalf("Cannot create auth client: %s", err)
	}

	c.clusterRoleBindings, err = c.authClient.ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot get ClusterRoleBindings: %s", err)
	}
	for _, clusterRoleBinding := range c.clusterRoleBindings.Items {
		log.Printf("Found ClusterRoleBindings %s/%s", clusterRoleBinding.RoleRef.Name, clusterRoleBinding.UserNames)
	}
}

func (c *ClusterRolesInspector) ClusterRoleBindingsForSA(namespace string, serviceAccount string) []authv1T.ClusterRoleBinding {
	systemServiceAccount := fmt.Sprintf("system:serviceaccount:%s:%s", namespace, serviceAccount)
	var clusterRoleBindings []authv1T.ClusterRoleBinding
	for _, clusterRoleBinding := range c.clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
			if strings.Compare(subject.Kind, "ServiceAccount") == 0 {
				for _, rbUserName := range clusterRoleBinding.UserNames {
					if strings.Compare(rbUserName, systemServiceAccount) == 0 && strings.Compare(namespace, subject.Namespace) == 0 {
						clusterRoleBindings = append(clusterRoleBindings, clusterRoleBinding)
						break
					}
				}
			}
		}
	}

	log.Printf("Found %d matching ClusterRoleBindings for %s/%s", len(clusterRoleBindings), namespace, serviceAccount)
	for _, clusterRoleBinding := range clusterRoleBindings {
		log.Printf("%s ", clusterRoleBinding.Name)
	}
	return clusterRoleBindings
}

func (c *ClusterRolesInspector) connectCluster() error {
	kubeconfig := ""

	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	//Load config for Openshift's go-client from kubeconfig file
	var err error
	c.config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	return err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
