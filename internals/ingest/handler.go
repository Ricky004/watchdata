package ingest


import (
	"io"
	"net/http"

	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
)

func HandleOTLPLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	req := plogotlp.NewExportRequest()
	if err := req.UnmarshalProto(body); err != nil {
		http.Error(w, "invalid protobuf payload", http.StatusUnsupportedMediaType)
		return
	}

	logs := req.Logs()

	w.WriteHeader(http.StatusAccepted)
}