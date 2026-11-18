//go:build linux

package repositories

// isProcNetSupported en Linux los sockets se leen de /proc/net y su proceso de /proc/<pid>/fd.
const isProcNetSupported = true
