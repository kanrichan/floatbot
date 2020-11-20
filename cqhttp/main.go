package cqhttp

var FirstStart bool = true

var XQPath = PathExecute()
var AppPath = XQPath + "data/app/onebot-yaya/"
var ImagePath = XQPath + "data/image/onebot-yaya/"
var RecordPath = XQPath + "data/record/onebot-yaya/"
var VideoPath = XQPath + "data/video/onebot-yaya/"

func init() {
}

func Main() {
}

func onStart() {
	if FirstStart {
		CreatePath(AppPath)
		CreatePath(ImagePath)
		CreatePath(RecordPath)
		CreatePath(VideoPath)
		INFO("夜夜は世界一かわいい")
		Conf = Load(AppPath + "config.yml")
		if Conf == nil {
			ERROR("晚安~")
			return
		}
		WSCInit(Conf)

		go WSCStarts()
	}
	FirstStart = false
}

func onDisable() {
	//
}
