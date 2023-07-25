package entities

import "github.com/idanya/evm-cli/clients/nodes"

type NodeClientGenerator = func() nodes.NodeClient
