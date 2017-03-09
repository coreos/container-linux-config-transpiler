package types

import (
	"errors"

	"github.com/coreos/go-semver/semver"
	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrFlannelTooOld      = errors.New("invalid flannel version (too old)")
	ErrFlannelMinorTooNew = errors.New("flannel minor version too new. Only options available in the previous minor version will be supported")
	OldestFlannelVersion  = *semver.New("0.5.0")
)

type Flannel struct {
	Version FlannelVersion `yaml:"version"`
	Options
}

type flannelCommon Flannel

type FlannelVersion semver.Version

func (v *FlannelVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := semver.Version(*v)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*v = FlannelVersion(t)
	return nil
}

func (fv FlannelVersion) Validate() report.Report {
	v := semver.Version(fv)
	switch {
	case v.LessThan(OldestFlannelVersion):
		return report.ReportFromError(ErrFlannelTooOld, report.EntryError)
	case v.Major == 0 && fv.Minor > 7:
		return report.ReportFromError(ErrFlannelMinorTooNew, report.EntryWarning)
	}
	return report.Report{}
}

func (fv FlannelVersion) String() string {
	return semver.Version(fv).String()
}

func (flannel *Flannel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	t := flannelCommon(*flannel)
	if err := unmarshal(&t); err != nil {
		return err
	}
	*flannel = Flannel(t)

	v := semver.Version(flannel.Version)
	if v.Major == 0 && v.Minor >= 7 {
		o := Flannel0_7{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	} else if v.Major == 0 && v.Minor == 6 {
		o := Flannel0_6{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	} else if v.Major == 0 && v.Minor == 5 {
		o := Flannel0_5{}
		if err := unmarshal(&o); err != nil {
			return err
		}
		flannel.Options = o
	}
	return nil
}

// Flannel0_7 represents flannel options for version 0.7.x. Don't embed Flannel0_6 because
// the yaml parser doesn't handle embedded structs
type Flannel0_7 struct {
	EtcdUsername  string `yaml:"etcd_username"   env:"FLANNELD_ETCD_USERNAME"`
	EtcdPassword  string `yaml:"etcd_password"   env:"FLANNELD_ETCD_PASSWORD"`
	EtcdEndpoints string `yaml:"etcd_endpoints"  env:"FLANNELD_ETCD_ENDPOINTS"`
	EtcdCAFile    string `yaml:"etcd_cafile"     env:"FLANNELD_ETCD_CAFILE"`
	EtcdCertFile  string `yaml:"etcd_certfile"   env:"FLANNELD_ETCD_CERTFILE"`
	EtcdKeyFile   string `yaml:"etcd_keyfile"    env:"FLANNELD_ETCD_KEYFILE"`
	EtcdPrefix    string `yaml:"etcd_prefix"     env:"FLANNELD_ETCD_PREFIX"`
	IPMasq        string `yaml:"ip_masq"         env:"FLANNELD_IP_MASQ"`
	SubnetFile    string `yaml:"subnet_file"     env:"FLANNELD_SUBNET_FILE"`
	Iface         string `yaml:"interface"       env:"FLANNELD_IFACE"`
	PublicIP      string `yaml:"public_ip"       env:"FLANNELD_PUBLIC_IP"`
	KubeSubnetMgr bool   `yaml:"kube_subnet_mgr" env:"FLANNEL_KUBE_SUBNET_MGR"`
}

type Flannel0_6 struct {
	EtcdUsername  string `yaml:"etcd_username"  env:"FLANNELD_ETCD_USERNAME"`
	EtcdPassword  string `yaml:"etcd_password"  env:"FLANNELD_ETCD_PASSWORD"`
	EtcdEndpoints string `yaml:"etcd_endpoints" env:"FLANNELD_ETCD_ENDPOINTS"`
	EtcdCAFile    string `yaml:"etcd_cafile"    env:"FLANNELD_ETCD_CAFILE"`
	EtcdCertFile  string `yaml:"etcd_certfile"  env:"FLANNELD_ETCD_CERTFILE"`
	EtcdKeyFile   string `yaml:"etcd_keyfile"   env:"FLANNELD_ETCD_KEYFILE"`
	EtcdPrefix    string `yaml:"etcd_prefix"    env:"FLANNELD_ETCD_PREFIX"`
	IPMasq        string `yaml:"ip_masq"        env:"FLANNELD_IP_MASQ"`
	SubnetFile    string `yaml:"subnet_file"    env:"FLANNELD_SUBNET_FILE"`
	Iface         string `yaml:"interface"      env:"FLANNELD_IFACE"`
	PublicIP      string `yaml:"public_ip"      env:"FLANNELD_PUBLIC_IP"`
}

type Flannel0_5 struct {
	EtcdEndpoints string `yaml:"etcd_endpoints" env:"FLANNELD_ETCD_ENDPOINTS"`
	EtcdCAFile    string `yaml:"etcd_cafile"    env:"FLANNELD_ETCD_CAFILE"`
	EtcdCertFile  string `yaml:"etcd_certfile"  env:"FLANNELD_ETCD_CERTFILE"`
	EtcdKeyFile   string `yaml:"etcd_keyfile"   env:"FLANNELD_ETCD_KEYFILE"`
	EtcdPrefix    string `yaml:"etcd_prefix"    env:"FLANNELD_ETCD_PREFIX"`
	IPMasq        string `yaml:"ip_masq"        env:"FLANNELD_IP_MASQ"`
	SubnetFile    string `yaml:"subnet_file"    env:"FLANNELD_SUBNET_FILE"`
	Iface         string `yaml:"interface"      env:"FLANNELD_IFACE"`
	PublicIP      string `yaml:"public_ip"      env:"FLANNELD_PUBLIC_IP"`
}
