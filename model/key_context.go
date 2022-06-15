package model

type KeyContext string

const (
	KeyContextToken   KeyContext = "token_payload"
	KeyContextTraceID KeyContext = "trace_id"
	KeyContextSpanID  KeyContext = "span_id"
)
