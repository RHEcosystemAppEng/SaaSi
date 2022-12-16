package installer

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/install-builder/pkg/config"
)

type Installer struct {
	appConfig       *config.ApplicationConfig
	installerConfig *config.InstallerConfig
}

func NewInstallerFromConfig(appConfig *config.ApplicationConfig, installerConfig *config.InstallerConfig) *Installer {
	installer := Installer{appConfig: appConfig, installerConfig: installerConfig}

	return &installer
}

func (i *Installer) BuildInstallerWithMove2Kube() {
	for _, ns := range i.appConfig.Application.Namespaces {
		log.Printf("Creating installer for NS %s with move2kube", ns.Name)

		outputFolder := i.installerConfig.OutputFolderForNS(ns.Name)

		RunCommand("move2kube", "plan", "--source", outputFolder, "--name", ns.Name)
		RunCommand("move2kube", "transform", "--qa-skip", "true", "--output", i.installerConfig.InstallerFolder())
	}
}

func (i *Installer) UpdateKustomize() {
	i.updateKustomizeBase()
	i.updateKustomizeOverlays()
}

func (i *Installer) updateKustomizeBase() {
	for _, ns := range i.appConfig.Application.Namespaces {
		log.Printf("Updating kustomize base for NS %s", ns.Name)
		baseFolder := i.installerConfig.BaseKustomizeFolderForNS(ns.Name)
		baseKustomization := i.installerConfig.KustomizationFileFrom(baseFolder)
		text := "generatorOptions:\n" +
			"  disableNameSuffixHash: true\n" +
			"configMapGenerator:"
		AppendToFile(baseKustomization, text)

		paramsFolder := i.installerConfig.KustomizeParamsFolderForNS(ns.Name)
		os.Rename(i.installerConfig.TmpParamsFolderForNS(ns.Name), paramsFolder)

		err := filepath.WalkDir(paramsFolder,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() && path != paramsFolder {
					configMap := d.Name()
					log.Printf("Creating configMapGenerator for %s", configMap)
					text = "\n" +
						"- name: %s\n" +
						"  behavior: create\n" +
						"  files:"
					AppendToFile(baseKustomization, text, configMap)

					err = filepath.WalkDir(path,
						func(path string, d fs.DirEntry, err error) error {
							if err != nil {
								return err
							}
							if !d.IsDir() {
								keyName := d.Name()
								text = "\n" +
									"  - params/%s/%s"
								AppendToFile(baseKustomization, text, configMap, keyName)
							}
							return nil
						})
				}
				return err
			})
		if err != nil {
			log.Fatalf("Cannot customize the base kustomize configuration: %s", err)
		}
	}
}
func (i *Installer) updateKustomizeOverlays() {
	for _, ns := range i.appConfig.Application.Namespaces {
		log.Printf("Updating kustomize overlays for NS %s", ns.Name)
		overlaysFolder := i.installerConfig.KustomizeOverlaysFolderForNS(ns.Name)
		paramsFolder := i.installerConfig.KustomizeParamsFolderForNS(ns.Name)

		err := filepath.WalkDir(overlaysFolder,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() && path != overlaysFolder {
					log.Printf("Updating overlay %s", d.Name())
					secretsFolder := filepath.Join(path, config.SecretsFolder)
					os.Rename(i.installerConfig.TmpSecretsFolderForNS(ns.Name), secretsFolder)
					kustomization := i.installerConfig.KustomizationFileFrom(path)

					text := "generatorOptions:\n" +
						"  disableNameSuffixHash: true\n" +
						"configMapGenerator:"
					AppendToFile(kustomization, text)
					os.Create(filepath.Join(path, "custom.env"))
					err = filepath.WalkDir(paramsFolder,
						func(path string, d fs.DirEntry, err error) error {
							if err != nil {
								return err
							}
							if d.IsDir() && path != paramsFolder {
								configMap := d.Name()
								log.Printf("Creating configMapGenerator for %s", configMap)
								text = "\n" +
									"- name: %s\n" +
									"  behavior: merge\n" +
									"  envs:\n" +
									"  - custom.env"
								AppendToFile(kustomization, text, configMap)
							}
							return nil
						})
					if err == nil {
						text := "\nsecretGenerator:"
						AppendToFile(kustomization, text)
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
									AppendToFile(kustomization, text, secret, config.SecretsFolder, d.Name())
								}
								return nil
							})
					}
				}
				return err
			})
		if err != nil {
			log.Fatalf("Cannot customize the kustomize overlays: %s", err)
		}
	}
}

/*
    touch ${overlaysKustomizeFolder}/${o}/custom.env

    for configmap in ${paramsFolder}/*
    do
      configmap=$(basename ${configmap})
      log "Creating configMapGenerator for ${configmap}"

      echo -n "
- name: ${configmap}
  behavior: merge
  envs:
  - custom.env" >> ${overlaysKustomizeFolder}/${o}/kustomization.yaml
    done
  done
*/
