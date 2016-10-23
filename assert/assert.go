package assert

import "log"

func Assert(assertion bool, message string) {
	if !assertion {
		log.Fatal(message)
	}
}
