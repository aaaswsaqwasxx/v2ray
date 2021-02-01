package api

import (
	"v2ray.com/core/main/commands/base"
)

// CmdAPI calls an API in an V2Ray process
var CmdAPI = &base.Command{
	UsageLine: "{{.Exec}} api",
	Short:     "Call V2Ray API",
	Long: `{{.Exec}} {{.LongName}} provides tools to manipulate V2Ray via its API.
`,
	Commands: []*base.Command{
		cmdBalancerCheck,
		cmdBalancerInfo,
		cmdBalancerOverride,
		cmdAddInbounds,
		cmdAddOutbounds,
		cmdRemoveInbounds,
		cmdRemoveOutbounds,
		cmdStats,
		cmdRestartLogger,
	},
}
