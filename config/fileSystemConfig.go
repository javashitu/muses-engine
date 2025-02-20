package config

import (
	"embed"
	"log"

	"gopkg.in/yaml.v3"
)

var FileSystemConfig FileSystemConf

type FileSystemConf struct {
	UploadConfig   UploadConf         `yaml:"upload"`
	DownloadConfig DownloadConf       `yaml:"download"`
	VideoManager   VideoManagerConfig `yaml:"videoManager"`
}

type UploadConf struct {
	FileNameFormat string `yaml:"fileNameFormat"`
	FilePathFormat string `yaml:"filePathFormat"`
}

type DownloadConf struct {
	VideoDownloadDir string `yaml:"videoDownloadDir"`
}

type VideoManagerConfig struct {
	ScriptDir            string `yaml:"scriptDir"`
	TranscodeVideoScript string `yaml:"transcodeVideoScript"`
}

//go:embed fileSystemConfig.yaml
var fsf embed.FS

func init() {
	configFileSystem()
}

func configFileSystem() {
	file, err := fsf.ReadFile("fileSystemConfig.yaml")
	if err != nil {
		log.Panicf("read file system config error,maybe file not exits -> %s ", err)
	}
	err = yaml.Unmarshal(file, &FileSystemConfig)
	if err != nil {
		log.Panicf("bind file system config error,maybe not yaml format -> %s ", err)
	}
}

//======================================================================对象方法======================================================================

func (videoManager VideoManagerConfig) FmtTrscodeVideoCmd(inputPath string, outputPath string) string {
	scriptPath := videoManager.ScriptDir + "/" + videoManager.TranscodeVideoScript
	return scriptPath + " " + inputPath + " " + outputPath
}
