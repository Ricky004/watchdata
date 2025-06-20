"use client"

import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { ScrollArea } from "@/components/ui/scroll-area"

type FacetCategory = {
  title: string
  values: { label: string; count: string }[]
}

type FacetFiltersProps = {
  selectedLevels: string[]
  setSelectedLevels: (levels: string[]) => void
}

const FACETS: FacetCategory[] = [
  {
    title: "Source",
    values: [
      { label: "java", count: "1.08G" },
      { label: "lambda", count: "708M" },
      { label: "nodejs", count: "168M" },
      { label: "dotnet", count: "114M" },
      { label: "ruby", count: "92.3M" },
      { label: "python", count: "68.6M" },
      { label: "cloudtrail", count: "63.8M" },
      { label: "vpc", count: "16.8M" },
      { label: "mongodb", count: "13.9M" },
    ],
  },
  {
    title: "Status",
    values: [
      { label: "DEBUG", count: "13.0k" },
      { label: "INFO", count: "9.1k" },
      { label: "WARN", count: "6.7k" },
      { label: "ERROR", count: "6.7k" },
      { label: "CRITICAL", count: "6.7k" },
    ],
  },
]

export default function FacetFilters({ selectedLevels, setSelectedLevels }: FacetFiltersProps) {
  const toggle = (label: string) => {
    if (selectedLevels.includes(label)) {
      setSelectedLevels(selectedLevels.filter((l: string) => l !== label));
    } else {
      setSelectedLevels([...selectedLevels, label]);
    }
  }

  return (
    <div className="w-[280px] border p-4 shadow-md h-130 overflow-hidden">
      <div className="mb-2 font-semibold text-lg">Facets Filter</div>
      <Accordion type="multiple" defaultValue={FACETS.map(f => f.title)}>
        {FACETS.map((facet) => (
          <AccordionItem key={facet.title} value={facet.title}>
            <AccordionTrigger>{facet.title}</AccordionTrigger>
            <AccordionContent>
              <Input
                  placeholder="Filter values"
                  className="mb-2"
                />
              <ScrollArea className="h-[210px] pr-2">
                {facet.values.map((val) => (
                  <div
                    key={val.label}
                    className="flex items-center space-x-2 mb-1"
                  >
                    <Checkbox
                      id={val.label}
                      checked={selectedLevels.includes(val.label)}
                      onCheckedChange={() => toggle(val.label)}
                    />
                    <label
                      htmlFor={val.label}
                      className="text-sm flex justify-between w-full"
                    >
                      <span>{val.label}</span>
                      <span className="text-muted-foreground">{val.count}</span>
                    </label>
                  </div>
                ))}
              </ScrollArea>
            </AccordionContent>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  )
}
