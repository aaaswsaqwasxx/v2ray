package api

import (
	"fmt"

	handlerService "github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	"github.com/v2fly/v2ray-core/v4/common/cmdarg"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/v2fly/v2ray-core/v4/infra/conf/serial"
	"github.com/v2fly/v2ray-core/v4/main/commands/base"
)

var cmdAddOutbounds = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api ado [--server=127.0.0.1:8080] <c1.json> [c2.json]...",
	Short:       "add outbounds",
	Long: `
Add outbounds to V2Ray.

Arguments:

	-s, -server <server:port>
		The API server address. Default 127.0.0.1:8080

	-t, -timeout <seconds>
		Timeout seconds to call API. Default 3

Example:

    {{.Exec}} {{.LongName}} --server=127.0.0.1:8080 c1.json c2.json
`,
	Run: executeAddOutbounds,
}

func executeAddOutbounds(cmd *base.Command, args []string) {
	setSharedFlags(cmd)
	cmd.Flag.Parse(args)
	unnamedArgs := cmd.Flag.Args()
	if len(unnamedArgs) == 0 {
		fmt.Println("Reading from STDIN")
		unnamedArgs = []string{"stdin:"}
	}

	outs := make([]conf.OutboundDetourConfig, 0)
	for _, arg := range unnamedArgs {
		r, err := cmdarg.LoadArg(arg)
		if err != nil {
			base.Fatalf("failed to load %s: %s", arg, err)
		}
		conf, err := serial.DecodeJSONConfig(r)
		if err != nil {
			base.Fatalf("failed to decode %s: %s", arg, err)
		}
		outs = append(outs, conf.OutboundConfigs...)
	}
	if len(outs) == 0 {
		base.Fatalf("no valid outbound found")
	}

	conn, ctx, close := dialAPIServer()
	defer close()

	client := handlerService.NewHandlerServiceClient(conn)
	for _, out := range outs {
		fmt.Println("adding:", out.Tag)
		o, err := out.Build()
		if err != nil {
			base.Fatalf("failed to build conf: %s", err)
		}
		r := &handlerService.AddOutboundRequest{
			Outbound: o,
		}
		_, err = client.AddOutbound(ctx, r)
		if err != nil {
			base.Fatalf("failed to add outbound: %s", err)
		}
	}
}
