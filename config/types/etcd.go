// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"errors"
	"fmt"

	"github.com/coreos/go-semver/semver"
	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate/report"
)

var (
	EtcdVersionTooOld      = errors.New("Etcd version specified is not valid (too old)")
	EtcdMinorVersionTooNew = errors.New("Etcd minor version specified is too new, only options available in the previous minor version will be accepted")
	EtcdMajorVersionTooNew = errors.New("Etcd version is not valid (too new)")
	OldestEtcd             = *semver.New("2.3.0")
)

// Options can be the options for any Etcd version
type Options interface{}

type etcdCommon Etcd

type EtcdVersion semver.Version

func (e *EtcdVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := semver.Version(*e)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*e = EtcdVersion(t)
	return nil
}

func (e EtcdVersion) Validate() report.Report {
	v := semver.Version(e)
	switch {
	case v.LessThan(OldestEtcd):
		return report.ReportFromError(EtcdVersionTooOld, report.EntryError)
	case v.Major == 2 && v.Minor > 3:
		fallthrough
	case v.Major == 3 && v.Minor > 1:
		return report.ReportFromError(EtcdMinorVersionTooNew, report.EntryWarning)
	case v.Major > 3:
		return report.ReportFromError(EtcdMajorVersionTooNew, report.EntryError)
	}
	return report.Report{}
}

func (e EtcdVersion) String() string {
	return semver.Version(e).String()
}

// Etcd is a stub for yaml unmarshalling that figures out which
// of the other Etcd structs to use and unmarshals to that. Options needs
// to be an embedded type so that the structure of the yaml tree matches the
// structure of the go config tree
type Etcd struct {
	Version EtcdVersion `yaml:"version"`
	Options
}

func (etcd *Etcd) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := etcdCommon(*etcd)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*etcd = Etcd(t)

	version := semver.Version(etcd.Version)
	if version.Major == 2 && version.Minor >= 3 {
		o := Etcd2{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	} else if version.Major == 3 && version.Minor == 0 {
		o := Etcd3_0{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	} else if version.Major == 3 && version.Minor >= 1 {
		o := Etcd3_1{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		etcd.Options = o
	}
	return nil
}

func init() {
	register2_0(func(in Config, out ignTypes.Config) (ignTypes.Config, report.Report) {
		if in.Etcd != nil {
			out.Systemd.Units = append(out.Systemd.Units, ignTypes.SystemdUnit{
				Name:   "etcd-member.service",
				Enable: true,
				DropIns: []ignTypes.SystemdUnitDropIn{{
					Name:     "20-clct-etcd-member.conf",
					Contents: etcdContents(*in.Etcd),
				}},
			})
		}
		return out, report.Report{}
	})
}

// etcdContents creates the string containing the systemd drop in for etcd-member
func etcdContents(etcd Etcd) string {
	vars := getEnvVars(etcd.Options)
	// Add the tag
	vars = append(vars, fmt.Sprintf("ETCD_IMAGE_TAG=v%s", etcd.Version))
	return serviceContentsFromEnvVars(vars)
}

type Etcd3_0 struct {
	Name                     string `yaml:"name"                        env:"ETCD_NAME"`
	DataDir                  string `yaml:"data_dir"                    env:"ETCD_DATA_DIR"`
	WalDir                   string `yaml:"wal_dir"                     env:"ETCD_WAL_DIR"`
	SnapshotCount            int    `yaml:"snapshot_count"              env:"ETCD_SNAPSHOT_COUNT"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"          env:"ETCD_HEARTBEAT_INTERVAL"`
	ElectionTimeout          int    `yaml:"election_timeout"            env:"ETCD_ELECTION_TIMEOUT"`
	ListenPeerUrls           string `yaml:"listen_peer_urls"            env:"ETCD_LISTEN_PEER_URLS"`
	ListenClientUrls         string `yaml:"listen_client_urls"          env:"ETCD_LISTEN_CLIENT_URLS"`
	MaxSnapshots             int    `yaml:"max_snapshots"               env:"ETCD_MAX_SNAPSHOTS"`
	MaxWals                  int    `yaml:"max_wals"                    env:"ETCD_MAX_WALS"`
	Cors                     string `yaml:"cors"                        env:"ETCD_CORS"`
	InitialAdvertisePeerUrls string `yaml:"initial_advertise_peer_urls" env:"ETCD_INITIAL_ADVERTISE_PEER_URLS"`
	InitialCluster           string `yaml:"initial_cluster"             env:"ETCD_INITIAL_CLUSTER"`
	InitialClusterState      string `yaml:"initial_cluster_state"       env:"ETCD_INITIAL_CLUSTER_STATE"`
	InitialClusterToken      string `yaml:"initial_cluster_token"       env:"ETCD_INITIAL_CLUSTER_TOKEN"`
	AdvertiseClientUrls      string `yaml:"advertise_client_urls"       env:"ETCD_ADVERTISE_CLIENT_URLS"`
	Discovery                string `yaml:"discovery"                   env:"ETCD_DISCOVERY"`
	DiscoverySrv             string `yaml:"discovery_srv"               env:"ETCD_DISCOVERY_SRV"`
	DiscoveryFallback        string `yaml:"discovery_fallback"          env:"ETCD_DISCOVERY_FALLBACK"`
	DiscoveryProxy           string `yaml:"discovery_proxy"             env:"ETCD_DISCOVERY_PROXY"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"       env:"ETCD_STRICT_RECONFIG_CHECK"`
	AutoCompactionRetention  int    `yaml:"auto_compaction_retention"   env:"ETCD_AUTO_COMPACTION_RETENTION"`
	Proxy                    string `yaml:"proxy"                       env:"ETCD_PROXY"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"          env:"ETCD_PROXY_FAILURE_WAIT"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"      env:"ETCD_PROXY_REFRESH_INTERVAL"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"          env:"ETCD_PROXY_DIAL_TIMEOUT"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"         env:"ETCD_PROXY_WRITE_TIMEOUT"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"          env:"ETCD_PROXY_READ_TIMEOUT"`
	CaFile                   string `yaml:"ca_file"                     env:"ETCD_CA_FILE"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                   env:"ETCD_CERT_FILE"`
	KeyFile                  string `yaml:"key_file"                    env:"ETCD_KEY_FILE"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"            env:"ETCD_CLIENT_CERT_AUTH"`
	TrustedCaFile            string `yaml:"trusted_ca_file"             env:"ETCD_TRUSTED_CA_FILE"`
	AutoTls                  bool   `yaml:"auto_tls"                    env:"ETCD_AUTO_TLS"`
	PeerCaFile               string `yaml:"peer_ca_file"                env:"ETCD_PEER_CA_FILE"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"              env:"ETCD_PEER_CERT_FILE"`
	PeerKeyFile              string `yaml:"peer_key_file"               env:"ETCD_PEER_KEY_FILE"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"       env:"ETCD_PEER_CLIENT_CERT_AUTH"`
	PeerTrustedCaFile        string `yaml:"peer_trusted_ca_file"        env:"ETCD_PEER_TRUSTED_CA_FILE"`
	PeerAutoTls              bool   `yaml:"peer_auto_tls"               env:"ETCD_PEER_AUTO_TLS"`
	Debug                    bool   `yaml:"debug"                       env:"ETCD_DEBUG"`
	LogPackageLevels         string `yaml:"log_package_levels"          env:"ETCD_LOG_PACKAGE_LEVELS"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"           env:"ETCD_FORCE_NEW_CLUSTER"`
}

type Etcd3_1 struct {
	Name                     string `yaml:"name"                        env:"ETCD_NAME"`
	DataDir                  string `yaml:"data_dir"                    env:"ETCD_DATA_DIR"`
	WalDir                   string `yaml:"wal_dir"                     env:"ETCD_WAL_DIR"`
	SnapshotCount            int    `yaml:"snapshot_count"              env:"ETCD_SNAPSHOT_COUNT"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"          env:"ETCD_HEARTBEAT_INTERVAL"`
	ElectionTimeout          int    `yaml:"election_timeout"            env:"ETCD_ELECTION_TIMEOUT"`
	ListenPeerUrls           string `yaml:"listen_peer_urls"            env:"ETCD_LISTEN_PEER_URLS"`
	ListenClientUrls         string `yaml:"listen_client_urls"          env:"ETCD_LISTEN_CLIENT_URLS"`
	MaxSnapshots             int    `yaml:"max_snapshots"               env:"ETCD_MAX_SNAPSHOTS"`
	MaxWals                  int    `yaml:"max_wals"                    env:"ETCD_MAX_WALS"`
	Cors                     string `yaml:"cors"                        env:"ETCD_CORS"`
	InitialAdvertisePeerUrls string `yaml:"initial_advertise_peer_urls" env:"ETCD_INITIAL_ADVERTISE_PEER_URLS"`
	InitialCluster           string `yaml:"initial_cluster"             env:"ETCD_INITIAL_CLUSTER"`
	InitialClusterState      string `yaml:"initial_cluster_state"       env:"ETCD_INITIAL_CLUSTER_STATE"`
	InitialClusterToken      string `yaml:"initial_cluster_token"       env:"ETCD_INITIAL_CLUSTER_TOKEN"`
	AdvertiseClientUrls      string `yaml:"advertise_client_urls"       env:"ETCD_ADVERTISE_CLIENT_URLS"`
	Discovery                string `yaml:"discovery"                   env:"ETCD_DISCOVERY"`
	DiscoverySrv             string `yaml:"discovery_srv"               env:"ETCD_DISCOVERY_SRV"`
	DiscoveryFallback        string `yaml:"discovery_fallback"          env:"ETCD_DISCOVERY_FALLBACK"`
	DiscoveryProxy           string `yaml:"discovery_proxy"             env:"ETCD_DISCOVERY_PROXY"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"       env:"ETCD_STRICT_RECONFIG_CHECK"`
	AutoCompactionRetention  int    `yaml:"auto_compaction_retention"   env:"ETCD_AUTO_COMPACTION_RETENTION"`
	Proxy                    string `yaml:"proxy"                       env:"ETCD_PROXY"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"          env:"ETCD_PROXY_FAILURE_WAIT"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"      env:"ETCD_PROXY_REFRESH_INTERVAL"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"          env:"ETCD_PROXY_DIAL_TIMEOUT"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"         env:"ETCD_PROXY_WRITE_TIMEOUT"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"          env:"ETCD_PROXY_READ_TIMEOUT"`
	CaFile                   string `yaml:"ca_file"                     env:"ETCD_CA_FILE"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                   env:"ETCD_CERT_FILE"`
	KeyFile                  string `yaml:"key_file"                    env:"ETCD_KEY_FILE"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"            env:"ETCD_CLIENT_CERT_AUTH"`
	TrustedCaFile            string `yaml:"trusted_ca_file"             env:"ETCD_TRUSTED_CA_FILE"`
	AutoTls                  bool   `yaml:"auto_tls"                    env:"ETCD_AUTO_TLS"`
	PeerCaFile               string `yaml:"peer_ca_file"                env:"ETCD_PEER_CA_FILE"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"              env:"ETCD_PEER_CERT_FILE"`
	PeerKeyFile              string `yaml:"peer_key_file"               env:"ETCD_PEER_KEY_FILE"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"       env:"ETCD_PEER_CLIENT_CERT_AUTH"`
	PeerTrustedCaFile        string `yaml:"peer_trusted_ca_file"        env:"ETCD_PEER_TRUSTED_CA_FILE"`
	PeerAutoTls              bool   `yaml:"peer_auto_tls"               env:"ETCD_PEER_AUTO_TLS"`
	Debug                    bool   `yaml:"debug"                       env:"ETCD_DEBUG"`
	LogPackageLevels         string `yaml:"log_package_levels"          env:"ETCD_LOG_PACKAGE_LEVELS"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"           env:"ETCD_FORCE_NEW_CLUSTER"`
	Metrics                  string `yaml:"metrics"                     env:"ETCD_METRICS"`
	LogOutput                string `yaml:"log_output"                  env:"ETCD_LOG_OUTPUT"`
}

type Etcd2 struct {
	AdvertiseClientURLs      string `yaml:"advertise_client_urls"         env:"ETCD_ADVERTISE_CLIENT_URLS"`
	CAFile                   string `yaml:"ca_file"                       env:"ETCD_CA_FILE"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file"                     env:"ETCD_CERT_FILE"`
	ClientCertAuth           bool   `yaml:"client_cert_auth"              env:"ETCD_CLIENT_CERT_AUTH"`
	CorsOrigins              string `yaml:"cors"                          env:"ETCD_CORS"`
	DataDir                  string `yaml:"data_dir"                      env:"ETCD_DATA_DIR"`
	Debug                    bool   `yaml:"debug"                         env:"ETCD_DEBUG"`
	Discovery                string `yaml:"discovery"                     env:"ETCD_DISCOVERY"`
	DiscoveryFallback        string `yaml:"discovery_fallback"            env:"ETCD_DISCOVERY_FALLBACK"`
	DiscoverySRV             string `yaml:"discovery_srv"                 env:"ETCD_DISCOVERY_SRV"`
	DiscoveryProxy           string `yaml:"discovery_proxy"               env:"ETCD_DISCOVERY_PROXY"`
	ElectionTimeout          int    `yaml:"election_timeout"              env:"ETCD_ELECTION_TIMEOUT"`
	EnablePprof              bool   `yaml:"enable_pprof"                  env:"ETCD_ENABLE_PPROF"`
	ForceNewCluster          bool   `yaml:"force_new_cluster"             env:"ETCD_FORCE_NEW_CLUSTER"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval"            env:"ETCD_HEARTBEAT_INTERVAL"`
	InitialAdvertisePeerURLs string `yaml:"initial_advertise_peer_urls"   env:"ETCD_INITIAL_ADVERTISE_PEER_URLS"`
	InitialCluster           string `yaml:"initial_cluster"               env:"ETCD_INITIAL_CLUSTER"`
	InitialClusterState      string `yaml:"initial_cluster_state"         env:"ETCD_INITIAL_CLUSTER_STATE"`
	InitialClusterToken      string `yaml:"initial_cluster_token"         env:"ETCD_INITIAL_CLUSTER_TOKEN"`
	KeyFile                  string `yaml:"key_file"                      env:"ETCD_KEY_FILE"`
	ListenClientURLs         string `yaml:"listen_client_urls"            env:"ETCD_LISTEN_CLIENT_URLS"`
	ListenPeerURLs           string `yaml:"listen_peer_urls"              env:"ETCD_LISTEN_PEER_URLS"`
	LogPackageLevels         string `yaml:"log_package_levels"            env:"ETCD_LOG_PACKAGE_LEVELS"`
	MaxSnapshots             int    `yaml:"max_snapshots"                 env:"ETCD_MAX_SNAPSHOTS"`
	MaxWALs                  int    `yaml:"max_wals"                      env:"ETCD_MAX_WALS"`
	Name                     string `yaml:"name"                          env:"ETCD_NAME"`
	PeerCAFile               string `yaml:"peer_ca_file"                  env:"ETCD_PEER_CA_FILE"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file"                env:"ETCD_PEER_CERT_FILE"`
	PeerKeyFile              string `yaml:"peer_key_file"                 env:"ETCD_PEER_KEY_FILE"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth"         env:"ETCD_PEER_CLIENT_CERT_AUTH"`
	PeerTrustedCAFile        string `yaml:"peer_trusted_ca_file"          env:"ETCD_PEER_TRUSTED_CA_FILE"`
	Proxy                    string `yaml:"proxy"                         env:"ETCD_PROXY"                       valid:"^(on|off|readonly)$"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout"            env:"ETCD_PROXY_DIAL_TIMEOUT"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait"            env:"ETCD_PROXY_FAILURE_WAIT"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout"            env:"ETCD_PROXY_READ_TIMEOUT"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval"        env:"ETCD_PROXY_REFRESH_INTERVAL"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout"           env:"ETCD_PROXY_WRITE_TIMEOUT"`
	SnapshotCount            int    `yaml:"snapshot_count"                env:"ETCD_SNAPSHOT_COUNT"`
	StrictReconfigCheck      bool   `yaml:"strict_reconfig_check"         env:"ETCD_STRICT_RECONFIG_CHECK"`
	TrustedCAFile            string `yaml:"trusted_ca_file"               env:"ETCD_TRUSTED_CA_FILE"`
	WalDir                   string `yaml:"wal_dir"                       env:"ETCD_WAL_DIR"`
}
