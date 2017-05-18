package main

const (
	debugNo       = 0
	debugMinimal  = 1
	debugNormal   = 2
	debugDetailed = 3
)

var debugLev int

func shouldDebug(myDebugLevel int) bool {
	return myDebugLevel <= debugLev
}
