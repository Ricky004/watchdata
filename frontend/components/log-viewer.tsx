"use client"

import { ScrollArea } from "@/components/ui/scroll-area"
import { useEffect, useState, useRef } from "react"
import { getTop10Logs, getLogsSince } from "@/api/logs"
import React from "react"
import { useLiveLogs } from "@/hooks/use-live-logs"
import LiveToggleButton from "./live-toggle-btn"

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

export default function LogViewer() {
  const [live, setLive] = useState(true)
  const [paused, setPaused] = useState(false)
  const [autoScroll, setAutoScroll] = useState(true)
  const [logs, setLogs] = useState<Log[]>([])
  const [expandedLogs, setExpandedLogs] = useState<Set<string>>(new Set())
  const lastSeenTimestamp = useRef<string | null>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  const { logs: liveLogs } = useLiveLogs(live)

  // Load initial logs
  useEffect(() => {
    getTop10Logs().then((initialLogs) => {
      setLogs(initialLogs)
    })
  }, [])

  // Update last seen timestamp
  useEffect(() => {
    if (logs.length > 0) {
      lastSeenTimestamp.current = logs[logs.length - 1].timestamp
    }
  }, [logs])

  // Append live logs if not paused
  useEffect(() => {
    if (!live || paused || liveLogs.length === 0) return
    setLogs((prev) => {
      const existingTimestamps = new Set(prev.map(log => log.timestamp))
      const newUnique = liveLogs.filter(log => !existingTimestamps.has(log.timestamp))
      return [...prev, ...newUnique]
    })
  }, [liveLogs, live, paused])

  // Fetch missed logs when turning live back on
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

  // Scroll to bottom when new logs arrive
  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [logs, autoScroll])

  // Toggle expanded JSON view
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
      {/* Control Buttons */}
      <div className="absolute p-3 right-0 flex gap-3 z-10">
        <LiveToggleButton live={live} toggleLive={() => setLive(prev => !prev)} />
        <button
          onClick={() => setPaused(p => !p)}
          className={`px-4 py-2 rounded text-sm ${paused ? "bg-yellow-500" : "bg-green-600"} text-white`}
        >
          {paused ? "Resume" : "Pause"}
        </button>
        <button
          onClick={() => setAutoScroll(a => !a)}
          className={`px-4 py-2 rounded text-sm ${autoScroll ? "bg-blue-500" : "bg-gray-500"} text-white`}
        >
          {autoScroll ? "Auto-Scroll: ON" : "Auto-Scroll: OFF"}
        </button>
      </div>

      {/* Log Viewer */}
      <div className="absolute bottom-0 left-0 w-full">
        <div className="mb-0.5 p-1 border bg-gray-100 dark:bg-slate-700">
          <h2 className="text-sm font-semibold">Log view</h2>
        </div>
        <ScrollArea className="h-110 w-full border">
          <div className="min-w-full" ref={scrollRef}>
            {logs.map((log, i) => {
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
