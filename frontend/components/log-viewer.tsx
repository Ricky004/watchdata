"use client"

import { ScrollArea } from "@/components/ui/scroll-area"
import { useEffect, useState, useRef } from "react"
import { getTop10Logs, getLogsSince, getLogsInTimeRanges } from "@/api/logs"
import { useLiveLogs } from "@/hooks/use-live-logs"
import { Log } from "./types/log-type"
import LiveToggleButton from "./live-toggle-btn"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Timer, Download } from 'lucide-react';
import { SidebarTrigger } from "./ui/sidebar"
import { Separator } from "./ui/separator"
import FacetFilters from "./facet-filters"
import { Button } from "./ui/button"
import React from "react"
import LogChart from "./graphs/log-volume"


const severityColorMap: Record<string, string> = {
  DEBUG: 'text-blue-500',
  INFO: 'text-blue-500',
  WARN: 'text-yellow-500',
  ERROR: 'text-red-500',
  CRITICAL: 'text-red-700 font-bold',
}

function getSeverityClass(severity: string): string {
  return severityColorMap[severity.toUpperCase()] || 'text-gray-500'
}

const TIME_RANGES = [
  { label: "15m", value: 15 * 60 },
  { label: "1h", value: 60 * 60 },
  { label: "4h", value: 4 * 60 * 60 },
  { label: "24h", value: 24 * 60 * 60 },
];

export default function LogViewer() {
  const [live, setLive] = useState(true)
  const [paused, setPaused] = useState(false)
  const [autoScroll] = useState(true)
  const [logs, setLogs] = useState<Log[]>([])
  const [expandedLogs, setExpandedLogs] = useState<Set<string>>(new Set())
  const [selectedRange, setSelectedRange] = useState<number | null>(null)
  const [selected, setSelected] = useState<string[]>([])
  const lastSeenTimestamp = useRef<string | null>(null)
  const scrollAreaRef = useRef<HTMLDivElement>(null)

  const { logs: liveLogs } = useLiveLogs(live)

  useEffect(() => {
    getTop10Logs().then((initialLogs) => {
      setLogs(initialLogs)
    })
  }, [])

  useEffect(() => {
    if (logs.length > 0) {
      lastSeenTimestamp.current = logs[logs.length - 1].timestamp
    }
  }, [logs])

  useEffect(() => {
    if (!live || paused || liveLogs.length === 0) return
    setLogs((prev) => {
      const existingTimestamps = new Set(prev.map(log => log.timestamp))
      const newUnique = liveLogs.filter(log => !existingTimestamps.has(log.timestamp))
      return [...prev, ...newUnique]
    })
  }, [liveLogs, live, paused])

  useEffect(() => {
    if (live && lastSeenTimestamp.current) {
      getLogsSince(lastSeenTimestamp.current).then((missedLogs) => {
        setLogs((prev) => {
          const existingTimestamps = new Set(prev.map(log => log.timestamp))
          const newUnique = missedLogs.filter(log => !existingTimestamps.has(log.timestamp))
          return [...prev, ...newUnique]
        })
      })
    }
  }, [live])

  // Fixed auto-scroll effect
  useEffect(() => {
    if (autoScroll && scrollAreaRef.current) {
      // Find the actual scrollable element within ScrollArea
      const scrollViewport = scrollAreaRef.current.querySelector('[data-radix-scroll-area-viewport]') as HTMLElement
      if (scrollViewport) {
        scrollViewport.scrollTop = scrollViewport.scrollHeight
      }
    }
  }, [logs, autoScroll])

  useEffect(() => {
    if (live) {
      setSelectedRange(null)
    }
  }, [live])

  useEffect(() => {
    if (live || selectedRange === null) return;
    const now = Math.floor(Date.now() / 1000);
    const start = now - selectedRange;

    getLogsInTimeRanges(start, now).then(setLogs).catch(console.error);
  }, [selectedRange, live])

  const toggleLog = (id: string) => {
    setExpandedLogs(prev => {
      const newSet = new Set(prev)
      if (newSet.has(id)) {
        newSet.delete(id)
      } else {
        newSet.add(id)
      }
      return newSet
    })
  }

  // Filter logs based on selected levels and source.name
  const filteredLogs = logs.filter(log => {
    const sourceName =
      log.resource?.attributes?.find(
        (a: { key: string; value: string }) => a.key === "source.name"
      )?.value || "";
    const serviceName =
      log.resource?.attributes?.find(
        (a: { key: string; value: string }) => a.key === "service.name"
      )?.value || "";

    const severity = log.severity_text?.toUpperCase();

    const severityMatch = selected.length === 0 || selected.includes(severity);
    const sourceMatch = selected.length === 0 || selected.includes(sourceName);
    const serviceMatch = selected.length === 0 || selected.includes(serviceName);

    return severityMatch || sourceMatch || serviceMatch
  });

  return (
    <div>
      <div className="flex gap-1 justify-between">
        <div className="m-2">
          <SidebarTrigger className="-ml-1" />
        </div>
        <div className="flex">
          {/* Time range */}
          <div className="m-2">
            <Select
              value={selectedRange !== null ? selectedRange.toString() : ""}
              onValueChange={(value) => {
                setSelectedRange(Number(value))
                setLive(false)
              }}
            >
              <SelectTrigger className="w-[160px]">
                <SelectValue placeholder="Select a range" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>time range</SelectLabel>
                  {TIME_RANGES.map((range) => (
                    <SelectItem
                      key={range.value}
                      value={range.value.toString()}
                    >
                      <Timer /> Last {range.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          {/* Control Buttons */}
          <div className="flex gap-2 mt-2 mb-2">
            <LiveToggleButton live={live} toggleLive={() => setLive(prev => !prev)} />
            <button
              onClick={() => setPaused(p => !p)}
              className={`px-4 py-2 rounded text-sm ${paused ? "bg-yellow-500" : "bg-green-600"} text-white`}
            >
              {paused ? "Resume" : "Pause"}
            </button>

          </div>
        </div>
      </div>
      <Separator />

      <div className="mt-20">
        <LogChart logs={logs} />
      </div>

      <div className="">
        {/* Sidebar (Facet Filters) */}
        <div className="bg-white dark:bg-slate-900 absolute bottom-0">
          <FacetFilters
            selected={selected}
            setSelected={setSelected}
          />
        </div>

        {/* Main Content (Log Viewer) */}
        <div className="bg-gray-50 dark:bg-slate-800 absolute bottom-0 left-[280px] right-0">
          <div className="flex justify-end p-1 border bg-white dark:bg-slate-700">
            <Button><Download /> Download as CSV</Button>
          </div>
          {/* Header row */}
          <div className="p-2 border-b bg-gray-100 dark:bg-slate-700">
            <div className="grid grid-cols-14 gap-2 text-xs font-semibold text-gray-700 dark:text-gray-300">
              <div className="col-span-3">DATE</div>
              <div className="col-span-2">SEVERITY LEVEL</div>
              <div className="col-span-2">SERVICE</div>
              <div className="col-span-7">CONTENT</div>
            </div>
          </div>

          {/* Scrollable Logs */}
          <ScrollArea className="h-110 w-full border" ref={scrollAreaRef}>
            <div className="min-w-full">
              {filteredLogs.map((log, i) => {
                const logId = `${log.timestamp}-${i}`
                const isExpanded = expandedLogs.has(logId)
                const latestTimestamp = logs.reduce((latest, log) => {
                  return new Date(log.timestamp) > new Date(latest) ? log.timestamp : latest
                }, logs[0]?.timestamp || "")
                const isLatestLog = log.timestamp === latestTimestamp

                return (
                  <React.Fragment key={logId}>
                    <div className="w-full font-mono text-sm border-b border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-800 whitespace-pre-wrap relative">
                      {/* Active Bar */}
                      <div
                        className={`absolute left-0 top-0 w-1 h-full ${isLatestLog ? "bg-red-500" : "bg-sky-200"
                          }`}
                      />

                      <div className="px-4 py-2 ml-2">
                        <div className="grid grid-cols-14 gap-2 items-center">
                          <div className="col-span-3">
                            <span className="text-gray-500">{log.timestamp}</span>
                          </div>
                          <div className="col-span-2">
                            <span
                              className={`${getSeverityClass(
                                log.severity_text
                              )} font-semibold`}
                            >
                              {log.severity_text.toUpperCase()}
                            </span>
                          </div>
                          <div className="col-span-2">
                            <span className="text-gray-500">
                              {
                                log.resource?.attributes?.find(
                                  (a: { key: string; value: string }) =>
                                    a.key === "service.name"
                                )?.value || "—"
                              }
                            </span>
                          </div>
                          <div className="col-span-6">
                            <span className="text-gray-500 dark:text-gray-200">
                              {log.body}
                            </span>
                            <span className="text-gray-400 text-xs ml-2">
                              (trace_id={log.trace_id} span_id={log.span_id})
                            </span>
                          </div>
                          <div className="col-span-1 text-right">
                            <button
                              onClick={() => toggleLog(logId)}
                              className="text-blue-500 hover:underline text-xs"
                            >
                              {isExpanded ? "Hide JSON" : "View JSON"}
                            </button>
                          </div>
                        </div>

                        {isExpanded && (
                          <pre className="mt-2 p-2 rounded bg-gray-100 dark:bg-slate-700 text-xs overflow-x-auto">
                            {JSON.stringify(log, null, 2)}
                          </pre>
                        )}
                      </div>
                    </div>
                  </React.Fragment>
                )
              })}
            </div>
          </ScrollArea>
        </div>
      </div>

    </div>
  )
}