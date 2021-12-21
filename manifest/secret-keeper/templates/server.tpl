{{- define "server.name" -}}
{{- printf "%s-server" (include "common.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "server.labels" -}}
app: raven-server
{{- if .Values.server.metadata.labels }}
{{- toYaml .Values.server.metadata.labels }}
{{- end }}
{{ include "common.labels" .}}
{{- end }}

{{- define "server.podLabels" -}}
{{ include "server.labels" . }}
{{- if .Values.server.metadata.podLabels }}
{{- toYaml .Values.server.metadata.podLabels }}
{{- end }}
{{- end }}

{{- define "server.annotations" -}}
{{- if .Values.server.metadata.annotations }}
{{- toYaml .Values.server.metadata.annotations }}
{{- end -}}
{{- end }}

{{- define "server.podAnnotations" -}}
{{- if .Values.server.metadata.annotations }}
{{- toYaml .Values.server.metadata.annotations }}
{{- end -}}
{{- if .Values.server.metadata.podAnnotations }}
{{- toYaml .Values.server.metadata.podAnnotations }}
{{- end -}}
{{- end }}

{{- define "server.matchLabels" -}}
app: secret-keeper-server
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}