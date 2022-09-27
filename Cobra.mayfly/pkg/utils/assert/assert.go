package assert

import "fmt"

func IsTrue(condition bool, panicMsg string, params ...interface{}) {
	if !condition {
		if len(params) != 0 {
			panic(fmt.Sprintf(panicMsg, params...))
		}
		panic(panicMsg)
	}
}
