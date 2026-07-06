package services

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync/atomic"

	"github.com/disintegration/imaging"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type SettingsService struct {
	settings atomic.Pointer[models.InstanceSettings]
	Config   config.Config
}

func NewSettingsService(config config.Config) *SettingsService {
	return &SettingsService{Config: config}
}

func (s *SettingsService) Bootstrap(ctx context.Context, q db.Querier) error {
	logrus.Infof("Bootstrapping settings from config")
	err := q.Queries().InitSettings(ctx, biomedb.InitSettingsParams{
		AppName:                s.Config.Instance.AppName,
		AppSubtitle:            s.Config.Instance.AppSubtitle,
		AppDescription:         s.Config.Instance.AppDescription,
		IsPublic:               s.Config.Instance.IsPublic,
		AccountRequestsEnabled: s.Config.Instance.AccountRequestsEnabled,
		AdminEmail:             s.Config.Instance.AdminEmail,
		MailFromAddress:        s.Config.Instance.MailFromAddress,
		MailFromName:           s.Config.Instance.MailFromName,
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap settings: %v", err)
	}
	return s.Reload(ctx, q)
}

func (s *SettingsService) GetSettings() models.InstanceSettings {
	return *s.settings.Load()
}

func (s *SettingsService) Reload(ctx context.Context, q db.Querier) error {
	settingsDB, err := q.Queries().GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to reload settings: %v", err)
	}
	settings := models.SettingsFromDB(settingsDB)
	s.settings.Store(&settings)
	return nil
}

func (s *SettingsService) UpdateInstanceSettings(ctx context.Context, q db.Querier, input models.InstanceSettingsUpdate) error {
	_, err := q.Queries().UpdateInstanceSettings(ctx, input.ToParams())
	if err == nil {
		err = s.Reload(ctx, q)
	}
	return err
}

func (s *SettingsService) TestSMTPConnection(ctx context.Context) (bool, error) {
	dialer := gomail.NewDialer(s.Config.SMTP.SMTPHost, int(s.Config.SMTP.SMTPPort), s.Config.SMTP.SMTPUser, s.Config.SMTP.SMTPPassword)
	closer, err := dialer.Dial()
	if err != nil {
		return false, err
	}
	closer.Close()
	return true, nil
}

func (s *SettingsService) TogglePublicAccess(ctx context.Context, q db.Querier, isPublic bool) error {
	_, err := q.Queries().UpdateInstanceSettings(ctx, biomedb.UpdateInstanceSettingsParams{
		IsPublic: &isPublic,
	})
	if err == nil {
		err = s.Reload(ctx, q)
	}
	return err
}

func (s *SettingsService) TogglePublicRegistration(ctx context.Context, q db.Querier, enabled bool) error {
	_, err := q.Queries().UpdateInstanceSettings(ctx, biomedb.UpdateInstanceSettingsParams{
		AccountRequestsEnabled: &enabled,
	})
	if err == nil {
		err = s.Reload(ctx, q)
	}
	return err
}

func (s *SettingsService) SetAppIcon(ctx context.Context, icon image.Image) error {
	resizedImg := imaging.Resize(icon, 300, 300, imaging.Lanczos)

	writer, err := os.Create("assets/app_icon.png")
	if err != nil {
		return err
	}
	defer writer.Close()
	err = png.Encode(writer, resizedImg)
	if err != nil {
		return err
	}

	return nil
}
