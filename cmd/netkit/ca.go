package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/biancarosa/netkit/internal/proxy"
)

func runCA(args []string) error {
	flags := flag.NewFlagSet("ca", flag.ContinueOnError)
	certPath := flags.String("cert", "netkit-ca.pem", "CA certificate output (share with clients)")
	keyPath := flags.String("key", "netkit-ca-key.pem", "Private key output (keep on proxy)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	cert, key, err := proxy.GenerateCA()
	if err != nil {
		return err
	}
	// Exclusive creation prevents accidental replacement of an already trusted CA.
	keyFile, err := os.OpenFile(*keyPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	if _, err = keyFile.Write(key); err != nil {
		_ = keyFile.Close()
		return err
	}
	if err = keyFile.Close(); err != nil {
		return err
	}
	certFile, err := os.OpenFile(*certPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		_ = os.Remove(*keyPath)
		return err
	}
	if _, err = certFile.Write(cert); err != nil {
		_ = certFile.Close()
		return err
	}
	if err = certFile.Close(); err != nil {
		return err
	}
	fmt.Printf("CA certificate: %s\nPrivate key: %s\nTrust only the certificate in proxy clients.\n", *certPath, *keyPath)
	return nil
}
