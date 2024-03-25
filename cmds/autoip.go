package cmds

import (
	"toolkit/logs"
	"toolkit/service"
	"toolkit/system"

	"github.com/valord577/clix"
)

const (
	envAutoipLanRecID = "TOOLKIT_AUTOIP_LAN_RECORD_ID"
	envAutoipWanRecID = "TOOLKIT_AUTOIP_WAN_RECORD_ID"
)

var AutoIp = &clix.Command{
	Name: "autoip",

	Summary: "Service of DDNS",
	Run: func(*clix.Command, []string) (err error) {
		lanRecID := system.GetEnvString(envAutoipLanRecID)
		if len(lanRecID) > 0 {
			var ip string
			if ip, err = service.GetLanIp(); err != nil {
				return
			}
			logs.Infof("lan ip: %s", ip)
			if err = service.DynamicDNS(ip, lanRecID); err != nil {
				return
			}
		}

		wanRecID := system.GetEnvString(envAutoipWanRecID)
		if len(wanRecID) > 0 {
			var ip string
			if ip, err = service.GetWanIp(); err != nil {
				return
			}
			logs.Infof("wan ip: %s", ip)
			if err = service.DynamicDNS(ip, wanRecID); err != nil {
				return
			}
		}
		return
	},
}
