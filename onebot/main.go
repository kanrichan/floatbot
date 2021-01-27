package onebot

var FirstStart bool = true

var XQPath = PathExecute()
var AppPath = XQPath + "OneBot/"
var ImagePath = XQPath + "OneBot/image/"
var RecordPath = XQPath + "OneBot/record/"
var VideoPath = XQPath + "OneBot/video/"
var CachePath = XQPath + "OneBot/cache/"

var PicPool = PicsCache{Max: 1000}

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
		go Conf.runDB()
		go Conf.runOnebot()
		apiMap.Register(&apiMap.this)
	}
	FirstStart = false
}

func onDisable() {
}
