package p2pv2

const (
	PeerNameSize int = 16
	PeerIDSize   int = 16
)

const (
	P2PMsgTypeRequestIDForPublicNodeCheck uint8 = 1 // 对方询问我的 PeerID 用于判断是否为公网，答复后即可立即断开连接
	P2PMsgTypeReportIdKeepConnectAsPeer   uint8 = 2 // 报告我的 Port + PeerName + PeerID ，请求希望作为持久节点来连接
	P2PMsgTypeAnswerIdKeepConnectAsPeer   uint8 = 3 // 回复我的 PeerName +PeerID 同意作为持久连接

	P2PMsgTypeRequestNearestPublicNodes uint8 = 4 // 对方请求公网节点列表， 答复后可立即断开连接

	P2PMsgTypePing uint8 = 5 // ping
	P2PMsgTypePong uint8 = 6 // pong

	// 客户上层消息
	P2PMsgTypeCustomer uint8 = 255
)