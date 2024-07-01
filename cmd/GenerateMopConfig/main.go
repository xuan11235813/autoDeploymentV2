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
	var inputSliceMap []map[string]string
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

		for _, v := range configs {

			item := map[string]string{
				"RadarID":             v.RadarID,
				"Comment":             v.Comment,
				"Stake":               strconv.FormatFloat(v.Position.X, 'g', -1, 64),
				"Direction":           strconv.FormatInt(int64(Btoi(v.IsZH2HK)), 10),
				"SoftMaxIncoming":     "0",
				"CoordinateLongitude": "0",
				"CoordinateLatitude":  "0",
				"ViewPositionX":       strconv.FormatFloat(v.Position.X, 'g', -1, 64),
				"ViewPositionY":       strconv.FormatFloat(v.Position.Y, 'g', -1, 64),
				"ViewPositionZ":       "0",
				"AngleDeg":            strconv.FormatFloat(v.Angle, 'g', -1, 64),
				"InTunnel":            "0",
			}

			inputSliceMap = append(inputSliceMap, item)
		}
	}
	var header = []string{"RadarID", "Comment", "Stake", "Direction", "SoftMaxIncoming", "CoordinateLongitude", "CoordinateLatitude", "ViewPositionX", "ViewPositionY", "ViewPositionZ", "AngleDeg", "InTunnel"}
	MapToCSVFile(inputSliceMap, "radarMop.csv", header)

	a, _ := json.Marshal(radarPosConfigs)
	fmt.Println(string(a))

}
