package pkg

import (
	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	meshconfig "github.com/kubeedge/kubeedge/edgemesh/pkg/config"
	"github.com/kubeedge/kubeedge/edgemesh/pkg/constant"
	"github.com/kubeedge/kubeedge/edgemesh/pkg/proxy"
	"github.com/kubeedge/kubeedge/edgemesh/pkg/server"
)

//EdgeMesh defines EdgeMesh object structure
type EdgeMesh struct {
}

// Register register edgemesh
func Register() {
	meshconfig.InitConfigure()
	core.Register(&EdgeMesh{})
}

//Name returns the name of EdgeMesh module
func (em *EdgeMesh) Name() string {
	return constant.ModuleNameEdgeMesh
}

//Group returns EdgeMesh group
func (em *EdgeMesh) Group() string {
	return modules.MeshGroup
}

//Start sets context and starts the controller
func (em *EdgeMesh) Start() {
	proxy.Init()
	go server.Start()
	// we need watch message to update the cache of instances
	for {
		select {
		case <-beehiveContext.Done():
			klog.Warning("EdgeMesh Stop")
			return
		default:
		}
		msg, err := beehiveContext.Receive(constant.ModuleNameEdgeMesh)
		if err != nil {
			klog.Warningf("edgemesh receive msg error %v", err)
			continue
		}
		klog.V(4).Infof("edgemesh get message: %v", msg)
		proxy.MsgProcess(msg)
	}
}
