package e

//定义错误码
const (
	//基本错误
	SUCCESS        = 200 //success
	ERROR          = 500 //fail
	INVALID_PARAMS = 400 //请求参数错误
	//标签类错误
	ERROR_EXIST_TAG       = 10001 //已存在该标签名称
	ERROR_EXIST_TAG_FAIL  = 10002 //获取已存在标签失败
	ERROR_NOT_EXIST_TAG   = 10003 //该标签不存在
	ERROR_GET_TAGS_FAIL   = 10004 //获取所有标签失败
	ERROR_COUNT_TAG_FAIL  = 10005 //统计标签失败
	ERROR_ADD_TAG_FAIL    = 10006 //新增标签失败
	ERROR_EDIT_TAG_FAIL   = 10007 //修改标签失败
	ERROR_DELETE_TAG_FAIL = 10008 //删除标签失败

	//文章类错误
	ERROR_NOT_EXIST_ARTICLE        = 10011 //该文章不存在
	ERROR_CHECK_EXIST_ARTICLE_FAIL = 10012 //检查文章是否存在失败
	ERROR_ADD_ARTICLE_FAIL         = 10013 //新增文章失败
	ERROR_DELETE_ARTICLE_FAIL      = 10014 //删除文章失败
	ERROR_EDIT_ARTICLE_FAIL        = 10015 //修改文章失败
	ERROR_COUNT_ARTICLE_FAIL       = 10016 //统计文章失败
	ERROR_GET_ARTICLES_FAIL        = 10017 //获取多个文章失败
	ERROR_GET_ARTICLE_FAIL         = 10018 //获取单个文章失败
	ERROR_GEN_ARTICLE_POSTER_FAIL  = 10019 //生成文章海报失败
	//评论类错误
	ERROR_NOT_EXIST_COMMENT = 10020
	//验证类错误
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001 //Token鉴权失败
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002 //Token已过期
	ERROR_AUTH_TOKEN               = 20003 //Token生成失败
	ERROR_AUTH_GET_USER_FAIL       = 20004 //Token获取用户失败
	ERROR_AUTH                     = 20005 //Token错误
	//上传错误
	ERROR_UPLOAD_SAVE_IMAGE_FAIL    = 30001
	ERROR_UPLOAD_CHECK_IMAGE_FAIL   = 30002
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT = 30003
	//数据库类错误
	ERROR_NOT_CONNECT_FAIL = 40001  //数据链接错误

)

//定义错误消息
var MsgFlags = map[int]string{
	SUCCESS:               "success",
	ERROR:                 "fail",
	INVALID_PARAMS:        "请求参数错误",
	ERROR_EXIST_TAG:       "已存在该标签名称",
	ERROR_EXIST_TAG_FAIL:  "获取已存在标签失败",
	ERROR_NOT_EXIST_TAG:   "该标签不存在",
	ERROR_GET_TAGS_FAIL:   "获取所有标签失败",
	ERROR_COUNT_TAG_FAIL:  "统计标签失败",
	ERROR_ADD_TAG_FAIL:    "新增标签失败",
	ERROR_EDIT_TAG_FAIL:   "修改标签失败",
	ERROR_DELETE_TAG_FAIL: "删除标签失败",

	ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	ERROR_ADD_ARTICLE_FAIL:          "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:       "删除文章失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已过期",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH_GET_USER_FAIL:        "Token获取用户失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",
}

//获取错误消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
