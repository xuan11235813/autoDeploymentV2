package comm

type ServerConfigurations struct {
	Port               int
	IPAddress          string
	ExtraPort          int
	LocalImplementPath string
}

type RadarType struct {
	TypeNum         int
	IsTunnel        int
	IncomingLaneNum int
	OutgoingLaneNum int
	StartIncoming   int
	StartOutgoing   int
	RadarDirection  int
}

type ProjectConfiguration struct {
	ProjectNum            int
	ProjectName           string
	ProjectStartStakeMark string
}

type Config struct {
	Server       ServerConfigurations
	Project      ProjectConfiguration
	RadarTypeVec []RadarType
}

/*
	type NodeConfig struct {
		DeviceID           int
		IpAddress          string
		UserName           string
		Password           string
		StakeMark          string
		Can0Type           int
		Can1Type           int
		Can2Type           int
		Can3Type           int
		Can0ChessboardFile string
		Can1ChessboardFile string
		Can2ChessboardFile string
		Can3ChessboardFile string
	}
*/
type NodeConfig struct {
	DeviceID           int
	IpAddress          string
	UserName           string
	Password           string
	StakeMark          string
	Can0Type           int
	Can1Type           int
	Can2Type           int
	Can3Type           int
	Can0ChessboardFile string
	Can1ChessboardFile string
	Can2ChessboardFile string
	Can3ChessboardFile string
	Net0Ip             string
	Net1Ip             string
	Net2Ip             string
	Net3Ip             string
	Net0PortIn         string
	Net1PortIn         string
	Net2PortIn         string
	Net3PortIn         string
	Net0Type           int
	Net1Type           int
	Net2Type           int
	Net3Type           int
	Net0PortOut        string
	Net1PortOut        string
	Net2PortOut        string
	Net3PortOut        string
}

type RadarPos struct {
	X float64
	Y float64
	Z float64
}

type RadarPosConfig struct {
	Position       RadarPos
	Angle          float64
	RadarID        string
	DenyLaneChange bool
	Comment        string
	IsZH2HK        bool
	Direction      bool

	RadarTypeItem RadarType
}

type RadarPortTransfer struct {
	InputIp    string
	InputPort  string
	OutputPort string
	OutputIp   string
}
