package repositories

// DiskUsageReaderSyscallRepository ocupacion del sistema de ficheros donde esta un path. La llamada al sistema
// (statfs) va en disk_usage_statfs_linux.go; en el resto de sistemas, disk_usage_statfs_other.go devuelve "no
// soportado" para que el agente compile igual en Windows y macOS.
type DiskUsageReaderSyscallRepository struct{}

func DiskUsageReaderSyscallRepositoryGetInstance() *DiskUsageReaderSyscallRepository {
    return &DiskUsageReaderSyscallRepository{}
}

// GetDiskUsageBytes total, usados y disponibles para un usuario normal (sin la reserva de root), en bytes. Error si
// el path no existe o el sistema no esta soportado.
func (r *DiskUsageReaderSyscallRepository) GetDiskUsageBytes(path string) (float64, float64, float64, error) {
    return getStatfsBytes(path)
}
