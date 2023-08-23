package cmd

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var (
	threshold float64
	keyword   string
	sleepTime int
)

var quickTopCmd = &cobra.Command{
	Use:     "quick_top",
	Short:   "便捷查看进程对资源的占用",
	Example: "quick_top  ",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	quickTopCmd.PersistentFlags().Float64VarP(&threshold, "threshold", "t", 80, "cpu占用的阈值 默认是80")
	quickTopCmd.PersistentFlags().StringVarP(&keyword, "keyword", "k", "", "进程关键词,默认为空")
	quickTopCmd.PersistentFlags().IntVarP(&sleepTime, "sleep_time", "s", 5, "执行频率")
}

type processInfo struct {
	PID        int32
	Name       string
	CPUPercent float64
	CWD        string
}

func clear() {
	// 清屏
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func run() {
	for {
		clear()

		// 获取所有进程
		allProcesses, err := process.Processes()
		if err != nil {
			log.Fatal(err)
		}

		var processList []processInfo

		// 遍历所有进程
		for _, p := range allProcesses {
			// 获取进程的名称
			name, err := p.Name()
			if err != nil {
				log.Printf("Failed to get name for PID %d: %v", p.Pid, err)
				continue
			}

			// 检查是否指定了关键词，或进程名称包含关键词
			if keyword == "" || strings.Contains(name, keyword) {
				// 获取进程的CPU占用率
				cpuPercent, err := p.CPUPercent()
				if err != nil {
					log.Printf("Failed to get CPU percent for PID %d: %v", p.Pid, err)
					continue
				}

				// 判断是否是高CPU占用进程
				if cpuPercent >= threshold {
					// 获取进程的工作目录
					cwd, err := p.Cwd()
					if err != nil {
						log.Printf("Failed to get current working directory for PID %d: %v", p.Pid, err)
						continue
					}

					// 添加进程信息到列表
					processList = append(processList, processInfo{
						PID:        p.Pid,
						Name:       name,
						CPUPercent: cpuPercent,
						CWD:        cwd,
					})
				}
			}
		}

		// 根据CPU占用率从高到低排序
		sort.Slice(processList, func(i, j int) bool {
			return processList[i].CPUPercent > processList[j].CPUPercent
		})

		// 打印进程信息
		for _, p := range processList {
			fmt.Printf("PID: %d, Name: %s, CWD: %s, CPU Percent: %.2f\n", p.PID, p.Name, p.CWD, p.CPUPercent)
		}

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
