{{/*
Expand the name of the chart.
*/}}
{{- define "rmq-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rmq-demo.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "rmq-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "rmq-demo.labels" -}}
helm.sh/chart: {{ include "rmq-demo.chart" . }}
{{ include "rmq-demo.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "rmq-demo.selectorLabels" -}}
app.kubernetes.io/name: {{ include "rmq-demo.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "rmq-demo.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "rmq-demo.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Env variables for connectivity
*/}}
{{- define "rmq.credentials" -}}
- name: SECRET_DEMO_PORT
  value: {{ .Values.port | default "5672" | quote }}
- name: SECRET_DEMO_SERVER
  value: {{ .Values.server | default "rabbitmq.default.svc.cluster.local" }}
- name: SECRET_DEMO_USERNAME
  value: {{ .Values.username | default "guest" }}
- name: SECRET_DEMO_PASSWORD
  valueFrom:
    secretKeyRef:
     name: rabbitmq
     key: password
{{- end }}
