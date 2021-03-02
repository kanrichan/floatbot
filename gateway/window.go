package gateway

import (
	"fmt"
	core "onebot/core/xianqu"
)

func OnSetting(ctx *core.Context) {
	OnDisable(ctx)
	core.XQApiCallMessageBox(fmt.Sprintf("等个好心人写UI，修改配置请到 %sconfig.yml\n", core.OneBotPath))
}
