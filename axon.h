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

#include <stdio.h>
#include <stdlib.h>
#include <semaphore.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdbool.h>

// Stores all state information for communication between
// two processes.
struct Axon {
	char * system_name; // Stores a name, specified when initializing 
	                    // an instance. Should be unique. Used as the
	                    // prefix to the names of the four semaphores
	                    // used for communication.

	bool is_controller;   // TRUE when initialized as controller,
	                      // FALSE when initialized as child.

	sem_t * controllerSent; // Waited on by the child and triggered 
	                        // by the controller when a message is sent.
	sem_t * childReceived;  // Waited on by the controller after a 
	                        // message is sent. The child will trigger
	                        // this when it receives the message.
	                        // This is used so that the controller can
	                        // wait for proper message receipt before
	                        // continuing execution.

	sem_t * childSent;      // Waited on by the parent to receive messages
	                        // from the child.
	sem_t * controllerReceived; // Waited on by the child to ensure that
	                            // messages are received by the parent 
	                            // before execution continues.

	int fd;            // File descriptor attached to the shared memory
	                   // that is used to pass message contents.

	int * messageCode; // Stores the number used to identify the type
	                   // of message being sent. Meaning is used defined.

	int * messageSize; // This is populated with the size of the message 
	                   // everytime a message is passed.
	int * messageKind; // This is populated with the message type every
	                   // time a message is passed. See the #define 
	                   // statements at the top of the file.
};

// This function allocates shared memory of the specified
// size and returns a pointer to it. Since the same function
// needs to be called in other processes using the same name
// to get access to the shared memory, a name paramter must 
// be passed.
void * mallocShared(size_t size, char * name);

// Used to reallocate memory associated with the file
// descriptor that is assigned to the instance. This
// memory is used for passing messages between the 
// controller and the child.
void * reallocShared(size_t size, int fd);

// Constructs the name that should be used to identify
// the file descriptor for the shared memory used as
// a location for addresses being passed between 
// processes.
// YES, I know the name of this function is confusing.
char * getMessageFDNameLocationFDName(char * base);

// Constructs the name that should be used to identify
// the file descriptor for the shared memory used as
// a location for message size indicators being passed
// between processes.
char * getMessageSizeFDName(char * base);

// Constructs the name that should be used to identify
// the file descriptor for the shared memory used as
// a location for message code indicators being passed
// between processes.
char * getMessageCodeFDName(char * base);

// Constructs the name that should be used to identify
// the file descriptor for the shared memory used as
// a location for message type indicators being passed
// between processes.
char * getMessageTypeFDName(char * base);

// ----------------------------------------------
// Initialization functions
// ----------------------------------------------
// These two functions need to be called in the 
// controller application and the child application.
// 
// Once called, both processes will be able to access
// the shared memory for communication and the
// semphores for synchronization. 


// Called in the controlling program.
// Will then wait for the child process to call createChildInstance with 
// the same name argument. CreateChildInstance will trigger the
// childReceived semaphore and the controller instance will 
// know that the child has started.
// 
// parameters:
//     - name: user defined unique string
//             must be the same in both the controller and child processes
struct Axon * createControllerInstance(char * name);

// Called by the child process.
//
// parameters:
//     - name: user defined unique string
//             must be the same in both the controller and child processes
//
struct Axon * createChildInstance(char * name);

// Sends a message.
// Can be called on either a child or controller, doesn't matter.
// Internally the function will allocate some shared memory and 
// copy the message to it before triggering a semaphore. The caller
// is responsible for deallocating the message that they pass in.
// This function will halt execution until the receiver confirms
// that they have received the message.
void sendMessage(struct Axon * instance, void * message, int code, int length, int kind);

// Halts until revceiving a message. When a message is received, it 
// will be copied from shared memory into local memory. The returned
// pointer is the responsibility of the caller to free.
void * recvMessage(struct Axon * instance, int * code, int * length, int * kind);

// Removes all semaphores, frees shared memory and 
// unlinks shared memory file descriptors. Call this in
// the controller program before it exits. Do not call
// in the child program. More than one call might cause
// a problem.
void destroyInstance(struct Axon * instance);