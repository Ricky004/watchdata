import { Log } from "@/components/types/log-type";

export async function getLogs() {
  const res = await fetch('http://localhost:8080/v1/logs')
  if (!res.ok) throw new Error('Failed to fetch logs');
  return res.json();
}

export async function getLogsSince(timestamp: string): Promise<Log[]> {
  const res = await fetch(`http://localhost:8080/v1/logs/since?timestamp=${encodeURIComponent(timestamp)}`)
  if (!res.ok) throw new Error("Failed to fetch logs since")
  return res.json()
}

export async function getLogsInTimeRanges(start: number, end: number): Promise<Log[]> {
  const res = await fetch(`http://localhost:8080/v1/logs/timerange?start=${start}&end=${end}`)
  if (!res.ok) throw new Error("Failed to fetch logs in time range")
  return res.json()
}
