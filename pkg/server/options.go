/*
Copyright 2025 Codenotary Inc. All rights reserved.

SPDX-License-Identifier: BUSL-1.1
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://mariadb.com/bsl11/

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codenotary/immudb/embedded/logger"
	"github.com/codenotary/immudb/pkg/database"
	"github.com/codenotary/immudb/pkg/replication"
	"github.com/codenotary/immudb/pkg/server/sessions"

	"github.com/codenotary/immudb/pkg/stream"

	"github.com/codenotary/immudb/pkg/auth"
)

const SystemDBName = "systemdb"
const DefaultDBName = "mydatabase"

// Options server options list
type Options struct {
	Dir                         string
	Network                     string
	Address                     string
	Port                        int
	Config                      string
	Pidfile                     string
	LogDir                      string
	Logfile                     string
	LogAccess                   bool
	LogRotationSize             int
	LogRotationAge              time.Duration
	AutoCert                    bool
	TLSConfig                   *tls.Config
	auth                        bool
	MaxRecvMsgSize              int
	MaxResultSize               int
	NoHistograms                bool
	Detached                    bool
	MetricsServer               bool
	MetricsServerPort           int
	WebServer                   bool
	WebServerPort               int
	DevMode                     bool
	AdminPassword               string `json:"-"`
	ForceAdminPassword          bool
	systemAdminDBName           string
	defaultDBName               string
	listener                    net.Listener
	usingCustomListener         bool
	maintenance                 bool
	SigningKey                  string
	synced                      bool
	RemoteStorageOptions        *RemoteStorageOptions
	StreamChunkSize             int
	TokenExpiryTimeMin          int
	PgsqlServer                 bool
	PgsqlServerPort             int
	ReplicationOptions          *ReplicationOptions
	SessionsOptions             *sessions.Options
	PProf                       bool
	LogFormat                   string
	GRPCReflectionServerEnabled bool
	SwaggerUIEnabled            bool
	LogRequestMetadata          bool
	MaxActiveDatabases          int
}

type RemoteStorageOptions struct {
	S3Storage             bool
	S3RoleEnabled         bool
	S3Role                string
	S3Endpoint            string
	S3AccessKeyID         string
	S3SecretKey           string `json:"-"`
	S3BucketName          string
	S3Location            string
	S3PathPrefix          string
	S3ExternalIdentifier  bool
	S3InstanceMetadataURL string
}

type ReplicationOptions struct {
	IsReplica                    bool
	SyncReplication              bool
	SyncAcks                     int    // only if !IsReplica && SyncReplication
	PrimaryHost                  string // only if IsReplica
	PrimaryPort                  int    // only if IsReplica
	PrimaryUsername              string // only if IsReplica
	PrimaryPassword              string // only if IsReplica
	PrefetchTxBufferSize         int    // only if IsReplica
	ReplicationCommitConcurrency int    // only if IsReplica
	AllowTxDiscarding            bool   // only if IsReplica
	SkipIntegrityCheck           bool   // only if IsReplica
	WaitForIndexing              bool   // only if IsReplica
}

// DefaultOptions returns default server options
func DefaultOptions() *Options {
	return &Options{
		Dir:                         "./data",
		Network:                     "tcp",
		Address:                     "0.0.0.0",
		Port:                        3322,
		Config:                      "configs/immudb.toml",
		Pidfile:                     "",
		Logfile:                     "",
		AutoCert:                    false,
		TLSConfig:                   nil,
		auth:                        true,
		MaxRecvMsgSize:              1024 * 1024 * 32, // 32Mb
		MaxResultSize:               database.MaxKeyScanLimit,
		NoHistograms:                false,
		Detached:                    false,
		MetricsServer:               true,
		MetricsServerPort:           9497,
		WebServer:                   true,
		WebServerPort:               8080,
		DevMode:                     false,
		AdminPassword:               auth.SysAdminPassword,
		ForceAdminPassword:          false,
		systemAdminDBName:           SystemDBName,
		defaultDBName:               DefaultDBName,
		usingCustomListener:         false,
		maintenance:                 false,
		synced:                      true,
		RemoteStorageOptions:        DefaultRemoteStorageOptions(),
		StreamChunkSize:             stream.DefaultChunkSize,
		TokenExpiryTimeMin:          1440,
		PgsqlServer:                 false,
		PgsqlServerPort:             5432,
		ReplicationOptions:          DefaultReplicationOptions(),
		SessionsOptions:             sessions.DefaultOptions(),
		PProf:                       false,
		GRPCReflectionServerEnabled: true,
		SwaggerUIEnabled:            true,
		LogRequestMetadata:          false,
		LogDir:                      "immulog",
		LogAccess:                   false,
		MaxActiveDatabases:          100,
	}
}

func DefaultRemoteStorageOptions() *RemoteStorageOptions {
	return &RemoteStorageOptions{
		S3Storage: false,
	}
}

func DefaultReplicationOptions() *ReplicationOptions {
	return &ReplicationOptions{
		IsReplica:                    false,
		SyncAcks:                     0,
		PrefetchTxBufferSize:         replication.DefaultPrefetchTxBufferSize,
		ReplicationCommitConcurrency: replication.DefaultReplicationCommitConcurrency,
	}
}

// WithDir sets dir
func (o *Options) WithDir(dir string) *Options {
	o.Dir = dir
	return o
}

// WithNetwork sets network
func (o *Options) WithNetwork(network string) *Options {
	o.Network = network
	return o
}

// WithAddress sets address
func (o *Options) WithAddress(address string) *Options {
	o.Address = address
	return o
}

// WithPort sets port
func (o *Options) WithPort(port int) *Options {
	o.Port = port
	return o
}

// WithConfig sets config file name
func (o *Options) WithConfig(config string) *Options {
	o.Config = config
	return o
}

// WithPidfile sets pid file
func (o *Options) WithPidfile(pidfile string) *Options {
	o.Pidfile = pidfile
	return o
}

// WithLogDir sets LogDir
func (o *Options) WithLogDir(dir string) *Options {
	o.LogDir = dir
	return o
}

// WithLogfile sets logfile
func (o *Options) WithLogfile(logfile string) *Options {
	o.Logfile = logfile
	return o
}

// WithLogRotationSize sets the log rotation size
func (o *Options) WithLogRotationSize(size int) *Options {
	o.LogRotationSize = size
	return o
}

// WithLogRotationAge sets the log rotation age
func (o *Options) WithLogRotationAge(age time.Duration) *Options {
	o.LogRotationAge = age
	return o
}

// WithLogAccess sets the log rotation age
func (o *Options) WithLogAccess(enabled bool) *Options {
	o.LogAccess = enabled
	return o
}

func (o *Options) WithLogFormat(logFormat string) *Options {
	o.LogFormat = logFormat
	return o
}

// WithTLS sets tls config
func (o *Options) WithTLS(tls *tls.Config) *Options {
	o.TLSConfig = tls
	return o
}

// WithAuth sets auth
// Deprecated: WithAuth will be removed in future release
func (o *Options) WithAuth(authEnabled bool) *Options {
	o.auth = authEnabled
	return o
}

func (o *Options) WithMaxRecvMsgSize(maxRecvMsgSize int) *Options {
	o.MaxRecvMsgSize = maxRecvMsgSize
	return o
}

// WithMaxResultSize sets the maximum number of results returned by any unary rpc method
func (o *Options) WithMaxResultSize(maxResultSize int) *Options {
	o.MaxResultSize = maxResultSize
	return o
}

// GetAuth gets auth
// Deprecated: GetAuth will be removed in future release
func (o *Options) GetAuth() bool {
	return o.auth
}

// WithNoHistograms disables collection of histograms metrics (e.g. query durations)
func (o *Options) WithNoHistograms(noHistograms bool) *Options {
	o.NoHistograms = noHistograms
	return o
}

// WithDetached sets immudb to be run in background
func (o *Options) WithDetached(detached bool) *Options {
	o.Detached = detached
	return o
}

// Bind returns bind address
func (o *Options) Bind() string {
	return o.Address + ":" + strconv.Itoa(o.Port)
}

// MetricsBind return metrics bind address
func (o *Options) MetricsBind() string {
	return o.Address + ":" + strconv.Itoa(o.MetricsServerPort)
}

// WebBind return bind address for the Web API/console
func (o *Options) WebBind() string {
	return o.Address + ":" + strconv.Itoa(o.WebServerPort)
}

// IsJSONLogger returns if the log format is json
func (o *Options) IsJSONLogger() bool {
	return o.LogFormat == logger.LogFormatJSON
}

// IsFileLogger returns if the log format is to a file
func (o *Options) IsFileLogger() bool {
	return o.Logfile != ""
}

// String print options
func (o *Options) String() string {
	rightPad := func(k string, v interface{}) string {
		return fmt.Sprintf("%-17s: %v", k, v)
	}
	opts := make([]string, 0, 17)
	opts = append(opts, "================ Config ================")
	opts = append(opts, rightPad("Data dir", o.Dir))
	opts = append(opts, rightPad("Address", fmt.Sprintf("%s:%d", o.Address, o.Port)))

	if o.MetricsServer {
		opts = append(opts, rightPad("Metrics address", fmt.Sprintf("%s:%d/metrics", o.Address, o.MetricsServerPort)))
		if o.PProf {
			opts = append(opts, rightPad("pprof enabled", "true"))
		}
	}

	repOpts := o.ReplicationOptions
	syncReplication := repOpts != nil && repOpts.SyncReplication
	isReplica := repOpts != nil && repOpts.IsReplica

	opts = append(opts, rightPad("Sync replication", syncReplication))

	if syncReplication && !isReplica {
		opts = append(opts, rightPad("Sync acks", repOpts.SyncAcks))
	}

	if isReplica {
		opts = append(opts, rightPad("Replica of", fmt.Sprintf("%s:%d", repOpts.PrimaryHost, repOpts.PrimaryPort)))
	}

	if o.Config != "" {
		opts = append(opts, rightPad("Config file", o.Config))
	}
	if o.Pidfile != "" {
		opts = append(opts, rightPad("PID file", o.Pidfile))
	}
	if o.Logfile != "" {
		opts = append(opts, rightPad("Log file", o.Logfile))
	}
	if o.LogFormat != "" {
		opts = append(opts, rightPad("Log format", o.LogFormat))
	}
	opts = append(opts, rightPad("Max recv msg size", o.MaxRecvMsgSize))
	opts = append(opts, rightPad("Auth enabled", o.auth))
	opts = append(opts, rightPad("Dev mode", o.DevMode))
	opts = append(opts, rightPad("Default database", o.defaultDBName))
	opts = append(opts, rightPad("Maintenance mode", o.maintenance))
	opts = append(opts, rightPad("Synced mode", o.synced))
	if o.SigningKey != "" {
		opts = append(opts, rightPad("Signing key", o.SigningKey))
	}
	if o.RemoteStorageOptions.S3Storage {
		opts = append(opts, "S3 storage")
		if o.RemoteStorageOptions.S3RoleEnabled {
			opts = append(opts, rightPad("   role auth", o.RemoteStorageOptions.S3RoleEnabled))
			opts = append(opts, rightPad("   role name", o.RemoteStorageOptions.S3Role))
		}
		opts = append(opts, rightPad("   endpoint", o.RemoteStorageOptions.S3Endpoint))
		opts = append(opts, rightPad("   bucket name", o.RemoteStorageOptions.S3BucketName))
		if o.RemoteStorageOptions.S3Location != "" {
			opts = append(opts, rightPad("   location", o.RemoteStorageOptions.S3Location))
		}
		opts = append(opts, rightPad("   prefix", o.RemoteStorageOptions.S3PathPrefix))
		opts = append(opts, rightPad("   external id", o.RemoteStorageOptions.S3ExternalIdentifier))
		opts = append(opts, rightPad("   metadata url", o.RemoteStorageOptions.S3InstanceMetadataURL))
	}
	if o.AdminPassword == auth.SysAdminPassword {
		opts = append(opts, "----------------------------------------")
		opts = append(opts, "Superadmin default credentials")
		opts = append(opts, rightPad("   Username", auth.SysAdminUsername))
		opts = append(opts, rightPad("   Password", auth.SysAdminPassword))
	}
	opts = append(opts, "========================================")
	return strings.Join(opts, "\n")
}

// WithMetricsServer ...
func (o *Options) WithMetricsServer(metricsServer bool) *Options {
	o.MetricsServer = metricsServer
	return o
}

// MetricsPort set Prometheus end-point port
func (o *Options) WithMetricsServerPort(port int) *Options {
	o.MetricsServerPort = port
	return o
}

// WithWebServer ...
func (o *Options) WithWebServer(webServer bool) *Options {
	o.WebServer = webServer
	return o
}

// WithWebServerPort ...
func (o *Options) WithWebServerPort(port int) *Options {
	o.WebServerPort = port
	return o
}

// WithDevMode ...
func (o *Options) WithDevMode(devMode bool) *Options {
	o.DevMode = devMode
	return o
}

// WithAdminPassword ...
func (o *Options) WithAdminPassword(adminPassword string) *Options {
	o.AdminPassword = adminPassword
	return o
}

// WithForceAdminPassword ...
func (o *Options) WithForceAdminPassword(forceAdminPassword bool) *Options {
	o.ForceAdminPassword = forceAdminPassword
	return o
}

// GetSystemAdminDBName returns the System database name
func (o *Options) GetSystemAdminDBName() string {
	return o.systemAdminDBName
}

// GetDefaultDBName returns the default database name
func (o *Options) GetDefaultDBName() string {
	return o.defaultDBName
}

// WithListener used usually to pass a bufered listener for testing purposes
func (o *Options) WithListener(lis net.Listener) *Options {
	o.listener = lis
	o.usingCustomListener = true
	return o
}

// WithMaintenance sets maintenance mode
func (o *Options) WithMaintenance(m bool) *Options {
	o.maintenance = m
	return o
}

// GetMaintenance gets maintenance mode
func (o *Options) GetMaintenance() bool {
	return o.maintenance
}

// WithSynced sets synced mode
func (o *Options) WithSynced(synced bool) *Options {
	o.synced = synced
	return o
}

// GetSynced gets synced mode
func (o *Options) GetSynced() bool {
	return o.synced
}

// WithSigningKey sets signature private key
func (o *Options) WithSigningKey(signingKey string) *Options {
	o.SigningKey = signingKey
	return o
}

// WithStreamChunkSize set the chunk size
func (o *Options) WithStreamChunkSize(streamChunkSize int) *Options {
	o.StreamChunkSize = streamChunkSize
	return o
}

// WithTokenExpiryTime set authentication token expiration time in minutes
func (o *Options) WithTokenExpiryTime(tokenExpiryTimeMin int) *Options {
	o.TokenExpiryTimeMin = tokenExpiryTimeMin
	return o
}

// PgsqlServerPort enable or disable pgsql server
func (o *Options) WithPgsqlServer(enable bool) *Options {
	o.PgsqlServer = enable
	return o
}

// PgsqlServerPort sets pgdsql server port
func (o *Options) WithPgsqlServerPort(port int) *Options {
	o.PgsqlServerPort = port
	return o
}

func (o *Options) WithRemoteStorageOptions(remoteStorageOptions *RemoteStorageOptions) *Options {
	o.RemoteStorageOptions = remoteStorageOptions
	return o
}

func (o *Options) WithReplicationOptions(replicationOptions *ReplicationOptions) *Options {
	o.ReplicationOptions = replicationOptions
	return o
}

func (o *Options) WithSessionOptions(options *sessions.Options) *Options {
	o.SessionsOptions = options
	return o
}

func (o *Options) WithPProf(pprof bool) *Options {
	o.PProf = pprof
	return o
}

func (o *Options) WithGRPCReflectionServerEnabled(enabled bool) *Options {
	o.GRPCReflectionServerEnabled = enabled
	return o
}

func (o *Options) WithSwaggerUIEnabled(enabled bool) *Options {
	o.SwaggerUIEnabled = enabled
	return o
}

func (o *Options) WithLogRequestMetadata(enabled bool) *Options {
	o.LogRequestMetadata = enabled
	return o
}

func (o *Options) WithMaxActiveDatabases(n int) *Options {
	o.MaxActiveDatabases = n
	return o
}

// RemoteStorageOptions

func (opts *RemoteStorageOptions) WithS3Storage(S3Storage bool) *RemoteStorageOptions {
	opts.S3Storage = S3Storage
	return opts
}

func (opts *RemoteStorageOptions) WithS3RoleEnabled(S3RoleEnabled bool) *RemoteStorageOptions {
	opts.S3RoleEnabled = S3RoleEnabled
	return opts
}

func (opts *RemoteStorageOptions) WithS3Role(S3Role string) *RemoteStorageOptions {
	opts.S3Role = S3Role
	return opts
}

func (opts *RemoteStorageOptions) WithS3Endpoint(s3Endpoint string) *RemoteStorageOptions {
	opts.S3Endpoint = s3Endpoint
	return opts
}

func (opts *RemoteStorageOptions) WithS3AccessKeyID(s3AccessKeyID string) *RemoteStorageOptions {
	opts.S3AccessKeyID = s3AccessKeyID
	return opts
}

func (opts *RemoteStorageOptions) WithS3SecretKey(s3SecretKey string) *RemoteStorageOptions {
	opts.S3SecretKey = s3SecretKey
	return opts
}

func (opts *RemoteStorageOptions) WithS3BucketName(s3BucketName string) *RemoteStorageOptions {
	opts.S3BucketName = s3BucketName
	return opts
}

func (opts *RemoteStorageOptions) WithS3Location(s3Location string) *RemoteStorageOptions {
	opts.S3Location = s3Location
	return opts
}

func (opts *RemoteStorageOptions) WithS3PathPrefix(s3PathPrefix string) *RemoteStorageOptions {
	opts.S3PathPrefix = s3PathPrefix
	return opts
}

func (opts *RemoteStorageOptions) WithS3ExternalIdentifier(s3ExternalIdentifier bool) *RemoteStorageOptions {
	opts.S3ExternalIdentifier = s3ExternalIdentifier
	return opts
}

func (opts *RemoteStorageOptions) WithS3InstanceMetadataURL(url string) *RemoteStorageOptions {
	opts.S3InstanceMetadataURL = url
	return opts
}

// ReplicationOptions

func (opts *ReplicationOptions) WithIsReplica(isReplica bool) *ReplicationOptions {
	opts.IsReplica = isReplica
	return opts
}

func (opts *ReplicationOptions) WithSyncReplication(syncReplication bool) *ReplicationOptions {
	opts.SyncReplication = syncReplication
	return opts
}

func (opts *ReplicationOptions) WithSyncAcks(syncAcks int) *ReplicationOptions {
	opts.SyncAcks = syncAcks
	return opts
}

func (opts *ReplicationOptions) WithPrimaryHost(primaryHost string) *ReplicationOptions {
	opts.PrimaryHost = primaryHost
	return opts
}

func (opts *ReplicationOptions) WithPrimaryPort(primaryPort int) *ReplicationOptions {
	opts.PrimaryPort = primaryPort
	return opts
}

func (opts *ReplicationOptions) WithPrimaryUsername(primaryUsername string) *ReplicationOptions {
	opts.PrimaryUsername = primaryUsername
	return opts
}

func (opts *ReplicationOptions) WithPrimaryPassword(primaryPassword string) *ReplicationOptions {
	opts.PrimaryPassword = primaryPassword
	return opts
}

func (opts *ReplicationOptions) WithPrefetchTxBufferSize(prefetchTxBufferSize int) *ReplicationOptions {
	opts.PrefetchTxBufferSize = prefetchTxBufferSize
	return opts
}

func (opts *ReplicationOptions) WithReplicationCommitConcurrency(replicationCommitConcurrency int) *ReplicationOptions {
	opts.ReplicationCommitConcurrency = replicationCommitConcurrency
	return opts
}

func (opts *ReplicationOptions) WithAllowTxDiscarding(allowTxDiscarding bool) *ReplicationOptions {
	opts.AllowTxDiscarding = allowTxDiscarding
	return opts
}

func (opts *ReplicationOptions) WithSkipIntegrityCheck(skipIntegrityCheck bool) *ReplicationOptions {
	opts.SkipIntegrityCheck = skipIntegrityCheck
	return opts
}

func (opts *ReplicationOptions) WithWaitForIndexing(waitForIndexingç bool) *ReplicationOptions {
	opts.WaitForIndexing = waitForIndexingç
	return opts
}
