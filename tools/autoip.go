package tools

import (
	"log/slog"

	"toolkit/system"
	"toolkit/tools/internal/autoip"

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
			if ip, err = autoip.GetLanIp(); err != nil {
				return
			}
			slog.Info("lan ip: " + ip)
			if err = autoip.DynamicDNS(ip, lanRecID); err != nil {
				return
			}
		}

		wanRecID := system.GetEnvString(envAutoipWanRecID)
		if len(wanRecID) > 0 {
			var ip string
			if ip, err = autoip.GetWanIp(); err != nil {
				return
			}
			slog.Info("wan ip: " + ip)
			if err = autoip.DynamicDNS(ip, wanRecID); err != nil {
				return
			}
		}
		return
	},
}
