package main

import (
	"MiniPrograms/api"
	"MiniPrograms/responsity/cache"
	"MiniPrograms/responsity/dao"
	"MiniPrograms/responsity/model"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var (
	username string
	password string
)

func main() {
	// 初始化配置
	initConfig()
	db, err := dao.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	s := newService(db)

	// 启动 Gin 服务
	g := gin.Default()
	g.Use(corsHdl())
	g.POST("/checkStatus", s.CheckStatus) // 检查项目状态
	g.PUT("/setStatus", s.SetStatus)      // 设置项目状态

	// 启动服务
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// 读取配置文件
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	username = viper.GetString("username")
	password = viper.GetString("password")
}

type Service struct {
	DAO   *dao.MiniProgramsDAO
	cache *cache.Cache
}

func newService(db *gorm.DB) *Service {
	return &Service{
		DAO:   dao.NewMiniProgramsDAO(db),
		cache: cache.NewCache(),
	}
}

// CheckStatus 检查项目状态
func (s *Service) CheckStatus(ctx *gin.Context) {
	var req api.CheckStatusReq
	if err := ctx.ShouldBind(&req); err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, 5001, err.Error())
		return
	}

	p, err := s.getMiniProgram(req.Name)
	if err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, 5001, "不存在的项目名称")
		return
	}

	ctx.JSON(http.StatusOK, api.Resp{
		Code: 0,
		Msg:  "获取成功!",
		Data: api.CheckStatusResp{Status: p.Status},
	})
}

// SetStatus 设置项目状态
func (s *Service) SetStatus(ctx *gin.Context) {
	var req api.SetStatusReq
	if err := ctx.ShouldBind(&req); err != nil {
		s.errorResponse(ctx, http.StatusBadRequest, 4001, err.Error())
		return
	}

	if req.Username != username || req.Password != password {
		s.errorResponse(ctx, http.StatusUnauthorized, 4001, "认证失败!")
		return
	}

	p, err := s.getOrCreateMiniProgram(req.ProgramsName, req.Status)
	if err != nil {
		s.errorResponse(ctx, http.StatusInternalServerError, 5002, fmt.Sprintf("保存失败: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, api.Resp{
		Msg:  fmt.Sprintf("%s 设置成功!", p.Name),
		Code: 0,
	})
}

// getMiniProgram 从缓存或数据库获取项目
func (s *Service) getMiniProgram(name string) (*model.MiniPrograms, error) {
	if p, ok := s.cache.Load(name); ok {
		return p, nil
	}

	p, err := s.DAO.Find(name)
	if err != nil {
		return nil, err
	}

	s.cache.Store(name, p)
	return p, nil
}

// getOrCreateMiniProgram 获取或创建项目
func (s *Service) getOrCreateMiniProgram(name string, status bool) (*model.MiniPrograms, error) {
	p, err := s.DAO.Find(name)
	if err != nil {
		// 如果项目不存在，创建新项目
		p = &model.MiniPrograms{Name: name, Status: status}
		if err := s.DAO.Save(*p); err != nil {
			return nil, err
		}
		s.cache.Store(name, p)
		return p, nil
	}

	// 更新状态并保存
	p.Status = status
	if err := s.DAO.Save(*p); err != nil {
		return nil, err
	}

	s.cache.Store(name, p)
	return p, nil
}

// errorResponse 通用错误响应
func (s *Service) errorResponse(ctx *gin.Context, httpCode, code int, msg string) {
	ctx.JSON(httpCode, api.Resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// corsHdl 配置跨域
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           12 * time.Hour,
	})
}
