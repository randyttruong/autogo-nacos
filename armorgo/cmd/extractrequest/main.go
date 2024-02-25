// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The extractrequest command runs the extractrequest analyzer.
package main

import (
	extractrequest "autoarmor/armorgo"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(extractrequest.Analyzer)
	//multichecker.Main(extractrequest.Analyzer)
	//fmt.Println("hello")
	//extractrequest.Output(extractrequest.ManifestData)
}

