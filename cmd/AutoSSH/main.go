package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "autoDeploy/mySSH"
)

func main() {
	/* variables */
	var operationLines []string
	var nodeOperations []NodeOperationItem
	var lines []string
	var nodes []NodeItem
	argsWithProg := os.Args
	operationFile := "operation.csv"
	if len(argsWithProg) > 1 {
		operationFile = argsWithProg[1]
		fmt.Printf(operationFile + "\n")
	}
	/* download data files */
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0700)
	}
	/* log file */
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0700)
	}
	dt := time.Now().Format("2006-01-02-15_04_05")
	fmt.Println(dt)
	logFile, err := os.Create("log/log" + dt + ".log")
	if err != nil {
		panic("Initialize log file failed.")
	}
	defer logFile.Close()
	/* domain file */
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	localNodeDomainFilePath := local + "/nodeDomain.txt"

	nodeDomainFile, err := os.Open(localNodeDomainFilePath)

	if err != nil {
		fmt.Println("There is no node domain file. Please recheck.")
		panic(err)
	}
	defer nodeDomainFile.Close()

	scanner := bufio.NewScanner(nodeDomainFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	nodeDomainFile.Close()

	localOperationFilePath := local + "/" + operationFile
	nodeOperationFile, err := os.Open(localOperationFilePath)

	if err != nil {
		fmt.Println("There is no operation file. Please recheck")
		panic(err)
	}
	defer nodeOperationFile.Close()

	scanner = bufio.NewScanner(nodeOperationFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		operationLines = append(operationLines, scanner.Text())
	}

	for i, v := range lines {
		nodeString := strings.Split(v, ",")
		if len(nodeString) == 5 {
			var currNodeItem NodeItem
			currNodeItem.NodeIndex = strings.Join(strings.Fields(nodeString[0]), "")
			currNodeItem.IPaddress = strings.Join(strings.Fields(nodeString[1]), "")
			currNodeItem.UserName = strings.Join(strings.Fields(nodeString[2]), "")
			currNodeItem.Password = strings.Join(strings.Fields(nodeString[3]), "")
			//currNodeItem.AbsolutePath = strings.Join(strings.Fields(nodeString[4]), "")
			currNodeItem.AbsolutePath = GetRealNameFromPattern(strings.Join(strings.Fields(nodeString[4]), ""), currNodeItem.NodeIndex)

			nodes = append(nodes, currNodeItem)
		} else {
			fmt.Printf("wrong input %d, : %s\n", i, v)
		}
	}
	if len(nodes) == 0 {
		panic("no suitable node domain files")
	}

	for i, v := range operationLines {
		operationString := strings.SplitN(v, ",", 2)
		if len(operationString) == 2 {
			var currOperationItem NodeOperationItem
			currOperationItem.OperationName = strings.TrimSpace(operationString[0])
			currOperationItem.OperationContent = strings.TrimSpace(operationString[1])
			nodeOperations = append(nodeOperations, currOperationItem)
		} else {
			fmt.Printf("wrong input %d, : %s\n", i, v)
			panic("no suitable operations, stop the program")
		}
	}
	if len(nodes) == 0 {
		panic("incorrect operation files, stop the program")
	} else {
		for i, v := range nodeOperations {
			if !StringInSlice(v.OperationName, LegalOperationName) {
				fmt.Printf("incorrect operation names: %s in line %d\n", v.OperationName, i)
				panic("opps!")
			}
			if v.OperationName == "copy" {
				localFile := local + v.OperationContent
				if _, err := os.Stat(localFile); os.IsNotExist(err) {
					fmt.Printf("No file in current path: %s in line %d\n", v.OperationContent, i)
					panic("opps!")
				}
			}

			if v.OperationName == "copyN" {
				indexPos := strings.Index(v.OperationContent, "%")
				if indexPos > -1 {
					preName := v.OperationContent[:indexPos]
					afterName := v.OperationContent[indexPos+1:]
					refinedName := preName + afterName
					nodeOperations[i].OperationRefinedContent = refinedName
				}
				for _, currNode := range nodes {
					realFile := GetRealNameFromPattern(v.OperationContent, currNode.NodeIndex)
					localFile := local + realFile
					if _, err := os.Stat(localFile); os.IsNotExist(err) {
						fmt.Printf("No file in current path: %s in line %d\n", localFile, i)
						panic("opps!")
					}
				}

			}

			if v.OperationName == "getN" {
				indexPos := strings.Index(v.OperationContent, "%")
				if indexPos > -1 {
					preName := v.OperationContent[:indexPos]
					afterName := v.OperationContent[indexPos+1:]
					nodeOperations[i].OperationPrefix = preName
					nodeOperations[i].OperationPostfix = afterName
					nodeOperations[i].IsSubstitute = true
				} else {
					nodeOperations[i].OperationPrefix = ""
					nodeOperations[i].OperationPostfix = v.OperationContent
					nodeOperations[i].IsSubstitute = false
				}
			}
		}
	}

	for _, v := range nodes {
		fmt.Printf("===============Implement node: %s, node index: %s =============\n", v.IPaddress, v.NodeIndex)
		fmt.Fprintf(logFile, "===============Implement node: %s, node index: %s =============\n", v.IPaddress, v.NodeIndex)
		var testOperation bool = true

		for _, opeItem := range nodeOperations {
			if !testOperation {
				break
			}
			fmt.Printf("current operation: %s; with detail: %s\n", opeItem.OperationName, opeItem.OperationContent)
			fmt.Fprintf(logFile, "current operation: %s; with detail: %s\n", opeItem.OperationName, opeItem.OperationContent)
			switch opeItem.OperationName {
			case "copy":
				{
					localFile := local + opeItem.OperationContent
					err := TransferFile(v, localFile, "", logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			case "copyN":
				{
					realFile := GetRealNameFromPattern(opeItem.OperationContent, v.NodeIndex)
					localFile := local + realFile
					destName := filepath.Base(opeItem.OperationRefinedContent)
					err := TransferFile(v, localFile, destName, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			case "command":
				{
					err := DirectImplement(v, opeItem.OperationContent, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s\n", v.IPaddress, opeItem.OperationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s\n", v.IPaddress, opeItem.OperationContent)
					} else {
						testOperation = false
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}
				}
			case "getN":
				{
					targetPathName := ""
					if !opeItem.IsSubstitute {
						targetPathName = opeItem.OperationPrefix + opeItem.OperationPostfix
					} else {
						targetPathName = opeItem.OperationPrefix + v.NodeIndex + opeItem.OperationPostfix
					}
					localCopyedFileName := local + "/data" + "/r" + v.NodeIndex + filepath.Base(targetPathName)
					fmt.Printf("targetFile: %s\n localFileName: %s\n", targetPathName, localCopyedFileName)
					err := DownloadFile(v, targetPathName, localCopyedFileName, logFile)
					if err == nil {
						fmt.Printf("Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
						fmt.Fprintf(logFile, "Success in node: %s with operation: %s : %s\n", v.IPaddress, opeItem.OperationName, opeItem.OperationContent)
					} else {
						fmt.Printf("\n")
						fmt.Fprintf(logFile, "\n")
					}

				}
			}
		}
		fmt.Printf("\n")
		fmt.Fprintf(logFile, "\n")
	}

}
