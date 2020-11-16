package cqhttp

var XQPath = PathExecute()
var AppPath = XQPath + "data/app/onebot-yaya/"
var ImagePath = XQPath + "data/image/onebot-yaya/"
var RecordPath = XQPath + "data/record/onebot-yaya/"

func init() {
}

func Main() {
}

func OnStart() {
	CreatePath(AppPath)
	CreatePath(ImagePath)
	CreatePath(RecordPath)
	INFO("夜夜は世界一かわいい")
	Conf = Load(AppPath + "config.yml")
	if Conf == nil {
		ERROR("晚安~")
		return
	}
	WSCInit(Conf)
	go WSCStarts()
}
