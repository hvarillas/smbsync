package smb

import (
	"archive/zip"
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

func startCopy(fs *smb2.Share, fileName, localBasePath, remoteBasePath string, deleteAfter, zippy bool) error {
	localFilePath := filepath.Join(localBasePath, fileName)
	remoteFilePath := filepath.Join(remoteBasePath, fileName)

	if zippy {
		logger.Sugar.Infof("Comprimiendo archivo: %s", fileName)
		zipFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".zip"
		zipFilePath := filepath.Join(localBasePath, zipFileName)
		remoteFilePath = filepath.Join(remoteBasePath, zipFileName)

		zipFile, err := os.Create(zipFilePath)
		if err != nil {
			return fmt.Errorf("failed to create zip file: %w", err)
		}

		zipWriter := zip.NewWriter(zipFile)

		sourceFile, err := os.Open(localFilePath)
		if err != nil {
			zipFile.Close()
			return fmt.Errorf("failed to open source file: %w", err)
		}

		fileInfo, err := sourceFile.Stat()
		if err != nil {
			sourceFile.Close()
			zipFile.Close()
			return fmt.Errorf("failed to get file info: %w", err)
		}

		writer, err := zipWriter.Create(fileName)
		if err != nil {
			sourceFile.Close()
			zipFile.Close()
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		bar := progressbar.NewOptions64(
			fileInfo.Size(),
			progressbar.OptionSetDescription("Comprimiendo..."),
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

		_, err = io.Copy(io.MultiWriter(writer, bar), sourceFile)
		
		sourceFile.Close()
		zipWriter.Close()
		zipFile.Close()
		
		if err != nil {
			return fmt.Errorf("failed to copy to zip: %w", err)
		}

		localFilePath = zipFilePath
		logger.Sugar.Infof("Archivo comprimido exitosamente: %s", zipFileName)
	}

	logger.Sugar.Infof("Iniciando copia de archivo: %s", filepath.Base(localFilePath))
	logger.Sugar.Debugf("Ruta local: %s -> Ruta remota: %s", localFilePath, remoteFilePath)

	var sourceHashSum []byte
	var localFile *os.File
	var remoteFile *smb2.File
	
	err := func() error {
		var err error
		localFile, err = os.Open(localFilePath)
		if err != nil {
			logger.Sugar.Errorf("Error al abrir archivo local %s: %v", localFilePath, err)
			return fmt.Errorf("could not open local file %s: %w", localFilePath, err)
		}

		fileInfo, err := localFile.Stat()
		if err != nil {
			logger.Sugar.Errorf("Error al obtener información del archivo local: %v", err)
			return fmt.Errorf("could not get local file info: %w", err)
		}
		fileSize := fileInfo.Size()
		logger.Sugar.Infof("Tamaño del archivo %s: %d bytes (%.2f MB)", fileName, fileSize, float64(fileSize)/(1024*1024))

		remoteFile, err = fs.Create(remoteFilePath)
		if err != nil {
			logger.Sugar.Errorf("Error al crear archivo remoto %s: %v", remoteFilePath, err)
			return fmt.Errorf("could not create remote file %s: %w", remoteFilePath, err)
		}

		logger.Sugar.Info("Fase: Copiando archivo y calculando hash SHA256")
		sourceHash := sha256.New()

		bar := progressbar.NewOptions64(
			fileSize,
			progressbar.OptionSetDescription("Copiando..."),
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

		destWriter := io.MultiWriter(remoteFile, sourceHash, bar)

		if _, err := io.Copy(destWriter, localFile); err != nil {
			logger.Sugar.Errorf("Error durante la copia del archivo %s: %v", fileName, err)
			return fmt.Errorf("file copy failed: %w", err)
		}

		logger.Sugar.Infof("Copia completada para %s (%d bytes transferidos)", fileName, fileSize)
		sourceHashSum = sourceHash.Sum(nil)
		logger.Sugar.Debugf("Hash SHA256 del archivo origen: %x", sourceHashSum)
		return nil
	}()

	if localFile != nil {
		localFile.Close()
	}
	if remoteFile != nil {
		remoteFile.Close()
	}

	if err != nil {
		return err
	}

	return verifyIntegrity(fs, remoteFilePath, sourceHashSum, fileName, localBasePath, deleteAfter, zippy)
}
