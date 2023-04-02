package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/katallaxie/pkg/utils/files"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Folder string
	Alpaca bool
	LLaMA  bool
	Size   string
}

var cfg = &config{}

// Downloader is an interface for downloading files.
type Downloader interface {
	Download(ctx context.Context, path string, size string) error
}

var root = &cobra.Command{
	Use: "install",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	root.PersistentFlags().StringVar(&cfg.Folder, "folder", "", "folder")
	root.PersistentFlags().BoolVar(&cfg.Alpaca, "alpha", false, "Alpaca")
	root.PersistentFlags().BoolVar(&cfg.LLaMA, "llama", true, "llama")
	root.PersistentFlags().StringVar(&cfg.Size, "size", "7B", "size")
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

type llama struct{}

// Download downloads the files.
func (l *llama) Download(ctx context.Context, p string, size string) error {
	var num = map[string]int{
		"7B":  1,
		"13B": 2,
		"30B": 4,
		"65B": 8,
	}

	ff := []string{}

	err := files.MkdirAll(p, 0o755)
	if err != nil {
		return nil
	}

	for i := 0; i <= num[size]; i++ {
		ff = append(ff, fmt.Sprintf("consolidated.0%d.pth", i))
	}

	for _, f := range ff {
		out, err := os.OpenFile(filepath.Join(p, f), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer out.Close()

		req, err := http.NewRequest("GET", fmt.Sprintf("https://agi.gpt4.org/llama/LLaMA/%s/%s", size, f), nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		io.Copy(io.MultiWriter(out, bar), resp.Body)
	}

	return nil
}

type alpaca struct{}

// Download downloads the files.
func (a *alpaca) Download(ctx context.Context, p string, size string) error {
	return nil
}

func run(ctx context.Context) error {
	var d Downloader

	d = &llama{}

	if cfg.Alpaca {
		d = &alpaca{}
	}

	err := d.Download(ctx, filepath.Join(cfg.Folder), cfg.Size)
	if err != nil {
		return err
	}

	return nil
}
