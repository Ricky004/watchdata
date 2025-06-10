import { ModeToggle } from "@/components/dark-mode";
import LogRecord from "@/components/log-record";


export default function Home() {
  return (
    <main className="p-3">
      <ModeToggle />
      <LogRecord />
    </main>
  );
}
