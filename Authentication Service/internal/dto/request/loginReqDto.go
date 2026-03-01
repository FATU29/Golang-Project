package request

type LoginReqDto struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	IpAddress  *string `json:"ip_address"`
	DeviceInfo *string `json:"device_info"`
}
