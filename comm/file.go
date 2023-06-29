package comm

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateTestDirectory(directoryName string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	if _, err := os.Stat(path + "/" + directoryName); os.IsNotExist(err) {
		if err := os.Mkdir(path+"/"+directoryName, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Print("fuck")
	}

}

func CreateAbsoluteDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// CSVFileToMap  reads csv file into slice of map
// slice is the line number
// map[string]string where key is column name
func CSVFileToMap(filePath string) (returnMap []map[string]string, err error) {

	// read csv file
	csvfile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	header := []string{} // holds first row (header)
	for lineNum, record := range rawCSVdata {

		// for first row, build the header slice
		if lineNum == 0 {
			for i := 0; i < len(record); i++ {
				header = append(header, strings.TrimSpace(record[i]))
			}

		} else {
			// for each cell, map[string]string k=header v=value
			line := map[string]string{}
			for i := 0; i < len(record); i++ {
				line[header[i]] = record[i]
			}
			returnMap = append(returnMap, line)
		}
	}

	return
}

// MapToCSVFile  writes slice of map into csv file
// filterFields filters to only the fields in the slice, and maintains order when writing to file
func MapToCSVFile(inputSliceMap []map[string]string, filePath string, filterFields []string) (err error) {

	var headers []string  // slice of each header field
	var line []string     // slice of each line field
	var csvLine string    // string of line converted to csv
	var CSVContent string // final output of csv containing header and lines

	// iter over slice to get all possible keys (csv header) in the maps
	// using empty Map[string]struct{} to get UNIQUE Keys; no value needed
	var headerMap = make(map[string]struct{})
	for _, record := range inputSliceMap {
		for k := range record {
			headerMap[k] = struct{}{}
		}
	}

	// convert unique headersMap to slice
	for headerValue := range headerMap {
		headers = append(headers, headerValue)
	}

	// filter to filteredFields and maintain order
	var filteredHeaders []string
	if len(filterFields) > 0 {
		for _, filterField := range filterFields {
			for _, headerValue := range headers {
				if filterField == headerValue {
					filteredHeaders = append(filteredHeaders, headerValue)
				}
			}
		}
	} else {
		filteredHeaders = append(filteredHeaders, headers...)
		sort.Strings(filteredHeaders) // alpha sort headers
	}

	// write headers as the first line
	csvLine, _ = WriteAsCSV(filteredHeaders)
	CSVContent += csvLine + "\n"

	// iter over inputSliceMap to get values for each map
	// maintain order provided in header slice
	// write to csv
	for _, record := range inputSliceMap {
		line = []string{}

		// lines
		for k := range filteredHeaders {
			line = append(line, record[filteredHeaders[k]])
		}
		csvLine, _ = WriteAsCSV(line)
		CSVContent += csvLine + "\n"
	}

	// make the dir incase it's not there
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// write out the csv contents to file
	ioutil.WriteFile(filePath, []byte(CSVContent), os.FileMode(0644))
	if err != nil {
		return err
	}

	return
}

func WriteAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func TestCopy() {
	Copy("./config/fuck.txt", "../fuck1.txt")
}

func GetConfiguration(configuration *Config) {
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

func TestViper() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("../config/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model

	fmt.Println("Reading variables using the model..")
	for _, radarItem := range configuration.RadarTypeVec {
		fmt.Println("type: ", radarItem.TypeNum)
		fmt.Println("isTunnel: ", radarItem.IsTunnel)
		fmt.Println("incoming lane number: ", radarItem.IncomingLaneNum)
	}
	fmt.Println("projectNum: ", configuration.Project.ProjectNum)
	fmt.Println("extraPort: ", configuration.Server.ExtraPort)
	fmt.Println("localPath: ", configuration.Server.LocalImplementPath)
}

func TestCSVFile() {
	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		fmt.Print(err.Error())
	}
	for _, item := range result {
		var node NodeConfig
		node.UserName = item["user_name"]
		node.StakeMark = item["stake_mark"]
		node.Password = item["password"]
		node.IpAddress = item["ip"]
		node.DeviceID, err = strconv.Atoi(item["device_id"])
		Check(err)
		node.Can0Type, err = strconv.Atoi(item["can0_type"])
		Check(err)
		node.Can0ChessboardFile = item["can0_file"]
		node.Can1Type, err = strconv.Atoi(item["can1_type"])
		Check(err)
		node.Can1ChessboardFile = item["can1_file"]
		node.Can2Type, err = strconv.Atoi(item["can2_type"])
		Check(err)
		node.Can2ChessboardFile = item["can2_file"]
		node.Can3Type, err = strconv.Atoi(item["can3_type"])
		Check(err)
		node.Can3ChessboardFile = item["can3_file"]
		var path string = "./test/" + strconv.Itoa(node.DeviceID) + "vec.config"
		CreateRadarConfigFile(path, node, configuration.RadarTypeVec, configuration.Server, configuration.Project)
	}
}

func Split(r rune) bool {
	return r == 'k' || r == 'K' || r == '+'
}

func TransformStakeMarkToDistance(stakeMark string) (distance float64) {
	distance = 0.0
	stakeMark = regexp.MustCompile(`[^0-9 +.]+`).ReplaceAllString(stakeMark, "")

	sVec := strings.Split(stakeMark, "+")
	var flag bool = true
	var km, hm float64

	km, err := strconv.ParseFloat(strings.TrimSpace(sVec[0]), 32)
	if err != nil {
		flag = false
	}
	hm, err = strconv.ParseFloat(strings.TrimSpace(sVec[1]), 32)
	if err != nil {
		flag = false
	}
	if flag {
		distance = km*1000 + hm
	} else {
		distance = 0
	}
	return
}

func GenerateRadarPosFromNode(node NodeConfig, radarTypes []RadarType, server ServerConfigurations, project ProjectConfiguration) (radarPosConfigs []RadarPosConfig) {
	var radarNum = 1
	var currUpRiver = false
	if node.DeviceID <= 64 {
		currUpRiver = true
	} else if node.DeviceID <= 128 {
		currUpRiver = false
	} else {
		currUpRiver = false
	}
	if node.Can0Type > 0 {
		var configItem RadarPosConfig
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can0Type {
				radarTypeTemp = radarItem
				break
			}
		}
		if radarTypeTemp.RadarDirection == 1 {
			configItem.Angle = 0
		} else {
			configItem.Angle = math.Pi
		}
		configItem.Comment = project.ProjectName
		configItem.DenyLaneChange = false
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		configItem.RadarID = strconv.Itoa(radarID)
		configItem.Position.X = -TransformStakeMarkToDistance(project.ProjectStartStakeMark) + TransformStakeMarkToDistance(node.StakeMark)
		configItem.Position.Y = 0
		configItem.Position.Z = 0
		if radarTypeTemp.TypeNum == 3 {
			configItem.IsZH2HK = true
		} else {
			configItem.IsZH2HK = currUpRiver
		}
		radarPosConfigs = append(radarPosConfigs, configItem)
		radarNum = radarNum + 1
	}
	if node.Can1Type > 0 {
		var configItem RadarPosConfig
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can1Type {
				radarTypeTemp = radarItem
				break
			}
		}
		if radarTypeTemp.RadarDirection == 1 {
			configItem.Angle = 0
		} else {
			configItem.Angle = math.Pi
		}
		configItem.Comment = project.ProjectName
		configItem.DenyLaneChange = false
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		configItem.RadarID = strconv.Itoa(radarID)
		configItem.Position.X = -TransformStakeMarkToDistance(project.ProjectStartStakeMark) + TransformStakeMarkToDistance(node.StakeMark)
		configItem.Position.Y = 0
		configItem.Position.Z = 0
		if radarTypeTemp.TypeNum == 3 {
			configItem.IsZH2HK = true
		} else {
			configItem.IsZH2HK = currUpRiver
		}
		radarPosConfigs = append(radarPosConfigs, configItem)
		radarNum = radarNum + 1
	}
	if node.Can2Type > 0 {
		var configItem RadarPosConfig
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can2Type {
				radarTypeTemp = radarItem
				break
			}
		}
		if radarTypeTemp.RadarDirection == 1 {
			configItem.Angle = 0
		} else {
			configItem.Angle = math.Pi
		}
		configItem.Comment = project.ProjectName
		configItem.DenyLaneChange = false
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		configItem.RadarID = strconv.Itoa(radarID)
		configItem.Position.X = -TransformStakeMarkToDistance(project.ProjectStartStakeMark) + TransformStakeMarkToDistance(node.StakeMark)
		configItem.Position.Y = 0
		configItem.Position.Z = 0
		if radarTypeTemp.TypeNum == 3 {
			configItem.IsZH2HK = true
		} else {
			configItem.IsZH2HK = currUpRiver
		}
		radarPosConfigs = append(radarPosConfigs, configItem)
		radarNum = radarNum + 1
	}
	if node.Can3Type > 0 {
		var configItem RadarPosConfig
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can3Type {
				radarTypeTemp = radarItem
				break
			}
		}
		if radarTypeTemp.RadarDirection == 1 {
			configItem.Angle = 0
		} else {
			configItem.Angle = math.Pi
		}
		configItem.Comment = project.ProjectName
		configItem.DenyLaneChange = false
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		configItem.RadarID = strconv.Itoa(radarID)
		configItem.Position.X = -TransformStakeMarkToDistance(project.ProjectStartStakeMark) + TransformStakeMarkToDistance(node.StakeMark)
		configItem.Position.Y = 0
		configItem.Position.Z = 0
		if radarTypeTemp.TypeNum == 3 {
			configItem.IsZH2HK = true
		} else {
			configItem.IsZH2HK = currUpRiver
		}
		radarPosConfigs = append(radarPosConfigs, configItem)
	}
	return
}

func CreateRadarConfigFile(path string, node NodeConfig, radarTypes []RadarType, server ServerConfigurations, project ProjectConfiguration) {
	f, err := os.Create(path)
	Check(err)
	defer f.Close()
	f.WriteString("remoteIpAddress=" + server.IPAddress + "\n")
	f.WriteString("remotePortNum=" + strconv.Itoa(server.Port) + "\n")
	f.WriteString("remotePortExtra=" + strconv.Itoa(server.ExtraPort) + "\n")
	f.WriteString("sectionMinSpeed=10\n")
	f.WriteString("sectionMaxSpeed=31\n")
	f.WriteString("rangeXMax=260\n")
	f.WriteString("samplePoolNum=6000\n")

	var deviceNum = 0
	var radarNum = 1
	if node.Can0Type > 0 {
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can0Type {
				radarTypeTemp = radarItem
				break
			}
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "isTunnel=" + strconv.Itoa(radarTypeTemp.IsTunnel) + "\n")
		if radarTypeTemp.IsTunnel == 1 {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=50" + "\n")
		} else {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=130" + "\n")
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "Range=260" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanName=can0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "USBIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "RadarModuleName=40821.module" + "\n")
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		f.WriteString("device" + strconv.Itoa(deviceNum) + "GlobalRadarID=" + strconv.Itoa(radarID) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IncomingLaneNum=" + strconv.Itoa(radarTypeTemp.IncomingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "OutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.OutgoingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartOutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.StartOutgoing) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartIncomingLaneNum=" + strconv.Itoa(radarTypeTemp.StartIncoming) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IsDriveRight=1" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "ChessboardFileName=" + node.Can0ChessboardFile + "\n")

		deviceNum = deviceNum + 1
		radarNum = radarNum + 1
	}
	if node.Can1Type > 0 {
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can1Type {
				radarTypeTemp = radarItem
				break
			}
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "isTunnel=" + strconv.Itoa(radarTypeTemp.IsTunnel) + "\n")
		if radarTypeTemp.IsTunnel == 1 {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=50" + "\n")
		} else {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=130" + "\n")
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "Range=260" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanName=can1" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "USBIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "RadarModuleName=40821.module" + "\n")
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		f.WriteString("device" + strconv.Itoa(deviceNum) + "GlobalRadarID=" + strconv.Itoa(radarID) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IncomingLaneNum=" + strconv.Itoa(radarTypeTemp.IncomingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "OutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.OutgoingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartOutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.StartOutgoing) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartIncomingLaneNum=" + strconv.Itoa(radarTypeTemp.StartIncoming) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IsDriveRight=1" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "ChessboardFileName=" + node.Can1ChessboardFile + "\n")
		deviceNum = deviceNum + 1
		radarNum = radarNum + 1
	}
	if node.Can2Type > 0 {
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can2Type {
				radarTypeTemp = radarItem
				break
			}
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "isTunnel=" + strconv.Itoa(radarTypeTemp.IsTunnel) + "\n")
		if radarTypeTemp.IsTunnel == 1 {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=50" + "\n")
		} else {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=130" + "\n")
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "Range=260" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanName=can2" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "USBIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "RadarModuleName=40821.module" + "\n")
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		f.WriteString("device" + strconv.Itoa(deviceNum) + "GlobalRadarID=" + strconv.Itoa(radarID) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IncomingLaneNum=" + strconv.Itoa(radarTypeTemp.IncomingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "OutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.OutgoingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartOutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.StartOutgoing) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartIncomingLaneNum=" + strconv.Itoa(radarTypeTemp.StartIncoming) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IsDriveRight=1" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "ChessboardFileName=" + node.Can2ChessboardFile + "\n")
		deviceNum = deviceNum + 1
		radarNum = radarNum + 1
	}
	if node.Can3Type > 0 {
		var radarTypeTemp RadarType
		for _, radarItem := range radarTypes {
			if radarItem.TypeNum == node.Can3Type {
				radarTypeTemp = radarItem
				break
			}
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "isTunnel=" + strconv.Itoa(radarTypeTemp.IsTunnel) + "\n")
		if radarTypeTemp.IsTunnel == 1 {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=50" + "\n")
		} else {
			f.WriteString("device" + strconv.Itoa(deviceNum) + "KeyCenterLong=130" + "\n")
		}
		f.WriteString("device" + strconv.Itoa(deviceNum) + "Range=260" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanName=can3" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "USBIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "CanIndex=0" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "RadarModuleName=40821.module" + "\n")
		var radarID int = 100000000 + project.ProjectNum*100000 + node.DeviceID*100 + radarNum
		f.WriteString("device" + strconv.Itoa(deviceNum) + "GlobalRadarID=" + strconv.Itoa(radarID) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IncomingLaneNum=" + strconv.Itoa(radarTypeTemp.IncomingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "OutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.OutgoingLaneNum) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartOutgoingLaneNum=" + strconv.Itoa(radarTypeTemp.StartOutgoing) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "StartIncomingLaneNum=" + strconv.Itoa(radarTypeTemp.StartIncoming) + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "IsDriveRight=1" + "\n")
		f.WriteString("device" + strconv.Itoa(deviceNum) + "ChessboardFileName=" + node.Can3ChessboardFile + "\n")
	}
}
