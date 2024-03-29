/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"kubevault.dev/cli/pkg/cmds"

	"github.com/spf13/cobra/doc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gomodules.xyz/runtime"
)

var (
	tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Reference | KubeVault CLI
description: KubeVault CLI Reference
menu:
  docs_{{ "{{ .version }}" }}:
    identifier: reference-cli
    name: KubeVault CLI
    weight: 30
    parent: reference
menu_name: docs_{{ "{{ .version }}" }}
---
`))

	_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  docs_{{ "{{ .version }}" }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: reference-cli
{{- if .RootCmd }}
    weight: 0
{{ end }}
menu_name: docs_{{ "{{ .version }}" }}
section_menu_id: reference
{{- if .RootCmd }}
url: /docs/{{ "{{ .version }}" }}/reference/cli/
aliases:
- /docs/{{ "{{ .version }}" }}/reference/cli/{{ .ID }}/
{{- end }}
---
`))
)

func docsDir() string {
	if dir, ok := os.LookupEnv("DOCS_ROOT"); ok {
		return dir
	}
	return runtime.GOPath() + "/src/kubevault.dev/docs"
}

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	rootCmd := cmds.NewRootCmd()
	dir := filepath.Join(docsDir(), "docs", "reference", "cli")
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		log.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		data := struct {
			ID      string
			Name    string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			cases.Title(language.English).String(strings.Replace(base, "_", " ", -1)),
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			log.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/cli/" + name
	}
	err = doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)
	if err != nil {
		log.Fatalln(err)
	}

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{}{})
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		log.Fatalln(err)
	}
}
