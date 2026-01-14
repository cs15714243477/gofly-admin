package datacenter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gofly/app/common/ffmpegtool"
	"gofly/utils/extend/uploads"
	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"strings"
	"time"
)

// 文件管理
type Upfile struct{ NoNeedAuths []string }

func init() {
	fpath := Upfile{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 上传admin端图片 系统默认附件attachment表，不限制空间
func (api *Upfile) Upload(c *gf.GinCtx) {
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
	var byteContainer []byte = make([]byte, 1000000)
	fileContent.Read(byteContainer)
	m_d5 := md5.New()
	m_d5.Write(byteContainer)
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	//查找该用户是否传过
	attachment, _ := gf.Model("admin_attachment").
		Where("sha1", sha1_str).Fields("id,pid,name,title,type,url,filesize,mimetype,cover_url").Find()
	if attachment != nil { //文件是否已经存在
		//更新到最前面
		maxId, _ := gf.Model("admin_attachment").Order("weigh desc").Value("id")
		if maxId != nil {
			gf.Model("admin_attachment").Data(map[string]interface{}{"weigh": maxId.Int() + 1, "pid": pid}).Where("id", attachment["id"]).Update()
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
			"type":       ftype,
			"pid":        pid,
			"sha1":       sha1_str,
			"title":      filename_arr[0],
			"name":       file.Filename,
			"url":        url,       //附件路径
			"cover_url":  cover_url, //封面
			"storage":    location,
			"createtime": nowTime,
			"filesize":   file.Size,
			"mimetype":   file.Header["Content-Type"][0],
		}
		//保存到数据库
		file_id, _ := gf.Model("admin_attachment").Data(fileData).InsertAndGetId()
		//更新排序
		gf.Model("admin_attachment").Where("id", file_id).Data(map[string]interface{}{"weigh": file_id}).Update()
		fileData["id"] = file_id
		//处理预览url地址
		fileData["url"] = gf.GetFullUrl(url)
		gf.Success().SetMsg("文件上传成功").SetData(fileData).Regin(c)
	}
}
