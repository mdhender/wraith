// MIT License
//
// Copyright (c) 2021 Pace Software Ltd
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package definition

// GreeterService makes nice greetings.
type GreeterService interface {
	// Greet makes a greeting.
	Greet(GreetRequest) GreetResponse
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
	// Name is the person to greet.
	// example: "Mat Ryer"
	Name string
}

// GreetResponse is the response object containing a
// person's greeting.
type GreetResponse struct {
	// Greeting is the greeting that was generated.
	// example: "Hello Mat Ryer"
	Greeting string
}
