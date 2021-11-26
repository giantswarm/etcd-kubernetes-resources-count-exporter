{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "etcd-kubernetes-resources-count-exporter.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "etcd-kubernetes-resources-count-exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "etcd-kubernetes-resources-count-exporter.labels" -}}
app: {{ include "etcd-kubernetes-resources-count-exporter.name" . | quote }}
{{ include "etcd-kubernetes-resources-count-exporter.selectorLabels" . }}
app.giantswarm.io/branch: {{ .Values.project.branch | quote }}
app.giantswarm.io/commit: {{ .Values.project.commit | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
helm.sh/chart: {{ include "etcd-kubernetes-resources-count-exporter.chart" . | quote }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "etcd-kubernetes-resources-count-exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "etcd-kubernetes-resources-count-exporter.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end -}}
