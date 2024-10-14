package main

import (
	_ "embed"
)

//go:embed index.html
var IndexHTML string

//go:embed bundle.js
var BundleJS string
