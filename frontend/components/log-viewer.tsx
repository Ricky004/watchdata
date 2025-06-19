"use client"

import { ScrollArea } from "@/components/ui/scroll-area"
import { useEffect, useState, useRef } from "react"
import { getTop10Logs, getLogsSince, getLogsInTimeRanges } from "@/api/logs"
import React from "react"
import { useLiveLogs } from "@/hooks/use-live-logs"
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
import { Timer } from 'lucide-react';

export type Log = {
  timestamp: string;
  observed_time: string;
  severity_number: string;
  severity_text: string;
  body: string;
  attributes: string;
  resource: string;
  trace_id: string;
  span_id: string;
  trace_flags: string;
  flags: string;
  dropped_attributes_count: string;
};

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
  const [selectedLevels, setSelectedLevels] = useState<string[]>([])
  const lastSeenTimestamp = useRef<string | null>(null)
  const scrollAreaRef = useRef<HTMLDivElement>(null)

  const { logs: liveLogs } = useLiveLogs(live)

  const allLevels = ["DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"]

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

  return (
    <div>
      <div className="flex gap-1 justify-end">

        {/* Time range */}
        <div className="mt-3 mr-2">
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
        <div className="pt-3 pr-2 flex gap-2">
          <LiveToggleButton live={live} toggleLive={() => setLive(prev => !prev)} />
          <button
            onClick={() => setPaused(p => !p)}
            className={`px-4 py-2 rounded text-sm ${paused ? "bg-yellow-500" : "bg-green-600"} text-white`}
          >
            {paused ? "Resume" : "Pause"}
          </button>

        </div>
      </div>

      <div className="flex gap-2 mt-10 p-2">
        {allLevels.map(level => (
          <button
            key={level}
            onClick={() => {
              setSelectedLevels(prev =>
                prev.includes(level)
                  ? prev.filter(l => l !== level)
                  : [...prev, level]
              )
            }}
            className={`px-3 py-1 rounded text-sm border ${selectedLevels.includes(level)
              ? `${severityColorMap[level]} border-black font-bold`
              : "bg-gray-200 text-gray-800"
              }`}
          >
            {level}
          </button>
        ))}
      </div>

      {/* Log Viewer */}
      <div className="absolute bottom-0 left-0 w-full">
        <div className="mb-0.5 p-1 border bg-gray-100 dark:bg-slate-700">
          <h2 className="text-sm font-semibold">Log view</h2>
        </div>
        <ScrollArea className="h-110 w-full border" ref={scrollAreaRef}>
          <div className="min-w-full">
            {logs
              .filter(log => selectedLevels.length === 0 || selectedLevels.includes(log.severity_text.toUpperCase()))
              .map((log, i) => {
                const logId = `${log.timestamp}-${i}`
                const isExpanded = expandedLogs.has(logId)
                return (
                  <React.Fragment key={logId}>
                    <div className="w-full font-mono text-sm px-4 py-2 border-b border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-800 whitespace-pre-wrap">
                      <div className="flex justify-between items-center">
                        <div>
                          <span className="text-gray-500">[{log.timestamp}]</span>{' '}
                          <span className={`${getSeverityClass(log.severity_text)} font-semibold`}>
                            {log.severity_text.toUpperCase()}:
                          </span>{' '}
                          <span className="text-gray-500 dark:text-gray-200">{log.body}</span>{' '}
                          <span className="text-gray-400">
                            (trace_id={log.trace_id} span_id={log.span_id})
                          </span>
                        </div>
                        <button
                          onClick={() => toggleLog(logId)}
                          className="text-blue-500 hover:underline text-xs ml-4"
                        >
                          {isExpanded ? "Hide JSON" : "View JSON"}
                        </button>
                      </div>

                      {isExpanded && (
                        <pre className="mt-2 p-2 rounded bg-gray-100 dark:bg-slate-700 text-xs overflow-x-auto">
                          {JSON.stringify(log, null, 2)}
                        </pre>
                      )}
                    </div>
                  </React.Fragment>
                )
              })}
          </div>
        </ScrollArea>
      </div>
    </div>
  )
}