package model

type Changes struct {
	Added    []FileSnapshot
	Removed  []FileSnapshot
	Modified []FileSnapshot
}
