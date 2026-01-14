package upload

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gofly/app/common/ffmpegtool"
	"gofly/utils/extend/uploads"
	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gconv"
	"strings"
	"time"
)

// 通用文件上传-小程序、app、h5等
type Index struct{ NoNeedAuths []string }

func init() {
	fpath := Index{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 业务端通用上传文件总接口
// 请求头添加 Businessid=当sass系统时记录(默认1账号)，filetype=附件类型(默认图片)
func (api *Index) UpFile(c *gf.GinCtx) {
	var businessID any = c.GetHeader("Businessid") //从请求头获取businessID判断是那个服务端传过来
	if businessID == "" {                          //找不到在去登录token获取
		businessID, _ = c.Get("businessID") //当前用户businessID(saas账号ID)
	}
	if businessID == "" { //找不到在去登录token获取
		businessID = 1 //默认单服务系统
	}
	//判断是否为saas结构是则增加上传文件限制
	if isExist, _ := gf.Model("admin_account").Exist(); isExist {
		//判断存储空间是否超出
		var usesize interface{}
		var fileSize interface{}
		usesize, _ = gf.Model("business_attachment").Where("business_id", businessID).Where("type", 0).Sum("filesize")
		if usesize == nil {
			usesize = 0
		}
		fileSize, _ = gf.Model("business_account").Where("id", businessID).Value("fileSize")
		if fileSize == nil {
			fileSize = 0
		}
		if gconv.Int(usesize) >= gconv.Int(fileSize) {
			gf.Failed().SetMsg("您的存储空间已满，请您先去购买存储空间！").Regin(c)
			return
		}
	}
	// 单个文件
	pid := c.DefaultPostForm("pid", "")
	filetype := c.DefaultPostForm("filetype", "image") //文件类型
	file, err := c.FormFile("file")
	if err != nil {
		gf.Failed().SetMsg("获取数据失败，").SetData(err).Regin(c)
		return
	}
	nowTime := time.Now().Unix() //当前时间
	//判断文件是否已经传过
	fileContent, _ := file.Open()
	defer fileContent.Close()
	var byteContainer []byte = make([]byte, 1000000)
	fileContent.Read(byteContainer)
	m_d5 := md5.New()
	m_d5.Write(byteContainer)
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	//查找该用户是否传过
	attachment, _ := gf.Model("business_attachment").Where("business_id", businessID).
		Where("sha1", sha1_str).Fields("id,pid,name,title,type,url,filesize,mimetype,cover_url as cover").Find()
	if attachment != nil { //文件是否已经存在
		//更新到最前面
		maxId, _ := gf.Model("business_attachment").Where("business_id", businessID).Order("weigh desc").Value("id")
		if maxId != nil {
			gf.Model("business_attachment").Data(map[string]interface{}{"weigh": maxId.Int() + 1, "pid": pid}).Where("id", attachment["id"]).Update()
		}
		gf.Success().SetMsg("文件已上传").SetData(attachment).Regin(c)
	} else {
		filename_arr := strings.Split(file.Filename, ".")
		location, _ := gcfg.Instance("upload").Get(c, "Type")
		//处理文件上传，bin返回地址
		url, cover_url, err := uploads.New().UploadFile(c, file)
		if err != nil {
			gf.Failed().SetMsg("上传文件失败").SetData(err).Regin(c)
			return
		}
		//文件类型
		var ftype int64 = 0
		switch filetype {
		case "video": //视频
			ftype = 2
			//使用ffmpeg生成是否封面
			videopath := fmt.Sprintf("./%s", url)
			pathroot := strings.Split(url, ".")
			imgpath := fmt.Sprintf("./%s", pathroot[0])
			fname, err := ffmpegtool.GetSnapshot(videopath, imgpath, 1)
			if err == nil {
				cover_url = fname
			}
		case "audio": //音频
			ftype = 3
		case "file": //附件类
			ftype = 4
		case "image": //图片
			ftype = 0
		default:
			ftype = 5
		}
		fileData := gf.Map{
			"business_id": businessID,
			"type":        ftype,
			"pid":         pid,
			"sha1":        sha1_str,
			"title":       filename_arr[0],
			"name":        file.Filename,
			"url":         url,       //附件路径
			"cover_url":   cover_url, //视频封面
			"storage":     location,
			"createtime":  nowTime,
			"filesize":    file.Size,
			"mimetype":    file.Header["Content-Type"][0],
		}
		//保存到数据库
		//保存数据
		file_id, _ := gf.Model("business_attachment").Data(fileData).InsertAndGetId()
		//更新排序
		gf.Model("business_attachment").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()

		fileData["id"] = file_id
		//处理预览url地址
		fileData["url"] = gf.GetFullUrl(url)
		gf.Success().SetMsg("文件上传成功").SetData(fileData).Regin(c)
	}
}
