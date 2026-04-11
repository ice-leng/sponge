package logic

import (
	"context"
	"errors"

	"github.com/go-dev-frame/sponge/internal/cache"
	"github.com/go-dev-frame/sponge/internal/dao"
	"github.com/go-dev-frame/sponge/internal/database"
	"github.com/go-dev-frame/sponge/internal/ecode"
	"github.com/go-dev-frame/sponge/internal/model"
	"github.com/go-dev-frame/sponge/internal/types"
	"github.com/go-dev-frame/sponge/pkg/copier"
)

// UserExampleLogic defining the handler interface
type UserExampleLogic interface {
	Create(ctx context.Context, request *types.CreateUserExampleRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdateUserExampleByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.UserExampleObjDetail, error)
	List(ctx context.Context, request *types.ListUserExamplesRequest) ([]*types.UserExampleObjDetail, int64, error)
}

type userExampleService struct {
	iDao dao.UserExampleDao
}

// NewUserExampleLogic creating the handler interface
func NewUserExampleLogic() UserExampleLogic {
	return NewUserExampleLogicByDAO(
		dao.NewUserExampleDao(
			database.GetDB(), // todo show db driver name here
			cache.NewUserExampleCache(database.GetCacheType()),
		),
	)
}

// NewUserExampleLogicByDAO creating the handler interface with injected dao, used for unit tests.
func NewUserExampleLogicByDAO(iDao dao.UserExampleDao) UserExampleLogic {
	return &userExampleService{iDao: iDao}
}

func (h userExampleService) Create(ctx context.Context, request *types.CreateUserExampleRequest) (uint64, error) {
	table := &model.UserExample{}
	err := copier.Copy(table, request)
	if err != nil {
		return 0, ecode.ErrCreateUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	err = h.iDao.Create(ctx, table)
	return table.ID, err
}

func (h userExampleService) DeleteByID(ctx context.Context, id uint64) error {
	return h.iDao.DeleteByID(ctx, id)
}

func (h userExampleService) UpdateByID(ctx context.Context, request *types.UpdateUserExampleByIDRequest) error {
	order := &model.UserExample{}
	err := copier.Copy(order, request)
	if err != nil {
		return ecode.ErrUpdateByIDUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return h.iDao.UpdateByID(ctx, order)
}

func (h userExampleService) GetByID(ctx context.Context, id uint64) (*types.UserExampleObjDetail, error) {
	result, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.NotFound.Err()
		}
		return nil, err
	}

	data, err := convertUserExample(result)
	if err != nil {
		return nil, ecode.ErrGetByIDUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func (h userExampleService) List(ctx context.Context, request *types.ListUserExamplesRequest) ([]*types.UserExampleObjDetail, int64, error) {
	orders, total, err := h.iDao.GetByColumns(ctx, &request.Params)
	if err != nil {
		return nil, 0, err
	}
	data, err := convertUserExamples(orders)
	if err != nil {
		return nil, 0, ecode.ErrListUserExample.Err()
	}
	return data, total, nil
}

func convertUserExample(userExample *model.UserExample) (*types.UserExampleObjDetail, error) {
	data := &types.UserExampleObjDetail{}
	err := copier.Copy(data, userExample)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserExamples(fromValues []*model.UserExample) ([]*types.UserExampleObjDetail, error) {
	toValues := []*types.UserExampleObjDetail{}
	for _, v := range fromValues {
		data, err := convertUserExample(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
