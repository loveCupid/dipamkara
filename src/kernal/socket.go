package kernal

import (
	"net"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func GetValidPort(ip string) int {
	ip += ":0"
	l, err := net.Listen("tcp", ip)
	ErrorPanic(err)

	ret := l.Addr().(*net.TCPAddr).Port

	l.Close()

	return ret
}

func GetValidIP() string {
	vv, err := net.Interfaces()
	ErrorPanic(err)

	for _, v := range vv {
		// ipnet, _ := v.(*net.IPNet)
		if ((v.Flags & net.FlagUp) != 0) && (v.Name == "wifi0") {
			addrs, _ := v.Addrs()
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						// fmt.Println(ipnet.IP.String())
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	panic("not found valid ip")
	return ""
}
