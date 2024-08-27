{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "neuvector.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "neuvector.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "neuvector.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Lookup secret.
*/}}
{{- define "neuvector.secrets.lookup" -}}
{{- $value := .defaultValue | toString | b64enc -}}
{{- printf "%s" $value -}}
{{- end -}}

{{- define "neuvector.controller.image" -}}
{{- if .Values.global.azure.enabled }}
  {{- printf "%s/%s:%s" .Values.global.azure.images.controller.registry .Values.global.azure.images.controller.image .Values.global.azure.images.controller.tag }}
{{- else }}
  {{- if eq .Values.registry "registry.neuvector.com" }}
    {{- if .Values.oem }}
      {{- printf "%s/%s/controller:%s" .Values.registry .Values.oem .Values.tag }}
    {{- else }}
      {{- printf "%s/controller:%s" .Values.registry .Values.tag }}
    {{- end }}
  {{- else }}
    {{- if .Values.controller.image.hash }}
      {{- printf "%s/%s@%s" .Values.registry .Values.controller.image.repository .Values.controller.image.hash }}
    {{- else }}
      {{- printf "%s/%s:%s" .Values.registry .Values.controller.image.repository .Values.tag }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end -}}
