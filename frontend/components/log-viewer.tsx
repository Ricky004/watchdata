"use client"

import { ScrollArea } from "@/components/ui/scroll-area"
import { useEffect, useState, useRef } from "react"
import { getTop10Logs, getLogsSince } from "@/api/logs"
import React from "react"
import { useLiveLogs } from "@/hooks/use-live-logs"
import { Button } from "./ui/button"

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

export default function LogRecord() {
  const [live, setLive] = useState(true)
  const [logs, setLogs] = useState<Log[]>([])
  const { logs: liveLogs } = useLiveLogs(live)
  const lastSeenTimestamp = useRef<string | null>(null)

  // Update last seen timestamp when logs change
  useEffect(() => {
    if (logs.length > 0) {
      lastSeenTimestamp.current = logs[logs.length - 1].timestamp
    }
  }, [logs])

  // On mount, load initial logs
  useEffect(() => {
    getTop10Logs().then((initialLogs) => {
      setLogs(initialLogs)
    })
  }, [])

  // When liveLogs change → append new ones
  useEffect(() => {
    if (live && liveLogs.length > 0) {
      setLogs((prev) => {
        const existingTimestamps = new Set(prev.map(log => log.timestamp))
        const newUnique = liveLogs.filter(log => !existingTimestamps.has(log.timestamp))
        return [...prev, ...newUnique]
      })
    }
  }, [liveLogs, live])

  // When live mode is turned ON again → fetch missing logs since last timestamp
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

  return (
    <div className="absolute bottom-0 left-0 w-full">
      <Button onClick={() => setLive((prev) => !prev)}>
        {live ? 'Live: on' : 'Live: off'}
      </Button>
      <div className="mb-0.5 p-1 border bg-gray-100 dark:bg-slate-700">
        <h2 className="text-sm font-semibold">Log view</h2>
      </div>
      <ScrollArea className="h-110 w-full border">
        <div className="min-w-full">
          {logs.map((log, i) => (
            <React.Fragment key={`${log.timestamp}-${i}`}>
              <div className="w-full font-mono text-sm px-4 py-2 border-b border-gray-300 dark:border-slate-600 bg-white dark:bg-slate-800 whitespace-pre-wrap">
                <span className="text-gray-500">[{log.timestamp}]</span>{' '}
                <span className={`${getSeverityClass(log.severity_text)} font-semibold`}>
                  {log.severity_text.toUpperCase()}:
                </span>{' '}
                <span className="text-gray-500 dark:text-gray-200">{log.body}</span>{' '}
                <span className="text-gray-400">
                  (trace_id={log.trace_id} span_id={log.span_id})
                </span>
              </div>
            </React.Fragment>
          ))}
        </div>
      </ScrollArea>
    </div>
  )
}
