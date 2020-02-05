apiVersion: v1
clusters:
- cluster:
{{- if eq .SkipTLSVerify true }}
    insecure-skip-tls-verify: true
{{- else if .Certificate }}
    certificate-authority-data: {{ .Certificate }}
{{- end}}
    server: {{ .APIServer }}
  name: helm
contexts:
- context:
    cluster: helm
{{- if .Namespace }}
    namespace: {{ .Namespace }}
{{- end }}
    user: {{ .ServiceAccount }}
  name: helm
current-context: "helm"
kind: Config
preferences: {}
users:
- name: {{ .ServiceAccount }}
  user:
{{- if .Token }}
    token: {{ .Token }}
{{- end }}
{{- if .ClientCertificate }}
    client-certificate-data: {{ .ClientCertificate }}
{{- end }}
{{- if .ClientKey }}
    client-key-data: {{ .ClientKey }}
{{- end }}
