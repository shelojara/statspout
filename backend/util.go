package backend

import "net"

// Creates a client for TCP (http) or Unix with the given address.
func createConn(http bool, address string) (net.Conn, error) {
	if http {
		return net.Dial("tcp", address)
	}

	return net.Dial("unix", address)
}

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

func sumTxBytesTotal(interfaces map[string]InterfaceStats) (sum uint32) {
	for _, i := range interfaces {
		sum += i.TxBytes
	}
	return
}

func sumRxBytesTotal(interfaces map[string]InterfaceStats) (sum uint32) {
	for _, i := range interfaces {
		sum += i.RxBytes
	}
	return
}
