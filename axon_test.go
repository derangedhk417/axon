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

package axon

import (
	"testing"
	"fmt"
	"os/exec"
	"os"
	"io"
	"bytes"
)

func TestInitializec(t *testing.T) {
	fmt.Println("[GO]     testing basic axon initialization . . .")
	fmt.Println("[GO]     \tinitializing Python Code")

	cmd := exec.Command("/usr/local/bin/python3", "axon_test.py")
	
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Start()

	if err != nil {
		t.Errorf("failed to start python\n")
	}

	fmt.Println("[GO]     \tpython code initialized")
	fmt.Println("[GO]     \tinitializing controller axon")

	channel, _ := AxonController("test2")

	fmt.Println("[GO]     \tcontroller initialized")

	channel.SendMessage([]byte{1, 2, 3, 4, 5}, 5, 8)

	cmd.Wait()

	fmt.Println("[GO]     \tdestroying contoller")
	channel.Destroy()
	fmt.Println("[GO]     \tcontroller destroyed")
}