package ezutil

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/config"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotisserie/eris"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// region Gin Utils

func GetPathParam[T any](ctx *gin.Context, key string) (T, bool, error) {
	var zero T

	paramValue, exists := ctx.Params.Get(key)
	if !exists {
		return zero, false, nil
	}

	parsedValue, err := Parse[T](paramValue)
	if err != nil {
		return zero, false, eris.Wrapf(err, "failed to parse parameter '%s'", key)
	}

	return parsedValue, true, nil
}

func BindRequest[T any](ctx *gin.Context, bindType binding.Binding) (T, error) {
	var zero T

	if err := ctx.ShouldBindWith(&zero, bindType); err != nil {
		return zero, eris.Wrapf(err, "failed to bind request with type %s", bindType.Name())
	}

	return zero, nil
}

func NewCorsMiddleware(corsConfig *cors.Config) gin.HandlerFunc {
	if corsConfig == nil {
		log.Println("CORS configuration is nil, using default settings")
		return cors.Default()
	}

	if err := corsConfig.Validate(); err != nil {
		log.Fatalf("invalid CORS configuration: %s", err.Error())
	}

	return cors.New(*corsConfig)
}

func NewAuthMiddleware(
	authStrategy string,
	tokenCheckFunc func(ctx *gin.Context, token string) (bool, map[string]any, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, errMsg, err := internal.ExtractToken(ctx, authStrategy)
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, "error extracting token"))
			ctx.Abort()
			return
		}
		if errMsg != "" {
			_ = ctx.Error(UnauthorizedError(errMsg))
			ctx.Abort()
			return
		}

		exists, data, err := tokenCheckFunc(ctx, token)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}
		if !exists {
			_ = ctx.Error(UnauthorizedError(config.MsgAuthUserNotFound))
			ctx.Abort()
			return
		}

		for key, val := range data {
			ctx.Set(key, val)
		}

		ctx.Next()
	}
}

func NewPermissionMiddleware(
	roleContextKey string,
	requiredPermission string,
	permissionMap map[string][]string,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString(roleContextKey)
		if role == "" {
			_ = ctx.Error(eris.Errorf("role not found in context or invalid type"))
			ctx.Abort()
			return
		}

		permissions, ok := permissionMap[role]
		if !ok {
			_ = ctx.Error(eris.Errorf("unknown role: %s", role))
			ctx.Abort()
			return
		}

		if !slices.Contains(permissions, requiredPermission) {
			_ = ctx.Error(ForbiddenError(config.MsgNoPermission))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func NewErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if err := ctx.Errors.Last(); err != nil {
			if originalErr, ok := err.Err.(AppError); ok {
				ctx.AbortWithStatusJSON(originalErr.HttpStatusCode, NewErrorResponse(originalErr))
				return
			}

			statusCode, appError := constructAppError(err)
			ctx.AbortWithStatusJSON(statusCode, NewErrorResponse(appError))
		}
	}
}

func constructAppError(err *gin.Error) (int, error) {
	originalErr := eris.Unwrap(err.Err)
	switch originalErr := originalErr.(type) {
	case validator.ValidationErrors:
		var errors []string
		for _, e := range originalErr {
			errors = append(errors, e.Error())
		}

		return http.StatusUnprocessableEntity, ValidationError(errors)
	case *json.SyntaxError:
		return http.StatusBadRequest, BadRequestError(config.MsgInvalidJson)
	default:
		// EOF error from json package is unexported
		if originalErr == io.EOF || originalErr.Error() == "EOF" {
			return http.StatusBadRequest, BadRequestError(config.MsgMissingBody)
		}

		log.Printf("unhandled error of type: %T\n", originalErr)
		log.Println(eris.ToString(err.Err, true))
		return http.StatusInternalServerError, InternalServerError()
	}
}

func GetFromContext[T any](ctx *gin.Context, key string) (T, error) {
	var zero T

	val, exists := ctx.Get(key)
	if !exists {
		return zero, eris.Errorf("value with key %s not found in context", key)
	}

	asserted, ok := val.(T)
	if !ok {
		return zero, eris.Errorf("error asserting value %s as type %T", val, zero)
	}

	return asserted, nil
}

type RouteConfig struct {
	Group    string
	Versions []RouteVersionConfig
	Handlers []gin.HandlerFunc
}

type RouteVersionConfig struct {
	Version  int
	Groups   []RouteGroupConfig
	Handlers []gin.HandlerFunc
}

type RouteGroupConfig struct {
	Group     string
	Endpoints []EndpointConfig
	Handlers  []gin.HandlerFunc
}

type EndpointConfig struct {
	Method   string
	Endpoint string
	Handlers []gin.HandlerFunc
}

func SetupRoutes(router *gin.Engine, routeConfigs []RouteConfig) {
	if router == nil {
		log.Fatal("Router cannot be nil")
	}

	for _, routeConfig := range routeConfigs {
		routeGroup := router.Group(routeConfig.Group, routeConfig.Handlers...)
		for _, versionConfig := range routeConfig.Versions {
			versionGroup := routeGroup.Group(fmt.Sprintf("/v%d", versionConfig.Version), versionConfig.Handlers...)
			for _, routeGroupConfig := range versionConfig.Groups {
				group := versionGroup.Group(routeGroupConfig.Group, routeGroupConfig.Handlers...)
				for _, endpointConfig := range routeGroupConfig.Endpoints {
					group.Handle(endpointConfig.Method, endpointConfig.Endpoint, endpointConfig.Handlers...)
				}
			}

		}
	}
}

// endregion

// region Gorm Utils

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * limit

		return db.Limit(limit).Offset(offset)
	}
}

func OrderBy(field string, ascending bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Basic validation to prevent SQL injection
		// Only allow alphanumeric characters, underscores, and dots for table.column
		if !internal.IsValidFieldName(field) {
			_ = db.AddError(eris.Errorf("invalid field name: %s", field))
			return db
		}

		if ascending {
			return db.Order(field + " ASC")
		}

		return db.Order(field + " DESC")
	}
}

func WhereBySpec[T any](spec T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&spec)
	}
}

func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}

		return db
	}
}

type Transactor interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

func NewTransactor(db *gorm.DB) Transactor {
	return &internal.GormTransactor{DB: db}
}

func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	return internal.GetTxFromContext(ctx)
}

func WithinTransaction(ctx context.Context, transactor Transactor, serviceFn func(ctx context.Context) error) error {
	ctx, err := transactor.Begin(ctx)
	if err != nil {
		return eris.Wrap(err, "error starting transaction")
	}
	defer transactor.Rollback(ctx)

	if err := serviceFn(ctx); err != nil {
		return eris.Wrap(err, "error executing service function")
	}

	return transactor.Commit(ctx)
}

// endregion

// region Slice Utils

func MapSlice[T any, U any](input []T, mapperFunc func(T) U) []U {
	output := make([]U, len(input))

	for i, v := range input {
		output[i] = mapperFunc(v)
	}

	return output
}

// endregion

// region Time Utils

func GetStartOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

func GetEndOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

// endregion

// region HTTP Utils

type AppError struct {
	Type           string `json:"type"`
	Message        string `json:"message"`
	HttpStatusCode int    `json:"-"`
	Details        any    `json:"details,omitempty"`
}

func (ae AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", ae.Type, ae.Message, ae.Details)
}

func InternalServerError() AppError {
	return AppError{
		Type:           "InternalServerError",
		Message:        "Undefined error occurred",
		HttpStatusCode: http.StatusInternalServerError,
	}
}

func ConflictError(details any) AppError {
	return AppError{
		Type:           "ConflictError",
		Message:        "Conflict with existing resource",
		HttpStatusCode: http.StatusConflict,
		Details:        details,
	}
}

func NotFoundError(details any) AppError {
	return AppError{
		Type:           "NotFoundError",
		Message:        "Requested resource is not found",
		HttpStatusCode: http.StatusNotFound,
		Details:        details,
	}
}

func UnauthorizedError(details any) AppError {
	return AppError{
		Type:           "UnauthorizedError",
		Message:        "Unauthorized access",
		HttpStatusCode: http.StatusUnauthorized,
		Details:        details,
	}
}

func ForbiddenError(details any) AppError {
	return AppError{
		Type:           "ForbiddenError",
		Message:        "Forbidden access",
		HttpStatusCode: http.StatusForbidden,
		Details:        details,
	}
}

func BadRequestError(details any) AppError {
	return AppError{
		Type:           "BadRequestError",
		Message:        "Request is not valid",
		HttpStatusCode: http.StatusBadRequest,
		Details:        details,
	}
}

func UnprocessableEntityError(details any) AppError {
	return AppError{
		Type:           "UnprocessableEntityError",
		Message:        "Request cannot be processed due to semantic errors",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}

func ValidationError(details any) AppError {
	return AppError{
		Type:           "ValidationError",
		Message:        "Failed to validate request",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}

type QueryOptions struct {
	Page  int `query:"page" binding:"required,min=1"`
	Limit int `query:"limit" binding:"required,min=1"`
}

type Pagination struct {
	TotalData   int  `json:"totalData"`
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
}

func (p *Pagination) IsZero() bool {
	return p.TotalData == 0 && p.CurrentPage == 0 && p.TotalPages == 0 && !p.HasNextPage && !p.HasPrevPage
}

type JSONResponse struct {
	Message    string     `json:"message"`
	Data       any        `json:"data,omitzero"`
	Errors     error      `json:"errors,omitempty"`
	Pagination Pagination `json:"pagination,omitzero"`
}

func NewResponse(message string) JSONResponse {
	return JSONResponse{
		Message: message,
	}
}

func NewErrorResponse(err error) any {
	return JSONResponse{
		Message: err.Error(),
		Errors:  err,
	}
}

func (jr JSONResponse) WithData(data any) JSONResponse {
	jr.Data = data
	return jr
}

func (jr JSONResponse) WithError(err error) JSONResponse {
	jr.Errors = err
	return jr
}

func (jr JSONResponse) WithPagination(queryOptions QueryOptions, totalData int) JSONResponse {
	totalPages := int(math.Ceil(float64(totalData) / float64(queryOptions.Limit)))

	jr.Pagination = Pagination{
		TotalData:   totalData,
		CurrentPage: queryOptions.Page,
		TotalPages:  totalPages,
		HasNextPage: queryOptions.Page < totalPages,
		HasPrevPage: queryOptions.Page > 1,
	}

	return jr
}

func RunServer(defaultConfigs Config, serverSetupFunc func(*Config) *http.Server) {
	configs := LoadConfig(defaultConfigs)
	srv := serverSetupFunc(configs)
	internal.ServeGracefully(srv, configs.App.Timeout)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Data map[string]any `json:"data"`
}

type JWTService interface {
	CreateToken(data map[string]any) (string, error)
	VerifyToken(tokenstr string) (JWTClaims, error)
}

type jwtServiceHS256 struct {
	issuer        string
	secretKey     string
	tokenDuration time.Duration
}

func NewJwtService(configs *Auth) JWTService {
	return &jwtServiceHS256{
		issuer:        configs.Issuer,
		secretKey:     configs.SecretKey,
		tokenDuration: configs.TokenDuration,
	}
}

func (j *jwtServiceHS256) CreateToken(data map[string]any) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Data: data,
		},
	)

	signed, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", eris.Wrap(err, "error signing token")
	}

	return signed, nil
}

func (j *jwtServiceHS256) VerifyToken(tokenstr string) (JWTClaims, error) {
	var claims JWTClaims

	_, err := jwt.ParseWithClaims(
		tokenstr,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, UnauthorizedError(config.MsgAuthExpiredToken)
		}

		return claims, eris.Wrap(err, "error parsing token")
	}

	return claims, nil
}

type HashService interface {
	Hash(val string) (string, error)
	CheckHash(hash, val string) (bool, error)
}

func NewHashService(cost int) HashService {
	return &internal.HashServiceBcrypt{Cost: cost}
}

// endregion

// region String Utils

func Parse[T any](value string) (T, error) {
	var parsed any
	var err error
	var zero T
	var parsedType string

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		parsed, err = strconv.Atoi(value)
		parsedType = "int"
	case bool:
		parsed, err = strconv.ParseBool(value)
		parsedType = "bool"
	case uuid.UUID:
		parsed, err = uuid.Parse(value)
		parsedType = "uuid"
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}

	if err != nil {
		return zero, eris.Wrapf(err, "failed to parse value '%s' as %s", value, parsedType)
	}

	return parsed.(T), nil
}

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", eris.New("length must be greater than 0")
	}

	randomBytes := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", eris.Wrap(err, "failed to generate random string")
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// endregion

// region Config Loader

type Config struct {
	App   *App
	Auth  *Auth
	SQLDB *SQLDB
	GORM  *gorm.DB
}

func LoadConfig(defaults Config) *Config {
	sqlDBConfig := loadSQLDBConfig()

	appDefaults := App{}
	if defaults.App != nil {
		appDefaults = *defaults.App
	}

	authDefaults := Auth{}
	if defaults.Auth != nil {
		authDefaults = *defaults.Auth
	}

	return &Config{
		App:   loadAppConfig(appDefaults),
		Auth:  loadAuthConfig(authDefaults),
		SQLDB: sqlDBConfig,
		GORM:  sqlDBConfig.openGormConnection(),
	}
}

type App struct {
	Env        string
	Port       string
	Timeout    time.Duration
	ClientUrls []string
	Timezone   string
}

func loadAppConfig(defaults App) *App {
	var loadedConfig App

	err := envconfig.Process("APP", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading app config: %s", err.Error())
	}

	if loadedConfig.Env == "" {
		loadedConfig.Env = defaults.Env
	}
	if loadedConfig.Port == "" {
		loadedConfig.Port = defaults.Port
	}
	// Validate port number
	if port, err := strconv.Atoi(loadedConfig.Port); err != nil || port < 1 || port > 65535 {
		log.Fatalf("invalid port number: %s", loadedConfig.Port)
	}
	if loadedConfig.Timeout == 0 {
		loadedConfig.Timeout = defaults.Timeout
	}
	if len(loadedConfig.ClientUrls) == 0 {
		loadedConfig.ClientUrls = defaults.ClientUrls
	}
	if loadedConfig.Timezone == "" {
		loadedConfig.Timezone = defaults.Timezone
	}
	// Validate timezone
	if _, err := time.LoadLocation(loadedConfig.Timezone); err != nil {
		log.Fatalf("invalid timezone: %s", loadedConfig.Timezone)
	}

	return &loadedConfig
}

type Auth struct {
	SecretKey      string
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
}

func loadAuthConfig(defaults Auth) *Auth {
	var loadedConfig Auth

	err := envconfig.Process("AUTH", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading auth config: %s", err.Error())
	}

	if loadedConfig.SecretKey == "" {
		loadedConfig.SecretKey = defaults.SecretKey
	}
	if loadedConfig.TokenDuration == 0 {
		loadedConfig.TokenDuration = defaults.TokenDuration
	}
	if loadedConfig.CookieDuration == 0 {
		loadedConfig.CookieDuration = defaults.CookieDuration
	}
	if loadedConfig.Issuer == "" {
		loadedConfig.Issuer = defaults.Issuer
	}
	if loadedConfig.URL == "" {
		loadedConfig.URL = defaults.URL
	}

	return &loadedConfig
}

type SQLDB struct {
	Host     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true"`
	Port     string `required:"true"`
	Driver   string `required:"true"`
}

func loadSQLDBConfig() *SQLDB {
	var loadedConfig SQLDB

	err := envconfig.Process("SQLDB", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading SQLDB config: %s", err.Error())
	}

	return &loadedConfig
}

func (sqldb *SQLDB) openGormConnection() *gorm.DB {
	db, err := gorm.Open(sqldb.getGormDialector(), &gorm.Config{})
	if err != nil {
		log.Fatalf("error opening GORM connection: %s", err.Error())
	}

	return db
}

func (sqldb *SQLDB) getGormDialector() gorm.Dialector {
	switch sqldb.Driver {
	case "mysql":
		return mysql.Open(sqldb.getDSN())
	case "postgres":
		return postgres.Open(sqldb.getDSN())
	default:
		log.Fatalf("unsupported SQLDB driver: %s", sqldb.Driver)
		return nil
	}
}

func (sqldb *SQLDB) getDSN() string {
	switch sqldb.Driver {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sqldb.User,
			sqldb.Password,
			sqldb.Host,
			sqldb.Port,
			sqldb.Name,
		)
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			sqldb.Host,
			sqldb.User,
			sqldb.Password,
			sqldb.Name,
			sqldb.Port,
		)
	default:
		log.Fatalf("unsupported SQLDB driver: %s", sqldb.Driver)
		return ""
	}
}

// endregion
