package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"runtime"
	"strings"
)

func usage() {
	fmt.Println("Usage: git-copy-tags <source-repo> <dest-repo> [-f]")
	fmt.Println("By default, the script is in \"dry run\" mode, which means that it only prints out what it would do, without actually doing it. If you are happy with the result, add -f.")
	os.Exit(1)
}

func shell(cmd string) ([]byte, error) {
	sh := "sh"
	c := "-c"
	if runtime.GOOS == "windows" {
		sh = "cmd"
		c = "/c"
	}

	return exec.Command(sh, c, cmd).CombinedOutput()
}

func exe(cmd string) string {
	result, err := shell(cmd)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return string(result)
}

func getTags() map[string]string {
	tags := exe("git tag")
	dict := make(map[string]string)

	for _, tag := range strings.Split(tags, "\n") {
		tag = strings.TrimSpace(tag)
		if len(tag) < 1 {
			continue
		}
		cmd := fmt.Sprintf("git rev-list --max-count=1 %s", tag)
		commit := strings.TrimSpace(exe(cmd))
		dict[tag] = commit
	}

	return dict
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	src := os.Args[1]
	dest := os.Args[2]
	force := false
	if len(os.Args) > 3 && os.Args[3] == "-f" {
		force = true
	}

	os.Chdir(src)
	src_tags := getTags()

	os.Chdir(dest)
	dest_tags := getTags()

	if !force {
		fmt.Println("Running dry, use -f to actually apply changes...")
	}

	for tag, commit := range src_tags {
		if _, ok := dest_tags[tag]; !ok {
			cmd := fmt.Sprintf("git rev-list --max-count=1 %s", commit)
			if _, err := shell(cmd); err == nil {
				if force {
					if _, err := shell(fmt.Sprintf("git tag %s %s\n", tag, commit)); err == nil {
						fmt.Printf("Tagged %s with %s", commit, tag)
					}else {
						fmt.Printf("Error while tagging %s with %s\n", commit, tag)
					}

				} else {
					fmt.Printf("Would tag %s with %s\n", commit, tag)
				}
			}
		}
	}
}
