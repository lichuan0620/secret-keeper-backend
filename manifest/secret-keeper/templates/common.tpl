{{/* Base name for all resources */}}
{{- define "common.name" -}}
{{- if contains .Chart.Name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/* Common labels */}}
{{- define "common.labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/* Return the proper image name */}}
{{- define "common.image" -}}
{{- $repositoryName := .Values.platform.imageRepository -}}
{{- $imageName := .Values.image.name -}}
{{- $tag := .Chart.AppVersion | toString -}}
{{- if .Values.platform.imageRegistry -}}
    {{- $registryName := .Values.platform.imageRegistry -}}
    {{- printf "%s/%s/%s:%s" $registryName $repositoryName $imageName $tag -}}
{{- else -}}
    {{- printf "%s/%s:%s" $repositoryName $imageName $tag -}}
{{- end -}}
{{- end -}}