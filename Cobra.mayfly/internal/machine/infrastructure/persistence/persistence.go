package persistence

import "learn_zinx/Cobra.mayfly/internal/machine/domain/repository"

var (
	machineRepo       repository.Machine       = newMachineRepo()
	machineFileRepo   repository.MachineFile   = newMachineFileRepo()
	machineScriptRepo repository.MachineScript = newMachineScriptRepo()
)

func GetMachineRepo() repository.Machine {
	return machineRepo
}

func GetMachineFileRepo() repository.MachineFile {
	return machineFileRepo
}

func GetMachineScriptRepo() repository.MachineScript {
	return machineScriptRepo
}
