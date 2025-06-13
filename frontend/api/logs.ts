export async function getTop10Logs() {
  const res = await fetch('http://localhost:8080/v1/logs');
  if (!res.ok) throw new Error('Failed to fetch logs');
  return res.json();
}
