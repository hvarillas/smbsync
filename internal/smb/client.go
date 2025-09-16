package smb

import (
	"fmt"
	"net"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/hvarillas/smbsync/internal/config"
	"github.com/hvarillas/smbsync/internal/logger"
)

func getSmbSession(user, password, smbHost string) (*smb2.Session, error) {
	logger.Sugar.Debugf("Estableciendo conexión TCP con %s:445", smbHost)

	conn, err := net.DialTimeout("tcp", smbHost+":445", 5*time.Second)
	if err != nil {
		logger.Sugar.Errorf("Error de conexión TCP: %v", err)
		return nil, fmt.Errorf("connection error: %w", err)
	}

	logger.Sugar.Debugf("Configurando autenticación NTLM para usuario: %s", user)
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
		},
	}

	logger.Sugar.Debug("Iniciando negociación SMB2")
	s, err := d.Dial(conn)
	if err != nil {
		logger.Sugar.Errorf("Error de autenticación SMB: %v", err)
		return nil, fmt.Errorf("SMB authentication error: %w", err)
	}

	logger.Sugar.Infof("Autenticación SMB exitosa para usuario: %s", user)
	return s, nil
}

func RunHeadless(cfg *config.Config) {
	logger.Sugar.Info("Iniciando en modo headless (sin TUI).")
	files := getRegexFiles(cfg.Regex, cfg.Path)
	if files == nil || len(files) == 0 {
		logger.Sugar.Warnf("No se encontraron archivos que coincidan con el patrón en: %s", cfg.Path)
		return
	}

	logger.Sugar.Infof("Encontrados %d archivos para sincronizar", len(files))
	session, err := getSmbSession(cfg.SMBUser, cfg.SMBPass, cfg.SMBHost)
	if err != nil {
		logger.Sugar.Fatalf("No se pudo establecer la sesión SMB: %v", err)
	}
	defer session.Logoff()

	share, err := session.Mount(cfg.Shared)
	if err != nil {
		logger.Sugar.Fatalf("No se pudo montar el recurso compartido '%s': %v", cfg.Shared, err)
	}
	defer share.Umount()

	for i, file := range files {
		logger.Sugar.Infof("Procesando archivo %d de %d: %s", i+1, len(files), file)
		if err := startCopy(share, file, cfg.Path, cfg.SharedPath, cfg.DeleteAfter, cfg.Zippy); err != nil {
			logger.Sugar.Errorf("Fallo al copiar %s: %v", file, err)
		} else {
			logger.Sugar.Infof("Archivo %s copiado y verificado exitosamente.", file)
		}
	}
	logger.Sugar.Info("Proceso de sincronización completado.")
}
