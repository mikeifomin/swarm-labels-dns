package docker

import (
	"context"
	//	"fmt"
	"github.com/docker/docker/api/types"
	dockerApi "github.com/docker/docker/client"
	"strings"
)

type Docker struct {
	clientApi *dockerApi.Client
	prefix    string
}

func NewFromEnv() Docker {

	clientApi, err := dockerApi.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return Docker{clientApi: clientApi, prefix: "domain_"}
}

func (d *Docker) FetchNodesAddrs() map[string]string {
	result := make(map[string]string)
	nodes, err := d.clientApi.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		result[node.Status.Addr] = node.ID
	}
	return result
}
func (d *Docker) clearLabels(labels *map[string]string) {
	for key, _ := range *labels {
		if strings.HasPrefix(key, d.prefix) {
			delete(*labels, key)
		}
	}
}
func (d *Docker) UpdateNodeLabels(nodeId string, domains []string) error {

	node, _, err := d.clientApi.NodeInspectWithRaw(context.Background(), nodeId)
	if err != nil {
		return err
	}
	if node.Spec.Annotations.Labels == nil {
		node.Spec.Annotations.Labels = make(map[string]string)
	} else {
		d.clearLabels(&node.Spec.Annotations.Labels)
	}
	for _, domain := range domains {
		node.Spec.Annotations.Labels[d.prefix+domain] = "true"
	}
	errUpd := d.clientApi.NodeUpdate(context.Background(), nodeId, node.Meta.Version, node.Spec)
	return errUpd
}
