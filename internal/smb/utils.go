package smb

import (
	"os"
	"regexp"

	"github.com/hvarillas/smbsync/internal/logger"
)

func getRegexFiles(regex, path string) []string {
	logger.Sugar.Debugf("Escaneando directorio '%s' con patrón regex: '%s'", path, regex)

	files, err := os.ReadDir(path)
	if err != nil {
		logger.Sugar.Errorf("Error al leer directorio local '%s': %v", path, err)
		return nil
	}

	re, err := regexp.Compile("(?i)" + regex)
	if err != nil {
		logger.Sugar.Errorf("Patrón regex inválido '%s': %v", regex, err)
		return nil
	}

	var matchingFiles []string
	for _, file := range files {
		if !file.IsDir() && re.MatchString(file.Name()) {
			matchingFiles = append(matchingFiles, file.Name())
			logger.Sugar.Debugf("Archivo encontrado: %s", file.Name())
		}
	}

	if len(matchingFiles) > 0 {
		logger.Sugar.Infof("Encontrados %d archivos que coinciden con el patrón en '%s'", len(matchingFiles), path)
	} else {
		logger.Sugar.Warnf("No se encontraron archivos que coincidan con el patrón '%s' en '%s'", regex, path)
	}

	return matchingFiles
}
