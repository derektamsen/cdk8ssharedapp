package cdk8ssharedapp

import (
	"fmt"
	"os"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/derektamsen/cdk8ssharedapp/imports/k8s"
)

type AppConfig struct {
	Name      string
	Namespace string
}

type ClusterProps struct {
	ClusterName string
	Image       string
}

type K8sClusters struct {
	Clusters *[]ClusterProps
}

func NewApp(appConfig *AppConfig, clusters *K8sClusters) error {
	// Remove any existing rendered manifests from dist
	// This ensures that the rendered output is exactly what we expect.
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	rendered_dir := fmt.Sprintf("%s/dist", dir)
	fmt.Printf("Removing rendered manifests from \"%s\"\n", rendered_dir)
	err = os.RemoveAll(rendered_dir)
	if err != nil {
		return err
	}

	// Generate the manifests for all clusters
	for _, v := range *clusters.Clusters {
		fmt.Printf("Generating manifests for %s\n", v.ClusterName)
		appProps := &cdk8s.AppProps{
			Outdir:              jsii.String(fmt.Sprintf("dist/%s", v.ClusterName)),
			OutputFileExtension: jsii.String(".yaml"),
			YamlOutputType:      cdk8s.YamlOutputType_FOLDER_PER_CHART_FILE_PER_RESOURCE,
		}

		app := cdk8s.NewApp(appProps)

		NewChart(app, appConfig.Name, appConfig.Namespace, appConfig.Name, v.Image)

		app.Synth()
	}

	return nil
}

func NewChart(scope constructs.Construct, id string, ns string, appLabel string, image string) cdk8s.Chart {

	chart := cdk8s.NewChart(scope, jsii.String(id), &cdk8s.ChartProps{
		Namespace: jsii.String(ns),
	})

	labels := map[string]*string{
		"app": jsii.String(appLabel),
	}

	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(3),
			Selector: &k8s.LabelSelector{
				MatchLabels: &labels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &labels,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name:  jsii.String("app-container"),
						Image: jsii.String(image),
						Ports: &[]*k8s.ContainerPort{{
							ContainerPort: jsii.Number(80),
						}},
					}},
					ServiceAccountName: jsii.String("service-account"),
				},
			},
		},
	})

	k8s.NewKubeServiceAccount(chart, jsii.String("service-account"), nil)

	return chart
}
