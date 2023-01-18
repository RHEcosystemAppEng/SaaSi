package packager

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/config"
)

func (pkg *ApplicationPkg) inspectMandatoryParams(ns config.Namespaces) {

	// find each param file in namespace template and inspect for unset mandatory params
	err = filepath.Walk(nsTmplDir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if file is not a directory and is of type "env", check file for mandatory params
		if !info.IsDir() && filepath.Ext(file) == PARAM_FILE_EXT {
			// if mandatory params still exist in file, save data
			if mandatoryParamsList := getUnsetMandatoryParams(file); len(mandatoryParamsList) > 0 {
				pkg.UnsetMandatoryParams[file] = mandatoryParamsList
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to reach %s parameter files to inspect for unset mandatory parameters, Error: %s", ns.Name, err)
	}
}

func getUnsetMandatoryParams(file string) []string {

	var mandatoryParams []string

	// read param file
	output, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to get unset mandatory parameters, could not read parameter file %s, Error: %s", file, err)
	}

	// parse param file output to string
	stringOutput := string(output)

	// check param file lines for mandatory placeholder
	fileLines := strings.Split(stringOutput, "\n")
	for _, line := range fileLines {
		// if mandatory placeholder exist in line, save line param
		if strings.Contains(line, MANDATORY_PLACEHOLDER) {
			mandatoryParams = append(mandatoryParams, line[:strings.Index(line, "=")])
		}
	}
	return mandatoryParams
}
