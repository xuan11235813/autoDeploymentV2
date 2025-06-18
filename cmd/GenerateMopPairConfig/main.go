package main

import (
	. "autoDeploy/comm"
	"math"
	"os"
	"slices"
	"strconv"
)

func main() {
	switch os.Args[1] {
	case "sztd":
		UpdateRadarConfigSZ()
	case "fsgs":
		UpdateRadarConfigFS()
	}

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
		} else if i.RadarTypeItem.TypeNum == 27 && j.RadarTypeItem.TypeNum == 28 {
			return 1
		} else if i.RadarTypeItem.TypeNum == 28 && j.RadarTypeItem.TypeNum == 27 {
			return -1
		} else if i.RadarTypeItem.TypeNum == 29 && j.RadarTypeItem.TypeNum == 26 {
			return -1
		} else if i.RadarTypeItem.TypeNum == 26 && j.RadarTypeItem.TypeNum == 29 {
			return 1
		} else {
			return -1
		}
	}
}

func UpdateRadarConfigFS() {
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
		node.Net0Type, _ = strconv.Atoi(item["radar0_type"])
		node.Net1Type, _ = strconv.Atoi(item["radar1_type"])
		node.Net2Type, _ = strconv.Atoi(item["radar2_type"])
		node.Net3Type, _ = strconv.Atoi(item["radar3_type"])
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

	var radarPosNormalRight []RadarPosConfig
	var radarPosNormalLeft []RadarPosConfig
	var radarPosTunnelRight []RadarPosConfig
	var radarPosTunnelLeft []RadarPosConfig
	for _, itemQue := range radarQueue {
		if itemQue[0].RadarTypeItem.TypeNum == 29 || itemQue[0].RadarTypeItem.TypeNum == 26 {
			radarPosNormalRight = append(radarPosNormalRight, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 28 || itemQue[0].RadarTypeItem.TypeNum == 27 {
			radarPosNormalLeft = append(radarPosNormalLeft, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 19 {
			radarPosTunnelRight = append(radarPosTunnelRight, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 18 {
			radarPosTunnelLeft = append(radarPosTunnelLeft, itemQue...)
		}
	}

	slices.SortFunc(radarPosNormalRight, compare)
	slices.SortFunc(radarPosNormalLeft, compare)
	slices.SortFunc(radarPosTunnelRight, compare)
	slices.SortFunc(radarPosTunnelLeft, compare)

	temp := radarPosNormalRight
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
	temp = radarPosNormalLeft
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
	temp = radarPosTunnelRight
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
	temp = radarPosTunnelLeft
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

func UpdateRadarConfigSZ() {

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

	var radarPosNormalRight []RadarPosConfig
	var radarPosNormalLeft []RadarPosConfig
	var radarPosTunnelRight []RadarPosConfig
	var radarPosTunnelLeft []RadarPosConfig
	for _, itemQue := range radarQueue {
		if itemQue[0].RadarTypeItem.TypeNum == 14 || itemQue[0].RadarTypeItem.TypeNum == 15 {
			radarPosNormalRight = append(radarPosNormalRight, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 17 || itemQue[0].RadarTypeItem.TypeNum == 16 {
			radarPosNormalLeft = append(radarPosNormalLeft, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 19 {
			radarPosTunnelRight = append(radarPosTunnelRight, itemQue...)
		} else if itemQue[0].RadarTypeItem.TypeNum == 18 {
			radarPosTunnelLeft = append(radarPosTunnelLeft, itemQue...)
		}
	}

	slices.SortFunc(radarPosNormalRight, compare)
	slices.SortFunc(radarPosNormalLeft, compare)
	slices.SortFunc(radarPosTunnelRight, compare)
	slices.SortFunc(radarPosTunnelLeft, compare)

	temp := radarPosNormalRight
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
	temp = radarPosNormalLeft
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
	temp = radarPosTunnelRight
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
	temp = radarPosTunnelLeft
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
