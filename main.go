package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/NYTimes/gziphandler"
)

func main() {
	// Fetch environment variables
	port := os.Getenv("PORT")

	// Set defaults for missing variables
	if port == "" {
		port = ":8442"
	}

	// Configure and start server
	s := &http.Server{
		Addr: port,
	}

	// Fetch envvars and set defaults
	cacheImageTime := os.Getenv("CACHE_IMAGE_TIME")
	cacheFontTime := os.Getenv("CACHE_FONT_TIME")
	webRoot := os.Getenv("DIRECTORY")
	if webRoot == "" {
		log.Printf("Using default directory")
		webRoot = "."
	}
	if cacheImageTime == "" {
		cacheImageTime = "max-age=7200"
	} else {
		log.Printf("Using custom image cache time")
		cacheImageTime = "max-age=" + cacheImageTime
	}
	if cacheFontTime == "" {
		cacheFontTime = "max-age=7200"
	} else {
		log.Printf("Using custom font cache time")
		cacheFontTime = "max-age=" + cacheFontTime
	}

	// Establish endpoint for files and use Gzip
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add HTTP headers
		var images = regexp.MustCompile(".(png)|.(jpg)|.(webp)|.(gif)|.(svg)")
		var fonts = regexp.MustCompile(".(ttf)|.(otf)|.(woff)|.(woff2)|.(eot)")
		if images.MatchString(r.URL.Path) {
			// Set headers
			defaultHeads(w)
			w.Header().Set("Cache-Control", cacheImageTime)
		} else if fonts.MatchString(r.URL.Path) {
			// Set headers
			defaultHeads(w)
			w.Header().Set("Cache-Control", cacheFontTime)
		} else {
			// Serve with default headers
			defaultHeads(w)
		}

		// Use GZip where possible
		fs := http.FileServer(http.Dir(webRoot))
		typ := mime.TypeByExtension(r.URL.Path)
		switch {
		case strings.HasPrefix(typ, "text/"):
			fallthrough
		case typ == "application/xml":
			fallthrough
		case typ == "":
			fs = gziphandler.GzipHandler(fs)
		}

		fs.ServeHTTP(w, r)
	})

	log.Printf("Inari is ready. Listening on port " + port + ".")
	log.Fatal(s.ListenAndServe())
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
