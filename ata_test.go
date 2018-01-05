// Copyright 2017-18 Daniel Swarbrick. All rights reserved.
// Use of this source code is governed by a GPL license that can be found in the LICENSE file.

package smart

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// ATA IDENTIFY DEVICE (command ECh)
var ataIdentifyData = [512]byte{
	0x40, 0x00, 0xff, 0x3f, 0x37, 0xc8, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3f, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x31, 0x53, 0x4d, 0x44, 0x45, 0x4e, 0x44, 0x41, 0x32, 0x31, 0x34, 0x33,
	0x36, 0x35, 0x20, 0x42, 0x20, 0x20, 0x20, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x58, 0x45,
	0x30, 0x54, 0x42, 0x44, 0x51, 0x36, 0x61, 0x53, 0x73, 0x6d, 0x6e, 0x75, 0x20, 0x67, 0x53, 0x53,
	0x20, 0x44, 0x34, 0x38, 0x20, 0x30, 0x56, 0x45, 0x20, 0x4f, 0x35, 0x37, 0x47, 0x30, 0x20, 0x42,
	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x01, 0x80,
	0x01, 0x40, 0x00, 0x2f, 0x00, 0x40, 0x00, 0x02, 0x00, 0x02, 0x07, 0x00, 0xf0, 0xff, 0x01, 0x00,
	0x3f, 0x00, 0x10, 0xfc, 0x3e, 0x00, 0x01, 0x01, 0xff, 0xff, 0xff, 0x0f, 0x00, 0x00, 0x07, 0x00,
	0x03, 0x00, 0x78, 0x00, 0x78, 0x00, 0x78, 0x00, 0x78, 0x00, 0x10, 0x0f, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1f, 0x00, 0x0e, 0x85, 0x46, 0x00, 0x6c, 0x00, 0x64, 0x00,
	0xfc, 0x03, 0x39, 0x00, 0x6b, 0x74, 0x01, 0x7d, 0x63, 0x41, 0x69, 0x74, 0x01, 0xbc, 0x63, 0x41,
	0x7f, 0x40, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0xfe, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x66, 0x54, 0x57, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x08, 0x00, 0x00, 0x40, 0x00, 0x00, 0x02, 0x50, 0x88, 0x53, 0x09, 0x50, 0x7f, 0x39,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x40,
	0x1c, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x29, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3d, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x10, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa5, 0x7f,
}

func TestATAIdentify(t *testing.T) {
	var d IdentifyDeviceData

	assert := assert.New(t)

	assert.Equal(uintptr(512), unsafe.Sizeof(d))
	binary.Read(bytes.NewBuffer(ataIdentifyData[:]), nativeEndian, &d)

	swapBytes(d.SerialNumber[:])
	swapBytes(d.FirmwareRevision[:])
	swapBytes(d.ModelNumber[:])

	assert.Equal("S1DMNEAD123456B     ", string(d.SerialNumber[:]))
	assert.Equal("EXT0DB6Q", string(d.FirmwareRevision[:]))
	assert.Equal("Samsung SSD 840 EVO 750GB               ", string(d.ModelNumber[:]))
	assert.Equal("5 002538 85009397f", d.getWWN())

	assert.Equal(uint16(1), d.RotationRate)
}
