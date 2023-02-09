package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qsr-mock-server/models"
	"strconv"
	"strings"
)

// RuleAdd
// @Summary 新增拦截规则
// @Description 拦截规则新增时的一些安全检查，不重复检查等
// @Tags 拦截规则相关接口
// @Accept application/json
// @Param Authorization header string true "ownerName"
// @Param proxyRulesModel body string false "新增的实体"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/addRule [POST]
func RuleAdd(ctx *gin.Context) {
	ownerNameValue := ctx.Request.Header.Get("ownerName")
	if ownerNameValue == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "需要包含当前用户信息",
		})
		return
	}

	var proxyRulesModel models.ProxyRulesModel
	proxyRulesModel.OwnerName = ownerNameValue
	if err := ctx.ShouldBind(&proxyRulesModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不合法",
		})
		return
	}
	fmt.Println(proxyRulesModel.ProxyHost, "-", proxyRulesModel.ProxyUrl, "-", proxyRulesModel.CustomizedResp)
	if proxyRulesModel.ProxyHost == "" || proxyRulesModel.ProxyUrl == "" || proxyRulesModel.CustomizedResp == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "预设的 主机地址 或 路由地址 或 预置返回信息 不能为空",
		})
		return
	}
	if proxyRulesModel.ProxyMethod != "ANY" && proxyRulesModel.ProxyMethod != http.MethodGet && proxyRulesModel.ProxyMethod != http.MethodPost {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ProxyMethod错误，暂不支持该请求方法",
		})
		return
	}
	proxyRulesModel.ProxyUrl = strings.Split(proxyRulesModel.ProxyUrl, "?")[0] // 如果 url中包含了? 就进行过滤掉
	if proxyRulesModel.ProxyMethod == "" {
		// ProxyMethod 如果 请求方法 为 空 就给默认 POST
		proxyRulesModel.ProxyMethod = http.MethodPost
	}
	bytes, err := json.Marshal(proxyRulesModel.ProxyConditions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "条件解析 ProxyConditionJson 为json发生错误",
		})
		return
	}

	proxyRulesModel.ProxyConditionJson = string(bytes)

	// 存入数据库
	db, rulesModel := models.CreateRules(&proxyRulesModel)
	if db == nil {
		// 插入失败，查询到了一个已经存在的数据
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "当前规则已经存在，规则id为" + strconv.Itoa(int(rulesModel.ID)),
		})
		return
	}

	if proxyRulesModel.ID == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "数据插入数据库发生意外，没有插入成功",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": proxyRulesModel,
		})
	}
}

// GetAll
// @Summary 获取当前所有的规则
// @Description 无差别获取，不根据用户进行隔离
// @Tags 拦截规则相关接口
// @Param Authorization header string true "ownerName"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/getAllRules [GET]
func GetAll(ctx *gin.Context) {
	// TODO: 根据用户进行隔离查询，加分页
	rulesModels, count := models.GetAllRules()
	if count <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": rulesModels,
		})
	}
}

// DeleteRuleById
// @Summary 删除一个拦截规则
// @Description 逻辑删除一个拦截规则
// @Tags 拦截规则相关接口
// @Accept application/json
// @Param Authorization header string true "ownerName"
// @Param id query int true "拦截规则的id"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/deleteRulesById [POST]
func DeleteRuleById(ctx *gin.Context) {
	// json中获取，id=xx
	bodyFromCtx, err := GetJSONBodyFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务端出错",
		})
		return
	}
	var tarId uint
	switch bodyFromCtx["id"].(type) {
	case string:
		id, err := strconv.Atoi(bodyFromCtx["id"].(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "id值类型错误",
			})
			return
		}
		tarId = uint(id)
	case float64:
		f := bodyFromCtx["id"].(float64)
		tarId = uint(f)
	case int:
		f := bodyFromCtx["id"].(int)
		tarId = uint(f)
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "id类型不正确",
		})
		return
	}

	rulesById := models.DeleteRulesById(tarId)
	if rulesById.Error == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "删除成功",
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败，检查这个id是否存在",
		})
		return
	}
}

// UpdateRuleById
// @Summary 根据id更新一个拦截规则
// @Description 更新一个拦截规则的所有信息
// @Tags 拦截规则相关接口
// @Accept application/json
// @Param Authorization header string true "ownerName"
// @Param id query int true "拦截规则的id"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/updateRulesById [POST]
func UpdateRuleById(ctx *gin.Context) {
	var proxyRulesModel models.ProxyRulesModel
	if err := ctx.ShouldBind(&proxyRulesModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不合法",
		})
		return
	}
	fmt.Println(proxyRulesModel.ProxyHost, "-", proxyRulesModel.ProxyUrl, "-", proxyRulesModel.CustomizedResp)
	if proxyRulesModel.ProxyHost == "" || proxyRulesModel.ProxyUrl == "" || proxyRulesModel.CustomizedResp == "" || proxyRulesModel.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "预设的 主机地址 或 路由地址 或 预置返回信息 或 ID 不能为空",
		})
		return
	}
	if proxyRulesModel.ProxyMethod != "ANY" && proxyRulesModel.ProxyMethod != http.MethodGet && proxyRulesModel.ProxyMethod != http.MethodPost {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "ProxyMethod错误，暂不支持该请求方法",
		})
		return
	}
	proxyRulesModel.ProxyUrl = strings.Split(proxyRulesModel.ProxyUrl, "?")[0] // 如果 url中包含了? 就进行过滤掉
	if proxyRulesModel.ProxyMethod == "" {
		// ProxyMethod 如果 请求方法 为 空 就给默认 POST
		proxyRulesModel.ProxyMethod = http.MethodPost
	}
	bytes, err := json.Marshal(proxyRulesModel.ProxyConditions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "条件解析 ProxyConditionJson 为json发生错误",
		})
		return
	}

	proxyRulesModel.ProxyConditionJson = string(bytes)

	// 更新数据库
	db := models.UpdateRuleById(proxyRulesModel)
	if db.Error == nil && db.RowsAffected == 1 { // 根据id更新，只能影响一行
		ctx.JSON(http.StatusOK, gin.H{
			"message": proxyRulesModel,
		})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失败",
		})
		return
	}

}

// GetAllByCons
// @Summary 根据指定的条件来获取一部分拦截规则
// @Description 根据拦截规则的其中一个或者多个参数来进行筛选，多个参数同时满足时才会返回
// @Tags 拦截规则相关接口
// @Accept application/json
// @Param Authorization header string true "ownerName"
// @Param id body string true "需要满足的多个条件json"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/getRuleByCons [GET]
func GetAllByCons(ctx *gin.Context) {
	// 使用 struct 查询时，GORM 只会查询非零字段，这意味着如果您的字段的值为0,或其他零值’'，则不会用于构建查询条件,要在查询条件中包含零值可以使用map
	var proxyRulesModel models.ProxyRulesModel
	if err := ctx.ShouldBindJSON(&proxyRulesModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不合法",
		})
		return
	}
	proxyRulesModels := models.GetRulesByCons(&proxyRulesModel)
	ctx.JSON(http.StatusOK, gin.H{
		"message": proxyRulesModels,
	})
}

// GetRuleById
// @Summary 根据指定的条件来获取一部分拦截规则
// @Description 根据拦截规则的其中一个或者多个参数来进行筛选，多个参数同时满足时才会返回
// @Tags 拦截规则相关接口
// @Accept application/json
// @Param Authorization header string true "ownerName"
// @Param id query int true "拦截规则的id"
// @Success 200 {string} json{"code","message"}
// @Router /mockRule/getRuleDetailById [GET]
func GetRuleById(ctx *gin.Context) {
	// json中获取，id=xx
	var proxyRulesModel models.ProxyRulesModel
	if err := ctx.ShouldBindJSON(&proxyRulesModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不合法",
		})
		return
	}
	tarId := proxyRulesModel.ID
	if tarId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "没有id这个重要的参数哦",
		})
		return
	}
	rulesById := models.GetRulesById(tarId)
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": rulesById,
	})
}
