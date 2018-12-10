package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetServiceFromClientEReturnsErrorForNonExistantService(t *testing.T) {
	t.Parallel()

	clientset, err := GetKubernetesClientE(t)
	require.NoError(t, err)
	_, err = GetServiceFromClientE(clientset, "default", "nginx-service")
	require.Error(t, err)
}

func TestGetServiceFromClientEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	service := GetService(t, uniqueID, "nginx-service")
	require.Equal(t, service.Name, "nginx-service")
	require.Equal(t, service.Namespace, uniqueID)
}

func TestWaitUntilServiceAvailableReturnsSuccessfullyOnNodePortType(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	WaitUntilServiceAvailable(t, uniqueID, "nginx-service", 10, 1*time.Second)
}

func TestGetServiceEndpointEReturnsAccessibleEndpointForNodePort(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	service := GetService(t, uniqueID, "nginx-service")
	endpoint := GetServiceEndpoint(t, service, 80)
	statusCode, _ := http_helper.HttpGet(t, fmt.Sprintf("http://%s", endpoint))
	require.Equal(t, statusCode, 200)
}

const EXAMPLE_DEPLOYMENT_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: %s
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.15.7
        ports:
        - containerPort: 80
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-service
  namespace: %s
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    targetPort: 80
    port: 80
  type: NodePort
`
