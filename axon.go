// Copyright 2018 Adam Robinson
// Permission is hereby granted, free of charge, to any person obtaining a copy of 
// this software and associated documentation files (the "Software"), to deal in the 
// Software without restriction, including without limitation the rights to use, copy, 
// modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, 
// and to permit persons to whom the Software is furnished to do so, subject to the 
// following conditions:
// The above copyright notice and this permission notice shall be included in all 
// copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A 
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT 
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF 
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE 
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// #include <malloc.h>

package axon


// #cgo CFLAGS: -fPIC -lm -lpthread -lrt
// #include <stdio.h>
// #include <stdlib.h>
// #include <semaphore.h>

// #include <fcntl.h>
// #include <string.h>
// #include <errno.h>
// #include <sys/stat.h>
// #include <sys/types.h>
// #include <sys/mman.h>
// #include <unistd.h>
// #include <stdbool.h>
// #include "axon.h"
import "C"

type Axon struct {
	c_Axon * C.struct_Axon
}

// Creates an instance of a controller object. The primary fields are stored in
// and managed by the C code.
func AxonController(name string) (Axon, error) {
	// This is intentionally not freed because the 
	// c code takes ownership of it.
	c_str_name := C.CString(name)
	c_axon_ptr := C.createControllerInstance(c_str_name)

	return Axon{c_Axon: c_axon_ptr}, nil
}

// Creates an instance of a child object. The primary fields are stored in and 
// managed by the C code.
func AxonChild(name string) (Axon, error) {
	c_str_name := C.CString(name)
	c_axon_ptr := C.createChildInstance(c_str_name)

	return Axon{c_Axon: c_axon_ptr}, nil
}

// Performs proper cleanup of the Axon object, using the C code.
func (self * Axon) Destroy() error {
	C.destroyInstance(self.c_Axon)

	return nil
}

// Sends the specified byte string to the other peer.
// 'code' is user defined, the receiver should interpret it.
// 'kind' is meant to indicate the data type, but can be interpeted by the 
//        receiver as desired.
func (self * Axon) SendMessage(msg []byte, code, kind int) error {
	c_bytes := C.CBytes(msg)
	length  := len(msg)
	C.sendMessage(
		self.c_Axon, 
		c_bytes, 
		C.int(code), 
		C.int(length), 
		C.int(kind),
	)

	return nil
}

// Receives a message as a byte string from the peer. This function will
// block until it receives a message.
func (self * Axon) RecvMessage() (msg []byte, code, kind int, err error) {
	c_code   := C.int(0)
	c_kind   := C.int(0)
	c_length := C.int(0)

	c_bytes := C.recvMessage(
		self.c_Axon, 
		(*C.int)(&c_code), 
		(*C.int)(&c_length), 
		(*C.int)(&c_kind),
	)

	bytes := C.GoBytes(c_bytes, c_length)

	return bytes, int(c_code), int(c_kind), nil
}




















