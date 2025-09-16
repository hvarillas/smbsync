package smb

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/hvarillas/smbsync/internal/logger"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func verifyIntegrity(fs *smb2.Share, remoteFilePath string, sourceHashSum []byte, fileName, localBasePath string, deleteAfter, zippy bool) error {
	logger.Sugar.Info("Fase: Verificación de integridad SHA256")
	
	copiedFile, err := fs.Open(remoteFilePath)
	if err != nil {
		logger.Sugar.Errorf("Error al reabrir archivo remoto para verificación: %v", err)
		return fmt.Errorf("could not reopen remote file for verification: %w", err)
	}

	copiedFileInfo, err := copiedFile.Stat()
	if err != nil {
		copiedFile.Close()
		logger.Sugar.Errorf("Error al obtener información del archivo remoto: %v", err)
		return fmt.Errorf("could not get remote file info: %w", err)
	}

	bar := progressbar.NewOptions64(
		copiedFileInfo.Size(),
		progressbar.OptionSetDescription("Calculando Hash..."),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: "░",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	destHash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(destHash, bar), copiedFile); err != nil {
		copiedFile.Close()
		logger.Sugar.Errorf("Error al calcular hash del archivo remoto: %v", err)
		return fmt.Errorf("failed to calculate remote file hash: %w", err)
	}

	copiedFile.Close()

	destHashSum := destHash.Sum(nil)
	logger.Sugar.Debugf("Hash SHA256 del archivo destino: %x", destHashSum)

	if !bytes.Equal(sourceHashSum, destHashSum) {
		logger.Sugar.Errorf("¡FALLO DE INTEGRIDAD! Los hashes no coinciden para %s", fileName)
		logger.Sugar.Errorf("Hash origen: %x", sourceHashSum)
		logger.Sugar.Errorf("Hash destino: %x", destHashSum)
		return fmt.Errorf("hash mismatch: file corruption likely")
	}

	logger.Sugar.Infof("✅ Archivo %s copiado y verificado exitosamente", fileName)
	logger.Sugar.Debugf("Verificación SHA256 exitosa - Hashes coinciden: %x", sourceHashSum)

	if deleteAfter {
		time.Sleep(100 * time.Millisecond)
		
		originalFileToDelete := filepath.Join(localBasePath, fileName)
		logger.Sugar.Infof("Eliminando archivo local original: %s", originalFileToDelete)
		if err := os.Remove(originalFileToDelete); err != nil {
			logger.Sugar.Errorf("Fallo al eliminar el archivo local original %s: %v", originalFileToDelete, err)
			return fmt.Errorf("failed to delete local file: %w", err)
		}
		logger.Sugar.Infof("Archivo local original %s eliminado.", originalFileToDelete)
		
		if zippy {
			zipFilePath := filepath.Join(localBasePath, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".zip")
			if zipFilePath != originalFileToDelete {
				if _, err := os.Stat(zipFilePath); err == nil {
					logger.Sugar.Infof("Eliminando archivo zip temporal: %s", zipFilePath)
					if err := os.Remove(zipFilePath); err != nil {
						logger.Sugar.Errorf("Fallo al eliminar el archivo zip temporal %s: %v", zipFilePath, err)
					} else {
						logger.Sugar.Infof("Archivo zip temporal %s eliminado.", zipFilePath)
					}
				}
			}
		}
	}

	return nil
}
