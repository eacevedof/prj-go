//go:build linux

package repositories

import "syscall"

// getStatfsBytes como df: usados = bloques - libres; disponibles = los que puede usar un usuario sin privilegios.
func getStatfsBytes(path string) (float64, float64, float64, error) {
    var stat syscall.Statfs_t
    if err := syscall.Statfs(path, &stat); err != nil {
        return 0, 0, 0, err
    }
    blockSize := float64(stat.Bsize)
    totalBytes := float64(stat.Blocks) * blockSize
    usedBytes := (float64(stat.Blocks) - float64(stat.Bfree)) * blockSize
    availableBytes := float64(stat.Bavail) * blockSize
    return totalBytes, usedBytes, availableBytes, nil
}
