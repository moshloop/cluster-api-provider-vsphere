/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package govmomi

import (
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/cloud/vsphere/constants"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/cloud/vsphere/context"
)

// Update updates the machine from the backend platform.
func Update(ctx *context.MachineContext) error {
	if ctx.MachineConfig.MachineRef == "" {
		return errors.Errorf("machine ref is empty while updating machine %q", ctx)
	}

	moRef := types.ManagedObjectReference{
		Type:  "VirtualMachine",
		Value: ctx.MachineConfig.MachineRef,
	}

	var obj mo.VirtualMachine
	if err := ctx.Session.RetrieveOne(ctx, moRef, []string{"name", "runtime"}, &obj); err != nil {
		return errors.Errorf("machine does not exist %q", ctx)
	}

	if obj.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn {
		return errors.Errorf("machine is not running %q: %v", ctx, obj.Runtime.PowerState)
	}

	if ctx.IPAddr() != "" {
		return nil
	}

	ctx.Logger.V(4).Info("waiting on machine's IP address")
	vm := object.NewVirtualMachine(ctx.Session.Client.Client, moRef)
	ipAddr, err := vm.WaitForIP(ctx)
	if err != nil {
		return errors.Wrapf(err, "error waiting for machine's IP address %q", ctx)
	}

	if ctx.Machine.Annotations == nil {
		ctx.Machine.Annotations = map[string]string{}
	}
	ctx.Machine.Annotations[constants.VmIpAnnotationKey] = ipAddr
	ctx.Logger.V(2).Info("discovered machine's IP address", "ip-addr", ipAddr)

	return nil
}
