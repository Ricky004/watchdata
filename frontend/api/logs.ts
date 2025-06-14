import { Log } from "@/components/log-viewer";

export async function getTop10Logs() {
  const res = await fetch('http://localhost:8080/v1/logs');
  if (!res.ok) throw new Error('Failed to fetch logs');
  return res.json();
}

export async function getLogsSince(timestamp: string): Promise<Log[]> {
  const res = await fetch(`http://localhost:8080/v1/logs/since?timestamp=${encodeURIComponent(timestamp)}`)
  if (!res.ok) throw new Error("Failed to fetch logs since")
  return res.json()
}
