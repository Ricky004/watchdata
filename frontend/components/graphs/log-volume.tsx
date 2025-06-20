"use client"

import * as echarts from 'echarts/core'
import { useEffect, useRef } from 'react'
import {
  TooltipComponent,
  GridComponent,
  TitleComponent,
  LegendComponent
} from 'echarts/components'
import { BarChart, LineChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  TooltipComponent,
  GridComponent,
  TitleComponent,
  LegendComponent,
  BarChart,
  LineChart,
  CanvasRenderer
])

type Props = {
  logs: { timestamp: string }[]  // minimal for count-based graph
}

export default function LogChart({ logs }: Props) {
  const chartRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!chartRef.current || logs.length === 0) return

    const chart = echarts.init(chartRef.current)

    // Group logs by hour
    const counts: Record<string, number> = {}
    logs.forEach(log => {
      const hour = new Date(log.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      counts[hour] = (counts[hour] || 0) + 1
    })

    const timeBuckets = Object.keys(counts)
    const logCounts = Object.values(counts)

    const option = {
        tooltip: {
          trigger: 'axis'
        },
        grid: {
          left: '2%',
          right: '2%',
          bottom: '5%',
          top: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: true, 
          data: timeBuckets,
          axisLabel: {
            fontSize: 10
          }
        },
        yAxis: {
          type: 'value',
        },
        series: [
          {
            name: 'Logs',
            type: 'bar',
            data: logCounts,
            barWidth: '60%', 
            itemStyle: {
              color: '#91caff'
            }
          }
        ]
      }

    chart.setOption(option)
    return() => chart.dispose()
  }, [logs])

  return <div ref={chartRef} style={{ width: '100%', height: '200px' }} />
}
