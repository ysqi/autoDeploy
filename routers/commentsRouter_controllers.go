package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/autoDeploy/controllers:GitHookController"] = append(beego.GlobalControllerRouter["github.com/ysqi/autoDeploy/controllers:GitHookController"],
		beego.ControllerComments{
			Method: "Payload",
			Router: `/payload`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
