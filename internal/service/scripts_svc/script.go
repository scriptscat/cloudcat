package scripts_svc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/scriptscat/cloudcat/pkg/scriptcat/plugin/window"
	"go.uber.org/zap/zapcore"

	"github.com/codfrm/cago/pkg/gogo"
	"github.com/scriptscat/cloudcat/internal/task/producer"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/robfig/cron/v3"
	api "github.com/scriptscat/cloudcat/internal/api/scripts"
	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"
	"github.com/scriptscat/cloudcat/internal/pkg/code"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"
	"github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/plugin"
	"go.uber.org/zap"
)

type ScriptSvc interface {
	// List 脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Install 安装脚本
	Install(ctx context.Context, req *api.InstallRequest) (*api.InstallResponse, error)
	// Get 获取脚本
	Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error)
	// Update 更新脚本
	Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error)
	// Delete 删除脚本
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
	// StorageList 值储存空间列表
	StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error)
	// Run 手动运行脚本
	Run(ctx context.Context, req *api.RunRequest) (*api.RunResponse, error)
	// Watch 监听脚本
	Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error)
	// Stop 停止脚本
	Stop(ctx context.Context, req *api.StopRequest) (*api.StopResponse, error)
}

type scriptSvc struct {
	sync.Mutex
	ctx       map[string]context.CancelFunc
	cron      *cron.Cron
	cronEntry map[string]cron.EntryID
}

var defaultScripts ScriptSvc

func Script() ScriptSvc {
	return defaultScripts
}

func NewScript(ctx context.Context) (ScriptSvc, error) {
	svc := &scriptSvc{
		cronEntry: make(map[string]cron.EntryID),
		ctx:       make(map[string]context.CancelFunc),
	}
	svc.cron = cron.New(cron.WithSeconds())
	svc.cron.Start()
	if err := gogo.Go(func(ctx context.Context) error {
		<-ctx.Done()
		svc.cron.Stop()
		logger.Ctx(ctx).Info("stop cron")
		return nil
	}, gogo.WithContext(ctx)); err != nil {
		return nil, err
	}
	// 初始化脚本运行环境
	scriptcat.RegisterRuntime(scriptcat.NewRuntime(
		logger.NewCtxLogger(logger.Default()),
		[]scriptcat.Plugin{
			window.NewBrowserPlugin(),
			plugin.NewGMPlugin(NewGMPluginFunc()),
		},
	))
	// 初始化运行脚本
	list, err := script_repo.Script().FindPage(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.State == script_entity.ScriptStateEnable {
			script := v
			_ = gogo.Go(func(ctx context.Context) error {
				return svc.run(ctx, script)
			})
		}
	}
	defaultScripts = svc
	return svc, nil
}

func (s *scriptSvc) runScript(ctx context.Context, script *script_entity.Script) error {
	s.Lock()
	ctx, cancel := context.WithCancel(ctx)
	s.ctx[script.ID] = cancel
	defer func(scriptId string) {
		s.Lock()
		defer s.Unlock()
		cancel()
		delete(s.ctx, scriptId)
	}(script.ID)
	s.Unlock()
	// 更新数据库中的状态
	script, err := script_repo.Script().Find(ctx, script.ID)
	if err != nil {
		return err
	}
	script.Status.SetRunStatus(script_entity.RunStateRunning)
	if err := script_repo.Script().Update(ctx, script); err != nil {
		return err
	}
	defer func() {
		script, err := script_repo.Script().Find(ctx, script.ID)
		if err == nil && script != nil {
			script.Status.SetRunStatus(script_entity.RunStateComplete)
			if err := script_repo.Script().Update(ctx, script); err != nil {
				logger.Ctx(ctx).Error("update script status error", zap.Error(err))
			}
		} else {
			logger.Ctx(ctx).Error("find script error", zap.Error(err))
		}
	}()
	with := logger.Ctx(ctx).
		With(zap.String("id", script.ID), zap.String("name", script.Name)).
		WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
			return nil
		}))
	if _, err := scriptcat.RuntimeCat().Run(logger.ContextWithLogger(ctx, with), &scriptcat.Script{
		ID:       script.ID,
		Code:     script.Code,
		Metadata: scriptcat.Metadata(script.Metadata),
	}); err != nil {
		with.Error("run script error", zap.Error(err))
		return err
	}
	return nil
}

func (s *scriptSvc) run(ctx context.Context, script *script_entity.Script) error {
	logger := logger.Ctx(ctx).
		With(zap.String("id", script.ID), zap.String("name", script.Name))
	// 判断是什么类型的脚本,如果是后台脚本,直接运行,定时脚本,添加定时任务
	if _, ok := script.Metadata["background"]; ok {
		if err := s.runScript(ctx, script); err != nil {
			logger.Error("run background script error", zap.Error(err))
			return err
		}
	} else if cron, ok := script.Crontab(); ok {
		if err := s.addCron(ctx, script, cron); err != nil {
			logger.Error("add cron error", zap.Error(err))
			return err
		}
	} else {
		logger.Error("script type error")
		return errors.New("script type error")
	}
	logger.Info("run script success")
	return nil
}

func (s *scriptSvc) addCron(ctx context.Context, script *script_entity.Script, c string) error {
	logger := logger.Ctx(ctx).
		With(zap.String("id", script.ID), zap.String("name", script.Name))
	var err error
	c, err = scriptcat.ConvCron(c)
	if err != nil {
		return err
	}
	cronEntry, err := s.cron.AddFunc(c, func() {
		err := s.runScript(ctx, script)
		if err != nil {
			logger.Error("run cron script error", zap.Error(err))
		} else {
			logger.Info("run cron script success")
		}
	})
	if err != nil {
		return err
	}
	s.Lock()
	s.cronEntry[script.ID] = cronEntry
	s.Unlock()
	return nil
}

// List 脚本列表
func (s *scriptSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, err := script_repo.Script().FindPage(ctx)
	if err != nil {
		return nil, err
	}
	resp := &api.ListResponse{
		List: make([]*api.Script, 0),
	}
	for _, v := range list {
		resp.List = append(resp.List, &api.Script{
			ID:           v.ID,
			Name:         v.Name,
			Metadata:     v.Metadata,
			SelfMetadata: v.SelfMetadata,
			Status:       v.Status,
			State:        v.State,
			Createtime:   v.Createtime,
			Updatetime:   v.Updatetime,
		})
	}
	return resp, nil
}

// Install 安装脚本
func (s *scriptSvc) Install(ctx context.Context, req *api.InstallRequest) (*api.InstallResponse, error) {
	resp := &api.InstallResponse{
		Scripts: make([]*api.Script, 0),
	}
	script, err := scriptcat.RuntimeCat().Parse(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	// 根据id判断是否已经存在
	model, err := script_repo.Script().Find(ctx, script.ID)
	if err != nil {
		return nil, err
	}
	// 如果存在则更新
	if model != nil {
		if err := model.Update(script); err != nil {
			return nil, err
		}
		if err := script_repo.Script().Update(ctx, model); err != nil {
			return nil, err
		}
		// 如果已经是开启,那么重新启动
		if model.State == script_entity.ScriptStateEnable {
			_ = gogo.Go(func(ctx context.Context) error {
				s.disable(model)
				return s.run(context.Background(), model)
			})
		}
	} else {
		model = &script_entity.Script{}
		if err := model.Create(script); err != nil {
			return nil, err
		}
		if err := script_repo.Script().Create(ctx, model); err != nil {
			return nil, err
		}
		// 开启
		_ = gogo.Go(func(ctx context.Context) error {
			return s.run(ctx, model)
		})
	}
	if err := producer.PublishScriptUpdate(ctx, model); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get 获取脚本
func (s *scriptSvc) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	script, err := script_repo.Script().FindByPrefixID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptNotFound)
	}
	return &api.GetResponse{
		Script: &api.Script{
			ID:         script.ID,
			Name:       script.Name,
			Code:       script.Code,
			Metadata:   script.Metadata,
			Status:     script.Status,
			State:      script.State,
			Createtime: script.Createtime,
			Updatetime: script.Updatetime,
		},
	}, nil
}

func (s *scriptSvc) stop(script *script_entity.Script) {
	s.Lock()
	defer s.Unlock()
	if cancel, ok := s.ctx[script.ID]; ok {
		cancel()
		delete(s.ctx, script.ID)
	}
}

func (s *scriptSvc) disable(script *script_entity.Script) {
	// 停止脚本
	s.Lock()
	defer s.Unlock()
	if cancel, ok := s.ctx[script.ID]; ok {
		cancel()
		delete(s.ctx, script.ID)
	}
	if eid, ok := s.cronEntry[script.ID]; ok {
		s.cron.Remove(eid)
		delete(s.cronEntry, script.ID)
	}
}

// Update 更新脚本
func (s *scriptSvc) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	// 查出脚本
	model, err := script_repo.Script().FindByPrefixID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptNotFound)
	}
	if req.Script.Code != "" {
		script, err := scriptcat.RuntimeCat().Parse(ctx, req.Script.Code)
		if err != nil {
			return nil, err
		}
		if err := model.Update(script); err != nil {
			return nil, err
		}
	}
	if model.State != req.Script.State {
		model.State = req.Script.State
		switch req.Script.State {
		case script_entity.ScriptStateEnable:
			_ = gogo.Go(func(ctx context.Context) error {
				return s.run(ctx, model)
			})
		case script_entity.ScriptStateDisable:
			go s.disable(model)
		default:
			return nil, i18n.NewError(ctx, code.ScriptStateError)
		}
	} else if model.Status.GetRunStatus() != req.Script.Status.GetRunStatus() {
		switch req.Script.Status.GetRunStatus() {
		case script_entity.RunStateRunning:
			_ = gogo.Go(func(ctx context.Context) error {
				return s.runScript(ctx, model)
			})
		case script_entity.RunStateComplete:
			go s.stop(model)
		default:
			return nil, i18n.NewError(ctx, code.ScriptRunStateError)
		}
	}
	// 更新信息
	if req.Script.Metadata != nil {
		model.Metadata = req.Script.Metadata
	}
	if req.Script.SelfMetadata != nil {
		model.SelfMetadata = req.Script.SelfMetadata
	}
	if req.Script.Status != nil {
		model.Status = req.Script.Status
	}
	if req.Script.State != "" {
		model.State = req.Script.State
	}
	model.Updatetime = time.Now().Unix()
	if err := script_repo.Script().Update(ctx, model); err != nil {
		return nil, err
	}
	if err := producer.PublishScriptUpdate(ctx, model); err != nil {
		return nil, err
	}
	return nil, nil
}

// Delete 删除脚本
func (s *scriptSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	// 查出脚本
	script, err := script_repo.Script().FindByPrefixID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptNotFound)
	}
	s.disable(script)
	if err := script_repo.Script().Delete(ctx, script.ID); err != nil {
		return nil, err
	}
	if err := producer.PublishScriptDelete(ctx, script); err != nil {
		return nil, err
	}
	return nil, nil
}

// StorageList 值储存空间列表
func (s *scriptSvc) StorageList(ctx context.Context, req *api.StorageListRequest) (*api.StorageListResponse, error) {
	list, err := script_repo.Script().StorageList(ctx)
	if err != nil {
		return nil, err
	}
	resp := &api.StorageListResponse{
		List: make([]*api.Storage, 0),
	}
	for _, v := range list {
		resp.List = append(resp.List, &api.Storage{
			Name:         v.Name,
			LinkScriptID: v.LinkScriptID,
		})
	}
	return resp, nil
}

// Run 手动运行脚本
func (s *scriptSvc) Run(ctx context.Context, req *api.RunRequest) (*api.RunResponse, error) {
	script, err := script_repo.Script().FindByPrefixID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptNotFound)
	}
	go func() {
		err = s.runScript(context.Background(), script)
		if err != nil {
			logger.Ctx(ctx).Error("run script error", zap.Error(err))
		}
	}()
	return nil, nil
}

// Watch 监听脚本
func (s *scriptSvc) Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error) {
	return nil, nil
}

// Stop 停止脚本
func (s *scriptSvc) Stop(ctx context.Context, req *api.StopRequest) (*api.StopResponse, error) {
	script, err := script_repo.Script().FindByPrefixID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptNotFound)
	}
	s.stop(script)
	return nil, nil
}
