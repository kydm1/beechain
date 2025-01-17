// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"syscall"
)

const ioctlReadTermios = syscall.TIOCGETA

// Termios functions describe a general terminal interface that is
// provided to control asynchronous communications ports.
// Go 1.2 doesn't include Termios for FreeBSD. This should be added in 1.3 and this could be merged with terminal_darwin.
type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}
