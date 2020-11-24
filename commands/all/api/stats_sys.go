package api

import (
	"context"
	"time"

	"google.golang.org/grpc"
	statsService "v2ray.com/core/app/stats/command"
	"v2ray.com/core/commands/base"
)

var cmdSysStats = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api statssys [--server=127.0.0.1:8080]",
	Short:       "Get system statistics",
	Long: `
Get sys statistics from V2Ray by calling its API. (timeout 3 seconds)

Arguments:

	-server=127.0.0.1:8080 
		The API server address. Default 127.0.0.1:8080
`,
	Run: executeSysStats,
}

func executeSysStats(cmd *base.Command, args []string) {
	setSharedFlags(cmd)
	cmd.Flag.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *apiServerAddrPtr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		base.Fatalf("failed to dial %s", *apiServerAddrPtr)
	}
	defer conn.Close()

	client := statsService.NewStatsServiceClient(conn)
	r := &statsService.SysStatsRequest{}
	resp, err := client.GetSysStats(ctx, r)
	if err != nil {
		base.Fatalf("failed to get sys stats: %s", err)
	}
	showResponese(responeseToString(resp))
}
