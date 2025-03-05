package atecc

import (
	"time"

	"periph.io/x/conn/v3/i2c"
)

// IfaceConfig is the configuration object for a device.
//
// Logical device configurations describe the device type and logical
// interface.
type IfaceConfig struct {
	// DeviceType affects how communication with the device is done.
	DeviceType DeviceType
	// I2C contains I²C specific configuration.
	I2C I2CConfig
	// WakeDelay defines the time to wait for the device before waking up.
	//
	// This represents the tWHI + tWLO and is configured based on device type.
	WakeDelay time.Duration
	// RxRetries is the number of retries to attempt when receiving data.
	RxRetries int
	// Debug is used for debug output.
	Debug Logger
}

type I2CConfig struct {
	Address uint16
	Bus     i2c.Bus
}

// ConfigATECCX08A_I2CDefault returns a default config for an ECCx08A device.
//
// TODO: re-think where we put bus, who owns it (who closes, do we have Close?)
func ConfigATECCX08A_I2CDefault(bus i2c.Bus) IfaceConfig {
	return IfaceConfig{
		DeviceType: DeviceATECC608,
		WakeDelay:  1500 * time.Microsecond,
		RxRetries:  20,
		I2C: I2CConfig{
			Address: 0x60,
			Bus:     bus,
		},
	}
}
