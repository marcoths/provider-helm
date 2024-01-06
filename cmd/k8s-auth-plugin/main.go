package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/crossplane-contrib/provider-helm/pkg/auth"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New(filepath.Base(os.Args[0]), "generate k8s auth tokens").DefaultEnvars()
	gke = app.Command("gke", "Generate GKE auth token")
)

func main() {
	ctx := context.Background()
	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case gke.FullCommand():
		err = auth.Gcp(ctx)
	}
	kingpin.FatalIfError(err, "Cannot generate k8s auth token")
}
