package ipdetails

import "encoding/json"

type IpRangeDetails struct {
	StartIP   string `csv:"start_ip"`
	EndIp     string `csv:"end_ip"`
	JoinKey   string `csv:"join_key"`
	IsHosting bool   `csv:"hosting"`
	IsProxy   bool   `csv:"proxy"`
	IsTor     bool   `csv:"tor"`
	IsVPN     bool   `csv:"vpn"`
	IsRelay   bool   `csv:"relay"`
	IsService bool   `csv:"service"`
}

func (m IpRangeDetails) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m IpRangeDetails) UnmarshalBinary(data []byte) error {
	var msg IpRangeDetails
	err := json.Unmarshal(data, &msg)
	m = msg
	return err
}
