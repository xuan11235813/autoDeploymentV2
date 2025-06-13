package main

import (
	. "autoDeploy/comm"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "canRadarServerLocal":
			UpdateRadarConfig()
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			CopyFiles(path + "/local/" + "outMultiple.lexe")
		case "netRadarServerLocal":
			updateRadarConfigNet()
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			if runtime.GOOS == "linux" {
				var paths []string
				paths = append(paths, path+"/local/"+"netMultiple.lexe")
				CopyNetFiles(paths)
			}
			if runtime.GOOS == "windows" {
				var paths []string
				paths = append(paths, path+"/local/"+"netMultiple.exe")
				paths = append(paths, path+"/local/"+"libgcc_s_dw2-1.dll")
				paths = append(paths, path+"/local/"+"libatomic-1.dll")
				paths = append(paths, path+"/local/"+"libgomp-1.dll")
				paths = append(paths, path+"/local/"+"libquadmath-0.dll")
				paths = append(paths, path+"/local/"+"libssp-0.dll")
				paths = append(paths, path+"/local/"+"libstdc++-6.dll")
				paths = append(paths, path+"/local/"+"libwinpthread-1.dll")
				CopyNetFiles(paths)
			}
		case "netPort":
			UpdateRadarConfigNetPortTransfer()
			updateRadarConfigNetRemote()
		case "onlyMultiple":
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			CopyFiles(path + "/local/" + "outMultiple.lexe")
		}

	}

}

func updateRadarConfigNet() {
	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
	}

	for _, item := range result {
		var node NodeConfig
		node.StakeMark = item["stake_mark"]

		node.Net0Ip = item["radar0_ip"]
		node.Net1Ip = item["radar1_ip"]
		node.Net2Ip = item["radar2_ip"]
		node.Net3Ip = item["radar3_ip"]
		node.Net0PortIn = item["radar0_port_in"]
		node.Net1PortIn = item["radar1_port_in"]
		node.Net2PortIn = item["radar2_port_in"]
		node.Net3PortIn = item["radar3_port_in"]
		node.Net0PortOut = item["radar0_port_out"]
		node.Net1PortOut = item["radar1_port_out"]
		node.Net2PortOut = item["radar2_port_out"]
		node.Net3PortOut = item["radar3_port_out"]
		node.Net0Type, _ = strconv.Atoi(item["radar0_type"])
		node.Net1Type, _ = strconv.Atoi(item["radar1_type"])
		node.Net2Type, _ = strconv.Atoi(item["radar2_type"])
		node.Net3Type, _ = strconv.Atoi(item["radar3_type"])
		node.DeviceID, _ = strconv.Atoi(item["device_id"])
		var path string = configuration.Server.LocalImplementPath + "r" + item["device_id"]
		CreateAbsoluteDirectory(path)
		CreateAbsoluteDirectory(path + "/config")
		CreateAbsoluteDirectory(path + "/data")
		CreateRadarConfigFileNet(path+"/config/vec.config", node, configuration.RadarTypeVec, configuration.Server, configuration.Project)
	}
}

func UpdateRadarConfig() {

	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
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
		var path string = configuration.Server.LocalImplementPath + "r" + item["device_id"]
		CreateAbsoluteDirectory(path)
		CreateAbsoluteDirectory(path + "/config")
		CreateAbsoluteDirectory(path + "/data")
		CreateRadarConfigFile(path+"/config/vec.config", node, configuration.RadarTypeVec, configuration.Server, configuration.Project)
	}
}

func CopyFiles(innerPath string) {
	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
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
		var path string = configuration.Server.LocalImplementPath + "r" + item["device_id"]
		CreateAbsoluteDirectory(path)
		Copy(innerPath, path+"/"+filepath.Base(innerPath))
	}
}

func UpdateRadarConfigNetPortTransfer() {

	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
	}
	var radarPortConfigs []RadarPortTransfer
	CreateAbsoluteDirectory("./portConfig")
	for _, item := range result {
		var node NodeConfig

		node.StakeMark = item["stake_mark"]
		node.DeviceID, err = strconv.Atoi(item["device_id"])
		Check(err)
		node.Net0Ip = item["radar0_ip"]
		node.Net1Ip = item["radar1_ip"]
		node.Net2Ip = item["radar2_ip"]
		node.Net3Ip = item["radar3_ip"]
		node.Net0PortIn = item["radar0_port_in"]
		node.Net1PortIn = item["radar1_port_in"]
		node.Net2PortIn = item["radar2_port_in"]
		node.Net3PortIn = item["radar3_port_in"]
		node.Net0PortOut = item["radar0_port_out"]
		node.Net1PortOut = item["radar1_port_out"]
		node.Net2PortOut = item["radar2_port_out"]
		node.Net3PortOut = item["radar3_port_out"]
		node.Net0Type, _ = strconv.Atoi(item["radar0_type"])
		node.Net1Type, _ = strconv.Atoi(item["radar1_type"])
		node.Net2Type, _ = strconv.Atoi(item["radar2_type"])
		node.Net3Type, _ = strconv.Atoi(item["radar3_type"])
		node.DeviceID, _ = strconv.Atoi(item["device_id"])

		var configs []RadarPortTransfer
		if node.Net0Ip != "" {
			var config RadarPortTransfer
			config.InputIp = node.Net0Ip
			config.InputPort = node.Net0PortIn
			config.OutputIp = "127.0.0.1"
			config.OutputPort = node.Net0PortOut
			configs = append(configs, config)
		}
		if node.Net1Ip != "" {
			var config RadarPortTransfer
			config.InputIp = node.Net1Ip
			config.InputPort = node.Net1PortIn
			config.OutputIp = "127.0.0.1"
			config.OutputPort = node.Net1PortOut
			configs = append(configs, config)

		}
		if node.Net2Ip != "" {
			var config RadarPortTransfer
			config.InputIp = node.Net2Ip
			config.InputPort = node.Net2PortIn
			config.OutputIp = "127.0.0.1"
			config.OutputPort = node.Net2PortOut
			configs = append(configs, config)

		}
		if node.Net3Ip != "" {
			var config RadarPortTransfer
			config.InputIp = node.Net3Ip
			config.InputPort = node.Net3PortIn
			config.OutputIp = "127.0.0.1"
			config.OutputPort = node.Net3PortOut
			configs = append(configs, config)

		}
		a, _ := json.Marshal(configs)
		_ = os.WriteFile("./portConfig/config"+strconv.Itoa(node.DeviceID)+".json", a, 0644)

		radarPortConfigs = append(radarPortConfigs, configs...)
	}
	a, _ := json.Marshal(radarPortConfigs)
	_ = os.WriteFile("./portConfig/config.json", a, 0644)
}

func CopyNetFiles(innerPaths []string) {
	var configuration Config
	GetConfiguration(&configuration)

	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
	}

	for _, item := range result {
		var path string = configuration.Server.LocalImplementPath + "r" + item["device_id"]
		CreateAbsoluteDirectory(path)
		for _, v := range innerPaths {
			Copy(v, path+"/"+filepath.Base(v))
		}

	}
}

func updateRadarConfigNetRemote() {
	var configuration Config
	GetConfiguration(&configuration)
	CreateAbsoluteDirectory("./portConfig")
	var result []map[string]string
	result, err := CSVFileToMap("./config/device.csv")
	if err != nil {
		result, err = CSVFileToMap("../config/device.csv")
		if err != nil {
			panic(err)
		}
	}

	for _, item := range result {
		var node NodeConfig
		node.StakeMark = item["stake_mark"]

		node.Net0Ip = item["radar0_ip"]
		node.Net1Ip = item["radar1_ip"]
		node.Net2Ip = item["radar2_ip"]
		node.Net3Ip = item["radar3_ip"]
		node.Net0PortIn = item["radar0_port_in"]
		node.Net1PortIn = item["radar1_port_in"]
		node.Net2PortIn = item["radar2_port_in"]
		node.Net3PortIn = item["radar3_port_in"]
		node.Net0PortOut = item["radar0_port_out"]
		node.Net1PortOut = item["radar1_port_out"]
		node.Net2PortOut = item["radar2_port_out"]
		node.Net3PortOut = item["radar3_port_out"]
		node.Net0Type, _ = strconv.Atoi(item["radar0_type"])
		node.Net1Type, _ = strconv.Atoi(item["radar1_type"])
		node.Net2Type, _ = strconv.Atoi(item["radar2_type"])
		node.Net3Type, _ = strconv.Atoi(item["radar3_type"])
		node.DeviceID, _ = strconv.Atoi(item["device_id"])
		CreateRadarConfigFileNet("./portConfig/vec"+item["device_id"]+".config", node, configuration.RadarTypeVec, configuration.Server, configuration.Project)
	}
}
