package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const userHZ = 100

type CPUStat struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

//type SoftIRQStat struct {
//	Hi          uint64
//	Timer       uint64
//	NetTx       uint64
//	NetRx       uint64
//	Block       uint64
//	BlockIoPoll uint64
//	Tasklet     uint64
//	Sched       uint64
//	Hrtimer     uint64
//	Rcu         uint64
//}

type Stat struct {
	// Boot time in seconds since the Epoch.
	//BootTime uint64
	// Summed up cpu statistics.
	CPUTotal CPUStat
	// Per-CPU statistics.
	CPU []CPUStat
	// Number of times interrupts were handled, which contains numbered and unnumbered IRQs.
	//IRQTotal uint64
	// Number of times a numbered IRQ was triggered.
	//IRQ []uint64
	// Number of times a context switch happened.
	//ContextSwitches uint64
	// Number of times a process was created.
	//ProcessCreated uint64
	// Number of processes currently running.
	//ProcessesRunning uint64
	// Number of processes currently blocked (waiting for IO).
	//ProcessesBlocked uint64
	// Number of times a softirq was scheduled.
	//SoftIRQTotal uint64
	// Detailed softirq statistics.
	//SoftIRQ SoftIRQStat
}

//func ReadFileNoStat(filename string) ([]byte, error) {
//	const maxBufferSize = 1024 * 512
//
//	f, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//
//	reader := io.LimitReader(f, maxBufferSize)
//	return ioutil.ReadAll(reader)
//}

func parseCPUStat(line string) (CPUStat, int64, error) {
	cpuStat := CPUStat{}
	var cpu string

	count, err := fmt.Sscanf(line, "%s %f %f %f %f %f %f %f %f %f %f",
		&cpu,
		&cpuStat.User, &cpuStat.Nice, &cpuStat.System, &cpuStat.Idle,
		&cpuStat.Iowait, &cpuStat.IRQ, &cpuStat.SoftIRQ, &cpuStat.Steal,
		&cpuStat.Guest, &cpuStat.GuestNice)

	if err != nil && err != io.EOF {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu): %w", line, err)
	}
	if count == 0 {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu): 0 elements parsed", line)
	}

	cpuStat.User /= userHZ
	cpuStat.Nice /= userHZ
	cpuStat.System /= userHZ
	cpuStat.Idle /= userHZ
	cpuStat.Iowait /= userHZ
	cpuStat.IRQ /= userHZ
	cpuStat.SoftIRQ /= userHZ
	cpuStat.Steal /= userHZ
	cpuStat.Guest /= userHZ
	cpuStat.GuestNice /= userHZ

	if cpu == "cpu" {
		return cpuStat, -1, nil
	}

	cpuID, err := strconv.ParseInt(cpu[3:], 10, 64)
	if err != nil {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu/cpuid): %w", line, err)
	}

	return cpuStat, cpuID, nil
}

//func parseSoftIRQStat(line string) (SoftIRQStat, uint64, error) {
//	softIRQStat := SoftIRQStat{}
//	var total uint64
//	var prefix string
//
//	_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d",
//		&prefix, &total,
//		&softIRQStat.Hi, &softIRQStat.Timer, &softIRQStat.NetTx, &softIRQStat.NetRx,
//		&softIRQStat.Block, &softIRQStat.BlockIoPoll,
//		&softIRQStat.Tasklet, &softIRQStat.Sched,
//		&softIRQStat.Hrtimer, &softIRQStat.Rcu)
//
//	if err != nil {
//		return SoftIRQStat{}, 0, fmt.Errorf("couldn't parse %q (softirq): %w", line, err)
//	}
//
//	return softIRQStat, total, nil
//}

func GetStat(data []byte) (stat Stat, err error) {
	//fileName := fs.proc.Path("stat")
	//data, err := ReadFileNoStat(fileName)
	//if err != nil {
	//	return Stat{}, err
	//}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		// require at least <key> <value>
		if len(parts) < 2 {
			continue
		}
		switch {
		//case parts[0] == "btime":
		//	if stat.BootTime, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (btime): %w", parts[1], err)
		//	}
		//case parts[0] == "intr":
		//	if stat.IRQTotal, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (intr): %w", parts[1], err)
		//	}
		//	numberedIRQs := parts[2:]
		//	stat.IRQ = make([]uint64, len(numberedIRQs))
		//	for i, count := range numberedIRQs {
		//		if stat.IRQ[i], err = strconv.ParseUint(count, 10, 64); err != nil {
		//			return Stat{}, fmt.Errorf("couldn't parse %q (intr%d): %w", count, i, err)
		//		}
		//	}
		//case parts[0] == "ctxt":
		//	if stat.ContextSwitches, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (ctxt): %w", parts[1], err)
		//	}
		//case parts[0] == "processes":
		//	if stat.ProcessCreated, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (processes): %w", parts[1], err)
		//	}
		//case parts[0] == "procs_running":
		//	if stat.ProcessesRunning, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (procs_running): %w", parts[1], err)
		//	}
		//case parts[0] == "procs_blocked":
		//	if stat.ProcessesBlocked, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
		//		return Stat{}, fmt.Errorf("couldn't parse %q (procs_blocked): %w", parts[1], err)
		//	}
		//case parts[0] == "softirq":
		//	softIRQStats, total, err := parseSoftIRQStat(line)
		//	if err != nil {
		//		return Stat{}, err
		//	}
		//	stat.SoftIRQTotal = total
		//	stat.SoftIRQ = softIRQStats
		case strings.HasPrefix(parts[0], "cpu"):
			cpuStat, cpuID, err := parseCPUStat(line)
			if err != nil {
				return Stat{}, err
			}
			if cpuID == -1 {
				stat.CPUTotal = cpuStat
			} else {
				for int64(len(stat.CPU)) <= cpuID {
					stat.CPU = append(stat.CPU, CPUStat{})
				}
				stat.CPU[cpuID] = cpuStat
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Stat{}, fmt.Errorf("couldn't parse: %w", err)
	}

	return stat, nil
}
