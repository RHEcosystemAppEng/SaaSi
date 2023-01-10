package main

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/packager"
)

func main() {
	// if len(os.Args) != 2 {
	// 	log.Fatal("Expected 1 argument, got ", len(os.Args)-1)
	// }

	// namespace := os.Args[1]
	// pretty.Printf("Deploying application to namespace %s", namespace)\

	pkgNs := "holdings"
	kustomizePath := "../install-builder/output/Infinity/installer/kustomize"

	pkg, err := packager.NewPkg(pkgNs)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pkg.GeneratePkgTemplate(kustomizePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pkg.InvokePkgCustomizations()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pkg.BuildPkg()
	if err != nil {
		fmt.Println(err)
		return
	}

	// DeployPkg()

}
