"use client"

import { useEffect, useState } from "react";
import { getLogs } from "@/api/logs";

export default function LogRecord() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    getLogs().then(setLogs);
  }, []);

  return (
      <div>
      <h1 className="text-xl font-bold mb-2">Log Records</h1>
      <ul>
        {logs.map((log, i) => (
          <li key={i}>{log}</li>
        ))}
      </ul>
    </div>
  );
}