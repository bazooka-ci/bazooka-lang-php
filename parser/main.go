package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	Php          = "php"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
	err := bazooka.LoadCryptoKeyFromFile("/bazooka-cryptokey")
	if err != nil {
		log.Fatal(err)
	}
}

type ConfigPhp struct {
	Base        bazooka.Config `yaml:",inline"`
	PhpVersions []string       `yaml:"php,omitempty"`
}

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &ConfigPhp{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}

	versions := conf.PhpVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"5.6"}
	}
	for i, version := range versions {
		if err := managePhpVersion(fmt.Sprintf("0%d", i), conf, version, ""); err != nil {
			log.Fatal(err)
		}
	}
	for i, image := range images {
		if err := managePhpVersion(fmt.Sprintf("1%d", i), conf, "", image); err != nil {
			log.Fatal(err)
		}
	}

}

func managePhpVersion(counter string, conf *ConfigPhp, version, image string) error {
	conf.PhpVersions = nil
	conf.Base.Image = nil

	setDefaultScript(conf)

	meta := map[string]string{}
	if len(version) > 0 {
		var err error
		image, err = resolvePhpImage(version)
		if err != nil {
			return err
		}
		meta[Php] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err := bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func resolvePhpImage(version string) (string, error) {
	//TODO extract this from db
	phpMap := map[string]string{
		"5.4": "bazooka/runner-php:5.4",
		"5.5": "bazooka/runner-php:5.5",
		"5.6": "bazooka/runner-php:5.6",
	}
	if val, ok := phpMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for PHP Runnner %s", version)
}

func setDefaultScript(conf *ConfigPhp) {
	if len(conf.Base.Script) == 0 {
		conf.Base.Script = []string{"phpunit --configuration /phpunit_$DB.xml --coverage-text"}
	}
}
