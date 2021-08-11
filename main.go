package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Fetch environment variables
	port := os.Getenv("PORT")

	// Set defaults for missing variables
	if port == "" {
		port = ":8442"
	}

	// Establish endpoint for files
	http.HandleFunc("/", router)

	// Configure and start server
	s := &http.Server{
		Addr: port,
	}
	log.Printf("Inari is ready. Listening on port " + port + ".")
	log.Fatal(s.ListenAndServe())
}

func router(w http.ResponseWriter, r *http.Request) {
	// Fetch envvars and set defaults
	cacheImageTime := os.Getenv("CACHE_IMAGE_TIME")
	cacheFontTime := os.Getenv("CACHE_FONT_TIME")
	webRoot := os.Getenv("DIRECTORY")
	if webRoot == "" {
		a, err := os.Getwd()
		if err != nil {
			log.Fatal("Couldn't fetch directory. More information may follow: ", err)
		}
		webRoot = a // Work around due to needing to define `err`
	}
	if cacheImageTime == "" {
		cacheImageTime = "max-age=7200"
	} else {
		cacheImageTime = "max-age=" + cacheImageTime
	}
	if cacheFontTime == "" {
		cacheFontTime = "max-age=7200"
	} else {
		cacheFontTime = "max-age=" + cacheFontTime
	}

	// Add HTTP headers
	var images = regexp.MustCompile(".(png)|.(jpg)|.(webp)|.(gif)|.(svg)")
	switch {
	case images.MatchString(r.URL.Path):
		// Set headers
		defaultHeads(w)
		w.Header().Set("Cache-Control", cacheImageTime)
	default:
		// Serve with default headers
		defaultHeads(w)
	}

	// Compress content and serve
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		// Set headers
		w.Header().Set("Content-Encoding", "gzip")

		// Compress data
		name := webRoot + r.URL.Path
		file, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		var b bytes.Buffer
		gzipWriter := gzip.NewWriter(&b, gzip.DefaultCompression)
		defer gzipWriter.Close()
		if _, err := gzipWriter.Write(data); err != nil {
			log.Fatal(err)
		}

		// Serve compressed content
		http.ServeContent(w, r, name, time.Now(), b.Bytes())
		return
	} else if strings.Contains(acceptEncoding, "defalte") {
		// Set headers
		w.Header().Set("Content-Encoding", "deflate")

		// Compress data
		name := webRoot + r.URL.Path
		file, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		var b bytes.Buffer
		flateWriter, _ := flate.NewWriter(&b, flate.DefaultCompression)
		defer flateWriter.Close()
		if _, err := flateWriter.Write(data); err != nil {
			log.Fatal(err)
		}

		// Serve compressed content
		http.ServeContent(w, r, name, time.Now(), b.Bytes())
		return
	} else {
		// Prepare files
		name := webRoot + r.URL.Path
		file, err := os.Open(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panicf("Unable to read file to display. More info may follow: ", err)
			return
		}
		defer file.Close()

		// Serve content
		http.ServeContent(w, r, name, time.Now(), file)
		return
	}
}

func defaultHeads(w http.ResponseWriter) {
	// Fetch envvars
	hstsDisable := os.Getenv("NOHSTS")
	unsafeFraming := os.Getenv("UNSAFE_FRAME")
	devel := os.Getenv("DEVELOPMENT")
	var noHSTS bool
	var frame string

	// Set defaults for envvars
	if hstsDisable == "true" {
		noHSTS = true
	}
	if unsafeFraming == "true" {
		frame = "ALLOW"
	} else {
		frame = "DENY"
	}

	// Setup HSTS
	if noHSTS == false {
		if devel == "" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000;preload;includeSubDomains")
		} else {
			w.Header().Set("Strict-Transport-Security", "max-age=3400;includeSubDomains")
		}
	}

	// Set default headers
	w.Header().Set("Permissions-Policy", "interest-cohort=()")
	w.Header().Set("Vary", "Origin")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", frame)
	w.Header().Set("X-XSS-Protection", "1; mode=block")
}
