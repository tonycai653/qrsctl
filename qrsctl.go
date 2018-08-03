package main

import (
	"fmt"
	"net/http"
	"os"

	"qiniu/api.v6/conf"

	log "github.com/qiniu/log.v1"
	"github.com/qiniu/reqid.v1"
	"qbox.us/shell/devtools/ver"
	"qbox.us/shell/shutil"
	"qbox.us/shell/shutil/acc"
	"qbox.us/shell/shutil/cdnmgr"
	"qbox.us/shell/shutil/env"
	"qbox.us/shell/shutil/pfop"
	"qbox.us/shell/shutil/pub"
	"qbox.us/shell/shutil/rs"
	"qbox.us/shell/shutil/rsf"
	"qbox.us/shell/shutil/uc"
)

var cmdHandlers = map[string]func(args []string, qbox *env.Env) int{
	"put":              rs.Put,
	"get":              rs.Get,
	"stat":             rs.Stat,
	"cat":              rs.Cat,
	"rm":               rs.Delete,
	"del":              rs.Delete,
	"mv":               rs.Move,
	"cp":               rs.Copy,
	"chgm":             rs.Chgm,
	"chtype":           rs.ChType,
	"chmeta":           rs.Chmeta,
	"mkbucket":         rs.Mkbucket,  //功能废弃
	"mkbucket2":        rs.Mkbucket2, //功能废弃
	"buckets":          rs.Buckets,
	"redirect":         rs.Redirect,
	"sync":             rs.PutDir,
	"share":            rs.Share,
	"appinfo":          uc.AppInfo,
	"newAccess":        uc.NewAccess,
	"bucketinfo":       uc.BucketInfo,
	"bucketinfov1":     uc.BucketInfoV1,
	"img":              uc.Image,
	"rule/add":         uc.RuleAdd,
	"rule/update":      uc.RuleUpdate,
	"rule/del":         uc.RuleDel,
	"rule/get":         uc.RuleGet,
	"image":            uc.Image, //别名
	"unimg":            uc.Unimage,
	"unimage":          uc.Unimage, //别名
	"protected":        pub.SetProtected,
	"sep":              pub.Separator, //别名
	"separator":        pub.Separator,
	"style":            pub.Style,
	"unstyle":          pub.Unstyle,
	"preferStyleAsKey": uc.PreferStyleAsKey,
	"info":             acc.UserInfo,
	"fopAuth":          uc.FopAuth,
	"antiLeechMode":    uc.AntiLeechMode,
	"private":          uc.Private,
	"tokenAntiLeech":   uc.TokenAntiLeech,
	"newMacKey":        uc.NewMacKey,
	"delMacKey":        uc.DeleteMacKey,
	"listprefix":       rsf.ListPrefix,
	"cdn/refresh":      cdnmgr.Refresh,
	"cdn/refreshdir":   cdnmgr.RefreshDir,
	"cdn/bandwidth":    cdnmgr.Bandwidth,
	"cdn/flux":         cdnmgr.Flux,
	"setKeyState":      uc.SetKeyState,
	"pfop":             pfop.Pfop,
	"listjobs":         pfop.ListJobs,
	"imgsft":           uc.ImgSFT,
	"maxage":           uc.MaxAge,
	"noIndexPage":      uc.NoIndexPage,
	"persistFop":       uc.PersistFop,
	"styleCopy":        uc.StyleCopy,

	"setMirrorRawQuery":     uc.SetMirrorRawQueryOption,
	"setMirrorCheckHeaders": uc.SetMirrorCheckHeaders,

	"sourceHeaders/set":    uc.SourceHeadersSet,
	"sourceHeaders/get":    uc.SourceHeadersGet,
	"sourceHeaders/delete": uc.SourceHeadersDelete,
}

func Help() {
	fmt.Print(`
Usage:
  qrsctl [-l|d|lan|it|-f <hostFile>] -v login <User> <Passwd>                                               - Login
  qrsctl [-l|d|lan|it|-f <hostFile>] -v info                                                                - Show user information
  qrsctl [-l|d|lan|it|-f <hostFile>] -v appinfo [<AppName>]                                                 - Get application info
  qrsctl [-l|d|lan|it|-f <hostFile>] -v put -c <Bucket> <Key> <SrcFile>                                     - Put file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v get <Bucket> <Key> <DestFile>                                       - Get file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v stat <Bucket> <Key>                                                 - Stat file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v cat <Bucket> <Key>                                                  - Cat file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v del <Bucket> <Key>                                                  - Delete a file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v mv <Bucket1:Key1> <Bucket2:Key2>                                    - Move file from Bucket1:Key1 to Bucket2:Key2
  qrsctl [-l|d|lan|it|-f <hostFile>] -v cp <Bucket1:Key1> <Bucket2:Key2>                                    - Copy file
  qrsctl [-l|d|lan|it|-f <hostFile>] -v chgm <Bucket> <key> <mimeType>                                      - Change MimeType
  qrsctl [-l|d|lan|it|-f <hostFile>] -v chtype <Bucket> <key> <type>                                        - Change file type, <type>:  0 (normal), 1 (line) 
  qrsctl [-l|d|lan|it|-f <hostFile>] -v listprefix <bucket> <prefix> [<limit>] [<marker>]                   - List files
  qrsctl [-l|d|lan|it|-f <hostFile>] -v mkbucket <BucketName> <Zone>                                        - Create a bucket, <Zone>:z0, z1, z2, na0
  qrsctl [-l|d|lan|it|-f <hostFile>] -v buckets <Shared>                                                    - List all buckets
  qrsctl [-l|d|lan|it|-f <hostFile>] -v share <Bucket> <Uid> <Permission>
                                        PermissionOptions: 1(RD) 2(RW) -1(Cancel Share)                     - Share bucket
  qrsctl [-l|d|lan|it|-f <hostFile>] -v listprefix <bucket> <prefix> [<limit>] [<marker>]                   - List files buckets
  qrsctl [-l|d|lan|it|-f <hostFile>] -v bucketinfo <Bucket>                                                 - Get bucket info
  qrsctl [-l|d|lan|it|-f <hostFile>] -v img <Bucket> <SrcUrl> [<SrcHost>] [<Expires>]                       - Image bucket with source
  qrsctl [-l|d|lan|it|-f <hostFile>] -v unimg <Bucket>                                                      - Unimage bucket
  qrsctl [-l|d|lan|it|-f <hostFile>] -v protected <Bucket> <Protected>                                      - Set bucket protected or not
  qrsctl [-l|d|lan|it|-f <hostFile>] -v separator <Bucket> <Sep>                                            - Set style separator
  qrsctl [-l|d|lan|it|-f <hostFile>] -v style <Bucket> <Name> <Style>                                       - Set style
  qrsctl [-l|d|lan|it|-f <hostFile>] -v unstyle <Bucket> <Name>                                             - Unset style
  qrsctl [-l|d|lan|it|-f <hostFile>] -v styleCopy <bucket_Src> <bucket_Dest>                                - Copy styles
  qrsctl [-l|d|lan|it|-f <hostFile>] -v private <Bucket> <Private>                                          - Set bucket private or not
  qrsctl [-l|d|lan|it|-f <hostFile>] -v imgsft  <Bucket> <imgsft>                                           - Set bucket image storage with fault tolerant
  qrsctl [-l|d|lan|it|-f <hostFile>] -v noIndexPage <Bucket> <0/1>                                          - Turn On/Off bucket index page
  qrsctl [-l|d|lan|it|-f <hostFile>] -v redirect <Bucket> <Key> <RedirectUrl> [<RedirectCode>]              - Redirect a key to an url

  qrsctl [-l|d|lan|it|-f <hostFile>] -v rule/add <Bucket> <ruleName> <prefix> <deleteAfterDays> <toLineAfterDays>
                                                                                                            - add bucket rule
  qrsctl [-l|d|lan|it|-f <hostFile>] -v rule/update <Bucket> <ruleName> <prefix> <deleteAfterDays> <toLineAfterDays>
                                                                                                            - update bucket rule
  qrsctl [-l|d|lan|it|-f <hostFile>] -v rule/del <Bucket> <ruleName>                                        - del bucket rule
  qrsctl [-l|d|lan|it|-f <hostFile>] -v rule/get <Bucket>                                                   - get bucket rule

  qrsctl [-l|d|lan|it|-f <hostFile>] -v sourceHeaders/set  <Bucket> <Header> <Value>                        - set source header
  qrsctl [-l|d|lan|it|-f <hostFile>] -v sourceHeaders/get  <Bucket>                                         - get source header
  qrsctl [-l|d|lan|it|-f <hostFile>] -v sourceHeaders/delete <Bucket> <Header>                              - delete source header

  qrsctl [-l|d|lan|it|-f <hostFile>] -v pfop <bucket> <key> <fops> [<notifyURL>] [<force>] [<pipeline>]     - Do pfop
  qrsctl [-l|d|lan|it|-f <hostFile>] -v listjobs <pipelineId> [<marker>] [<limit>]                          - List jobs of pfop

  qrsctl [-l|d|lan|it|-f <hostFile>] -v cdn/refresh <Url1> <Url2>...<UrlN>                                  - Refresh cdn cache for urls
  qrsctl [-l|d|lan|it|-f <hostFile>] -v cdn/bandwidth <domains> <start_date> <end_date> [<granularity>]     - Get bandwidth of domains
  qrsctl [-l|d|lan|it|-f <hostFile>] -v cdn/flux <domains> <start_date> <end_date> [<granularity>]          - Get traffic of domains
  qrsctl [-l|d|lan|it|-f <hostFile>] -v setMirrorRawQuery <Bucket> <1|0>                                    - set mirror Raw Query support

Authorization:
  1) qrsctl login <User> <Passwd>: and then remember the login token
  2) qrsctl login <AccessKey> <SecretKey>
  3) qrsctl -a <AccountConf>: provide an account config file that provide access_key & secret_key (in json format)
BuildVersion:
  qrsctl v`, ver.String, `
`)
}

func main() {

	log.SetOutputLevel(1)
	if len(os.Args) < 2 {
		Help()
		os.Exit(0)
	}
	args := os.Args[1:]
	cmd := args[0]

	var qbox *env.Env
	if cmd == "redirect" {
		if len(args) < 4 || len(args) > 5 {
			fmt.Println("redirect <Bucket> <Key> <RedirectUrl> [<RedirectCode>]")
			os.Exit(-2)
		}
	}
	qbox = shutil.GetEnv("qrsctl")
	if cmd == "login" {
		shutil.Login(args, qbox)
		return
	}

	if len(os.Args) < 2 {
		Help()
		os.Exit(-1)
	}

	conf.SetUser("qrsctl-" + ver.String)

	if f, ok := cmdHandlers[cmd]; ok {
		shutil.TryLogin(qbox)
		reqId := reqid.Gen()
		qbox.Transport = &ReqIdTransport{Transport: qbox.Transport, ReqId: reqId}
		rv := f(os.Args[1:], qbox)
		if rv != 0 {
			fmt.Println("RequestId: ", reqId)
		}
		os.Exit(rv)
	}

	fmt.Println("Unknown command:", cmd)
}

type ReqIdTransport struct {
	ReqId     string
	Transport http.RoundTripper
}

func (t *ReqIdTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Reqid", t.ReqId)
	return t.Transport.RoundTrip(req)
}
