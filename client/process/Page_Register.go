package process

import (
	"chatroom/public/errs"
	"errors"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var regId, regName, regPwd, regPwd_Sure *walk.LineEdit

func (this *View) RegisterView() {

	layout := VBox{}
	children := []Widget{
		// HSplitter{
		// 	Children: []Widget{
		// 		TextEdit{AssignTo: &inTE},
		// 		TextEdit{AssignTo: &outTE, ReadOnly: true},
		// 	},
		// },
		TextLabel{
			Text: "欢迎使用聊天室 - 注册账号",
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
					AssignTo:  &regId,
					MaxLength: 50,
				},
			},
		},
		HSplitter{
			Children: []Widget{
				Label{
					Text: "用户名",
				},
				LineEdit{
					AssignTo:  &regName,
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
					AssignTo:  &regPwd,
					MaxLength: 50,
				},
			},
		},
		HSplitter{
			Children: []Widget{
				Label{
					Text: "确认密码",
				},
				LineEdit{
					AssignTo:  &regPwd_Sure,
					MaxLength: 50,
				},
			},
		},
		PushButton{
			Text: "注册",
			OnClicked: func() {
				this.viewRegister()
			},
		},
	}

	this.mw_RegisterPage = MainWindow{
		Title:    "聊天室 - 注册账号",
		Size:     Size{400, 300},
		AssignTo: &this.mwAssign_RegisterPage,
		Layout:   layout,
		Children: children,
	}

	this.mw_RegisterPage.Run()

}

func (this *View) viewRegister() {

	regCheckErr := viewRegisterCheckCountPwd(regId.Text(), regName.Text(), regPwd.Text(), regPwd_Sure.Text())
	if regCheckErr != nil {
		messageBox("错误", regCheckErr.Error())
		return
	}

	registerErr := ClientProc.Register(regId.Text(), regName.Text(), regPwd.Text())
	if registerErr != nil {
		messageBox("错误", registerErr.Error())
		return
	}

	messageBox("注册成功", fmt.Sprintf("新用户 %s (%s) 注册成功, 请返回登录界面", regName.Text(), regId.Text()))
}

func viewRegisterCheckCountPwd(clientId string, name string, pwd string, pwd_sure string) (err error) {

	if clientId == "" {
		err = errs.CountEmpty
		return
	}
	if name == "" {
		err = errors.New("请输入用户名")
		return
	}
	if pwd == "" {
		err = errs.PwdEmpty
		return
	}
	if pwd_sure == "" {
		err = errs.PwdEmpty
		return
	}
	if pwd != pwd_sure {
		err = errors.New("两次输入密码不相同")
		return
	}
	return
}
