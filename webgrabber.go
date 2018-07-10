package main

import (
	"fmt"
	"log"

	"flag"
	"github.com/PuerkitoBio/goquery"
	. "github.com/gtechx/base/collections"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//var scheme string = "http://"
//var host string =  "www.critterai.org"
var gurldata *url.URL
var err error
var externaldir string = "/external/"
var urlset *Set
var outputdir string = "tmp/"

func GrabData(rawurl string) {
	//check if is ""
	if rawurl == "" || urlset.Has(rawurl) {
		return
	}
	urlset.Add(rawurl)
	urldata, err := url.Parse(rawurl)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//if external url , then return
	if urldata.Host != gurldata.Host {
		fmt.Println("external url and return")
		return
	}

	//fmt.Println(htmlstr)

	//save htmlstr
	fpath := urldata.Path
	if pos := strings.LastIndex(fpath, "."); pos < 0 {
		if fpath == "" || string(fpath[len(fpath)-1]) != "/" {
			fpath += "/"
		}
		fpath += "index.html"
	}
	if string(fpath[0]) != "/" {
		fpath = "./" + fpath
	} else {
		fpath = "." + fpath
	}
	// if PathExists(fpath) {
	// 	return
	// }
	fpath = strings.Replace(fpath, "//", "/", -1)
	fmt.Println("fpath", fpath)
	if !PathExists(fpath) {
		dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}
		// setup a http client
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		// set our socks5 as the dialer
		httpTransport.Dial = dialer.Dial

		resp, err := httpClient.Get(rawurl)

		if resp == nil {
			return
		}

		if resp.StatusCode != 200 {
			return
		}

		if err != nil {
			return
		}

		doc, err := goquery.NewDocumentFromResponse(resp)

		if err != nil {
			return
		}

		pathext := "./"
		// tmppath := fpath
		// if tmppath != "/" && tmppath != "\\" && tmppath != "./" {
		// 	if len(tmppath) >= 1 && (string(tmppath[0]) == "/" || string(tmppath[0]) == "\\") {
		// 		tmppath = tmppath[1:]
		// 	} else if len(tmppath) >= 1 && (string(tmppath[len(tmppath)-1]) == "/" || string(tmppath[len(tmppath)-1]) == "\\") {
		// 		tmppath = tmppath[0 : len(tmppath)-1]
		// 	}
		// }

		patharr := strings.Split(fpath, "/")

		if len(patharr) > 2 {
			for i := 1; i < len(patharr)-1; i++ {
				pathext += "../"
			}
		}

		pos := strings.LastIndex(urldata.Path, ".")
		posslash := strings.LastIndex(urldata.Path, "/")
		if pos > 0 && pos > posslash {
			urldata.Path = urldata.Path[:posslash+1]
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			href, _ := s.Attr("href")
			if href != "" {
				aurldata, _ := url.Parse(href)
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					//s.SetAttr("href", externaldir+aurldata.String())
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					//fragment := aurldata.Fragment
					//aurldata.Fragment = ""
					if aurldata.Path != "" || (aurldata.Path == "" && aurldata.Fragment == "") {
						pos := strings.LastIndex(aurldata.Path, ".")
						posslash := strings.LastIndex(aurldata.Path, "/")
						if pos < 0 || pos < posslash {
							aurldata.Path += "index.html"
						}
					}

					// if fragment != "" {
					// 	aurldata.Fragment = fragment
					// }
					//s.SetAttr("href", pathext+aurldata.String())
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					fmt.Println("surl", surl)
					s.SetAttr("href", surl)
				}
			}
		})

		doc.Find("frame").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			href, _ := s.Attr("src")
			if href != "" {
				aurldata, _ := url.Parse(href)
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					//s.SetAttr("href", externaldir+aurldata.String())
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					//fragment := aurldata.Fragment
					//aurldata.Fragment = ""
					if aurldata.Path != "" || (aurldata.Path == "" && aurldata.Fragment == "") {
						pos := strings.LastIndex(aurldata.Path, ".")
						posslash := strings.LastIndex(aurldata.Path, "/")
						if pos < 0 || pos < posslash {
							aurldata.Path += "index.html"
						}
					}
					// if fragment != "" {
					// 	aurldata.Fragment = fragment
					// }
					//s.SetAttr("href", pathext+aurldata.String())
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					s.SetAttr("src", surl)
				}
			}
		})

		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			src, _ := s.Attr("src")
			if src != "" {
				aurldata, _ := url.Parse(src)
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					aurldata.Host = ""
					aurldata.Scheme = ""
					pos := strings.LastIndex(aurldata.Path, "/")
					surl := strings.Replace(pathext+externaldir+aurldata.Path[pos+1:], "//", "/", -1)
					s.SetAttr("src", surl)
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					aurldata.Fragment = ""
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					s.SetAttr("src", surl)
				}
			}
		})

		doc.Find("link").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			href, _ := s.Attr("href")
			if href != "" {
				aurldata, _ := url.Parse(href)
				pos := strings.LastIndex(href, "?")
				if pos > 0 {
					if p := strings.LastIndex(aurldata.Path, "."); p <= 0 {
						aurldata.Path += ".css"
					}
				}
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					aurldata.Host = ""
					aurldata.Scheme = ""
					pos := strings.LastIndex(aurldata.Path, "/")
					surl := strings.Replace(pathext+externaldir+aurldata.Path[pos+1:], "//", "/", -1)
					s.SetAttr("href", surl)
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					aurldata.Fragment = ""
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					s.SetAttr("href", surl)
					//fmt.Println(aurldata.String())
				}
			}
		})

		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			src, _ := s.Attr("src")
			if src != "" {
				aurldata, _ := url.Parse(src)
				pos := strings.LastIndex(src, "?")
				if pos > 0 {
					if p := strings.LastIndex(aurldata.Path, "."); p <= 0 {
						aurldata.Path += ".js"
					}
				}
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					aurldata.Host = ""
					aurldata.Scheme = ""
					pos := strings.LastIndex(aurldata.Path, "/")
					surl := strings.Replace(pathext+externaldir+aurldata.Path[pos+1:], "//", "/", -1)
					s.SetAttr("src", surl)
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					aurldata.Fragment = ""
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					s.SetAttr("src", surl)
				}
			}
		})

		doc.Find("input").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			src, _ := s.Attr("src")
			if src != "" {
				aurldata, _ := url.Parse(src)
				if aurldata.Host != "" && aurldata.Host != urldata.Host {
					//strings.LastIndex(src, ".")
					aurldata.Host = ""
					aurldata.Scheme = ""
					pos := strings.LastIndex(aurldata.Path, "/")
					surl := strings.Replace(pathext+externaldir+aurldata.Path[pos+1:], "//", "/", -1)
					s.SetAttr("src", surl)
				} else {
					aurldata.Host = ""
					aurldata.Scheme = ""
					aurldata.RawQuery = ""
					aurldata.Fragment = ""
					surl := ""
					if aurldata.Path == "" || string(aurldata.Path[0:1]) != "/" {
						surl = aurldata.String()
						//s.SetAttr("href", aurldata.String())
					} else {
						surl = pathext + aurldata.String()
						//s.SetAttr("href", pathext+aurldata.String())
					}
					surl = strings.Replace(surl, "//", "/", -1)
					s.SetAttr("src", surl)
				}
			}
		})

		fmt.Println(fpath)
		htmlstr, _ := doc.Html()
		saveFile(fpath, htmlstr)
	}

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial

	resp, err := httpClient.Get(rawurl)

	if resp == nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		return
	}

	// Find script items
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, _ := s.Attr("src")
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		//s.SetText("111")
		//s.SetAttr("href","222")
		//need save
		if src != "" && src != "/" && src != "//" && src != "./" {
			aurldata, _ := url.Parse(src)
			fmt.Println("111111", src)
			if aurldata.Scheme != "" || aurldata.Host != "" {
				saveScriptFile(src)
			} else {
				aurldata.Host = urldata.Host
				aurldata.Scheme = urldata.Scheme
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					if urldata.Path != "" {
						pos := strings.LastIndex(urldata.Path, "/")
						if pos > 0 {
							aurldata.Path = urldata.Path[:pos] + "/" + aurldata.Path
						}
					} else {
						aurldata.Path = "/" + aurldata.Path
					}
				}
				//fmt.Println("111", aurldata.String())

				saveScriptFile(aurldata.String())
			}
		}
	})

	// Find link items
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, _ := s.Attr("href")
		//s.SetText("111")
		//s.SetAttr("href","222")
		//need save
		if href != "" && href != "/" && href != "//" && href != "./" {
			aurldata, _ := url.Parse(href)
			if aurldata.Scheme != "" || aurldata.Host != "" {
				saveCssFile(href)
			} else {
				aurldata.Host = urldata.Host
				aurldata.Scheme = urldata.Scheme
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					if urldata.Path != "" {
						pos := strings.LastIndex(urldata.Path, "/")
						if pos > 0 {
							aurldata.Path = urldata.Path[:pos] + "/" + aurldata.Path
						}
					} else {
						aurldata.Path = "/" + aurldata.Path
					}
				}
				saveCssFile(aurldata.String())
			}
		}
	})

	// Find img items
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, _ := s.Attr("src")
		// filepath := strings.Replace(src, "http://www.critterai.org", "", -1)
		// srcdoc, _ := goquery.NewDocument("http://www.critterai.org" + filepath)
		// saveAssetFile(filepath, srcdoc)
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		//s.SetText("111")
		//s.SetAttr("href","222")
		if src != "" && src != "/" && src != "//" && src != "./" {
			aurldata, _ := url.Parse(src)
			if aurldata.Scheme != "" || aurldata.Host != "" {
				saveAssetFile(src)
			} else {
				aurldata.Host = urldata.Host
				aurldata.Scheme = urldata.Scheme
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					if urldata.Path != "" {
						pos := strings.LastIndex(urldata.Path, "/")
						if pos > 0 {
							aurldata.Path = urldata.Path[:pos] + "/" + aurldata.Path
						}
					} else {
						aurldata.Path = "/" + aurldata.Path
					}
				}
				saveAssetFile(aurldata.String())
			}
		}
	})

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		src, _ := s.Attr("src")
		if src != "" && src != "/" && src != "//" && src != "./" {
			aurldata, _ := url.Parse(src)
			if aurldata.Scheme != "" || aurldata.Host != "" {
				saveAssetFile(src)
			} else {
				aurldata.Host = urldata.Host
				aurldata.Scheme = urldata.Scheme
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					if urldata.Path != "" {
						pos := strings.LastIndex(urldata.Path, "/")
						if pos > 0 {
							aurldata.Path = urldata.Path[:pos] + "/" + aurldata.Path
						}
					} else {
						aurldata.Path = "/" + aurldata.Path
					}
				}
				saveAssetFile(aurldata.String())
			}
		}
	})

	// Find a items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, _ := s.Attr("href")
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		//s.SetText("111")
		//s.SetAttr("href","222")
		//check if start with #
		if href != "" && href != "/" && href != "//" && href != "./" {
			if string(href[0]) != "#" {
				// if pos := strings.LastIndex(href, "#"); pos > 0 {
				// 	href = string(href[:pos])
				// }
				aurldata, _ := url.Parse(href)

				if aurldata.Scheme == "" && aurldata.Host == "" {
					aurldata.Scheme = urldata.Scheme
					aurldata.Host = urldata.Host
				}
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					fmt.Println("1...", urldata.Path)
					fmt.Println("2...", aurldata.Path)
					pos := strings.LastIndex(urldata.Path, ".")
					posslash := strings.LastIndex(urldata.Path, "/")
					if pos > 0 && pos > posslash {
						urldata.Path = urldata.Path[:posslash]
						fmt.Println("3...", urldata.Path)
					}
					// if string(urldata.Path[0] == "/"){
					// 	urldata.Path = urldata.Path[1:]
					// }
					if len(urldata.Path) > 2 && string(urldata.Path[len(urldata.Path)-1]) == "/" {
						urldata.Path = urldata.Path[:len(urldata.Path)-1]
					}
					parentpath := "/"
					if aurldata.Path != "" && string(aurldata.Path[0:2]) != "./" && string(aurldata.Path[0:2]) == ".." {
						pararr := strings.Split(urldata.Path, "/")
						n := checkpath(aurldata.Path)
						for i, str := range pararr {
							if i < len(pararr)-n {
								if str != "" {
									parentpath += str + "/"
								}
							} else {
								break
							}
						}
					} else if aurldata.Path != "" && string(aurldata.Path[0:2]) == "./" && string(aurldata.Path[3:5]) == ".." {
						pararr := strings.Split(urldata.Path, "/")
						n := checkpath(aurldata.Path[3:])
						for i, str := range pararr {
							if i < len(pararr)-n {
								if str != "" {
									parentpath += str + "/"
								}
							} else {
								break
							}
						}
					} else {
						parentpath = urldata.Path + "/"
					}
					aurldata.Path = parentpath + removedot(aurldata.Path)
					// if urldata.Path != "" && string(urldata.Path[len(urldata.Path)-1]) != "/" {
					// 	aurldata.Path = parentpath + "/" + aurldata.Path
					// } else if urldata.Path != "" {
					// 	aurldata.Path = parentpath + aurldata.Path
					// } else {
					// 	aurldata.Path = "/" + aurldata.Path
					// }
				}

				aurldata.Fragment = ""
				ext := getext(aurldata.Path)
				if ext != "" && ext != "html" && ext != "htm" && ext != "php" && ext != "jsp" && ext != "asp" && ext != "aspx" {
					aurldata.Path = strings.Replace(aurldata.Path, "/./", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "////", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "///", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "//", "/", -1)
					surl := aurldata.String()
					//surl = strings.Replace(surl, "//", "/", -1)
					fmt.Println("save a href asset", surl)

					if ext == "js" {
						saveScriptFile(surl)
					} else if ext == "css" {
						saveCssFile(surl)
					} else {
						saveAssetFile(surl)
					}
				} else {
					aurldata.Path = strings.Replace(aurldata.Path, "/./", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "////", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "///", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "//", "/", -1)
					surl := aurldata.String()
					//surl = strings.Replace(surl, "//", "/", -1)
					fmt.Println("grab sub page", surl)
					GrabData(surl)
				}
			}
		}
	})

	doc.Find("frame").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, _ := s.Attr("src")
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		//s.SetText("111")
		//s.SetAttr("href","222")
		//check if start with #
		if href != "" && href != "/" && href != "//" && href != "./" {
			if string(href[0]) != "#" {
				// if pos := strings.LastIndex(href, "#"); pos > 0 {
				// 	href = string(href[:pos])
				// }
				aurldata, _ := url.Parse(href)

				if aurldata.Scheme == "" && aurldata.Host == "" {
					aurldata.Scheme = urldata.Scheme
					aurldata.Host = urldata.Host
				}
				if aurldata.Path != "" && string(aurldata.Path[0]) != "/" {
					fmt.Println("1...", urldata.Path)
					fmt.Println("2...", aurldata.Path)
					pos := strings.LastIndex(urldata.Path, ".")
					posslash := strings.LastIndex(urldata.Path, "/")
					if pos > 0 && pos > posslash {
						urldata.Path = urldata.Path[:posslash]
						fmt.Println("3...", urldata.Path)
					}
					// if string(urldata.Path[0] == "/"){
					// 	urldata.Path = urldata.Path[1:]
					// }
					if len(urldata.Path) > 2 && string(urldata.Path[len(urldata.Path)-1]) == "/" {
						urldata.Path = urldata.Path[:len(urldata.Path)-1]
					}
					parentpath := "/"
					if aurldata.Path != "" && string(aurldata.Path[0:2]) != "./" && string(aurldata.Path[0:2]) == ".." {
						pararr := strings.Split(urldata.Path, "/")
						n := checkpath(aurldata.Path)
						for i, str := range pararr {
							if i < len(pararr)-n {
								if str != "" {
									parentpath += str + "/"
								}
							} else {
								break
							}
						}
					} else if aurldata.Path != "" && string(aurldata.Path[0:2]) == "./" && string(aurldata.Path[3:5]) == ".." {
						pararr := strings.Split(urldata.Path, "/")
						n := checkpath(aurldata.Path[3:])
						for i, str := range pararr {
							if i < len(pararr)-n {
								if str != "" {
									parentpath += str + "/"
								}
							} else {
								break
							}
						}
					} else {
						parentpath = urldata.Path + "/"
					}
					aurldata.Path = parentpath + removedot(aurldata.Path)
					// if urldata.Path != "" && string(urldata.Path[len(urldata.Path)-1]) != "/" {
					// 	aurldata.Path = parentpath + "/" + aurldata.Path
					// } else if urldata.Path != "" {
					// 	aurldata.Path = parentpath + aurldata.Path
					// } else {
					// 	aurldata.Path = "/" + aurldata.Path
					// }
				}

				aurldata.Fragment = ""
				ext := getext(aurldata.Path)
				if ext != "" && ext != "html" && ext != "htm" && ext != "php" && ext != "jsp" && ext != "asp" && ext != "aspx" {
					aurldata.Path = strings.Replace(aurldata.Path, "/./", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "////", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "///", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "//", "/", -1)
					surl := aurldata.String()
					//surl = strings.Replace(surl, "//", "/", -1)
					fmt.Println("save a href asset", surl)

					if ext == "js" {
						saveScriptFile(surl)
					} else if ext == "css" {
						saveCssFile(surl)
					} else {
						saveAssetFile(surl)
					}
				} else {
					aurldata.Path = strings.Replace(aurldata.Path, "/./", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "////", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "///", "/", -1)
					aurldata.Path = strings.Replace(aurldata.Path, "//", "/", -1)
					surl := aurldata.String()
					//surl = strings.Replace(surl, "//", "/", -1)
					fmt.Println("grab sub page", surl)
					GrabData(surl)
				}
			}
		}
	})
}

func getext(pathstr string) string {
	pos := strings.LastIndex(pathstr, ".")
	posslash := strings.LastIndex(pathstr, "/")
	if pos > 0 && pos > posslash {
		return pathstr[pos+1:]
	} else {
		return ""
	}
}

func removedot(pathstr string) string {
	if len(pathstr) < 2 {
		return pathstr
	}
	if string(pathstr[0:2]) == ".." {
		if string(pathstr[2]) == "/" && len(pathstr) > 3 {
			return removedot(pathstr[3:])
		} else if len(pathstr) > 3 {
			return pathstr[3:]
		} else {
			return ""
		}
	} else {
		if string(pathstr[0:2]) == "./" && len(pathstr) > 2 {
			return removedot(pathstr[2:])
		} else {
			return pathstr
		}
	}
}

func checkpath(pathstr string) int {
	if string(pathstr[0:2]) == ".." {
		if string(pathstr[2]) == "/" {
			return 1 + checkpath(pathstr[3:])
		} else {
			return 1
		}
	} else {
		return 0
	}
}

func saveScriptFile(rawurl string) {
	if rawurl == "" || urlset.Has(rawurl) {
		return
	}
	urlset.Add(rawurl)
	//fmt.Println("121", rawurl)
	//parse rawurl
	urldata, err := url.Parse(rawurl)

	if err != nil {
		return
	}

	if urldata.Path == "" || urldata.Path == "/" {
		return
	}
	// fmt.Println("Scheme:", urldata.Scheme)
	// fmt.Println("Opaque:", urldata.Opaque)
	// fmt.Println("User:", urldata.User)
	// fmt.Println("Host:", urldata.Host)
	// fmt.Println("Path:", urldata.Path)
	// fmt.Println("RawPath:", urldata.RawPath)
	// fmt.Println("ForceQuery:", urldata.ForceQuery)
	// fmt.Println("RawQuery:", urldata.RawQuery)
	// fmt.Println("Fragment:", urldata.Fragment)
	// fmt.Println("112", urldata.String())
	fpath := urldata.Path
	fmt.Println("1", fpath)
	if urldata.Scheme == "" {
		urldata.Scheme = "http"
	}
	fmt.Println("12", urldata.String())
	// doc, err := goquery.NewDocument(urldata.String())

	// if err != nil {
	// 	fmt.Println("error:", err.Error())
	// 	return
	// }

	// htmlstr, _ := doc.Html()

	//if start with http:// or https://

	//if external url
	if urldata.Host != gurldata.Host {
		//strings.LastIndex(src, ".")
		pos := strings.LastIndex(urldata.Path, "/")
		if pos < 0 {
			pos = 0
		}
		fpath = externaldir + urldata.Path[pos+1:]
	}

	pos := strings.LastIndex(fpath, "/")
	if p := strings.LastIndex(fpath, "."); p < 0 || p < pos {
		fpath += ".js"
	}

	if string(fpath[0]) != "/" {
		fpath = "./" + fpath
	} else {
		fpath = "." + fpath
	}

	fpath = outputdir + fpath
	pos = strings.LastIndex(fpath, "/")

	if PathExists(fpath) {
		return
	} else {
		newpos := strings.Index(fpath, "/")
		if pos != newpos && !PathExists(fpath[:pos+1]) {
			mkdir(fpath[:pos+1])
		}
	}

	//if start with /
	fmt.Println(fpath)
	f, err1 := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial

	resp, err := httpClient.Get(urldata.String())
	if err != nil {
		// handle error
		fmt.Println("error:", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//f.WriteString(htmlstr)
	f.Write(body)
}

func saveCssFile(rawurl string) {
	if rawurl == "" || urlset.Has(rawurl) {
		return
	}
	urlset.Add(rawurl)
	//parse rawurl
	urldata, err := url.Parse(rawurl)

	if err != nil {
		return
	}

	if urldata.Path == "" || urldata.Path == "/" {
		return
	}

	fpath := urldata.Path
	fmt.Println("1", fpath)
	if urldata.Scheme == "" {
		urldata.Scheme = "http"
	}
	// doc, err := goquery.NewDocument(urldata.String())

	// if err != nil {
	// 	fmt.Println("error:", err.Error())
	// 	return
	// }

	// htmlstr, _ := doc.Html()

	//if start with http:// or https://

	//if external url
	if urldata.Host != gurldata.Host {
		//strings.LastIndex(src, ".")
		pos := strings.LastIndex(urldata.Path, "/")
		if pos < 0 {
			pos = 0
		}
		fpath = externaldir + urldata.Path[pos:]
	}

	pos := strings.LastIndex(fpath, "/")
	if p := strings.LastIndex(fpath, "."); p < 0 || p < pos {
		fpath += ".css"
	}

	if string(fpath[0]) != "/" {
		fpath = "./" + fpath
	} else {
		fpath = "." + fpath
	}

	fpath = outputdir + fpath
	pos = strings.LastIndex(fpath, "/")
	// pos := strings.LastIndex(fpath, "/")
	// if p := strings.LastIndex(fpath, "."); p < 0 && p != 0 {
	// 	fpath += ".css"
	// }
	if PathExists(fpath) {
		return
	} else {
		newpos := strings.Index(fpath, "/")
		if pos != newpos && !PathExists(fpath[:pos+1]) {
			mkdir(fpath[:pos+1])
		}
	}

	//if start with /
	fmt.Println(fpath)
	f, err1 := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial

	resp, err := httpClient.Get(urldata.String())
	if err != nil {
		// handle error
		fmt.Println("error:", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//f.WriteString(htmlstr)
	f.Write(body)
}

func saveAssetFile(rawurl string) {
	if rawurl == "" || urlset.Has(rawurl) {
		return
	}
	urlset.Add(rawurl)
	//parse rawurl
	urldata, err := url.Parse(rawurl)

	if err != nil {
		return
	}

	if urldata.Path == "" || urldata.Path == "/" {
		return
	}

	fpath := urldata.Path
	fmt.Println("1", fpath)
	if urldata.Scheme == "" {
		urldata.Scheme = "http"
	}
	// doc, err := goquery.NewDocument(urldata.String())

	// if err != nil {
	// 	fmt.Println("error:", err.Error())
	// 	return
	// }

	// htmlstr, _ := doc.Html()

	//if start with http:// or https://

	//if external url
	if urldata.Host != gurldata.Host {
		//strings.LastIndex(src, ".")
		pos := strings.LastIndex(urldata.Path, "/")
		if pos < 0 {
			pos = 0
		}
		fpath = externaldir + urldata.Path[pos:]
	}

	if string(fpath[0]) != "/" {
		fpath = "./" + fpath
	} else {
		fpath = "." + fpath
	}

	fpath = outputdir + fpath
	pos := strings.LastIndex(fpath, "/")
	// if p := strings.LastIndex(fpath, "."); p < 0 {
	// 	fpath += ".js"
	// }
	if PathExists(fpath) {
		return
	} else {
		newpos := strings.Index(fpath, "/")
		if pos != newpos && !PathExists(fpath[:pos+1]) {
			mkdir(fpath[:pos+1])
		}
	}

	//if start with /
	fmt.Println(fpath)
	f, err1 := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial

	resp, err := httpClient.Get(urldata.String())
	if err != nil {
		// handle error
		fmt.Println("error:", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//f.WriteString(htmlstr)
	f.Write(body)
}

func saveFile(fpath string, content string) {
	fpath = outputdir + fpath
	pos := strings.LastIndex(fpath, "/")
	if pos < 0 {
		return
	} else {
		if !PathExists(fpath[:pos+1]) {
			mkdir(fpath[:pos+1])
		}
	}
	f, err1 := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	//f.WriteString(htmlstr)
	f.Write([]byte(content))
}

func mkdir(pathstr string) {
	err := os.MkdirAll(pathstr, 0777)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://www.critterai.org/projects/nmgen_study/overview.html#voxelization")
	if err != nil {
		log.Fatal(err)
	}

	htmlstr, _ := doc.Html()
	fmt.Println(htmlstr)

	fmt.Println(doc.Url.RawQuery)
	fmt.Println(doc.Url.Host)
	fmt.Println(doc.Url.Path)

	//if url.path is end with / then add .index.html
	//else if url.path is start with # then continue
	//else if url.path is start with http or https
	//then check if host is www.critterai.org
	//	if true then replace http:... with ""

	//save html to file url.path

	//recurse all href, src

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Text()
		title, _ := s.Attr("href")
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
		s.SetText("111")
		s.SetAttr("href", "222")
	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Text()
		title, _ := s.Attr("href")
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
		s.SetText("111")
	})

	sela := doc.Find("script")
	sela.Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Text()
		title, _ := s.Attr("src")
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Text()
		title, _ := s.Attr("href")
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Text()
		title, _ := s.Attr("src")
		fmt.Printf("Review %d: %s - %s\n", i, band, title[0])
	})
}

func main() {
	//ExampleScrape()
	//http://www.uml-diagrams.org/
	purl := flag.String("url", "https://www.baidu.com", "-url=")
	pdir := flag.String("outdir", "tmp/", "-outdir=")
	flag.Parse()
	rawurl := *purl //"http://discuzt.cr180.com/"

	outputdir = *pdir

	// if string(outputdir[0]) == "/" || string(outputdir[0]) == "\\" {
	// 	f
	// }
	outputdir, _ = filepath.Abs(outputdir)

	if string(outputdir[len(outputdir)-1]) != "/" {
		outputdir += "/"
	}

	if !PathExists(outputdir) {
		mkdir(outputdir)
	}

	urlset = NewSet()
	gurldata, err = url.Parse(rawurl)

	if err != nil {
		fmt.Println(rawurl + " parse error:" + err.Error())
		return
	}

	GrabData(rawurl)
}
