package client // import "github.com/docker/docker/client"

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// CheckpointCreate creates a checkpoint from the given container with the given name
func (cli *Client) CheckpointCreate(ctx context.Context, container string, options types.CheckpointCreateOptions) error {
	fmt.Println("Invoke Client.CheckpointCreate, which is to invoke containerd")
	resp, err := cli.post(ctx, "/containers/"+container+"/checkpoints", nil, options, nil)
	ensureReaderClosed(resp)
	return err
}
