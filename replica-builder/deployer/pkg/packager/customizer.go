package packager

import (
	"log"
	"fmt"
	"bytes"
    "io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/utils"
)

func (pkg *ApplicationPkg)invokeNsCustomizations(ns config.SourceNamespace) { 

	// validate kustomize cli
	utils.ValidateRequirements()

	// set kustomize.yaml
	pkg.customizeKustomize(ns)

	// set configMaps
	customizeParams(ns, CONFIGMAPS_DIR)

	// set secrets
	customizeParams(ns, SECRETS_DIR)
	

}

func (pkg *ApplicationPkg)customizeKustomize(ns config.SourceNamespace) {
	
	// set the namespace resource to target namespace
	cmd := exec.Command("kustomize", "edit", "set", "namespace", ns.Target)
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to set namespace resource in %s template", ns.Name)
	}

	// set a common annotation to uuid
	cmd = exec.Command("kustomize", "edit", "set", "annotation", COMMON_ANNOTATION_KEY + pkg.Uuid.String())
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to set uuid common annotations in %s template", ns.Name)
	}
}

func customizeParams(ns config.SourceNamespace, paramsDir string ) {
	
	// select param type
	if paramsDir == CONFIGMAPS_DIR {
		// customize configmap params
		for _, configMap := range ns.ConfigMaps {

			// define configmap filepath
			configMapsFilepath := filepath.Join(nsTmplDir, paramsDir, configMap.ConfigMap + ".env")
			// if configmap filepath exists, replace configmap params with custom values
			if utils.FileExists(configMapsFilepath) {
				// replace configmap params with custom values
				replaceParamValues(configMapsFilepath, configMap.Params)
			} else {
				log.Printf("WARNING: configmap \"%s\" does not exist", configMap.ConfigMap)
			}
		}
	} else {
		// customize secret params
		for _, secret := range ns.Secrets {

			// define secret filepath
			secretsMapFilepath := filepath.Join(nsTmplDir, paramsDir, secret.Secret + ".env")
			// if secret filepath exists, replace secret params with custom values
			if utils.FileExists(secretsMapFilepath) {
				// replace secret params with custom values
				replaceParamValues(secretsMapFilepath, secret.Params) 
			} else {
				log.Printf("WARNING: secret \"%s\" does not exist", secret.Secret)
			}
		}
	}
 }


 func replaceParamValues(file string, params []config.Params) {

	// read param file
	output, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Could not read file %s", file)
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
				log.Printf("WARNING: \"%s\" no such param exists in file %s", param.Name, file)
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
		log.Fatalf("Could not update file %s with custom params", file)
	} 
}