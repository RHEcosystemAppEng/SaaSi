package packager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
)

func (pkg *ApplicationPkg) invokeNsCustomizations(ns config.Namespaces) {

	// validate kustomize cli
	err = utils.ValidateRequirements(utils.KUSTOMIZE)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("%s command not found", utils.KUSTOMIZE)
		pkg.Error = err
		return
	}

	// set kustomize.yaml
	pkg.customizeKustomize(ns)
	if pkg.Error != nil {
		return
	}

	// set configMaps
	pkg.customizeParams(ns, CONFIGMAPS_DIR)
	if pkg.Error != nil {
		return
	}

	// set secrets
	pkg.customizeParams(ns, SECRETS_DIR)
	if pkg.Error != nil {
		return
	}

	// check if any unset mandatory params still exist in namespace template files
	pkg.inspectMandatoryParams(ns)

}

func (pkg *ApplicationPkg) customizeKustomize(ns config.Namespaces) {

	// set the namespace resource to target namespace
	cmd := exec.Command("kustomize", "edit", "set", "namespace", ns.Target)
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to set namespace resource in %s kustomize template", ns.Name)
		pkg.Error = err
		return
	}

	// set a common annotation to uuid
	cmd = exec.Command("kustomize", "edit", "set", "annotation", COMMON_ANNOTATION_KEY+pkg.Uuid.String())
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to set uuid common annotations in %s kustomize template", ns.Name)
		pkg.Error = err
	}
}

func (pkg *ApplicationPkg) customizeParams(ns config.Namespaces, paramsDir string) {

	// select param type
	if paramsDir == CONFIGMAPS_DIR {
		// customize configmap params
		for _, configMap := range ns.ConfigMaps {

			// define configmap filepath
			configMapFilepath := filepath.Join(nsTmplDir, paramsDir, configMap.ConfigMap+".env")
			// if configmap filepath exists, replace configmap params with custom values
			if utils.FileExists(configMapFilepath) {
				// replace configmap params with custom values
				pkg.replaceParamValues(configMapFilepath, configMap.Params)
				if pkg.Error != nil {
					return
				}
			} else {
				pkg.DeployerContext.GetLogger().Warningf("ConfigMap \"%s\" does not exist", configMapFilepath)
			}
		}
	} else {
		// customize secret params
		for _, secret := range ns.Secrets {

			// define secret filepath
			secretFilepath := filepath.Join(nsTmplDir, paramsDir, secret.Secret+".env")
			// if secret filepath exists, replace secret params with custom values
			if utils.FileExists(secretFilepath) {
				// replace secret params with custom values
				pkg.replaceParamValues(secretFilepath, secret.Params)
				if pkg.Error != nil {
					return
				}
			} else {
				pkg.DeployerContext.GetLogger().Warningf("Secret \"%s\" does not exist", secretFilepath)
			}
		}
	}
}

func (pkg *ApplicationPkg) replaceParamValues(file string, params []config.ApplicationParams) {

	// read param file
	output, err := ioutil.ReadFile(file)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Could not read file %s", file)
		pkg.Error = err
		return
	}

	// parse param file output to string
	stringOutput := string(output)

	for _, param := range params {

		// locate param with empty placeholder
		source := fmt.Sprintf("#%s=%s", param.Name, EMPTY_PLACEHOLDER)
		if !strings.Contains(stringOutput, source) {
			// if param with empty placeholder does not exist, locate param with mandatory placeholder
			source = fmt.Sprintf("%s=%s", param.Name, MANDATORY_PLACEHOLDER)
			if !strings.Contains(stringOutput, source) {
				pkg.DeployerContext.GetLogger().Warningf("\"%s\" no such param exists in file %s", param.Name, file)
				continue
			}
		}

		// replace param placeholder with custome value
		target := fmt.Sprintf("%s=%s", param.Name, param.Value)

		// make replacement
		output = bytes.Replace(output, []byte(source), []byte(target), -1)
	}

	// write changes to param file
	if err = ioutil.WriteFile(file, output, 0666); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Could not update file %s with custom params", file)
		pkg.Error = err
	}
}
