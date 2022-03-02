package policy

import (
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

type arguments []string

func (args arguments) Contain(values []string) bool {

	findRet := map[string]bool{}

	for _, value := range values {
		for _, a := range args {
			if a == value {
				findRet[value] = true
				break
			}
		}
	}

	if len(findRet) == len(values) {
		return true
	}

	return false
}

// Executor provides interfaces for testing policy.
type Executor interface {
	UpTerra() (bool, error)
}

type terraExecutor struct {
}

func (e *terraExecutor) UpTerra() (bool, error) {
	log.Println("UpTerra")
	return isProcessRunning("terrad", []string{"start"})
}

// func (e *maticExecutor) UpOracle() (bool, error) {
// 	log.Println("IsBorUp")
// 	return isProcessRunning("bor", []string{})
// }

func isProcessRunning(pName string, mustHaveArgs []string) (bool, error) {
	ps, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range ps {
		var (
			pName     string
			isRunning bool
			err       error
		)

		if pName, err = p.Name(); err != nil {
			continue
		}

		if pName != pName {
			continue
		}

		args, err := p.CmdlineSlice()
		if err != nil {
			continue
		}

		if !arguments(args).Contain(mustHaveArgs) {
			continue
		}

		if isRunning, err = p.IsRunning(); err != nil {
			return false, err
		}

		return isRunning, nil
	}

	return false, nil
}

// NewExecutor returns new executor.
func NewExecutor() Executor {
	return &terraExecutor{}
}
