package main

import "sync"

// -----------------------------
// In-Memory User Store
// -----------------------------
// This map holds UserDetail entries keyed by username.
// A mutex guards concurrent access.
var (
	userStore   = make(map[string]UserDetail)
	userStoreMu sync.RWMutex
)