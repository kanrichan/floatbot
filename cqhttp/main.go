package cqhttp

var FirstStart bool = true

var XQPath = PathExecute()
var AppPath = XQPath + "onebot/"
var ImagePath = XQPath + "onebot/image/"
var RecordPath = XQPath + "onebot/record/"
var VideoPath = XQPath + "onebot/video/"

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
		go Conf.runOnebot()
		go Conf.heartBeat()
	}
	FirstStart = false
}

func onDisable() {
	//
}
