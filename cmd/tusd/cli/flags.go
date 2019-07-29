package cli

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Flags struct {
	HttpHost            string
	HttpPort            string
	HttpSock            string
	MaxSize             int64  `default: 0`
	UploadDir           string `default: "./data"`
	StoreSize           int64  `default: 0`
	Basepath            string `default: "/files"`
	Timeout             int64  `default: 30000`
	S3Bucket            string
	S3ObjectPrefix      string
	S3Endpoint          string
	GCSBucket           string
	GCSObjectPrefix     string
	FileHooksDir        string
	HttpHooksEndpoint   string
	HttpHooksRetry      int `default: 3`
	HttpHooksBackoff    int `default: 1`
	HooksStopUploadCode int `default: 0`
	PluginHookPath      string
	ShowVersion         bool
	ExposeMetrics       bool   `default: true`
	MetricsPath         string `default: "/metrics"`
	BehindProxy         bool   `default: false`

	FileHooksInstalled bool
	HttpHooksInstalled bool
}

func ParseFlags() {

	flag.StringVar(&Flags.HttpHost, "host", "0.0.0.0", "Host to bind HTTP server to")
	flag.StringVar(&Flags.HttpPort, "port", "4000", "Port to bind HTTP server to")
	flag.BoolVar(&Flags.ShowVersion, "version", false, "Print tusd version information")

	flag.Parse()

	// now read and parse ENV vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = envconfig.Process("tusd", &Flags)
	if err != nil {
		log.Fatal("Error binding flags")
	}

	if Flags.FileHooksDir != "" {
		Flags.FileHooksDir, _ = filepath.Abs(Flags.FileHooksDir)
		Flags.FileHooksInstalled = true

		stdout.Printf("Using '%s' for hooks", Flags.FileHooksDir)
	}

	if Flags.HttpHooksEndpoint != "" {
		Flags.HttpHooksInstalled = true

		stdout.Printf("Using '%s' as the endpoint for hooks", Flags.HttpHooksEndpoint)
	}

	if Flags.UploadDir == "" && Flags.S3Bucket == "" {
		stderr.Fatalf("Either an upload directory (using -dir) or an AWS S3 Bucket " +
			"(using -s3-bucket) must be specified to start tusd but " +
			"neither flag was provided. Please consult `tusd -help` for " +
			"more information on these options.")
	}

	if Flags.GCSObjectPrefix != "" && strings.Contains(Flags.GCSObjectPrefix, "_") {
		stderr.Fatalf("gcs-object-prefix value (%s) can't contain underscore. "+
			"Please remove underscore from the value", Flags.GCSObjectPrefix)
	}
}
