package comm

import (
	"autoDeploy/mySSH"
)

func CommTest() {
	prefix := "data"
	CreateTestDirectory(prefix)
	//autoDeployment()
	TestViper()
	TestCopy()
	TestCSVFile()
	mySSH.TestOutput()
}
