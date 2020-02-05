package run

import (
	"errors"
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"io"
	"os"
	"text/template"
)

// InitKube is a step in a helm Plan that initializes the kubernetes config file.
type InitKube struct {
	*config
	templateFilename string
	configFilename   string
	template         *template.Template
	configFile       io.WriteCloser
	values           kubeValues
}

type kubeValues struct {
	SkipTLSVerify     bool
	Certificate       string
	APIServer         string
	Namespace         string
	ServiceAccount    string
	Token             string
	ClientCertificate string
	ClientKey         string
}

// NewInitKube creates a InitKube using the given Config and filepaths. No validation is performed at this time.
func NewInitKube(cfg env.Config, templateFile, configFile string) *InitKube {
	return &InitKube{
		config: newConfig(cfg),
		values: kubeValues{
			SkipTLSVerify:     cfg.SkipTLSVerify,
			Certificate:       cfg.Certificate,
			APIServer:         cfg.APIServer,
			Namespace:         cfg.Namespace,
			ServiceAccount:    cfg.ServiceAccount,
			Token:             cfg.KubeToken,
			ClientCertificate: cfg.ClientCertificate,
			ClientKey:         cfg.ClientKey,
		},
		templateFilename: templateFile,
		configFilename:   configFile,
	}
}

// Execute generates a kubernetes config file from drone-helm3's template.
func (i *InitKube) Execute() error {
	if i.debug {
		fmt.Fprintf(i.stderr, "writing kubeconfig file to %s\n", i.configFilename)
	}
	defer i.configFile.Close()
	return i.template.Execute(i.configFile, i.values)
}

// Prepare ensures all required configuration is present and that the config file is writable.
func (i *InitKube) Prepare() error {
	var err error

	if i.values.APIServer == "" {
		return errors.New("an API Server is needed to deploy")
	}
	if i.values.Token == "" && (i.values.ClientCertificate == "" && i.values.ClientKey == "") {
		return errors.New("token or client certs are needed to deploy")
	} else if i.values.ClientCertificate == "" || i.values.ClientKey == "" {
		return errors.New("both client certificate and key are needed to deploy")
	}

	if i.values.ServiceAccount == "" {
		i.values.ServiceAccount = "helm"
	}

	if i.debug {
		fmt.Fprintf(i.stderr, "loading kubeconfig template from %s\n", i.templateFilename)
	}
	i.template, err = template.ParseFiles(i.templateFilename)
	if err != nil {
		return fmt.Errorf("could not load kubeconfig template: %w", err)
	}

	if i.debug {
		if _, err := os.Stat(i.configFilename); err != nil {
			// non-nil err here isn't an actual error state; the kubeconfig just doesn't exist
			fmt.Fprint(i.stderr, "creating ")
		} else {
			fmt.Fprint(i.stderr, "truncating ")
		}
		fmt.Fprintf(i.stderr, "kubeconfig file at %s\n", i.configFilename)
	}

	i.configFile, err = os.Create(i.configFilename)
	if err != nil {
		return fmt.Errorf("could not open kubeconfig file for writing: %w", err)
	}
	return nil
}
