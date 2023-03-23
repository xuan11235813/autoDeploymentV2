package main

import (
	. "autoDeploy/comm"
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {

	UpdateRadarConfig()
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
	var radarPosConfigs []RadarPosConfig

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
		var configs []RadarPosConfig = GenerateRadarPosFromNode(node, configuration.RadarTypeVec, configuration.Server, configuration.Project)
		radarPosConfigs = append(radarPosConfigs, configs...)
	}

	a, _ := json.Marshal(radarPosConfigs)
	fmt.Println(string(a))

}
