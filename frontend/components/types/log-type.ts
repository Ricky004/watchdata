export type Log = {
  timestamp: string;
  observed_time: string;
  severity_number: string;
  severity_text: string;
  body: string;
  attributes: string;
  resource: {
    attributes?: { key: string; value: string }[];
    [key: string]: unknown;
  };
  trace_id: string;
  span_id: string;
  trace_flags: string;
  flags: string;
  dropped_attributes_count: string;
};