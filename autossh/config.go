package autossh

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"

	"toolkit/logs"
	"toolkit/system"
)

var c conf

type conf struct {
	Tunnel []*struct {
		Service   string `json:"service"`
		ListenOn  string `json:"listenOn"`
		ListenAt  string `json:"listenAt"`
		SshAlias  string `json:"sshAlias"`
		ForwardTo string `json:"forwardTo"`
	} `json:"tunnel"`
	SshConfig []*struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
		User    string `json:"user"`
		Auth    *struct {
			Pass string   `json:"pass"`
			Keys []string `json:"keys"`
		} `json:"auth"`
	} `json:"sshConfig"`
}

const envAutoSshConfig = "TOOLKIT_AUTOSSH_CONF_PATH"

func ReadConfInFile() error {
	file := system.GetEnvString(envAutoSshConfig)
	bs, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsoncIgnoreComments(bs), &c)
}

const (
	listenUnknown  listenType = "unknown"
	listenOnLocal  listenType = "local"
	listenOnRemote listenType = "remote"
)

type listenType string

func (t listenType) string() string {
	return string(t)
}

var listenTypeMapping = map[string]listenType{
	string(listenOnLocal):  listenOnLocal,
	string(listenOnRemote): listenOnRemote,
}

func getListenType(s string) listenType {
	t, ok := listenTypeMapping[s]
	if !ok {
		t = listenUnknown
	}
	return t
}

type sshConfig struct {
	Alias   string
	Address string
	Config  *ssh.ClientConfig
}

func sshAuthKeys(keys []string) []ssh.Signer {
	s := make([]ssh.Signer, 0, len(keys))
	for _, k := range keys {
		path := filepath.Clean(k)
		isAbs := filepath.IsAbs(path)
		if !isAbs {
			home, err := os.UserHomeDir()
			if err != nil {
				logs.Errorf("user home dir, err: %s", err.Error())
				continue
			}
			path = filepath.Join(home, ".ssh", path)
		}

		bs, err := os.ReadFile(path)
		if err != nil {
			logs.Errorf("read: %s, err: %s", path, err.Error())
			continue
		}
		signer, err := ssh.ParsePrivateKey(bs)
		if err != nil {
			logs.Errorf("ssh parse pk, err: %s", err.Error())
			continue
		}
		s = append(s, signer)
	}
	return s
}

func sshConfigMapping() map[string]*sshConfig {
	m := make(map[string]*sshConfig, len(c.SshConfig))
	for _, k := range c.SshConfig {
		if k == nil || k.Auth == nil {
			continue
		}
		pass := k.Auth.Pass
		keys := sshAuthKeys(k.Auth.Keys)
		if len(keys) < 1 && len(pass) < 1 {
			continue
		}

		m[k.Alias] = &sshConfig{
			Alias:   k.Alias,
			Address: k.Address,
			Config: &ssh.ClientConfig{
				User: k.User,
				Auth: []ssh.AuthMethod{
					ssh.Password(pass),
					ssh.PublicKeys(keys...),
				},
				Timeout:         10 * time.Second,
				BannerCallback:  func(string) (e error) { return },
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}
	}
	return m
}

type tunnel struct {
	Service   string
	ListenOn  listenType
	ListenAt  string
	ForwardTo string
	SshConfig *sshConfig
}

func tunnels() []*tunnel {
	sshConfigMap := sshConfigMapping()
	tuns := make([]*tunnel, 0, len(c.Tunnel))
	for _, k := range c.Tunnel {
		listenOn := getListenType(k.ListenOn)
		if listenOn == listenUnknown {
			continue
		}
		sshConfig, ok := sshConfigMap[k.SshAlias]
		if !ok || sshConfig == nil {
			continue
		}

		tuns = append(tuns,
			&tunnel{
				Service:   k.Service,
				ListenOn:  listenOn,
				ListenAt:  k.ListenAt,
				ForwardTo: k.ForwardTo,
				SshConfig: sshConfig,
			},
		)
	}
	return tuns
}
