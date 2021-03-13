package main

func Headers() map[string]string {
	headers := make(map[string]string)
	headers["Accept"] = "ext/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	// empty here
	headers["Accept-Encoding"] = ""
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	headers["Host"] = "music.163.com"
	headers["Cache-Control"] = "no-cache"
	headers["Connection"] = "keep-alive"
	headers["Pragma"] = "no-cache"
	headers["Origin"] = "https://music.163.com"
	headers["Accept"] = "ext/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	return headers
}
