package all

import (
	"bytes"
	"os"

	"google.golang.org/protobuf/proto"
	"v2ray.com/core/commands/base"
	"v2ray.com/core/infra/conf/merge"
	"v2ray.com/core/infra/conf/serial"
)

var cmdConvert = &base.Command{
	UsageLine: "{{.Exec}} convert [c1.json] [<url>.json] [dir1] ...",
	Short:     "Convert multiple json config to protobuf",
	Long: `
Merge multiple JSON config and convert to protobuf.

Arguments:

	-r
		Load confdir recursively.

Examples:

	{{.Exec}} {{.LongName}} c1.json c2.json 
	{{.Exec}} {{.LongName}} c1.json https://url.to/c2.json 
	{{.Exec}} {{.LongName}} "path/to/json_dir"

To learn how JSON files are merged, run "{{.Exec}} help merge"
`,
}

func init() {
	cmdConvert.Run = executeConvert // break init loop
}

var convertReadDirRecursively = cmdConvert.Flag.Bool("r", false, "")

func executeConvert(cmd *base.Command, args []string) {
	unnamed := cmd.Flag.Args()
	files := resolveFolderToFiles(unnamed, *convertReadDirRecursively)
	if len(files) == 0 {
		base.Fatalf("empty config list")
	}

	data, err := merge.JSONs(files)
	if err != nil {
		base.Fatalf("failed to load json: %s", err)
	}
	r := bytes.NewReader(data)
	cf, err := serial.DecodeJSONConfig(r)
	if err != nil {
		base.Fatalf("failed to decode json: %s", err)
	}

	pbConfig, err := cf.Build()
	if err != nil {
		base.Fatalf(err.Error())
	}

	bytesConfig, err := proto.Marshal(pbConfig)
	if err != nil {
		base.Fatalf("failed to marshal proto config: %s", err)
	}

	if _, err := os.Stdout.Write(bytesConfig); err != nil {
		base.Fatalf("failed to write proto config: %s", err)
	}
}