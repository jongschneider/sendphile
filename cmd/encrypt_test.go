/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"testing"
)

func Test_createDstFilepath(t *testing.T) {

	tests := []struct {
		name string
		fn   func() func()
		in   string
		want string
	}{
		{
			name: "empty directory",
			fn:   func() func() { return func() {} },
			in:   "test.txt",
			want: "enc_test.txt",
		},
		{
			name: "duplicate file",
			fn: func() func() {
				_, err := os.Create("./enc_test.txt")
				if err != nil {
					panic(err)
				}

				return func() {
					os.Remove("./enc_test.txt")
				}
			},
			in:   "./test.txt",
			want: "enc_test(1).txt",
		},
		{
			name: "duplicate (1) file",
			fn: func() func() {
				_, err := os.Create("./enc_test.txt")
				if err != nil {
					panic(err)
				}
				_, err = os.Create("./enc_test(1).txt")
				if err != nil {
					panic(err)
				}

				return func() {
					os.Remove("./enc_test.txt")
					os.Remove("./enc_test(1).txt")
				}
			},
			in:   "./test.txt",
			want: "enc_test(2).txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set up test
			defer tt.fn()()
			if got, _ := createDstFilepath(tt.in); got != tt.want {
				t.Errorf("createDstFilepath() = %v, want %v", got, tt.want)
			}

		})
	}
}
