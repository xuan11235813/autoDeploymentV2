package main

import (
	. "autoDeploy/comm"
	"math"
	"slices"
	"strconv"
)

func main() {

	UpdateRadarConfig()
}

func compare(i, j RadarPosConfig) int {
	if i.Position.X < j.Position.X {
		return -1
	} else if i.Position.X > j.Position.X {
		return 1
	} else {
		if i.RadarTypeItem.TypeNum == 14 && j.RadarTypeItem.TypeNum == 15 {
			return 1
		} else if i.RadarTypeItem.TypeNum == 15 && j.RadarTypeItem.TypeNum == 14 {
			return -1
		} else if i.RadarTypeItem.TypeNum == 16 && j.RadarTypeItem.TypeNum == 17 {
			return -1
		} else if i.RadarTypeItem.TypeNum == 17 && j.RadarTypeItem.TypeNum == 16 {
			return 1
		} else {
			return -1
		}
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
	var radarPosConfigs []RadarPosConfig
	var radarQueue [][]RadarPosConfig
	var radarPairResult []map[string]string

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
	for _, item := range radarPosConfigs {
		if len(radarQueue) == 0 {
			var radarQueueItem []RadarPosConfig
			radarQueueItem = append(radarQueueItem, item)
			radarQueue = append(radarQueue, radarQueueItem)
		} else {
			addFlag := false
			for idx, queueItem := range radarQueue {
				if queueItem[0].RadarTypeItem.TypeNum == item.RadarTypeItem.TypeNum {
					queueItem = append(queueItem, item)
					radarQueue[idx] = queueItem
					addFlag = true
				}
			}
			if !addFlag {
				var radarQueueItem []RadarPosConfig
				radarQueueItem = append(radarQueueItem, item)
				radarQueue = append(radarQueue, radarQueueItem)
			}
		}
	}

	var radarPosNormalSZ []RadarPosConfig
	var radarPosNormalZS []RadarPosConfig
	var radarPosTunnelSZ []RadarPosConfig
	var radarPosTunnelZS []RadarPosConfig
	for _, itemQue := range radarQueue {
		if itemQue[0].RadarTypeItem.TypeNum == 14 || itemQue[0].RadarTypeItem.TypeNum == 15 {
			radarPosNormalSZ = append(radarPosNormalSZ, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 17 || itemQue[0].RadarTypeItem.TypeNum == 16 {
			radarPosNormalZS = append(radarPosNormalZS, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 19 {
			radarPosTunnelSZ = append(radarPosTunnelSZ, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 18 {
			radarPosTunnelZS = append(radarPosTunnelZS, itemQue...)
		}
	}

	slices.SortFunc(radarPosNormalSZ, compare)
	slices.SortFunc(radarPosNormalZS, compare)
	slices.SortFunc(radarPosTunnelSZ, compare)
	slices.SortFunc(radarPosTunnelZS, compare)

	temp := radarPosNormalSZ
	for i := 0; i < len(temp)-1; i++ {
		item := map[string]string{}
		item["RadarID[0]"] = temp[i].RadarID
		item["RadarID[1]"] = temp[i+1].RadarID
		if math.Abs(temp[i].Position.X-temp[i+1].Position.X) <= 1 {
			item["Type"] = "0"
		} else {
			if temp[i].RadarTypeItem.RadarDirection == temp[i+1].RadarTypeItem.RadarDirection {
				item["Type"] = "2"
			} else {
				item["Type"] = "1"
			}
		}
		item["ApproximateDistance"] = ""
		item["Fixed"] = "FALSE"
		radarPairResult = append(radarPairResult, item)
	}
	temp = radarPosNormalZS
	for i := 0; i < len(temp)-1; i++ {
		item := map[string]string{}
		item["RadarID[0]"] = temp[i].RadarID
		item["RadarID[1]"] = temp[i+1].RadarID
		if math.Abs(temp[i].Position.X-temp[i+1].Position.X) <= 1 {
			item["Type"] = "0"
		} else {
			if temp[i].RadarTypeItem.RadarDirection == temp[i+1].RadarTypeItem.RadarDirection {
				item["Type"] = "2"
			} else {
				item["Type"] = "1"
			}
		}
		item["ApproximateDistance"] = ""
		item["Fixed"] = "FALSE"
		radarPairResult = append(radarPairResult, item)
	}
	temp = radarPosTunnelSZ
	for i := 0; i < len(temp)-1; i++ {
		item := map[string]string{}
		item["RadarID[0]"] = temp[i].RadarID
		item["RadarID[1]"] = temp[i+1].RadarID
		if math.Abs(temp[i].Position.X-temp[i+1].Position.X) <= 1 {
			item["Type"] = "0"
		} else {
			if temp[i].RadarTypeItem.RadarDirection == temp[i+1].RadarTypeItem.RadarDirection {
				item["Type"] = "2"
			} else {
				item["Type"] = "1"
			}
		}
		item["ApproximateDistance"] = ""
		item["Fixed"] = "FALSE"
		radarPairResult = append(radarPairResult, item)
	}
	temp = radarPosTunnelZS
	for i := 0; i < len(temp)-1; i++ {
		item := map[string]string{}
		item["RadarID[0]"] = temp[i].RadarID
		item["RadarID[1]"] = temp[i+1].RadarID
		if math.Abs(temp[i].Position.X-temp[i+1].Position.X) <= 1 {
			item["Type"] = "0"
		} else {
			if temp[i].RadarTypeItem.RadarDirection == temp[i+1].RadarTypeItem.RadarDirection {
				item["Type"] = "2"
			} else {
				item["Type"] = "1"
			}
		}
		item["ApproximateDistance"] = ""
		item["Fixed"] = "FALSE"
		radarPairResult = append(radarPairResult, item)
	}
	var header = []string{"RadarID[0]", "RadarID[1]", "Type", "ApproximateDistance", "Fixed"}
	MapToCSVFile(radarPairResult, "radarPair.csv", header)
}
