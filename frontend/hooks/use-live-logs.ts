import { useEffect, useRef, useState } from "react";
import { Log } from "@/components/log-viewer";

export function useLiveLogs(enabled: boolean) {
  const [logs, setLogs] = useState<Log[]>([]);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!enabled) return;

    socketRef.current = new WebSocket("ws://localhost:8080/ws");

    socketRef.current.onmessage = (event) => {
      const log = JSON.parse(event.data);
      setLogs((prev) => [log, ...prev.slice(0, 99)]); // keep last 100
    };

    socketRef.current.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    return () => {
      socketRef.current?.close();
    };
  }, [enabled]);

  return { logs };
}
