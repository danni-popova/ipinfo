package ipdetails

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
