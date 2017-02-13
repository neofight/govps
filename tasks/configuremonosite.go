package tasks

import (
	"encoding/xml"
	"fmt"
)

type ConfigureMonoSite struct {
}

func (task ConfigureMonoSite) Execute(cxt Context) error {

	file, err := cxt.VPS.DownloadFile("/etc/xsp4/debian.webapp")

	if err != nil {
		return fmt.Errorf("Failed to download Mono FastCGI configuration: %v", err)
	}

	var config apps

	err = xml.Unmarshal([]byte(file), &config)

	if err != nil {
		return fmt.Errorf("Failed to parse Mono FastCGI configuration: %v", err)
	}

	for _, app := range config.Apps {
		if app.Name == cxt.Domain {
			fmt.Println("Site found in Mono FastCGI configuration")
			return nil
		}
	}

	config.Apps = append(config.Apps, webApplication{cxt.Domain, cxt.Domain, 80, "/", "/var/www/" + cxt.Domain, true})

	data, err := xml.MarshalIndent(&config, "", "	")

	if err != nil {
		return fmt.Errorf("Failed to generate xml for Mono FastCGI configuration: %v", err)
	}

	cxt.VPS.UploadData(string(data), "/etc/xsp4/debian.webapp")

	if err != nil {
		return fmt.Errorf("Failed to upload Mono FastCGI configuration: %v", err)
	}

	fmt.Println("Site added to Mono FastCGI configuration")

	return nil
}

type apps struct {
	Apps []webApplication `xml:"web-application"`
}

type webApplication struct {
	Name    string `xml:"name"`
	VHost   string `xml:"vhost"`
	VPort   int    `xml:"vport"`
	VPath   string `xml:"vpath"`
	Path    string `xml:"path"`
	Enabled bool   `xml:"enabled"`
}
