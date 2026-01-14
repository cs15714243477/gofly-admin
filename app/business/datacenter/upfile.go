package datacenter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gofly/app/common/ffmpegtool"
	"gofly/utils/extend/uploads"
	"gofly/utils/gf"
	"gofly/utils/tools/gcfg"
	"gofly/utils/tools/gconv"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 文件管理
type Upfile struct{ NoNeedAuths []string }

func init() {
	fpath := Upfile{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// 上传文件总接口
func (api *Upfile) Upload(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID") //当前商户ID
	//判断是否为saas结构是则增加上传文件限制
	if isExist, _ := gf.Model("admin_account").Exist(); isExist {
		//判断存储空间是否超出
		usesize, err := gf.Model("business_attachment").Where("business_id", businessID).Where("type", 0).Sum("filesize")
		if err != nil {
			usesize = 0
		}
		var fileSizes int64 = 0
		fileSize, err := gf.Model("business_account").Where("id", businessID).Value("fileSize")
		if err == nil && fileSize != nil {
			fileSizes = gconv.Int64(fileSize)
		}
		if gconv.Int64(usesize) >= fileSizes {
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
			gf.Failed().SetMsg(err.Error()).SetData(err).Regin(c)
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
			"createtime":  time.Now().Unix(),
			"filesize":    file.Size,
			"mimetype":    file.Header["Content-Type"][0],
		}
		//保存到数据库
		file_id, _ := gf.Model("business_attachment").Data(fileData).InsertAndGetId()
		//更新排序
		gf.Model("business_attachment").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
		fileData["id"] = file_id
		//处理预览url地址
		fileData["url"] = gf.GetFullUrl(url)
		gf.Success().SetMsg("文件上传成功").SetData(fileData).Regin(c)
	}
}

// 2.上传文件到本地
func (api *Upfile) LocalFile(c *gf.GinCtx) {
	// 单个文件
	file, err := c.FormFile("file")
	if err != nil {
		gf.Failed().SetMsg("获取数据失败").SetData(err).Regin(c)
		return
	}
	businessID := c.GetInt64("businessID") //当前商户ID
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
			gf.Model("business_attachment").Data(map[string]interface{}{"weigh": maxId.Int() + 1, "pid": 1}).Where("id", attachment["id"]).Update()
		}
		gf.Success().SetMsg("文件已上传").SetData(attachment).Regin(c)
	} else {
		//处理文件上传，bin返回地址
		url, cover_url, err := uploads.New("local").UploadFile(c, file)
		if err != nil {
			gf.Failed().SetMsg("上传文件失败").SetData(err).Regin(c)
			return
		}
		filename_arr := strings.Split(file.Filename, ".")
		//文件类型
		fileData := gf.Map{
			"business_id": 1,
			"type":        0, //图片
			"pid":         1, //存储在默认文件夹
			"sha1":        sha1_str,
			"title":       filename_arr[0],
			"name":        file.Filename,
			"url":         url,       //附件路径
			"cover_url":   cover_url, //封面
			"storage":     "local",
			"createtime":  time.Now().Unix(),
			"filesize":    file.Size,
			"mimetype":    file.Header["Content-Type"][0],
		}
		//保存数据
		file_id, _ := gf.Model("business_attachment").Data(fileData).InsertAndGetId()
		//更新排序
		gf.Model("business_attachment").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
		//处理预览url地址
		gf.Success().SetMsg("文件上传成功").SetData(gf.GetFullUrl(url)).Regin(c)
	}
}

// 3.编辑器保存第三方图片到本地
func (api *Upfile) ThirdImage(c *gf.GinCtx) {
	businessID, _ := c.Get("businessID") //当前用户businessID(saas账号ID)
	if businessID == "" {                //找不到在去登录token获取
		businessID = 1 //默认单服务系统
	}
	params, _ := gf.RequestParam(c)
	if url, ok := params["url"]; !ok || url == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"result": false,
			"data": map[string]interface{}{
				"url": "",
			},
			"message": "地址无效",
		})
	} else {
		file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
		if _, err := os.Stat(file_path); err != nil {
			if !os.IsExist(err) {
				os.MkdirAll(file_path, os.ModePerm)
			}
		}
		nowTime := time.Now().Unix() //当前时间
		localPicName := fmt.Sprintf("%vlocal_%v", file_path, nowTime)
		imgtype, _ := gf.DownPic(gconv.String(params["url"]), localPicName)
		imgurl := fmt.Sprintf("%s%s", localPicName, imgtype)
		//判断文件是否存在
		path, _ := os.Getwd()
		filego_path := filepath.Join(path, imgurl)
		if _, err := os.Stat(filego_path); err == nil {
			Insertdata := map[string]interface{}{
				"business_id": businessID,
				"type":        0,
				"pid":         0,
				"title":       fmt.Sprintf("local_%v%v", nowTime, imgtype),
				"name":        fmt.Sprintf("local_%v", nowTime),
				"url":         imgurl,
				"storage":     "local",
				"createtime":  nowTime,
				"mimetype":    "image/png",
			}
			//保存数据
			gf.Model("business_attachment").Data(Insertdata).InsertAndGetId()
			gf.Success().SetMsg("文件上传成功").SetData(gf.GetLocalFullUrl(imgurl)).Regin(c)
		} else {
			c.JSON(200, gin.H{
				"code":   400,
				"result": false,
				"data": map[string]interface{}{
					"url": "",
				},
				"message": "保存失败",
			})
		}
	}
}
