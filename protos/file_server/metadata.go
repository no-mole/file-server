package file_server

import (
	"github.com/no-mole/neptune/registry"
)

func Metadata() *registry.Metadata {
	return &registry.Metadata{
		ServiceDesc: FileServerService_ServiceDesc,
		Namespace:   "biomind",
		Version:     "v1",
	}
}
