package model

import (
	"fmt"
	"runtime"
	"runtime/metrics"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/richavery/donk-cli/internal/ui/styles"
)

const resourceMonitorInterval = time.Second

type resourceSnapshotMsg struct {
	time       time.Time
	cpuSeconds float64
	ramPercent float64
}

func (m *UI) resourceMonitorCmd() tea.Cmd {
	return tea.Tick(resourceMonitorInterval, func(now time.Time) tea.Msg {
		var samples = []metrics.Sample{{Name: "/cpu/classes/total:cpu-seconds"}}
		metrics.Read(samples)

		var memory runtime.MemStats
		runtime.ReadMemStats(&memory)
		ramPercent := 0.0
		if memory.HeapSys > 0 {
			ramPercent = float64(memory.HeapAlloc) / float64(memory.HeapSys) * 100
		}

		return resourceSnapshotMsg{
			time:       now,
			cpuSeconds: samples[0].Value.Float64(),
			ramPercent: min(100, max(0, ramPercent)),
		}
	})
}

func resourceCPUPercent(previous, current resourceSnapshotMsg) float64 {
	elapsed := current.time.Sub(previous.time).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return min(100, max(0, (current.cpuSeconds-previous.cpuSeconds)/elapsed/float64(runtime.NumCPU())*100))
}

func (m *UI) drawResourceMonitor(scr uv.Screen, area uv.Rectangle) {
	if !m.resourceReady || area.Dx() < 10 || area.Dy() == 0 {
		return
	}

	margin := max(1, area.Dx()/10)
	width := max(0, area.Dx()-margin*2)
	cpuWidth := width * 40 / 100
	ramWidth := width * 40 / 100
	gapWidth := max(0, width-cpuWidth-ramWidth)
	cpu := renderResourceBar("CPU", m.resourceCPU, cpuWidth, m.com.Styles)
	ram := renderResourceBar("RAM", m.resourceRAM, ramWidth, m.com.Styles)
	gap := strings.Repeat(" ", gapWidth)
	line := strings.Repeat(" ", margin) + cpu + gap + ram
	uv.NewStyledString(lipgloss.NewStyle().Width(area.Dx()).Render(line)).Draw(scr, area)
}

func renderResourceBar(label string, percent float64, width int, styleSet *styles.Styles) string {
	if width <= 0 {
		return ""
	}
	value := fmt.Sprintf("%3.0f%%", percent)
	labelText := label + " "
	fixedWidth := lipgloss.Width(labelText) + lipgloss.Width(value) + 1
	barWidth := max(1, width-fixedWidth)
	filledWidth := min(barWidth, max(0, int(float64(barWidth)*percent/100)))
	emptyWidth := barWidth - filledWidth
	bar := styleSet.Status.ResourceFilled.Render(strings.Repeat("█", filledWidth)) +
		styleSet.Status.ResourceEmpty.Render(strings.Repeat("░", emptyWidth))
	return styleSet.Status.ResourceLabel.Render(labelText) + bar + " " + styleSet.Status.ResourceValue.Render(value)
}
