# kubevirt-auto-uuid
Mutating kubernetes webhook to provide a UUID to all kubevirt VirtualMachines upon create.

This is a workaround for the [stable UUID generation](https://github.com/kubevirt/kubevirt/blob/afbeb269840e984b6c49d984ed1e53ab6eb02302/pkg/virt-controller/watch/vm.go#L1830) applied by default on VirtualMachine objects.
