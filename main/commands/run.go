package commands

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"v2ray.com/core"
	"v2ray.com/core/commands/base"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/common/platform"
)

// CmdRun runs V2Ray with config
var CmdRun = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} run [-c config.json] [-d dir]",
	Short:       "Run V2Ray with config",
	Long: `
Run V2Ray with config.

Example:

	{{.Exec}} {{.LongName}} -c config.json
	{{.Exec}} {{.LongName}} -d path/to/json_dir

Arguments:

	-c, -config
		Config file for V2Ray. Multiple assign is accepted.

	-d, -confdir
		A dir with multiple json config. Multiple assign is accepted.

	-r
		Load confdir recursively.

	-format
		Format of input files. (default "json")

Use "{{.Exec}} help format-loader" for more information about format.
	`,
	Run: executeRun,
}

var (
	configFiles          cmdarg.Arg
	configDirs           cmdarg.Arg
	configFormat         *string
	configDirRecursively *bool
)

func setConfigFlags(cmd *base.Command) {
	configFormat = cmd.Flag.String("format", "json", "")
	configDirRecursively = cmd.Flag.Bool("r", false, "")

	cmd.Flag.Var(&configFiles, "config", "")
	cmd.Flag.Var(&configFiles, "c", "")
	cmd.Flag.Var(&configDirs, "confdir", "")
	cmd.Flag.Var(&configDirs, "d", "")
}

func executeRun(cmd *base.Command, args []string) {
	setConfigFlags(cmd)
	cmd.Flag.Parse(args)
	printVersion()
	server, err := startV2Ray()
	if err != nil {
		base.Fatalf("Failed to start: %s", err)
	}

	if err := server.Start(); err != nil {
		base.Fatalf("Failed to start: %s", err)
	}
	defer server.Close()

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
	}
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func dirExists(file string) bool {
	if file == "" {
		return false
	}
	info, err := os.Stat(file)
	return err == nil && info.IsDir()
}

func readConfDir(dirPath string, extension []string) cmdarg.Arg {
	confs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		base.Fatalf("failed to read dir %s: %s", dirPath, err)
	}
	files := make(cmdarg.Arg, 0)
	for _, f := range confs {
		ext := filepath.Ext(f.Name())
		for _, e := range extension {
			if strings.EqualFold(e, ext) {
				files.Set(filepath.Join(dirPath, f.Name()))
				break
			}
		}
	}
	return files
}

// getFolderFiles get files in the folder and it's children
func readConfDirRecursively(dirPath string, extension []string) cmdarg.Arg {
	files := make(cmdarg.Arg, 0)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		for _, e := range extension {
			if strings.EqualFold(e, ext) {
				files.Set(path)
				break
			}
		}
		return nil
	})
	if err != nil {
		base.Fatalf("failed to read dir %s: %s", dirPath, err)
	}
	return files
}

func getLoaderExtension() ([]string, error) {
	firstFile := ""
	if len(configFiles) > 0 {
		firstFile = configFiles[0]
	}
	loader, err := core.GetConfigLoader(*configFormat, firstFile)
	if err != nil {
		return nil, err
	}
	return loader.Extension, nil
}

func getConfigFilePath() cmdarg.Arg {
	extension, err := getLoaderExtension()
	if err != nil {
		base.Fatalf(err.Error())
	}
	dirReader := readConfDir
	if *configDirRecursively {
		dirReader = readConfDirRecursively
	}
	if len(configDirs) > 0 {
		for _, d := range configDirs {
			log.Println("Using confdir from arg:", d)
			configFiles = append(configFiles, dirReader(d, extension)...)
		}
	} else if envConfDir := platform.GetConfDirPath(); dirExists(envConfDir) {
		log.Println("Using confdir from env:", envConfDir)
		configFiles = append(configFiles, dirReader(envConfDir, extension)...)
	}
	if len(configFiles) > 0 {
		return configFiles
	}

	if len(configFiles) == 0 && len(configDirs) > 0 {
		base.Fatalf("no config file found with extension: %s", extension)
	}

	if workingDir, err := os.Getwd(); err == nil {
		configFile := filepath.Join(workingDir, "config.json")
		if fileExists(configFile) {
			log.Println("Using default config: ", configFile)
			return cmdarg.Arg{configFile}
		}
	}

	if configFile := platform.GetConfigurationPath(); fileExists(configFile) {
		log.Println("Using config from env: ", configFile)
		return cmdarg.Arg{configFile}
	}

	log.Println("Using config from STDIN")
	return cmdarg.Arg{"stdin:"}
}

func startV2Ray() (core.Server, error) {
	configFiles := getConfigFilePath()

	config, err := core.LoadConfig(*configFormat, configFiles[0], configFiles)
	if err != nil {
		return nil, newError("failed to read config files: [", configFiles.String(), "]").Base(err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, newError("failed to create server").Base(err)
	}

	return server, nil
}