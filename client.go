package gk8s

import (
	"encoding/base64"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	RestConfig *rest.Config
	APIClient  kubernetes.Interface
}

// FromRestConfig - Init K8s Client with rest configuration
func FromRestConfig(config *rest.Config) (*Config, error) {
	apiClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Config{
		RestConfig: config,
		APIClient:  apiClient,
	}, nil
}

// InCluster - Init K8s Client with in cluster credentials
func InCluster() (*Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return FromRestConfig(config)
}

// FromKubeConfig - Init K8s Client with kubeconfig file $USER/.kube/config
func FromKubeConfig() (*Config, error) {
	home := homedir.HomeDir()
	path := filepath.Join(home, ".kube", "config")
	return FromKubeConfigPath(path)
}

// FromKubeConfigPath - Init K8s Client with specific kubeconfig file
func FromKubeConfigPath(path string) (*Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}
	return FromRestConfig(config)
}

// FromEKS - Init K8s Client with AWS shared-credentials-file and cluster name
func FromEKS(clusterName string) (*Config, error) {
	sess, err := getSession()
	if err != nil {
		return nil, err
	}
	config, err := getEKSConfig(sess, clusterName)
	if err != nil {
		return nil, err
	}
	return FromRestConfig(config)
}

func getSession() (*session.Session, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func getEKSCluster(sess *session.Session, clusterName string) (*eks.Cluster, error) {
	svc := eks.New(sess)
	input := eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}
	output, err := svc.DescribeCluster(&input)
	if err != nil {
		return nil, err
	}
	return output.Cluster, nil
}

func getEKSConfig(sess *session.Session, clusterName string) (*rest.Config, error) {
	cluster, err := getEKSCluster(sess, clusterName)
	if err != nil {
		return nil, err
	}
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}
	opts := &token.GetTokenOptions{
		Session:   sess,
		ClusterID: aws.StringValue(cluster.Name),
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}
	return &rest.Config{
		Host:        aws.StringValue(cluster.Endpoint),
		BearerToken: tok.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: ca,
		},
	}, nil
}
