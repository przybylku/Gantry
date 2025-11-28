package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/przybylku/gantry/internal/builder"
	"github.com/przybylku/gantry/internal/gitops"
	"github.com/przybylku/gantry/internal/runtime"
	"github.com/spf13/cobra"
)

// Zaciaganie gita i deploy

// Flagi init
var gitaddr string
var name string
var internalPort string = "3000"
func init() {
	deployCmd.Flags().StringVarP(&gitaddr, "git", "g", "", "Git repository address")
	deployCmd.Flags().StringVarP(&name, "name", "n", "mysite", "Name of your site")
	deployCmd.Flags().StringVarP(&internalPort, "port", "p", "3000", "Internal port your app listens on")
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your static site quickly and easily",
	Long:  `A Fast and Flexible Static Site Generator`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println("Error creating Docker client:", err)
			return
		}
		defer cli.Close()


		fmt.Println("Deploying your static site named:", name, " from  repository: ", gitaddr)
		path, hash, err := gitops.CloneRepo(gitaddr)
		if err != nil {
			fmt.Println("Error cloning repository:", err)
			return
		}
		fmt.Println("Repository cloned to:", path, "with commit hash:", hash)
		
		imageName := "gantry-"+name
		fmt.Printf(" Commit: %s\n", hash)

		err = builder.BuildImage(ctx, cli, path, imageName, hash)
		if err != nil {
			fmt.Println("Error building image:", err)
			return
		}

		_ = cli.ContainerRemove(ctx, name, container.RemoveOptions{Force: true})
		
		
		
		containerID, err := runtime.RunContainer(ctx, cli, imageName, name, internalPort)
		if err != nil {
			fmt.Println("Error running container:", err)
			return
		}

		fmt.Println("Container started with ID:", containerID)
		fmt.Println("------------------------------------------------")
		fmt.Printf("✅ DEPLOYMENT SUCCESSFUL!\n")
		fmt.Printf("   App Name:  %s\n", name)
		fmt.Printf("   Container: %s\n", containerID[:12])
		fmt.Printf("   URL:       http://%s.localhost\n", name) // Port na sztywno z runtime.go
		fmt.Println("------------------------------------------------")
	},
}

