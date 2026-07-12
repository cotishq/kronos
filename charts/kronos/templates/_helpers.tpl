{{/*
Expand the name of the chart.
*/}}
{{- define "kronos.name" -}}
{{- .Chart.Name }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kronos.labels" -}}
app.kubernetes.io/name: {{ include "kronos.name" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}