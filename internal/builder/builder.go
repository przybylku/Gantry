package builder

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/client"
)
type buildLogLine struct {
	Stream string      `json:"stream"`      // To jest czysty tekst logu (np. "Step 1/5...")
	Error  string      `json:"error"`       // Tutaj Docker wrzuca błędy krytyczne
	ErrorDetail struct {
		Message string `json:"message"`
	} `json:"errorDetail"`
}
func BuildImage(ctx context.Context, cli *client.Client, workDir string, imageName string, commitHash string) error {
	buildContext, err := createTarContext(workDir)
	if err != nil {
		fmt.Errorf("Error taring context, %w", err)
		return err
	}

	tags := []string{
		fmt.Sprintf("%s:latest", imageName),
		fmt.Sprintf("%s:%s", imageName, commitHash),
	}
	fmt.Printf("🔨 Building image: %s (tags: %v)... \n", imageName, tags)
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:        tags,
		Remove:      true,
	}
	var res build.ImageBuildResponse
	res, err = cli.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		fmt.Errorf("Error building image, %w", err)
		return err
	}
	defer res.Body.Close()

	return printBuildLogs(res.Body)
}


func createTarContext(src string) (io.Reader, error){
	pr, pw := io.Pipe()
	go func() {
		tw := tar.NewWriter(pw)
		defer pw.Close()
		defer tw.Close()

		err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Create tar header
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// abs -> relative
			relPath, err := filepath.Rel(src, file)
			if err != nil {
				return err
			}

			header.Name = filepath.ToSlash(relPath)

			if header.Name == "." {
				return nil
			}

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err := io.Copy(tw, data); err != nil {
					data.Close()
					return err
				}
				data.Close()
			}
			return nil

		})
		if err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr, nil
}

func printBuildLogs(rd io.Reader) error {
	decoder := json.NewDecoder(rd)

	for {
		var line buildLogLine
		if err := decoder.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if line.Error != "" {
			return fmt.Errorf("docker build failed: %s", line.Error)
		}

		// Wypisz "czysty" tekst (Stream) bez cudzysłowów JSON-a
		// Docker wysyła czasem puste linie, pomijamy je dla estetyki
		if strings.TrimSpace(line.Stream) != "" {
			// fmt.Print(line.Stream)
		}
	}
	return nil
}