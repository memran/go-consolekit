package main

import (
	"fmt"
	"go-consolekit/console"
	"io"
	"net/http"
	"os"
)

type DownloadCommand struct{}

func (c *DownloadCommand) Name() string {
	return "download"
}

func (c *DownloadCommand) Description() string {
	return "Download a file with progress bar"
}

func (c *DownloadCommand) Configure(config *console.CommandConfig) {
	config.Argument("url").Required().Description("URL to download")
	config.Argument("output").Required().Description("Output file path")
	config.Option("resume").Shortcut("r").Description("Resume partial download")
}

func (c *DownloadCommand) Handle(ctx *console.Context) error {
	url := ctx.Arg("url")
	output := ctx.Arg("output")
	resume := ctx.Option("resume") == "true"

	ctx.Title(fmt.Sprintf("Downloading %s", url))

	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	total := int(resp.ContentLength)
	if total <= 0 {
		ctx.Warning("Unknown file size, downloading without progress")
		return downloadPlain(url, output)
	}

	ctx.Line(fmt.Sprintf("Size: %s (%d bytes)", formatBytes(total), total))

	var offset int64
	if resume {
		if info, err := os.Stat(output); err == nil {
			offset = info.Size()
			if offset >= int64(total) {
				ctx.Success("File already fully downloaded")
				return nil
			}
			ctx.Line(fmt.Sprintf("Resuming from byte %d", offset))
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	flag := os.O_CREATE | os.O_WRONLY
	if offset > 0 {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	f, err := os.OpenFile(output, flag, 0644)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	bar := ctx.Output().Progress(fmt.Sprintf("Downloading %s", output), total)

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return fmt.Errorf("write error: %w", werr)
			}
			bar.Add(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
	}

	bar.Finish()
	ctx.Success(fmt.Sprintf("Downloaded to %s", output))
	return nil
}

func downloadPlain(url, output string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func formatBytes(n int) string {
	switch {
	case n >= 1024*1024*1024:
		return fmt.Sprintf("%.1f GiB", float64(n)/(1024*1024*1024))
	case n >= 1024*1024:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1024*1024))
	case n >= 1024:
		return fmt.Sprintf("%.1f KiB", float64(n)/1024)
	default:
		return fmt.Sprintf("%d B", n)
	}
}

var _ console.Command = (*DownloadCommand)(nil)
