package installer

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/install-builder/pkg/config"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

type Installer struct {
	appConfig             *config.ApplicationConfig
	installerConfig       *config.InstallerConfig
	clusterRolesInspector *ClusterRolesInspector
}

func NewInstallerFromConfig(appConfig *config.ApplicationConfig, installerConfig *config.InstallerConfig, clusterRolesInspector *ClusterRolesInspector) *Installer {
	installer := Installer{appConfig: appConfig, installerConfig: installerConfig, clusterRolesInspector: clusterRolesInspector}

	return &installer
}

func (i *Installer) BuildKustomizeInstaller() {
	for _, ns := range i.appConfig.Application.Namespaces {
		log.Printf("Creating kustomize installer for NS %s", ns.Name)

		outputFolder := i.installerConfig.OutputFolderForNS(ns.Name)
		kustomizeFolder := i.installerConfig.BaseKustomizeFolderForNS(ns.Name)

		kustomization := filepath.Join(kustomizeFolder, config.KustomizationFile)
		os.Create(kustomization)
		AppendToFile(kustomization, "resources:")
		filepath.WalkDir(outputFolder, func(path string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
				yfile, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				decode := scheme.Codecs.UniversalDeserializer().Decode
				obj, gKV, err := decode(yfile, nil, nil)
				if err == nil {
					if gKV.Kind == "ServiceAccount" {
						serviceAccount := obj.(*v1.ServiceAccount)
						i.handleServiceAccount(kustomization, ns.Name, serviceAccount)
					}
				}

				// log.Printf("Moving %s to %s", d.Name(), kustomizeFolder)
				os.Rename(path, filepath.Join(kustomizeFolder, d.Name()))
				AppendToFile(kustomization, fmt.Sprintf("\n  - %s", d.Name()))
			}
			return nil
		})
		// RunCommand("move2kube", "plan", "--source", outputFolder, "--name", ns.Name)
		// RunCommand("move2kube", "transform", "--qa-skip", "true", "--output", i.installerConfig.InstallerFolder())
	}

	i.createKustomizeTemplate()
}

func (i *Installer) createKustomizeTemplate() {
	for _, ns := range i.appConfig.Application.Namespaces {
		log.Printf("Creating kustomize template for NS %s", ns.Name)
		templateFolder := i.installerConfig.KustomizeTemplateFolderForNS(ns.Name)

		paramsFolder := filepath.Join(templateFolder, config.ParamsFolder)
		os.Rename(i.installerConfig.TmpParamsFolderForNS(ns.Name), paramsFolder)
		secretsFolder := filepath.Join(templateFolder, config.SecretsFolder)
		os.Rename(i.installerConfig.TmpSecretsFolderForNS(ns.Name), secretsFolder)

		templateKustomization := i.installerConfig.KustomizationFileFrom(templateFolder)
		os.Create(templateKustomization)
		text := "resources:\n" +
			"  - ../base\n"
		AppendToFile(templateKustomization, text)

		text = "generatorOptions:\n" +
			"  disableNameSuffixHash: true\n" +
			"configMapGenerator:"
		AppendToFile(templateKustomization, text)
		err := filepath.WalkDir(paramsFolder,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					configMap := strings.Replace(d.Name(), ".env", "", 1)

					log.Printf("Creating configMapGenerator for %s", configMap)
					text = "\n" +
						"- name: %s\n" +
						"  behavior: merge\n" +
						"  envs:\n" +
						"  - %s/%s"
					AppendToFile(templateKustomization, text, configMap, config.ParamsFolder, d.Name())
				}
				return nil
			})
		if err == nil {
			text := "\nsecretGenerator:"
			AppendToFile(templateKustomization, text)
			err = filepath.WalkDir(secretsFolder,
				func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() {
						secret := strings.Replace(d.Name(), ".env", "", 1)
						log.Printf("Creating secretGenerator for %s", secret)
						text = "\n" +
							"- name: %s\n" +
							"  behavior: create\n" +
							"  envs:\n" +
							"  - %s/%s"
						AppendToFile(templateKustomization, text, secret, config.SecretsFolder, d.Name())
					}
					return nil
				})
			if err != nil {
				log.Fatalf("Cannot create kustomize template: %s", err)
			}
		}
	}
}

func (i *Installer) handleServiceAccount(kustomization string, namespace string, serviceAccount *v1.ServiceAccount) {
	log.Printf("Handling ServiceAccount %s", serviceAccount.Name)

	clusterRoleBindings := i.clusterRolesInspector.ClusterRoleBindingsForSA(namespace, serviceAccount.Name)

	for _, clusterRoleBinding := range clusterRoleBindings {
		// TODO: update CRB name
		yamlFile := fmt.Sprintf("%s-%s.yaml", "ClusterRoleBinding", clusterRoleBinding.Name)
		yamlPath := filepath.Join(i.installerConfig.BaseKustomizeFolderForNS(namespace), yamlFile)

		clusterRoleBindingSpec := rbacV1.ClusterRoleBinding{
			// TODO: These two are not collected by client-go API
			TypeMeta: metaV1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: metaV1.ObjectMeta{
				Name: clusterRoleBinding.Name,
			},
			RoleRef: rbacV1.RoleRef{
				Kind: clusterRoleBinding.RoleRef.Kind,
				Name: clusterRoleBinding.RoleRef.Name,
				// Do not copy the original namespace, will be overriden at install time
			},
			// TODO: API Group
			Subjects: []rbacV1.Subject{},
		}
		for _, subject := range clusterRoleBinding.Subjects {
			clusterRoleBindingSpec.Subjects = append(clusterRoleBindingSpec.Subjects, rbacV1.Subject{
				Kind:      subject.Kind,
				Name:      subject.Name,
				Namespace: subject.Namespace,
			})
			// TODO: API Group
		}

		log.Printf("Creating YAML %s for ClusterRoleBinding %s to assign role %s to ServiceAccount %s", yamlFile,
			clusterRoleBindingSpec.Name, clusterRoleBindingSpec.RoleRef.Name, serviceAccount.Name)
		newFile, err := os.Create(yamlPath)
		if err != nil {
			log.Fatal(err)
		}
		y := printers.YAMLPrinter{}
		defer newFile.Close()
		if err = y.PrintObj(&clusterRoleBindingSpec, newFile); err != nil {
			log.Fatal(err)
		}

		AppendToFile(kustomization, fmt.Sprintf("\n  - %s", yamlFile))
	}
}
