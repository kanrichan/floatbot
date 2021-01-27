package onebot

// runOnebot run all server in config
func (conf *Yaml) runOnebot() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[OneBot] OneBot =X=> =X=> Start Error: %v", err)
			WARN("[OneBot] OneBot ==> ==> Sleep")
		}
	}()
	for i := range conf.BotConfs {
		for j := range conf.BotConfs[i].WSSConf {
			if conf.BotConfs[i].WSSConf[j].Status == 0 && conf.BotConfs[i].WSSConf[j].Enable == true && conf.BotConfs[i].WSSConf[j].Host != "" {
				go conf.BotConfs[i].WSSConf[j].start()
				go conf.BotConfs[i].WSSConf[j].listen()
				go conf.BotConfs[i].WSSConf[j].send()
			}
		}
		for k := range conf.BotConfs[i].WSCConf {
			if conf.BotConfs[i].WSCConf[k].Status == 0 && conf.BotConfs[i].WSCConf[k].Enable == true && conf.BotConfs[i].WSCConf[k].Url != "" {
				go conf.BotConfs[i].WSCConf[k].listen()
				go conf.BotConfs[i].WSCConf[k].send()
			}
		}
		for l := range conf.BotConfs[i].HTTPConf {
			if conf.BotConfs[i].HTTPConf[l].Status == 0 && conf.BotConfs[i].HTTPConf[l].Enable == true {
				if conf.BotConfs[i].HTTPConf[l].Host != "" {
					go conf.BotConfs[i].HTTPConf[l].listen()
				}
				go conf.BotConfs[i].HTTPConf[l].send()
			}
		}
	}
}
