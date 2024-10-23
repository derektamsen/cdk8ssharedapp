package main

import (
	"example.com/cdk8s/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// type ClusterProps struct {
// 	ClusterName string
// }

// type K8sClusters struct {
// 	Clusters map[string]*ClusterProps
// }

func NewApp() error {
	appProps := &cdk8s.AppProps{
		Outdir:              jsii.String("dist/cluster"),
		OutputFileExtension: jsii.String(".yaml"),
		YamlOutputType:      cdk8s.YamlOutputType_FOLDER_PER_CHART_FILE_PER_RESOURCE,
	}

	app := cdk8s.NewApp(appProps)

	NewChart(app, "exemplar", "my-app", "my-app")

	app.Synth()

	return nil
}

func NewChart(scope constructs.Construct, id string, ns string, appLabel string) cdk8s.Chart {

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
						Image: jsii.String("nginx:1.19.10"),
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

func main() {
	NewApp()
}
