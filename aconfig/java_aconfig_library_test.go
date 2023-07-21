// Copyright 2023 Google Inc. All rights reserved.
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

package aconfig

import (
	"strings"
	"testing"

	"android/soong/android"
	"android/soong/java"
)

// Note: These tests cover the code in the java package. It'd be ideal of that code could
// be in the aconfig package.

// With the bp parameter that defines a my_module, make sure it has the LOCAL_ACONFIG_FILES entries
func runJavaAndroidMkTest(t *testing.T, bp string) {
	result := android.GroupFixturePreparers(
		PrepareForTestWithAconfigBuildComponents,
		java.PrepareForTestWithJavaDefaultModules).
		ExtendWithErrorHandler(android.FixtureExpectsNoErrors).
		RunTestWithBp(t, bp+`
			aconfig_declarations {
				name: "my_aconfig_declarations",
				package: "com.example.package",
				srcs: ["foo.aconfig"],
			}

			java_aconfig_library {
				name: "my_java_aconfig_library",
				aconfig_declarations: "my_aconfig_declarations",
			}
		`)

	module := result.ModuleForTests("my_module", "android_common").Module()

	entry := android.AndroidMkEntriesForTest(t, result.TestContext, module)[0]

	makeVar := entry.EntryMap["LOCAL_ACONFIG_FILES"]
	android.AssertIntEquals(t, "len(LOCAL_ACONFIG_FILES)", 1, len(makeVar))
	if !strings.HasSuffix(makeVar[0], "intermediate.pb") {
		t.Errorf("LOCAL_ACONFIG_FILES should end with /intermediates.pb, instead it is: %s", makeVar[0])
	}
}

func TestAndroidMkJavaLibrary(t *testing.T) {
	bp := `
		java_library {
			name: "my_module",
			srcs: [
				"src/foo.java",
			],
			static_libs: [
				"my_java_aconfig_library",
			],
			platform_apis: true,
		}
	`

	runJavaAndroidMkTest(t, bp)
}

func TestAndroidMkAndroidApp(t *testing.T) {
	bp := `
		android_app {
			name: "my_module",
			srcs: [
				"src/foo.java",
			],
			static_libs: [
				"my_java_aconfig_library",
			],
			platform_apis: true,
		}
	`

	runJavaAndroidMkTest(t, bp)
}

func TestAndroidMkBinary(t *testing.T) {
	bp := `
		java_binary {
			name: "my_module",
			srcs: [
				"src/foo.java",
			],
			static_libs: [
				"my_java_aconfig_library",
			],
			platform_apis: true,
			main_class: "foo",
		}
	`

	runJavaAndroidMkTest(t, bp)
}

func TestAndroidMkAndroidLibrary(t *testing.T) {
	bp := `
		android_library {
			name: "my_module",
			srcs: [
				"src/foo.java",
			],
			static_libs: [
				"my_java_aconfig_library",
			],
			platform_apis: true,
		}
	`

	runJavaAndroidMkTest(t, bp)
}

func TestAndroidMkBinaryThatLinksAgainstAar(t *testing.T) {
	// Tests AndroidLibrary's propagation of flags through JavaInfo
	bp := `
		android_library {
			name: "some_library",
			srcs: [
				"src/foo.java",
			],
			static_libs: [
				"my_java_aconfig_library",
			],
			platform_apis: true,
		}
		java_binary {
			name: "my_module",
			srcs: [
				"src/bar.java",
			],
			static_libs: [
				"some_library",
			],
			platform_apis: true,
			main_class: "foo",
		}
	`

	runJavaAndroidMkTest(t, bp)
}
