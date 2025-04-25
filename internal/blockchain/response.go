package blockchain

type Participant struct {
	CLContext CLContext `json:"cl_context"`
	CLType    string    `json:"cl_type"`
	ELContext ELContext `json:"el_context"`
	ELType    string    `json:"el_type"`
}

type CLContext struct {
	BeaconHTTPURL string `json:"beacon_http_url"`
}

type ELContext struct {
	RpcHttpUrl string `json:"rpc_http_url"`
	WsUrl      string `json:"ws_url"`
}

type ConfigResponse struct {
	AllParticipants       []Participant `json:"all_participants"`
	NetworkID             string        `json:"network_id"`
	FinalGenesisTimestamp string        `json:"final_genesis_timestamp"`
}
