package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

// shamelessly copied from skopeo
// https://github.com/containers/skopeo/blob/c3e66b5876d349e05947e841e40df40d4f19c2ec/cmd/skopeo/sync.go#L72
// registrySyncConfig contains information about a single registry, read from
// the source YAML file
type registrySyncConfig struct {
	Images           map[string][]string // Images map images name to slices with the images' references (tags, digests)
	ImagesByTagRegex map[string]string   `yaml:"images-by-tag-regex"` // Images map images name to regular expression with the images' tags
	ImagesBySemver   map[string]string   `yaml:"images-by-semver"`    // ImagesBySemver maps a repository to a semver constraint (e.g. '>=3.14') to match images' tags to
}

// sourceConfig contains all registries information read from the source YAML file
type sourceConfig map[string]registrySyncConfig

// newSourceConfig unmarshals the provided YAML file path to the sourceConfig type.
// It returns a new unmarshaled sourceConfig object and any error encountered.
func newSourceConfig(yamlFile string) (sourceConfig, error) {
	var cfg sourceConfig
	source, err := os.ReadFile(yamlFile)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("Failed to unmarshal %q: %w", yamlFile, err)
	}
	return cfg, nil
}

type buildEntry struct {
	Source          string `yaml:"source"`
	Destination     string `yaml:"destination"`
	BuildOnMainOnly bool   `yaml:"buildOnMainOnly"`
}

type builderConfig struct {
	Build []buildEntry `yaml: "build"`
}

type diffOutput struct {
	Images map[string][]string
}

type packageRecord struct {
	tags []string
	ghId []int64
}

type scanConfig struct {
	CreatedAt string                         `yaml:"created_at"`
	Images    []map[string]map[string]string // Images map images name to slices with the images' references (tags, digests)
}

func newScanConfig(yamlFile string) (scanConfig, error) {
	var cfg scanConfig
	source, err := os.ReadFile(yamlFile)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("Failed to unmarshal %q: %w", yamlFile, err)
	}
	return cfg, nil
}

func newBuilderConfig(yamlFile string) (builderConfig, error) {
	var cfg builderConfig
	source, err := os.ReadFile(yamlFile)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("Failed to unmarshal %q: %w", yamlFile, err)
	}
	return cfg, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	builderYamlPath     string = "config.yaml"
	syncYamlPath        string = "sync-ghcr.yml"
	scanYamlPath        string = "base-image-usage-report.yaml"
	org                        = "geonet"
	repository                 = "base-images"
	registryPrefix             = "ghcr.io/geonet"
	packageView                = "public"
	tagExclusionPattern        = "sha256-[a-fA-F0-9]{64}"
	keep_versions              = 5
	packageType                = "container"
)

func main() {
	//token := os.Getenv("GITHUB_TOKEN")

	//dryrun := flag.Bool("dryrun", true, "set to false to enable deletion")

	flag.Parse()

	//read artisinal docker build configuration
	builderCfg, err := newBuilderConfig(builderYamlPath)
	if err != nil {
		log.Fatalln(err)
	}

	//check for tags?
	//read skopeo sync.yml
	sourceCfg, err := newSourceConfig(syncYamlPath)
	if err != nil {
		log.Fatalln(err)
	}
	yaml_inv := make(map[string][]string, len(sourceCfg)+len(builderCfg.Build))
	for _, v := range builderCfg.Build {
		//I could do something clever here with the oci image library code but was lazy
		reference := strings.Split(v.Destination, ":")
		tags := make([]string, 1)
		tags[0] = reference[1]
		yaml_inv[reference[0]] = tags
	}

	scanCfg, err := newScanConfig(scanYamlPath)
	if err != nil {
		log.Fatalln(err)
	}
	//reformatting custom scanner report yaml to align with sourceCfg format
	scanImageList := func(scanRaw scanConfig) map[string][]string {
		imgList := make(map[string][]string, len(scanRaw.Images))
		for _, img := range scanRaw.Images {
			for i := range img {
				//find tags in name string
				//this could use the oci string parsing format code to be more clever
				if strings.Contains(i, ":") {
					image := strings.Split(i, ":")
					if len(image) == 2 {
						imgList[image[0]] = make([]string, 1)
						imgList[image[0]][0] = image[1]
					} else {
						tag := strings.Join(image[1:], ":")
						imgList[image[0]] = make([]string, 1)
						imgList[image[0]][0] = tag
					}
				} else {
					//assuming default tag of latest
					imgList[i] = make([]string, 1)
					imgList[i][0] = "latest"
				}
			}
		}
		return imgList
	}(scanCfg)

	fmt.Println(scanImageList)
	for k := range sourceCfg {
		for v := range sourceCfg[k].Images {
			name := strings.Split(v, "/")
			var basename string
			if len(name) > 1 {
				basename = name[1]
			} else {
				basename = name[0]
			}
			inventory_name := fmt.Sprintf("%s/%s/%s", registryPrefix, repository, basename)
			yaml_inv[inventory_name] = sourceCfg[k].Images[v]
		}
	}
	for i, tags := range scanImageList {
		inventory_name := fmt.Sprintf("%s/%s/%s", registryPrefix, repository, i)
		for _, t := range tags {
			if !slices.Contains(yaml_inv[inventory_name], t) {
				yaml_inv[inventory_name] = slices.Insert(yaml_inv[inventory_name], 0, t)
			}
		}
	}
	fmt.Println(yaml_inv)
	logfile, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	errorfile, err := os.Create("error.log")
	if err != nil {
		log.Fatal(err)
	}
	for image, tags := range yaml_inv {
		for _, v := range tags {
			cmd := exec.Command("/usr/bin/docker", "ghcr.io/anchore/grype:latest", fmt.Sprintf("%s:%s", image, v), "-o", "json")
			cmd.Stdout = logfile
			cmd.Stderr = errorfile
			err := cmd.Run()
			if err != nil {
				log.Println("failed to scan: ", fmt.Sprintf("podman:%s:%s", image, v))
			}
		}
	}
	defer logfile.Close()
	defer errorfile.Close()
}
