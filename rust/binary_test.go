// Copyright 2019 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rust

import (
	"strings"
	"testing"
)

// Test that the prefer_dynamic property is handled correctly.
func TestPreferDynamicBinary(t *testing.T) {
	ctx := testRust(t, `
		rust_binary_host {
			name: "fizz-buzz-dynamic",
			srcs: ["foo.rs"],
			prefer_dynamic: true,
		}

		rust_binary_host {
			name: "fizz-buzz",
			srcs: ["foo.rs"],
		}`)

	fizzBuzz := ctx.ModuleForTests("fizz-buzz", "linux_glibc_x86_64").Output("fizz-buzz")
	fizzBuzzDynamic := ctx.ModuleForTests("fizz-buzz-dynamic", "linux_glibc_x86_64").Output("fizz-buzz-dynamic")

	path := ctx.ModuleForTests("fizz-buzz", "linux_glibc_x86_64").Module().(*Module).HostToolPath()
	if g, w := path.String(), "/host/linux-x86/bin/fizz-buzz"; !strings.Contains(g, w) {
		t.Errorf("wrong host tool path, expected %q got %q", w, g)
	}

	// Do not compile binary modules with the --test flag.
	flags := fizzBuzzDynamic.Args["rustcFlags"]
	if strings.Contains(flags, "--test") {
		t.Errorf("extra --test flag, rustcFlags: %#v", flags)
	}
	if !strings.Contains(flags, "prefer-dynamic") {
		t.Errorf("missing prefer-dynamic flag, rustcFlags: %#v", flags)
	}

	flags = fizzBuzz.Args["rustcFlags"]
	if strings.Contains(flags, "--test") {
		t.Errorf("extra --test flag, rustcFlags: %#v", flags)
	}
	if strings.Contains(flags, "prefer-dynamic") {
		t.Errorf("unexpected prefer-dynamic flag, rustcFlags: %#v", flags)
	}
}