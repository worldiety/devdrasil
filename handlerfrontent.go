package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type frontendHandler struct {
	server *Devdrasil
}

func installFrontendHandler(server *Devdrasil) {
	handler := &frontendHandler{}

	//directly server the wwt from the working dir
	wwtDir := filepath.Join(server.cwd, "wwt", "wwt")
	assertDir(wwtDir)
	log.Printf("wwt served from %s\n", wwtDir)
	http.Handle("/wwt/", NoCache(http.StripPrefix("/wwt", http.FileServer(http.Dir(wwtDir)))))

	//serve the frontend javascript app
	appDir := filepath.Join(server.cwd, "devdrasil", "frontend")
	assertDir(appDir)
	log.Printf("frontend served from %s\n", appDir)
	http.Handle("/frontend/", NoCache(http.StripPrefix("/frontend", http.FileServer(http.Dir(appDir)))))

	//everything else redirects to the html root
	server.mux.HandleFunc("/", handler.handleRoot)
}

func (h *frontendHandler) handleRoot(writer http.ResponseWriter, request *http.Request) {

	sb := strings.Builder{}
	sb.WriteString(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"`)
	sb.WriteString("\n")

	sb.WriteString(`"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`)
	sb.WriteString("\n")

	sb.WriteString(`<html xmlns="http://www.w3.org/1999/xhtml">`)
	sb.WriteString("\n")

	sb.WriteString(`<head>`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="HandheldFriendly" content="true"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="apple-mobile-web-app-capable" content="yes"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="apple-mobile-web-app-status-bar-style" content="green"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="theme-color" content="green">`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="apple-mobile-web-app-title" content="demo wwt app">`)
	sb.WriteString("\n")

	sb.WriteString(`<meta name="viewport" content="width=device-width, minimum-scale=1.0,initial-scale=1 maximum-scale=1 user-scalable=0 minimal-ui"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<link rel="stylesheet" href="/wwt/mcw.min.css"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<link rel="stylesheet" href="/wwt/material_icons.css">`)
	sb.WriteString("\n")

	sb.WriteString(`<link rel="stylesheet" href="/frontend/custom.css">`)
	sb.WriteString("\n")

	sb.WriteString(`<script type="text/javascript" src="/wwt/mcw.min.js"></script>`)
	sb.WriteString("\n")

	sb.WriteString(`<link rel="stylesheet" href="/wwt/fixes.css"/>`)
	sb.WriteString("\n")

	sb.WriteString(`<script src="/frontend/app.js" type="module"></script>`)
	sb.WriteString("\n")

	sb.WriteString(`</head>`)
	sb.WriteString("\n")

	sb.WriteString(`<body>`)
	sb.WriteString("\n")

	sb.WriteString(`</body>`)
	sb.WriteString("\n")

	sb.WriteString(`</html>`)

	writer.Header().Add("Content-Type", "text/html;charset=utf-8")
	writer.Write([]byte(sb.String()))
}

func assertDir(dir string) {
	stat, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("folder expected in '%s' but is missing or inaccesible: %s\n", dir, err)
	}

	if !stat.IsDir() {
		log.Fatalf("folder '%s' is not a directory\n", dir)
	}
}
