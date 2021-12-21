{{- define "queue.name" -}}
{{- printf "%s-queue" (include "common.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "queue.labels" -}}
app: raven-queue
{{- if .Values.queue.metadata.labels }}
{{- toYaml .Values.queue.metadata.labels }}
{{- end }}
{{ include "common.labels" .}}
{{- end }}

{{- define "queue.podLabels" -}}
{{ include "queue.labels" . }}
{{- if .Values.queue.metadata.podLabels }}
{{- toYaml .Values.queue.metadata.podLabels }}
{{- end }}
{{- end }}

{{- define "queue.annotations" -}}
{{- if .Values.queue.metadata.annotations }}
{{- toYaml .Values.queue.metadata.annotations }}
{{- end -}}
{{- end }}

{{- define "queue.podAnnotations" -}}
{{- if .Values.queue.metadata.annotations }}
{{- toYaml .Values.queue.metadata.annotations }}
{{- end -}}
{{- if .Values.queue.metadata.podAnnotations }}
{{- toYaml .Values.queue.metadata.podAnnotations }}
{{- end -}}
{{- end }}

{{- define "queue.matchLabels" -}}
app: secret-keeper-queue
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}