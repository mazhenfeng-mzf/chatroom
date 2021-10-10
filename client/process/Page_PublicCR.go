package process

import (
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (this *View) publicCRPage() {
	layout := VBox{}
	children := []Widget{
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					AssignTo: &this.PcCR.PcMessageWindow,
					ReadOnly: true,
					VScroll:  true,
				},
				TextEdit{
					AssignTo: &this.PcCR.PcOnlineClient,
					ReadOnly: true,
					VScroll:  true,
				},
			},
		},
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					AssignTo: &this.PcCR.PcInput,
					OnKeyDown: func(key walk.Key) { //键盘事件
						if key == walk.KeyReturn { //回车键
							ClientProc.SendPCSms(this.PcCR.PcInput.Text())
						}
					},
				},
				PushButton{
					Text: "发送",
					OnClicked: func() {
						ClientProc.SendPCSms(this.PcCR.PcInput.Text())
					},
				},
			},
		},
	}

	this.PcCR.mw.Title = fmt.Sprintf("公共聊天室 - %s", ClientProc.Name)
	this.PcCR.mw.Size = Size{800, 600}
	this.PcCR.mw.Layout = layout
	this.PcCR.mw.Children = children
	this.PcCR.mw.AssignTo = &this.PcCR.mwAssign

	MyLOG.Log("主程序已进入公共聊天室, 让子程序完成填充内容")
	//MyLOG.Log("this.PcCR.PcOnlineClient: %#v", this.PcCR.PcOnlineClient)
	this.PcCR.mw.Run()

	MyLOG.Log("退出公共聊天室")
	ClientProc.ExitPC()
}
