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
}

type ProjectConfiguration struct {
	ProjectNum int
}

type Config struct {
	Server       ServerConfigurations
	Project      ProjectConfiguration
	RadarTypeVec []RadarType
}

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
