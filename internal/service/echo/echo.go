package service

import (
	"errors"

	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/echo"
	repository "github.com/lin-snow/ech0/internal/repository/echo"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

type EchoService struct {
	commonService  commonService.CommonServiceInterface
	echoRepository repository.EchoRepositoryInterface
}

func NewEchoService(commonService commonService.CommonServiceInterface, echoRepository repository.EchoRepositoryInterface) EchoServiceInterface {
	return &EchoService{
		commonService:  commonService,
		echoRepository: echoRepository,
	}
}

// PostEcho 创建新的Echo
func (echoService *EchoService) PostEcho(userid uint, newEcho *model.Echo) error {
	newEcho.UserID = userid

	user, err := echoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查Extension内容
	if newEcho.Extension != "" && newEcho.ExtensionType != "" {
		switch newEcho.ExtensionType {
		case model.Extension_MUSIC:
			// 处理音乐链接 (暂无)
		case model.Extension_VIDEO:
			// 处理视频链接 (暂无)
		case model.Extension_GITHUBPROJ:
			// 处理GitHub项目的链接
			newEcho.Extension = httpUtil.TrimURL(newEcho.Extension)
		case model.Extension_WEBSITE:
			// 处理网站链接 (暂无)
		}
	} else {
		newEcho.Extension = ""
		newEcho.ExtensionType = ""
	}

	newEcho.Username = user.Username

	for i := range newEcho.Images {
		if newEcho.Images[i].ImageURL == "" {
			newEcho.Images[i].ImageSource = ""
		}
	}

	if newEcho.Content == "" && len(newEcho.Images) == 0 && (newEcho.Extension == "" || newEcho.ExtensionType == "") {
		return errors.New(commonModel.ECHO_CAN_NOT_BE_EMPTY)
	}

	return echoService.echoRepository.CreateEcho(newEcho)
}

// GetEchosByPage 获取Echo列表，支持分页
func (echoService *EchoService) GetEchosByPage(userid uint, pageQueryDto commonModel.PageQueryDto) (commonModel.PageQueryResult[[]model.Echo], error) {
	// 参数校验
	if pageQueryDto.Page < 1 {
		pageQueryDto.Page = 1
	}
	if pageQueryDto.PageSize < 1 || pageQueryDto.PageSize > 100 {
		pageQueryDto.PageSize = 10
	}

	//管理员登陆则支持查看隐私数据，否则不允许
	showPrivate := false
	if userid == authModel.NO_USER_LOGINED {
		showPrivate = false
	} else {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return commonModel.PageQueryResult[[]model.Echo]{}, err
		}
		if !user.IsAdmin {
			showPrivate = false
		}
		showPrivate = true
	}

	echosByPage, total := echoService.echoRepository.GetEchosByPage(pageQueryDto.Page, pageQueryDto.PageSize, pageQueryDto.Search, showPrivate)
	result := commonModel.PageQueryResult[[]model.Echo]{
		Items: echosByPage,
		Total: total,
	}

	return result, nil
}

// DeleteEchoById 删除指定ID的Echo
func (echoService *EchoService) DeleteEchoById(userid, id uint) error {
	user, err := echoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查该Echo是否存在图片
	echo, err := echoService.echoRepository.GetEchosById(id)
	if err != nil {
		return err
	}
	if echo == nil {
		return errors.New(commonModel.ECHO_NOT_FOUND)
	}

	// 删除Echo中的图片
	if len(echo.Images) > 0 {
		for _, img := range echo.Images {
			if err := echoService.commonService.DirectDeleteImage(img.ImageURL, img.ImageSource); err != nil {
				return err
			}
		}
	}

	return echoService.echoRepository.DeleteEchoById(id)
}

// GetTodayEchos 获取今天的Echo列表
func (echoService *EchoService) GetTodayEchos(userid uint) ([]model.Echo, error) {
	//管理员登陆则支持查看隐私数据，否则不允许
	showPrivate := false
	if userid == authModel.NO_USER_LOGINED {
		showPrivate = false
	} else {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return nil, err
		}
		if !user.IsAdmin {
			showPrivate = false
		}
		showPrivate = true
	}

	// 获取当日发布的Echos
	todayEchos := echoService.echoRepository.GetTodayEchos(showPrivate)

	return todayEchos, nil
}

// UpdateEcho 更新指定ID的Echo
func (echoService *EchoService) UpdateEcho(userid uint, echo *model.Echo) error {
	user, err := echoService.commonService.CommonGetUserByUserId(userid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(commonModel.NO_PERMISSION_DENIED)
	}

	// 检查Extension内容
	if echo.Extension != "" && echo.ExtensionType != "" {
		switch echo.ExtensionType {
		case model.Extension_MUSIC:
			// 处理音乐链接 (暂无)
		case model.Extension_VIDEO:
			// 处理视频链接 (暂无)
		case model.Extension_GITHUBPROJ:
			echo.Extension = httpUtil.TrimURL(echo.Extension)
		case model.Extension_WEBSITE:
			// 处理网站链接 (暂无)
		}
	} else {
		echo.Extension = ""
		echo.ExtensionType = ""
	}

	// 处理无效图片
	for i := range echo.Images {
		if echo.Images[i].ImageURL == "" {
			echo.Images[i].ImageSource = ""
			echo.Images[i].ImageURL = ""
		}
		// 确保外键正确设置
		echo.Images[i].MessageID = echo.ID
	}

	// 检查是否为空
	if echo.Content == "" && len(echo.Images) == 0 && (echo.Extension == "" || echo.ExtensionType == "") {
		return errors.New(commonModel.ECHO_CAN_NOT_BE_EMPTY)
	}

	return echoService.echoRepository.UpdateEcho(echo)
}

// LikeEcho 点赞指定ID的Echo
func (echoService *EchoService) LikeEcho(id uint) error {
	return echoService.echoRepository.LikeEcho(id)
}

// GetEchoById 获取指定 ID 的 Echo
func (echoService *EchoService) GetEchoById(id uint) (*model.Echo, error) {
	var echo *model.Echo

	echo, err := echoService.echoRepository.GetEchosById(id)
	if err != nil {
		return nil, err
	}

	if echo != nil && echo.Private == true {
		// 不允许通过ID获取私密Echo
		return nil, errors.New(commonModel.ECHO_NOT_FOUND)
	}

	return echo, nil
}
