package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/sirupsen/logrus"
)

type Installer struct {
	appContext *AppContext

	sccToBeReplacedByNS map[string][]SccForSA
	logger              *logrus.Logger
}

type SccForSA struct {
	serviceAccountName string
	sccName            string
}

func NewInstallerFromConfig(appContext *AppContext) *Installer {
	installer := Installer{appContext: appContext}

	installer.logger = appContext.Logger
	installer.sccToBeReplacedByNS = make(map[string][]SccForSA)
	return &installer
}

func (i *Installer) BuildKustomizeInstaller() error {
	for _, ns := range i.appContext.AppConfig.Namespaces {
		i.logger.Infof("Creating kustomize installer for NS %s", ns.Name)

		outputFolder := i.appContext.OutputFolderForNS(ns.Name)
		kustomizeFolder := i.appContext.BaseKustomizeFolderForNS(ns.Name)

		kustomization := filepath.Join(kustomizeFolder, KustomizationFile)
		os.Create(kustomization)
		utils.AppendToFile(kustomization, "resources:")
		filepath.WalkDir(outputFolder, func(path string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
				i.logger.Debugf("Moving %s to %s", d.Name(), kustomizeFolder)
				os.Rename(path, filepath.Join(kustomizeFolder, d.Name()))
				utils.AppendToFile(kustomization, fmt.Sprintf("\n  - %s", d.Name()))
			}
			return nil
		})
	}

	err := i.createKustomizeTemplate()
	return err
}

func (i *Installer) createKustomizeTemplate() error {
	for _, ns := range i.appContext.AppConfig.Namespaces {
		i.logger.Infof("Creating kustomize template for NS %s", ns.Name)
		templateFolder := i.appContext.KustomizeTemplateFolderForNS(ns.Name)

		paramsFolder := filepath.Join(templateFolder, ParamsFolder)
		os.Rename(i.appContext.TmpParamsFolderForNS(ns.Name), paramsFolder)
		secretsFolder := filepath.Join(templateFolder, SecretsFolder)
		os.Rename(i.appContext.TmpSecretsFolderForNS(ns.Name), secretsFolder)

		templateKustomization := i.appContext.KustomizationFileFrom(templateFolder)
		os.Create(templateKustomization)
		text := "resources:\n" +
			"  - ../base\n"
		err := utils.AppendToFile(templateKustomization, text)
		if err != nil {
			return err
		}

		text = "generatorOptions:\n" +
			"  disableNameSuffixHash: true\n" +
			"configMapGenerator:"
		err = utils.AppendToFile(templateKustomization, text)
		if err != nil {
			return err
		}
		err = filepath.WalkDir(paramsFolder,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					configMap := strings.Replace(d.Name(), ".env", "", 1)

					i.logger.Infof("Creating configMapGenerator for %s", configMap)
					text = "\n" +
						"- name: %s\n" +
						"  behavior: merge\n" +
						"  envs:\n" +
						"  - %s/%s"
					utils.AppendToFile(templateKustomization, text, configMap, ParamsFolder, d.Name())
				}
				return nil
			})
		if err != nil {
			return err
		}
		text = "\nsecretGenerator:"
		utils.AppendToFile(templateKustomization, text)
		err = filepath.WalkDir(secretsFolder,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					secret := strings.Replace(d.Name(), ".env", "", 1)
					i.logger.Infof("Creating secretGenerator for %s", secret)
					text = "\n" +
						"- name: %s\n" +
						"  behavior: create\n" +
						"  envs:\n" +
						"  - %s/%s"
					utils.AppendToFile(templateKustomization, text, secret, SecretsFolder, d.Name())
				}
				return nil
			})
		if err != nil {
			i.logger.Errorf("Cannot create kustomize template: %s", err)
			return err
		}

		if len(i.sccToBeReplacedByNS[ns.Name]) > 0 {
			text := "\nreplacements:"
			utils.AppendToFile(templateKustomization, text)

			for _, sccForSA := range i.sccToBeReplacedByNS[ns.Name] {
				text = "\n" +
					"- source:\n" +
					"    kind: ServiceAccount\n" +
					"    name: %s\n" +
					"    fieldPath: metadata.namespace\n" +
					"  targets:\n" +
					"  - select:\n" +
					"      kind: SecurityContextConstraints\n" +
					"      name: %s\n" +
					"    fieldPaths:\n" +
					"    - users.*\n" +
					"    options:\n" +
					"      delimiter: \":\"\n" +
					"      index: 2\n"
				utils.AppendToFile(templateKustomization, text, sccForSA.serviceAccountName, sccForSA.sccName)
			}
		}
	}
	return nil
}
