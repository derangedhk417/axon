# Copyright 2018 Adam Robinson
# Permission is hereby granted, free of charge, to any person obtaining a copy of 
# this software and associated documentation files (the "Software"), to deal in the 
# Software without restriction, including without limitation the rights to use, copy, 
# modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, 
# and to permit persons to whom the Software is furnished to do so, subject to the 
# following conditions:
# The above copyright notice and this permission notice shall be included in all 
# copies or substantial portions of the Software.
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A 
# PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT 
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF 
# CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE 
# OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
import numpy           as np
import numpy.ctypeslib as ctl
import ctypes


libname = 'axonlib.so'
libdir  = './'
lib     = ctl.load_library(libname, libdir)

createControllerInstance = lib.createControllerInstance
createControllerInstance.argtypes = [ctypes.c_char_p]
createControllerInstance.restype  = ctypes.c_void_p

createChildInstance = lib.createChildInstance
createChildInstance.argtypes = [ctypes.c_char_p]
createChildInstance.restype  = ctypes.c_void_p

destroyInstance = lib.destroyInstance
destroyInstance.argtypes = [ctypes.c_void_p]
destroyInstance.restype  = None

sendMessage = lib.sendMessage
sendMessage.argtypes = [
	ctypes.c_void_p, 
	ctypes.c_char_p, 
	ctypes.c_int, 
	ctypes.c_int, 
	ctypes.c_int
]
sendMessage.restype  = None

recvMessage = lib.recvMessage
recvMessage.argtypes = [
	ctypes.c_void_p, 
	ctypes.c_int, 
	ctypes.c_int, 
	ctypes.c_int
]
recvMessage.restype  = ctypes.c_char_p

class Axon:
	def __init__(self, name, is_child=False):
		self.child  = is_child
		self.name   = name

		if self.child:
			self.c_Axon = createChildInstance(self.name.encode())
		else:
			self.c_Axon = createControllerInstance(self.name.encode())

	def SendMessage(self, _bytes, code, kind):
		sendMessage(
			self.c_Axon, 
			ctypes.from_buffer_copy(_bytes), 
			code, 
			len(_bytes), 
			kind
		)

	def RecvMessage(self):
		code   = ctypes.c_int()
		length = ctypes.c_int()
		kind   = ctypes.c_int()

		_bytes = recvMessage(
			self.c_Axon, 
			ctypes.byref(code),
			ctypes.byref(length),
			ctypes.byref(kind)
		)

		bytes_type = ctypes.c_byte * length
		return bytearray(bytes_type(_bytes)), code, kind

	def __del__(self):
		destroyInstance(self.c_Axon)
		