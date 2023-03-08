package doubangroupjs

import (
	"github.com/Tiandapaopao/crawler/collect"
)

var DoubangroupJSTask = &collect.TaskModle{
	Property: collect.Property{
		Name:     "js_find_douban_sun_room",
		WaitTime: 2,
		MaxDepth: 5,
		Cookie:   "__utma=30149280.609234316.1675058944.1676360711.1676443276.12; __utmb=30149280.2.10.1676443276; __utmc=30149280; __utmt=1; __utmv=30149280.19775; __utmz=30149280.1676013678.3.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/; push_doumail_num=0; push_noty_num=0; __gpi=UID=00000bb0bfe53e79:T=1675058959:RT=1676443268:S=ALNI_MYAKZt9MDINtpmS9LfugzL2iJA7uw; ap_v=0,6.0; _pk_id.100001.8cb4=14ce418015a1a446.1676013677.10.1676443265.1676360727.; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1676443265%2C%22https%3A%2F%2Ftime.geekbang.org%2F%22%5D; _pk_ses.100001.8cb4=*; ck=2kCv; ct=y; douban-fav-remind=1; dbcl2=\"197752134:rDFJ+M8i7fk\"; __gads=ID=87e1dea389dd72db-220535d3b2d9003f:T=1676013679:RT=1676013679:S=ALNI_MaT031PxpmNS81fN7vDuyMHqvkbGA; __yadk_uid=Hgf0RfykrqGvtLG6Yejj19BgmSf8ssDr; viewed=\"26416768_1007305_35720728\"; ll=\"118318\"; bid=xYSZqPEUrRo",
	},
	Root: `
 		var arr = new Array();
  		for (var i = 25; i <= 25; i+=25) {
 			var obj = {
 			   Url: "https://www.douban.com/group/szsh/discussion?start=" + i,
 			   Priority: 1,
 			   RuleName: "解析网站URL",
 			   Method: "GET",
 		   };
 			arr.push(obj);
 		};
 		console.log(arr[0].Url);
 		AddJsReq(arr);
 			`,
	Rules: []collect.RuleModle{
		{
			Name: "解析网站URL",
			ParseFunc: `
 			ctx.ParseJSReg("解析阳台房","(https://www.douban.com/group/topic/[0-9a-z]+/)\"[^>]*>([^<]+)</a>");
 			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
 			//console.log("parse output");
 			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
 			`,
		},
	},
}
