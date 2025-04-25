package main

import (
	"context"
	"flag"

	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/protofy-xyz/wread/fs"
)

func main() {
	url := flag.String("url", "", "target url")
	out := flag.String("out", "", "output path for file with findings")
	pull := flag.String("pull", "", "pull resources locally")
	timeout := flag.Duration("timeout", 15*time.Second, "max time to load site")
	flag.Parse()

	if *url == "" {
		println("")
		println("wread - v0")
		println("")
		println("usage:")
		flag.PrintDefaults()
		println("")
		os.Exit(0)
	}

	println("\nwread - v0\n")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, *timeout)
	defer cancel()

	var resources []string

	// network listener
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if resp, ok := ev.(*network.EventResponseReceived); ok {
			resources = append(resources, resp.Response.URL)
		}
	})

	if err := chromedp.Run(ctx,
		network.Enable(),
		network.SetCacheDisabled(true), // disable cacche to force resources load
		chromedp.Navigate(*url),
		chromedp.Sleep(5*time.Second),
	); err != nil {
		println("Error loading site")
		os.Exit(1)
	}

	// filter duplicates
	for _, r := range unique(resources) {
		fmt.Println(r)
	}

	if *out != "" {
		fs.WriteFile(*out, []byte(strings.Join(resources, "\n")))
	}

	if *pull != "" {
		root := fs.EnsureDir(*pull)
		for _, resource := range resources {
			res, err := http.Get(resource)
			if err != nil {
				fmt.Printf(" %s\n", err)
				os.Exit(1)
			}

			resBody, err := io.ReadAll(res.Body)

			fs.WriteFile(path.Join(root, fs.GetFilenameFromURL(resource)), resBody)
		}
		println("")
		println("Resources saved at", root)
		println("")
	}
}

func unique(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}
