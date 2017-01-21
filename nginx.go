package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func createNginxHostConfiguration(domain string) (string, error) {

	template, err := template.New("nginx").Parse(nginxTemplate)

	if err != nil {
		return "", fmt.Errorf("Unable to parse Nginx Host Configuration template: %v", err)
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, domain)

	return buffer.String(), nil
}

var nginxTemplate = `server {
        listen   80;
        server_name  {{.}};
        access_log   /var/log/nginx/{{.}}.access.log;
        root /var/www/{{.}}/;

        location / {
                index index.html index.htm default.aspx Default.aspx;
                fastcgi_index Home;
                fastcgi_pass 127.0.0.1:9000;
                include /etc/nginx/fastcgi_params;
        }
}`
