package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"

	"github.com/costinm/cert-ssh/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func Conf(keyn, def string) string {
	val := os.Getenv(keyn)
	if val == "" {
		val = def
	}
	return val
}

func main() {
	sshc := &ssh.Client{
		SSHCa: Conf("SSH_CA", "127.0.0.1:14023"),
		SSHD: Conf("SSH_GATE", "127.0.0.1:14022"),
		Namespace: Conf("WORKLOAD_NAMESPACE", "default"),

	}
	sshc.CertProvider = func(ctx context.Context, sshCA string) (gossh.Signer, error) {
		ephemeralPrivate, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		ephemeralSigner, _ := gossh.NewSignerFromKey(ephemeralPrivate)
		return ephemeralSigner, nil
	}
	// run-k8s helper can't start a debug ssh server if running ssh_signer -
	// no signer. Start one in-process, for debugging.
	err := sshc.Start()
	if err != nil {
		panic(err)
	}
}
