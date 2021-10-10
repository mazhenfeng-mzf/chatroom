package process

import (
	"chatroom/public/errs"
	"chatroom/public/tools"

	//. "chatroom/public/tools"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (this *View) LoginView() {
	//var inCount, inPwd *walk.LineEdit

	layout := VBox{}
	children := []Widget{
		// HSplitter{
		// 	Children: []Widget{
		// 		TextEdit{AssignTo: &inTE},
		// 		TextEdit{AssignTo: &outTE, ReadOnly: true},
		// 	},
		// },
		TextLabel{
			Text: "欢迎使用聊天室",
			// MaxSize: Size{
			// 	Width:  20,
			// 	Height: 20,
			// },
		},
		HSplitter{
			Children: []Widget{
				Label{
					Text: "账号",
				},
				LineEdit{
					AssignTo:  &this.LoginCount,
					MaxLength: 50,
				},
			},
		},
		HSplitter{
			Children: []Widget{
				Label{
					Text: "密码",
				},
				LineEdit{
					AssignTo:  &this.LoginPwd,
					MaxLength: 50,
					OnKeyDown: func(key walk.Key) { //键盘事件
						if key == walk.KeyReturn { //回车键
							this.viewLogin()
						}
					},
				},
			},
		},
		PushButton{
			Text: "登陆",
			OnClicked: func() {
				this.viewLogin()
			},
		},
		PushButton{
			Text: "注册",
			OnClicked: func() {
				this.RegisterView()
			},
		},
	}

	// MyLOG.Log("Before: %#v", this.LoginPwd)

	this.mw_LoginPage.Layout = layout
	this.mw_LoginPage.Children = children

	// if first {
	this.mw_LoginPage.Run()
	// } else {
	// 	this.Run(true)
	// }
}

func (this *View) viewLogin() {

	// MyLOG.Log("After: %#v", this.LoginPwd)

	tools.MyLOG.Log("登陆账号: %s", this.LoginCount.Text())

	loginCheckErr := viewLoginCheckCountPwd(this.LoginCount.Text(), this.LoginPwd.Text())
	if loginCheckErr != nil {
		messageBox("错误", loginCheckErr.Error())
		return
	}

	loginErr := ClientProc.Login(this.LoginCount.Text(), this.LoginPwd.Text())
	if loginErr != nil {
		messageBox("错误", loginErr.Error())
		return
	}
	this.mwAssign_LoginPage.Close()
	this.MainPage()
}

func viewLoginCheckCountPwd(count string, pwd string) (err error) {

	if count == "" {
		err = errs.CountEmpty
		return
	}
	if pwd == "" {
		err = errs.PwdEmpty
		return
	}
	return
}
