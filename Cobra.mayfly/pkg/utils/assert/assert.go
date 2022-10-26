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

func State(condition bool, panicMsg string, params ...interface{}) {
	IsTrue(condition, panicMsg, params...)
}

func NotEmpty(str string, panicMsg string, params ...interface{}) {
	IsTrue(str != "", panicMsg, params...)
}
