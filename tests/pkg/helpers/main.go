package helpers

import (
	"bytes"
	"fmt"
	"net"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
)

// HellowoldOpt type allows you to specify options for
// deploying a helloworld app template in a cluster.
type HelloworldOpt struct {
	Class         string `default:"nginx"`
	Identifier    string `default:"green-blue"`
	EnableModSec  string `default:"false"`
	Weight        string `default:"100"`
	ModSecSnippet string
	Hostname      string `example:"hostname.cloud-platform...."`
}

// CreateHelloWorldApp takes a HelloworldOpt type and KubectlOptions arguments
// to create a HelloWorld app in the environment of your choice.
func CreateHelloWorldApp(app *HelloworldOpt, opt *k8s.KubectlOptions) error {
	if app.Hostname == "" {
		return fmt.Errorf("helloworld app hostname must not be empty")
	}

	templateVars := map[string]interface{}{
		"ingress_annotations": map[string]string{
			"kubernetes.io/ingress.class":                     app.Class,
			"external-dns.alpha.kubernetes.io/aws-weight":     app.Weight,
			"external-dns.alpha.kubernetes.io/set-identifier": app.Identifier,
			"nginx.ingress.kubernetes.io/enable-modsecurity":  app.EnableModSec,
			"nginx.ingress.kubernetes.io/modsecurity-snippet": "|\n     SecRuleEngine On",
		},
		"host": app.Hostname,
	}

	tpl, err := TemplateFile("./fixtures/helloworld-deployment.yaml.tmpl", "outputTemplateDeployment.yaml.tmpl", templateVars)
	if err != nil {
		return fmt.Errorf("failed to create the helloworld template: %s", err)
	}

	err = k8s.KubectlApplyFromStringE(GinkgoT(), opt, tpl)
	if err != nil {
		return fmt.Errorf("failed to apply the helloworld template: %s", err)
	}

	return nil
}

// HttpStatusCode return the HTTP code for an endpoint
func HttpStatusCode(u string) (int, error) {
	client := resty.New()
	resp, err := client.R().EnableTrace().Get(u)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode(), nil
}

// DNSLookUp returns error if there is not DNS entry for an endpoint, used
// with retry library from terratest
func DNSLookUp(h string) (string, error) {
	if _, err := net.LookupIP(h); err != nil {
		return "", fmt.Errorf("DNS propagation hasn't happened'%w'", err)
	}

	return "", nil
}

// TemplateFile returns a string with the content of a template rendered
func TemplateFile(f string, n string, m interface{}) (string, error) {
	var b bytes.Buffer

	t, err := template.ParseFiles(f)
	if err != nil {
		return "", err
	}

	err = t.ExecuteTemplate(&b, n, m)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
