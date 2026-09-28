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
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap settings: %w", err)
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

// SaveSettings leaves cache publication to the caller after transaction commit.
func (s *SettingsService) SaveSettings(ctx context.Context, q db.Querier, input models.InstanceSettingsUpdate) error {
	_, err := q.Queries().UpdateInstanceSettings(ctx, input.ToParams())
	return err
}

func (s *SettingsService) UpdateInstanceSettings(ctx context.Context, q db.Querier, input models.InstanceSettingsUpdate) error {
	err := s.SaveSettings(ctx, q, input)
	if err == nil {
		err = s.Reload(ctx, q)
	}
	return err
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

func (s *SettingsService) GetDashboardMessage(ctx context.Context, q db.Querier) (*string, error) {
	settingsDB, err := q.Queries().GetSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %v", err)
	}
	return settingsDB.FrontpageMessageMD, nil
}

func (s *SettingsService) SetDashboardMessage(ctx context.Context, q db.Querier, message *string) error {
	_, err := q.Queries().SetDashboardMessage(ctx, message)
	if err == nil {
		err = s.Reload(ctx, q)
	}
	return err
}
