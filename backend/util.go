package backend

// taken from: https://github.com/portainer/portainer/blob/develop/app/components/stats/statsController.js#L177-L193
func calcCpuPercent(stats *ContainerStats) float64 {
	cpuPercent := 0.0

	cpuDelta := float64(stats.Cpu.Usage.Total - stats.PreCpu.Usage.Total)
	systemDelta := float64(stats.Cpu.SystemCpuUsage - stats.PreCpu.SystemCpuUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.Cpu.Usage.PerCpu)) * 100.0
	}

	return cpuPercent
}

func calcMemoryPercent(stats *ContainerStats) float64 {
	return float64(stats.Memory.Usage) * 100.0 / float64(stats.Memory.Limit)
}
