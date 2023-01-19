package export

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	authv1T "github.com/openshift/api/authorization/v1"
	secuv1T "github.com/openshift/api/security/v1"
	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	secuv1 "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterRolesInspector struct {
	config              *rest.Config
	authClient          *authv1.AuthorizationV1Client
	secuClient          *secuv1.SecurityV1Client
	clusterRoleBindings *authv1T.ClusterRoleBindingList
	sccs                *secuv1T.SecurityContextConstraintsList
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

	c.secuClient, err = secuv1.NewForConfig(c.config)
	if err != nil {
		log.Fatalf("Cannot create security client: %s", err)
	}

	c.clusterRoleBindings, err = c.authClient.ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot get ClusterRoleBindings: %s", err)
	}
	for _, clusterRoleBinding := range c.clusterRoleBindings.Items {
		log.Printf("Found ClusterRoleBinding %s/%s", clusterRoleBinding.RoleRef.Name, clusterRoleBinding.UserNames)
	}

	c.sccs, err = c.secuClient.SecurityContextConstraints().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot get SecurityContextConstraints: %s", err)
	}
	for _, scc := range c.sccs.Items {
		log.Printf("Found SecurityContextConstraints %s/%s", scc.Name, scc.Users)
	}
}

func (c *ClusterRolesInspector) ClusterRoleBindingsForSA(namespace string, serviceAccount string) []authv1T.ClusterRoleBinding {
	systemServiceAccount := SystemNameForSA(namespace, serviceAccount)
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

func (c *ClusterRolesInspector) SecurityContextConstraintsForSA(namespace string, serviceAccount string) []secuv1T.SecurityContextConstraints {
	systemServiceAccount := SystemNameForSA(namespace, serviceAccount)
	var sccs []secuv1T.SecurityContextConstraints
	for _, scc := range c.sccs.Items {
		for _, user := range scc.Users {
			if strings.Compare(user, systemServiceAccount) == 0 {
				sccs = append(sccs, scc)
				break
			}
		}
	}

	log.Printf("Found %d matching SecurityContextConstraints for %s/%s", len(sccs), namespace, serviceAccount)
	for _, scc := range sccs {
		log.Printf("%s ", scc.Name)
	}
	return sccs
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
