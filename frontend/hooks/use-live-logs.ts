import { useEffect, useRef, useState } from "react";
import { Log } from "@/components/types/log-type";

export function useLiveLogs(enabled: boolean) {
  const [logs, setLogs] = useState<Log[]>([]);
  const socketRef = useRef<WebSocket | null>(null);
  const lastTimestampRef = useRef<string | null>(null);

  // Update last timestamp whenever logs change
  useEffect(() => {
    if (logs.length > 0) {
      lastTimestampRef.current = logs[0].timestamp; 
    }
  }, [logs]);

  useEffect(() => {
    if (!enabled) return;

    // Step 1: Fetch missed logs since last known timestamp
    if (lastTimestampRef.current) {
      fetch(`/v1/logs/since?timestamp=${encodeURIComponent(lastTimestampRef.current)}`)
        .then((res) => res.json())
        .then((fetchedLogs: Log[]) => {
          // Prepend new logs (assumed ascending order from server)
          setLogs((prev) => [...fetchedLogs.reverse(), ...prev].slice(0, 100));
        })
        .catch((err) => {
          console.error("Failed to fetch missed logs:", err);
        });
    }

    // Step 2: Setup WebSocket
    const ws = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = ws;

    ws.onmessage = (event) => {
      const log: Log = JSON.parse(event.data);

      // Prevent duplicates (match by timestamp and trace_id/span_id)
      setLogs((prev) => {
        const exists = prev.some(
          (l) => l.timestamp === log.timestamp && l.trace_id === log.trace_id && l.span_id === log.span_id
        );
        if (exists) return prev;
        return [log, ...prev.slice(0, 99)];
      });
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    return () => {
      ws.close();
    };
  }, [enabled]);

  return { logs };
}