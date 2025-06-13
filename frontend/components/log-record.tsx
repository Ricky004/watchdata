"use client"

import { useEffect, useState } from "react";
import { getTop10Logs } from "@/api/logs";


type Log = {
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

export default function LogRecord() {
  const [logs, setLogs] = useState<Log[]>([]);

  useEffect(() => {
    getTop10Logs().then(setLogs);
  }, []);

  return (
    <div>
      <h1 className="text-xl font-bold mb-2">Log Records</h1>
      <ul>
        {logs.map((log, i) => (
          <li key={i} className="mb-2 p-2 border rounded bg-gray-100 dark:bg-slate-700">
            <p className="flex gap-4">
              <span>{log.timestamp}</span>
              <span>{log.severity_number}</span>
              <span>{log.severity_number}</span>
              <span>{log.severity_text}</span>
              <span>{log.body}</span>
              <span>{JSON.stringify(log.attributes)}</span>
              <span>{JSON.stringify(log.resource)}</span>
              <span>{log.trace_id}</span>
              <span>{log.span_id}</span>
              <span>{log.trace_flags}</span>
              <span>{log.flags}</span>
              <span>{log.dropped_attributes_count}</span>
            </p>
          </li>
        ))}
      </ul>
    </div>
  );
}