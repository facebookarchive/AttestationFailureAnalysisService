package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"libfb/go/afas"

	"github.com/immune-gmbh/AttestationFailureAnalysisService/pkg/observability"
	"github.com/immune-gmbh/AttestationFailureAnalysisService/server/controller"
	"github.com/immune-gmbh/AttestationFailureAnalysisService/server/controller/consts"
	"github.com/immune-gmbh/AttestationFailureAnalysisService/server/thrift"

	"github.com/facebookincubator/go-belt/tool/logger"
	fianoLog "github.com/linuxboot/fiano/pkg/log"
	"github.com/spf13/pflag"
)

const (
	diffFirmwareCacheSizeDefault     = 16
	storageCacheSizeDefault          = 1 << 30 // 1GiB
	reportHostConfigCacheSizeDefault = 1000
	rtpfwCacheSizeDefault            = 20
	rtpfwCacheEvictionTimeoutDefault = 24 * time.Hour
	apiCachePurgeTimeoutDefault      = time.Hour
	dataCacheSizeDefault             = 1000
)

func assertNoError(ctx context.Context, err error) {
	if err != nil {
		logger.FromCtx(ctx).Fatalf("%v", err)
	}
}

func usageExit() {
	pflag.Usage()
	os.Exit(2) // The default Go's exitcode on flag.Parse() problems
}

func main() {
	logLevel := logger.LevelInfo // the default value

	fbid := pflag.Int64("fbid", consts.DefaultFBID, "user ID to act from")
	hipsterACL := pflag.String("hipster-acl", afas.DefaultHipsterACL, "Hipster ACL name")
	pflag.Var(&logLevel, "log-level", "logging level")
	netPprofAddr := pflag.String("net-pprof-addr", "", "if non-empty then listens with net/http/pprof")
	thriftBindAddr := pflag.String("thrift-bind-addr", `:17545`, "the address to listen by thrift")
	tierBind := pflag.String("tier", afas.DefaultSMCTier, "tier to register on")
	rdbmsURL := pflag.String("tier-db", `xdb.afas`, "tier to RDBMS")
	tierSeRF := pflag.String("tier-serf", `serf_device_cache`, "SeRF device cache tier")
	manifoldBucket := pflag.String("manifold-bucket", `firmware_images`, "manifold bucket to storage firmware images")
	manifoldAPIKey := pflag.String("manifold-api-key", `firmware_images-key`, "manifold API key to storage firmware images")
	amountOfWorkers := pflag.Uint("workers", uint(runtime.NumCPU()), "amount of concurrent workers")
	workersQueue := pflag.Uint("workers-queue", uint(runtime.NumCPU())*10000, "maximal amount of requests permitted in the queue")
	cpuLoadLimit := pflag.Float64("cpu-load-limit", 0.8, "suspend accepting requests while fraction of busy CPU cycles is more than the specified number")
	rtpfwCacheSize := pflag.Int("rtp-cache-size", rtpfwCacheSizeDefault, "defines how many RTP table firmwares are stored in memory")
	rtpfwCacheEvictionTimeout := pflag.Duration(
		"rtp-table-eviction-timeout",
		rtpfwCacheEvictionTimeoutDefault,
		"defines eviction timeout based on element access after which rtp cache items are evicted",
	)
	apiCachePurgeTimeout := pflag.Duration(
		"api-cache-purge-timeout",
		apiCachePurgeTimeoutDefault,
		"defines API cache purge timeout",
	)
	storageCacheSize := pflag.Uint64("image-storage-cache-size", storageCacheSizeDefault, "defines the memory limit for the storage used to save images, analyzed by AFAS")
	diffFirmwareCacheSize := pflag.Int("diff-firmware-cache-size", diffFirmwareCacheSizeDefault, "defines how many DiffFirmware reports are stored in memory")
	reportHostConfigurationCacheSize := pflag.Int("report-host-config-cache-size", reportHostConfigCacheSizeDefault, "defines how many ReportHostConfiguration results are stored in memory")
	dataCacheSize := pflag.Int("data-cache-size", dataCacheSizeDefault, "defines the size of the cache for internally caclulated data objects like parsed firmware, measurements flow")
	pflag.Parse()
	if pflag.NArg() != 0 {
		usageExit()
	}

	ctx := observability.WithBelt(
		context.Background(),
		logLevel,
		"AFAS", true,
	)

	log := logger.FromCtx(ctx)

	if *netPprofAddr != "" {
		go func() {
			err := http.ListenAndServe(*netPprofAddr, nil)
			log.Errorf("unable to start listening for https/net/pprof: %v", err)
		}()
	}

	fianoLog.DefaultLogger = newFianoLogger(log.WithField("module", "fiano"))

	ctrl, err := controller.New(ctx,
		*fbid,
		*rtpfwCacheSize,
		*rtpfwCacheEvictionTimeout,
		*apiCachePurgeTimeout,
		*storageCacheSize,
		*diffFirmwareCacheSize,
		*reportHostConfigurationCacheSize,
		*dataCacheSize,
		*rdbmsURL, *tierSeRF,
		"firmware_measurements",
		"attestation_hosts_configurations",
		*manifoldBucket, *manifoldAPIKey,
	)
	assertNoError(ctx, err)
	log.Debugf("created a controller")

	srv, err := thrift.NewServer(
		ctx,
		*amountOfWorkers,
		*workersQueue,
		*cpuLoadLimit,
		*tierBind,
		*hipsterACL,
		ctrl,
	)
	assertNoError(ctx, err)
	log.Debugf("created a Thrift server")

	err = srv.Serve(ctx, *thriftBindAddr)
	assertNoError(ctx, err)
}
