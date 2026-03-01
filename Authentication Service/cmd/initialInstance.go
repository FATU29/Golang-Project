package main

import (
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/controller"
	"Authentication_Service/internal/cron"
	"Authentication_Service/internal/handler"
	"Authentication_Service/internal/repository"
	"Authentication_Service/internal/service"
	"Authentication_Service/pkg/email"
	"Authentication_Service/pkg/google"

	_ "Authentication_Service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type InitialInstance struct {
	App *App
}

func NewInitialInstance(app *App) *InitialInstance {
	return &InitialInstance{
		App: app,
	}
}

func (instance *InitialInstance) DefineRouter() {
	// Swagger UI: /swagger/index.html
	instance.App.Re.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	root := instance.App.Re.Group(constant.Root)

	userRepo := &repository.UserRepository{
		Db: instance.App.Db,
	}

	authRepo := &repository.AuthenticationRepository{
		Rd: instance.App.Rd,
		Db: instance.App.Db,
	}

	otpRepo := &repository.OtpRepository{
		Rd: instance.App.Rd,
	}

	refreshTokenRepo := &repository.RefreshTokenRepository{
		Db: instance.App.Db,
	}

	tokenGenerator := service.NewTokenGenerator(instance.App.Cfg.GetJwtSecret())

	var googleOAuth2 *google.OAuth2Config
	if instance.App.Cfg.GoogleClientID != "" && instance.App.Cfg.GoogleClientSecret != "" && instance.App.Cfg.GoogleRedirectURL != "" {
		googleOAuth2 = google.NewOAuth2Config(
			instance.App.Cfg.GoogleClientID,
			instance.App.Cfg.GoogleClientSecret,
			instance.App.Cfg.GoogleRedirectURL,
		)
	}

	authService := &service.AuthenticationService{
		AuthenticationRepository: authRepo,
		UserRepository:           userRepo,
		OtpRepository:            otpRepo,
		TokenGenerator:           tokenGenerator,
		RefreshTokenRepository:   refreshTokenRepo,
		GoogleOAuth2:             googleOAuth2,
	}

	authController := &controller.AuthenticationController{
		AuthenticationService: authService,
		Email: &email.BrevoEmail{
			Cfg: instance.App.Cfg,
		},
	}

	handler.NewAuthenticationHandler(authController, tokenGenerator).RouterList(root)

	// Cron: clean expired refresh tokens every hour
	if _, err := cron.StartRefreshTokenCleanup(refreshTokenRepo, "0 * * * *"); err != nil {
		panic("failed to start refresh token cleanup cron: " + err.Error())
	}
}
