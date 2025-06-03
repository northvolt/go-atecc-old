package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Scania-Goldcup/go-atecc/pkg/atecc"
	"github.com/peterbourgon/ff/v3/ffcli"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	defaultI2CAddress     = 0x60
	defaultDeviceIdentity = 0
)

func newATECC(ctx context.Context, c *rootConfig) (*atecc.Dev, io.Closer, error) {
	i2cAddress, err := getI2CAddress(c.addr, c.trustPlatformFormat)
	if err != nil {
		return nil, nil, err
	}

	if _, err = host.Init(); err != nil {
		return nil, nil, err
	}
	bus, err := i2creg.Open(strconv.Itoa(c.bus))
	if err != nil {
		return nil, nil, fmt.Errorf("atecc: failed to connect to bus: %w", err)
	}

	cfg := atecc.ConfigATECCX08A_I2CDefault(bus)
	cfg.Debug = newLogger(c.verbose)
	cfg.I2C.Address = i2cAddress
	d, err := atecc.NewI2CDev(ctx, cfg)
	return d, bus, err
}

func getI2CAddress(addrStr string, trustPlatformFormat bool) (uint16, error) {
	if addrStr == "" {
		return defaultI2CAddress, nil
	}
	addr, err := strconv.ParseUint(strings.TrimPrefix(addrStr, "0x"), 16, 16)
	if err != nil {
		return 0, err
	}

	if trustPlatformFormat {
		return uint16(addr >> 1), nil
	} else {
		return uint16(addr), nil
	}
}

func prettyHex(data []byte) string {
	return prettyHexIndent(data, "    ", "")
}

func prettyHexIndent(data []byte, prefix string, space string) string {
	var buf strings.Builder

	// prefix and space every 16 byte, and 2 hex, and one space/newline
	cols := 16
	size := (len(data)/cols+1)*(len(prefix)+len(space)+1) + len(data)*3
	buf.Grow(size)

	for i := range data {
		if i > 0 {
			switch i % cols {
			case 0:
				buf.WriteByte('\n')
			case cols / 2:
				buf.WriteByte(' ')
				buf.WriteString(space)
			default:
				buf.WriteByte(' ')
			}
		}
		if i%cols == 0 {
			buf.WriteString(prefix)
		}

		buf.WriteString(fmt.Sprintf("%02X", data[i:i+1]))
	}

	return buf.String()
}

func pemEncodePublicKey(pk crypto.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})), nil
}

func addLongHelp(cmd *ffcli.Command) *ffcli.Command {
	if cmd.LongHelp == "" {
		cmd.LongHelp = cmd.ShortHelp
	}

	cmd.LongHelp += ateccLongHelp

	return cmd
}

func newLogger(verbose bool) atecc.Logger {
	if verbose {
		return log.New(os.Stderr, "", 0)
	} else {
		return nil
	}
}
