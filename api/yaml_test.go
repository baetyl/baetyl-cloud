package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	mf "github.com/baetyl/baetyl-cloud/v2/mock/facade"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/service"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testSecret = `
apiVersion: v1
kind: Secret
metadata:
  name: dcell
  labels:
    secret: dcell
  annotations:
    secret: dcell
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
type: Opaque`
	updateSecret = `
apiVersion: v1
kind: Secret
metadata:
  name: dcell
  labels:
    secret: dcell
  annotations:
    secret: dcell
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
  123: YWJj
type: Opaque`

	testRegistry = `
apiVersion: v1
data:
  .dockerconfigjson: eyJhdXRocyI6eyJET0NLRVJfUkVHSVNUUllfU0VSVkVSIjp7InVzZXJuYW1lIjoiRE9DS0VSX1VTRVIiLCJwYXNzd29yZCI6IkRPQ0tFUl9QQVNTV09SRCIsImVtYWlsIjoiRE9DS0VSX0VNQUlMIiwiYXV0aCI6IlJFOURTMFZTWDFWVFJWSTZSRTlEUzBWU1gxQkJVMU5YVDFKRSJ9fX0=
kind: Secret
metadata:
  name: myregistrykey
  namespace: default
type: kubernetes.io/dockerconfigjson`
	updateRegistry = `
apiVersion: v1
data:
  .dockerconfigjson: eyJhdXRocyI6eyJET0NLRVJfUkVHSVNUUllfU0VSVkVSIjp7InVzZXJuYW1lIjoiRE9DS0VSX1VTRVIiLCJwYXNzd29yZCI6IkRPQ0tFUl9QQVNTV09SRDEyMyIsImVtYWlsIjoiRE9DS0VSX0VNQUlMIiwiYXV0aCI6IlJFOURTMFZTWDFWVFJWSTZSRTlEUzBWU1gxQkJVMU5YVDFKRSJ9fX0=
kind: Secret
metadata:
  name: myregistrykey
  namespace: default
type: kubernetes.io/dockerconfigjson`

	testCert = `apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURpekNDQW5PZ0F3SUJBZ0lRWENzeE02UnpjcTJVUk5ROHBINjZKekFOQmdrcWhraUc5dzBCQVFzRkFEQVYKTVJNd0VRWURWUVFERXdwcmRXSmxjbTVsZEdWek1CNFhEVEl4TVRJeU16QTJOREF5TmxvWERUSXlNVEl5TXpBMgpOREF5Tmxvd0xURXJNQ2tHQTFVRUF4TWlZbUZsZEhsc0xYZGxZbWh2YjJzdGMyVnlkbWxqWlM1a1pXWmhkV3gwCkxuTjJZekNDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFNeEdhays0QXd6ZDI2VWMKMzhvM05scVp2RWcrUFJqRkRvV1F0NnB1Y3ZDdFo0MUdGWXdFODlaMm1tT1RSeXlma1d1Mi95aWgzbkMwYTdvSwpxMmF5dzJzYTExRWU3TkZYMUY2MzB3UHNmZ1A0aG1LZmdIbEd6N05hTHU1SWQyNTNtZWtwZmI0QUs1UnVkNmZWCjEyOXJ1WXJUSDFLSHl3cWdkTGdqWFlVS0Q2VGhPYk1hTW5vVG5haTVnR3Evd1VMVVR0cy8zcUNaTU53ajlrNkoKR2hwWTQ3cFFhck1MVU1tWTl4TzExcDZNVzFIelRFQ2VQM0FQMUtVenZUMFd5U0kyb2NxdENIMlN1bUl5aW5uUwpEcnVTU3Jxb0pKOFE0MHhwdjROTjdqYm55M0k2TDZyWERJV2kwUmFZQ2pTaHc3TGlZSEVtcENBWEIxSHl3VmswCnZlNTFYRXNDQXdFQUFhT0J2akNCdXpBT0JnTlZIUThCQWY4RUJBTUNCYUF3RXdZRFZSMGxCQXd3Q2dZSUt3WUIKQlFVSEF3RXdEQVlEVlIwVEFRSC9CQUl3QURBZkJnTlZIU01FR0RBV2dCUWx0aXVhem9leW5IZVdJQXJWVC9EKwpEeGM1V2pCbEJnTlZIUkVFWGpCY2doWmlZV1YwZVd3dGQyVmlhRzl2YXkxelpYSjJhV05sZ2g1aVlXVjBlV3d0CmQyVmlhRzl2YXkxelpYSjJhV05sTG1SbFptRjFiSFNDSW1KaFpYUjViQzEzWldKb2IyOXJMWE5sY25acFkyVXUKWkdWbVlYVnNkQzV6ZG1Nd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFHbkVMMzFwaitZdkFUa3N3L09pK2dBTwp5a05uYVhlZ2dLMEtnMW5vZmhXNXlhajBHNU9yVHFyVTBsYm5SVDEwTUgyKzZpcG9JSEVPNEliSTRnUEo0Z2JqCllNb2JQdTJGeDN0TVd2SStTcEs1NEJMT3FlZk5VMEJPV2pwSU5Vcis2MGl1OFZaNDhhYnVLZ0FjUmJSNktiQTIKRlpFN2VsZ0JHbnJ6ZUh1NHlNdEx2VEI2VUNzMFZnL0YvRkdVWjJ1ZnU3bEM5dFVFc3c0U3ZFSEp6ZEZLRnBNOAp3cGZOQU5WOWVEL1dlWDJRNEY5ci9NZm1XUFdIS2pJWTlTRzN1c29VUjlKZzgvZlZFaXBnRlRYQ3NlSnFEb1RECk9hN0J6TnVBZ3NlOVVGUmswc2ZGZnJkM1N4WkVmNHhZVzkwcmtIU3crL0c3Ni9RS3BYeTNaV0FFTXVnZmxsaz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBekVacVQ3Z0RETjNicFJ6ZnlqYzJXcG04U0Q0OUdNVU9oWkMzcW01eThLMW5qVVlWCmpBVHoxbmFhWTVOSExKK1JhN2IvS0tIZWNMUnJ1Z3FyWnJMRGF4clhVUjdzMFZmVVhyZlRBK3grQS9pR1lwK0EKZVViUHMxb3U3a2gzYm5lWjZTbDl2Z0FybEc1M3A5WFhiMnU1aXRNZlVvZkxDcUIwdUNOZGhRb1BwT0U1c3hveQplaE9kcUxtQWFyL0JRdFJPMnovZW9Ka3czQ1AyVG9rYUdsamp1bEJxc3d0UXlaajNFN1hXbm94YlVmTk1RSjQvCmNBL1VwVE85UFJiSklqYWh5cTBJZlpLNllqS0tlZElPdTVKS3VxZ2tueERqVEdtL2cwM3VOdWZMY2pvdnF0Y00KaGFMUkZwZ0tOS0hEc3VKZ2NTYWtJQmNIVWZMQldUUzk3blZjU3dJREFRQUJBb0lCQVFDeG5uU29ObzlvYTZ5bAp5QktMQ1RFWTNGNUx6dHBmSkZFNU1CbVRkeEE0Vi84WFp4NHJMczg2NXN2ZDAxenEzeUNhTVhkeVJEVmZlSEhIClJhbTkxYWg0QTlHL05vMmloYVVpYXVKdm9mMzh3K2RONy9UTS94WndrL1VFdWp3bThKWUNtRkkrbUhWTVRqVjYKMlVUSEhEc3NDK0ZMYU1uU3hxeFd4R0YzNGNTeWpwb3VWNkdxTmY3eDBCa2llb1NieEZTd1cwa3ZreGJsdnpiZQpBTGpmcW1JcTJuWlo2OFZISTRncGR4dlIvU09ENjFRVzc1QU1wUGplbytzOHluaWZ3dDV0Tzh1dU5CRkE2NmJLCjlqRFFuc096RXR5aXZlN1AvWC9BME5YMEdNUFVHOHRBZjFMMFNlV0E4WmtwQ3RsU05KZlVLTWIwM0g3Y2hQVDgKQVhlL0R2MkJBb0dCQU94dytTUVhHQ0lzL3NNUmZ4YjZBS2loSHdDc1BieG1Yck8rNlZFK1F5VHlvMGFYdWxVbQo2UjZxNkNnRlhSWitEL2NKWEpSTUZmYmczK3dlVmtvZkdDeE5rZHl5bDVVUTdHczlVRXd0OXJZMy9PQ1Z5Zjc3Ck4rc0NpbkY0Z0E3amlzRVRSUGhVUUlBVk1wSmlDaHRvTTBTNDBVeXZwMzBTczNDZVBwcTZaczBMQW9HQkFOMHMKUlJnaXlZNjRmbHlrZGl0ZEZubnpocG1Gc3U3NlVHb0ZlYlAzMzdQSk51THJqRWxQUnZTSUtlczlzVUVsN0xRNQpnajg2TjhZY2txRmJyZWVaVkZsYlZNbGZpRjhpaFpacjl2cy9pNFhXeE5aK3ZlM2tEZngrTlpQYmhycmpBbS9JCkFDNXIzTWZNaWNlZEt5c2xlNVd0TjJkVGRBQ1NsUlVLTWxMTmRMWEJBb0dBWHN3YzE5OTZpWmxJdTZVME0xNGgKRFhzc0Z2VDMrNlYvcXNtTWVrcGdXVnYvSXJxS3RzRlhEaml2dy93Q2lwWVlpSTkwVXZEK2pYRXoxbE9EZlV4aQpRTUVKRGxkOGR3UEdCbWthM0xCQkRtWDhPWDlVOGFwL2pQWUQwK0xnVlJmZDlmTm4zN2pIODVLTUtDeXVxTFpxCmQ4OHgrM0VoMGYvQmVoRzRRQWtrVm1rQ2dZRUFzTlRxVVVmTyt1c0xMS3JaU0FaZktCWEtzZ2d4YmR4NFdxd1MKQ0EvUXJYL2RBRVR2bnRWaGw3VWVQdFRPV1pZbTBGbUNoMmJXblBEUFUyOW5kVm9rRkdWdlBxbkE4TDg3SzI4YQp3dnFsWk5hMy9mN0xmOTNzU01ubnNGVytQTUd2ZXd2ZkNUNTRBTTdLQWV6cFRNL2xKV0NlZ1dBNXlSTnBXcThTCldSMm5pSUVDZ1lBMTVUUjZtaXlFdW1Xbnp4MDBURk0zZDNTTC9LejFGSkRSd1BPam45Ri8vRFVKZ0J0RFErQWwKalc5MTFrekpvZjRBMjhVd0FKMFhQbzd5Zkx2bEpobWlPU1BhNi9ja0lsZ0E4S0pFZ2l4RExROHk5QjZqajFoQgpWZk84VVJSL3FUZG9uMTd6L1lHQ09BUzhzWTQ2UEova2J2R1BiUVpOOGFZaGFRSGZidCtUSWc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
kind: Secret
metadata:
  namespace: default
  name: baetyl-tls-secret
type: kubernetes.io/tls`
	updateCert = `apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN2akNDQWFZQ0NRQ3JxajNTckI4WnhqQU5CZ2txaGtpRzl3MEJBUXNGQURBaE1Rc3dDUVlEVlFRR0V3SkQKVGpFU01CQUdBMVVFQXd3SlkyRjBkR3hsTFdOaE1CNFhEVEl3TURrd01qQTRNVEUxTmxvWERUTXdNRGd6TVRBNApNVEUxTmxvd0lURUxNQWtHQTFVRUJoTUNRMDR4RWpBUUJnTlZCQU1NQ1dOaGRIUnNaUzFqWVRDQ0FTSXdEUVlKCktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQUtYZUwzUVZEcmJsdnIrR0Y1S2lWQmNzOWZEZDdaTk4KQlg5MStVbjA5TitxaEF1RWdSYkdROGpRRitXNFZ5cENqRjlLQUZyZTE2TUtRckZMYWFRb1pJbTN3RURiVUd6ZwpUazFRdlZDUkxYbi8zblJGY3o4Y2hHL0l3YXQrZmFYbnQ4MmFGZWVSTU1UZTdscXRhRXY4MmZvdXEwcDNINmlqCmN5UC9aQzZNc1hXdWdvSndBYXlLUk1IVExhbHJQRzlXUDlLdnJQL2VLTUhPSTg4cWNVMWtzRnk2OWQ3TzZkbVEKeldpZThaZHdDL004eWNQVFNtUE9QUWlKZ0JydUJUdkNmbXVjNVV4SXFDakNzOXp1bzFMMldWNlVTMEhRT2dEbQpTbm9xTThEVjZpaVhEY01mSU81R1U4azhhb3A2dnpub2hvbmo2L1hNQkphbnZDcE1EK0tTRmFVQ0F3RUFBVEFOCkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQVBVVUZnT09lbm1COGVleGV0VDJvOEhLRmIxM2plUUZFaWFzdU95TTQKZ2I4dTFMMGUzT1VITFFpNENKbWw3c3p3MTBzZGM0L3g2d3F0eEpkOEsxMW0xeEo4TnhxbXNtR0dlZHpHdElZdQpsWjZ4dGhjZ3B2d1NvRHFWUXI0RllZWWpDbXVyTnNhWUZ1N0dWU3FQQ2dMMFNiaXFObWZMSzVzS2Y5QlZrTWx0CmNETnNWbS96eE50bVJGeVRNOHVHWDNSb0Eva2x0eG5FbEFPcVlKRUFjZVBNUTRJb1RoZmlmL3FsM1kxRDZKNXUKTzBvN0lNNXRsRDhVaFlmQVlWSC94eUt4ZHY2SzIzekhnYnVGNFFXckJFbllleTBlcXpuTys1Yk5WN2J6MHBrYwp0dTVud3c1UmRqQ3o0VWtzMDhQMkdObVpqTE84MU1nWWtoUjdCOXdpM0tETnhnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBcGQ0dmRCVU90dVcrdjRZWGtxSlVGeXoxOE4zdGswMEZmM1g1U2ZUMDM2cUVDNFNCCkZzWkR5TkFYNWJoWEtrS01YMG9BV3Q3WG93cENzVXRwcENoa2liZkFRTnRRYk9CT1RWQzlVSkV0ZWYvZWRFVnoKUHh5RWI4akJxMzU5cGVlM3pab1Y1NUV3eE43dVdxMW9TL3paK2k2clNuY2ZxS056SS85a0xveXhkYTZDZ25BQgpySXBFd2RNdHFXczhiMVkvMHErcy85NG93YzRqenlweFRXU3dYTHIxM3M3cDJaRE5hSjd4bDNBTDh6ekp3OU5LClk4NDlDSW1BR3U0Rk84SithNXpsVEVpb0tNS3ozTzZqVXZaWlhwUkxRZEE2QU9aS2Vpb3p3TlhxS0pjTnd4OGcKN2taVHlUeHFpbnEvT2VpR2llUHI5Y3dFbHFlOEtrd1A0cElWcFFJREFRQUJBb0lCQVFDRkd2Z1p2NHcvV2I3cApFMEozZWF6aHJFTHhPQ2NldmdCYmVPREVhTDdaZm96WWNVem1hZFNib2VLTGhwTHNadHNlM05QTUdHZ1RmbmhtCnJvM29Ia0lRQWxWVnRxbWp0WjBnamxwZC9TTHhkRk9nR3R1UkdlRnRrejFYMGZvaTJRQzNEWi9tWkswdVQzZ1gKYkhEMkNjTWk4YkNqNFZTV2tCUW1IeHpWL2pHcXJVVmtrRWJmWGpldWprWG1jdjlYLzhTajFVdm5JSUVDV0s0cgpIVDZsUm95NzV0V2ErL3dmY05xK29GRWZ3MjBWeVNTcCtmaDdPMHN3WG9oZFRkdFp4WVhhRUFYakU4cHRkb2RXCjgvc0tBc1BiamZubXFzdGtmcTNpV2wzdklta05jeXVPeEJ0YXIwbUw1N2tIZkVPU21UUUpFc05oWlFOcFNERkkKWU5kS1g3OUJBb0dCQU5JR2tVdWdDeFgxOXV6Y1JKY0dKdWhZaURkUStFdzduak41bnFqaE1SblYwMjNpRkxpdgozeGFVeEJVTWdEUVgvaUhZRU5DL1oxdUpacjRmL2dDTXk0WmU5TnFtKytwOVFFeXczenlLRGZYdTQ4NmxLZHl4ClhTTGlCWlE2a0drc0dCSExNcXd3elo5Ykw5NDJjaFZqcG5PM2tuL1JQdDZYNTd0OGlGMUpSRUtUQW9HQkFNb3QKR21ReUtzMDVwVk1Bdm1rN1pQRDNWMGZuMlBpZ094TFJVdUdSQThJaEdETlI3RkJXa2J1VXQxa25pVjExUkNKKwpMMyt4RkE5MzdGcHFxelN6S29kdTZtWFV3dEt4TUs1S2lQZHNLUXFVbU1pdU9ueXphUHBDK2l6OXFlNjVPVkE2ClcrOUJXazljbU41WHF5YUlVRGRDYjVTTEtPUGhrTUVmc2tEcDROSG5Bb0dBYlBtMWFDMEpzNEpsZGk4d2M4QmcKYmN5S0dWR3RGRGtXOUJTVjY0QzFMbmVSZGdHSnlPNlFiYklSTCs3RmtzSWtQY0ZUc0V5d1A0SEN5c0hrMUxvNQpYR1ptM0JFcXcxZnNCaDc4SmZob0dBUzFOV0xqbnJ4MDNBVzA2VjJkMHNSclZNZy9hYk1FN2p1dFVicWtaVTdJCmJtQ0E1a3RYT0w1UElpd1N3WHlqcTNzQ2dZRUFwRVpmc2xnOUJRSTQvaWVWa0NYdGtBbzV4amh4eVJ0UXhLcUgKSUxkWENXOGduZHFNSEg4cTdQTWF3M3RubHlQSW1BcFdCL2hYWjNZMit3Uy9WaFBhazY4aEVGci9ibmtCS0MxeAorekRNYkVkdm1XaFFKN0VUdEgybGo5Y1JNK01XMmNTQm5QZEtMVC85Q25UTG9ZU1RRVU5mTEtDaU9mKzNRZVRDClR4SjZWYk1DZ1lFQWk5WXAvTTRhYWV1QTZVaWVTTkxaY1FwUUx2VkxEdEhnOHV5ZUQ3eUlIL2w2S2ZUMVg1bU8KcXlEUHk1Wm9lYU5GTHJjcXVJcmpkMFo1SGNaUUYwRzMvZ1ZsZ3FsN0FHTWlXZ0syd1BtQ1Y2SGdHaVpFZHVhMgpPcVB4SnVaYlAzYUFOeTE5Vm41aytXOHdFODFOZWg4UmpPMkNDOGYwd2tTeEhhckJLREp4UjM4PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==
kind: Secret
metadata:
  namespace: default
  name: baetyl-tls-secret
type: kubernetes.io/tls`
)

var (
	commonCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: common-cm
  labels:
    abc: abc
data:
  example.property.1: hello
  example.property.2: world
  conf.yaml: |-
    property.1: value-1
    property.2: value-2
    property.3: value-3`
	commonCfgUpdate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: common-cm
  labels:
    abc: abc
data:
  example.property.1: hello1
  example.property.2: world1
  conf.yaml: |-
    property.1: value-11
    property.2: value-22
    property.3: value-33`

	objectTypeCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: object-cm
data:
  123.jpg: |-
    type: object
    source: baidubos
    account: current
    bucket: bie-document
    object: test-images/1.jpg`

	objectTypeCfgUpdate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: object-cm
data:
  123.jpg: |-
    type: object
    source: baidubos
    account: current
    bucket: bie-document
    object: test-images/123.jpg`

	bosCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: bos-cm
data:
  download1616485040271.zip: |-
    type: object
    source: baidubos
    account: other
    endpoint: http://bj-bos-sandbox.baidu-int.com
    bucket: baetyl-test1616485040271
    object: download1616485040271.zip
    unpack: zip
    ak: 618466d309734d1b908590ec2ee46932
    sk: 9891e1e22a34463ea1d3d5c673e2bfd0`

	httpCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: http-cm
data:
  tf_mnist.zip: |-
    type: object
    source: http
    url: https://doc.bce.baidu.com/bce-documentation/BIE/tf_mnist.zip
    object: tf_mnist.zip
    unpack: zip`

	imageCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: image-cm
  labels:
    baetyl-config-type: baetyl-image
data:
  address: nginx:latest`

	programCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: program-cm
  labels:
    baetyl-config-type: baetyl-program
data:
  linux-amd64: |-
    type: object
    source: baidubos
    account: current
    bucket: bie-document
    object: Easyedge-SDK/program1.zip
    platform: linux-amd64
    unpack: zip`

	functionCfg = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: function-cm
  labels:
    baetyl-config-type: baetyl-function
    baetyl-function: ""
data:
  cfc-new-demo.zip: |-
    type: function
    source: baidubos
    bucket: baetyl-cloud-1cd2d7790b6f4347bbeb3ecee54eca6e
    object: 9741a152f1282a514d72804927cc3f4c71d49aae86684195c2930c63c56f784b/cfc-new-demo.zip
    function: cfc-new-demo
    handler: index.add2
    runtime: python3
    version: 1
    unpack: zip`

	functionCfgUpdate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: function-cm
  labels:
    baetyl-config-type: baetyl-function
    baetyl-function: ""
data:
  cfc-new-demo.zip: |-
    type: function
    source: baidubos
    bucket: baetyl-cloud-1cd2d7790b6f4347bbeb3ecee54eca6e
    object: 9741a152f1282a514d72804927cc3f4c71d49aae86684195c2930c63c56f784b/cfc-new-demo.zip
    function: cfc-new-demo
    handler: index.sum
    runtime: python3
    version: 1
    unpack: zip`
)

func initYamlAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }

	yaml := router.Group("v1/yaml")
	{
		yaml.POST("", mockIM, common.Wrapper(api.CreateYamlResource))
		yaml.PUT("", mockIM, common.Wrapper(api.UpdateYamlResource))
		yaml.DELETE("", mockIM, common.Wrapper(api.DeleteYamlResource))
	}

	return api, router, mockCtl
}

func TestAPI_CreateSecret(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create secret
	expectSecret := &specV1.Secret{
		Name:      "dcell",
		Namespace: "default",
		Labels: map[string]string{
			"secret": "dcell",
		},
		Annotations: map[string]string{
			"secret": "dcell",
		},
		Data: map[string][]byte{
			"password": []byte("1f2d1e2e67df"),
			"username": []byte("admin"),
		},
	}
	sSecret.EXPECT().Get("default", "dcell", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateSecret("default", expectSecret).Return(expectSecret, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "secret.yaml")
	io.Copy(fw, strings.NewReader(testSecret))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create registry
	expectRegistry := &specV1.Secret{
		Name:      "myregistrykey",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"password": []byte("DOCKER_PASSWORD"),
			"username": []byte("DOCKER_USER"),
			"address":  []byte("DOCKER_REGISTRY_SERVER"),
		},
	}
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateSecret("default", expectRegistry).Return(expectRegistry, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "registry.yaml")
	io.Copy(fw, strings.NewReader(testRegistry))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create cert
	expectCert := &specV1.Secret{
		Name:      "baetyl-tls-secret",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAzEZqT7gDDN3bpRzfyjc2Wpm8SD49GMUOhZC3qm5y8K1njUYV\njATz1naaY5NHLJ+Ra7b/KKHecLRrugqrZrLDaxrXUR7s0VfUXrfTA+x+A/iGYp+A\neUbPs1ou7kh3bneZ6Sl9vgArlG53p9XXb2u5itMfUofLCqB0uCNdhQoPpOE5sxoy\nehOdqLmAar/BQtRO2z/eoJkw3CP2TokaGljjulBqswtQyZj3E7XWnoxbUfNMQJ4/\ncA/UpTO9PRbJIjahyq0IfZK6YjKKedIOu5JKuqgknxDjTGm/g03uNufLcjovqtcM\nhaLRFpgKNKHDsuJgcSakIBcHUfLBWTS97nVcSwIDAQABAoIBAQCxnnSoNo9oa6yl\nyBKLCTEY3F5LztpfJFE5MBmTdxA4V/8XZx4rLs865svd01zq3yCaMXdyRDVfeHHH\nRam91ah4A9G/No2ihaUiauJvof38w+dN7/TM/xZwk/UEujwm8JYCmFI+mHVMTjV6\n2UTHHDssC+FLaMnSxqxWxGF34cSyjpouV6GqNf7x0BkieoSbxFSwW0kvkxblvzbe\nALjfqmIq2nZZ68VHI4gpdxvR/SOD61QW75AMpPjeo+s8ynifwt5tO8uuNBFA66bK\n9jDQnsOzEtyive7P/X/A0NX0GMPUG8tAf1L0SeWA8ZkpCtlSNJfUKMb03H7chPT8\nAXe/Dv2BAoGBAOxw+SQXGCIs/sMRfxb6AKihHwCsPbxmXrO+6VE+QyTyo0aXulUm\n6R6q6CgFXRZ+D/cJXJRMFfbg3+weVkofGCxNkdyyl5UQ7Gs9UEwt9rY3/OCVyf77\nN+sCinF4gA7jisETRPhUQIAVMpJiChtoM0S40Uyvp30Ss3CePpq6Zs0LAoGBAN0s\nRRgiyY64flykditdFnnzhpmFsu76UGoFebP337PJNuLrjElPRvSIKes9sUEl7LQ5\ngj86N8YckqFbreeZVFlbVMlfiF8ihZZr9vs/i4XWxNZ+ve3kDfx+NZPbhrrjAm/I\nAC5r3MfMicedKysle5WtN2dTdACSlRUKMlLNdLXBAoGAXswc1996iZlIu6U0M14h\nDXssFvT3+6V/qsmMekpgWVv/IrqKtsFXDjivw/wCipYYiI90UvD+jXEz1lODfUxi\nQMEJDld8dwPGBmka3LBBDmX8OX9U8ap/jPYD0+LgVRfd9fNn37jH85KMKCyuqLZq\nd88x+3Eh0f/BehG4QAkkVmkCgYEAsNTqUUfO+usLLKrZSAZfKBXKsggxbdx4WqwS\nCA/QrX/dAETvntVhl7UePtTOWZYm0FmCh2bWnPDPU29ndVokFGVvPqnA8L87K28a\nwvqlZNa3/f7Lf93sSMnnsFW+PMGvewvfCT54AM7KAezpTM/lJWCegWA5yRNpWq8S\nWR2niIECgYA15TR6miyEumWnzx00TFM3d3SL/Kz1FJDRwPOjn9F//DUJgBtDQ+Al\njW911kzJof4A28UwAJ0XPo7yfLvlJhmiOSPa6/ckIlgA8KJEgixDLQ8y9B6jj1hB\nVfO8URR/qTdon17z/YGCOAS8sY46PJ/kbvGPbQZN8aYhaQHfbt+TIg==\n-----END RSA PRIVATE KEY-----\n"),
			"certificate":        []byte("-----BEGIN CERTIFICATE-----\nMIIDizCCAnOgAwIBAgIQXCsxM6Rzcq2URNQ8pH66JzANBgkqhkiG9w0BAQsFADAV\nMRMwEQYDVQQDEwprdWJlcm5ldGVzMB4XDTIxMTIyMzA2NDAyNloXDTIyMTIyMzA2\nNDAyNlowLTErMCkGA1UEAxMiYmFldHlsLXdlYmhvb2stc2VydmljZS5kZWZhdWx0\nLnN2YzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMxGak+4Awzd26Uc\n38o3NlqZvEg+PRjFDoWQt6pucvCtZ41GFYwE89Z2mmOTRyyfkWu2/yih3nC0a7oK\nq2ayw2sa11Ee7NFX1F630wPsfgP4hmKfgHlGz7NaLu5Id253mekpfb4AK5Rud6fV\n129ruYrTH1KHywqgdLgjXYUKD6ThObMaMnoTnai5gGq/wULUTts/3qCZMNwj9k6J\nGhpY47pQarMLUMmY9xO11p6MW1HzTECeP3AP1KUzvT0WySI2ocqtCH2SumIyinnS\nDruSSrqoJJ8Q40xpv4NN7jbny3I6L6rXDIWi0RaYCjShw7LiYHEmpCAXB1HywVk0\nve51XEsCAwEAAaOBvjCBuzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYB\nBQUHAwEwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBQltiuazoeynHeWIArVT/D+\nDxc5WjBlBgNVHREEXjBcghZiYWV0eWwtd2ViaG9vay1zZXJ2aWNlgh5iYWV0eWwt\nd2ViaG9vay1zZXJ2aWNlLmRlZmF1bHSCImJhZXR5bC13ZWJob29rLXNlcnZpY2Uu\nZGVmYXVsdC5zdmMwDQYJKoZIhvcNAQELBQADggEBAGnEL31pj+YvATksw/Oi+gAO\nykNnaXeggK0Kg1nofhW5yaj0G5OrTqrU0lbnRT10MH2+6ipoIHEO4IbI4gPJ4gbj\nYMobPu2Fx3tMWvI+SpK54BLOqefNU0BOWjpINUr+60iu8VZ48abuKgAcRbR6KbA2\nFZE7elgBGnrzeHu4yMtLvTB6UCs0Vg/F/FGUZ2ufu7lC9tUEsw4SvEHJzdFKFpM8\nwpfNANV9eD/WeX2Q4F9r/MfmWPWHKjIY9SG3usoUR9Jg8/fVEipgFTXCseJqDoTD\nOa7BzNuAgse9UFRk0sfFfrd3SxZEf4xYW90rkHSw+/G76/QKpXy3ZWAEMugfllk=\n-----END CERTIFICATE-----\n"),
			"signatureAlgorithm": []byte("SHA256-RSA"),
			"effectiveTime":      []byte("2021-12-23 06:40:26 +0000 UTC"),
			"expiredTime":        []byte("2022-12-23 06:40:26 +0000 UTC"),
			"serialNumber":       []byte("122513242306731234686214994664217295399"),
			"issuer":             []byte("kubernetes"),
			"fingerPrint":        []byte("A7:ED:24:3A:93:D0:66:43:E9:BB:E3:B4:A6:91:41:4B:61:10:BD:C3:DB:11:4D:7C:6B:A8:F9:1E:A8:F6:CE:92"),
		},
	}
	sSecret.EXPECT().Get("default", "baetyl-tls-secret", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateSecret("default", expectCert).Return(expectCert, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "cert.yaml")
	io.Copy(fw, strings.NewReader(testCert))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateSecret(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create secret
	se := &specV1.Secret{
		Name:      "dcell",
		Namespace: "default",
		Labels: map[string]string{
			"secret": "dcell",
		},
		Annotations: map[string]string{
			"secret": "dcell",
		},
		Data: map[string][]byte{
			"password": []byte("1f2d1e2e67df"),
			"username": []byte("admin"),
		},
	}
	se_updated := &specV1.Secret{
		Name:      "dcell",
		Namespace: "default",
		Labels: map[string]string{
			"secret": "dcell",
		},
		Annotations: map[string]string{
			"secret": "dcell",
		},
		Data: map[string][]byte{
			"password": []byte("1f2d1e2e67df"),
			"username": []byte("admin"),
			"123":      []byte("abc"),
		},
	}
	sSecret.EXPECT().Get("default", "dcell", "").Return(se, nil)
	sFacade.EXPECT().UpdateSecret("default", gomock.Any()).Return(se_updated, nil)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "secret.yaml")
	io.Copy(fw, strings.NewReader(updateSecret))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateRegistry(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create cert
	registry := &specV1.Secret{
		Name:      "myregistrykey",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"password": []byte("DOCKER_PASSWORD"),
			"username": []byte("DOCKER_USER"),
			"address":  []byte("DOCKER_REGISTRY_SERVER"),
		},
	}
	registryUpdated := &specV1.Secret{
		Name:      "myregistrykey",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"password": []byte("DOCKER_PASSWORD123"),
			"username": []byte("DOCKER_USER"),
			"address":  []byte("DOCKER_REGISTRY_SERVER"),
		},
	}
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(registry, nil)
	sFacade.EXPECT().UpdateSecret("default", gomock.Any()).Return(registryUpdated, nil)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "registry.yaml")
	io.Copy(fw, strings.NewReader(updateRegistry))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateCert(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create cert
	cert := &specV1.Secret{
		Name:      "baetyl-tls-secret",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAzEZqT7gDDN3bpRzfyjc2Wpm8SD49GMUOhZC3qm5y8K1njUYV\njATz1naaY5NHLJ+Ra7b/KKHecLRrugqrZrLDaxrXUR7s0VfUXrfTA+x+A/iGYp+A\neUbPs1ou7kh3bneZ6Sl9vgArlG53p9XXb2u5itMfUofLCqB0uCNdhQoPpOE5sxoy\nehOdqLmAar/BQtRO2z/eoJkw3CP2TokaGljjulBqswtQyZj3E7XWnoxbUfNMQJ4/\ncA/UpTO9PRbJIjahyq0IfZK6YjKKedIOu5JKuqgknxDjTGm/g03uNufLcjovqtcM\nhaLRFpgKNKHDsuJgcSakIBcHUfLBWTS97nVcSwIDAQABAoIBAQCxnnSoNo9oa6yl\nyBKLCTEY3F5LztpfJFE5MBmTdxA4V/8XZx4rLs865svd01zq3yCaMXdyRDVfeHHH\nRam91ah4A9G/No2ihaUiauJvof38w+dN7/TM/xZwk/UEujwm8JYCmFI+mHVMTjV6\n2UTHHDssC+FLaMnSxqxWxGF34cSyjpouV6GqNf7x0BkieoSbxFSwW0kvkxblvzbe\nALjfqmIq2nZZ68VHI4gpdxvR/SOD61QW75AMpPjeo+s8ynifwt5tO8uuNBFA66bK\n9jDQnsOzEtyive7P/X/A0NX0GMPUG8tAf1L0SeWA8ZkpCtlSNJfUKMb03H7chPT8\nAXe/Dv2BAoGBAOxw+SQXGCIs/sMRfxb6AKihHwCsPbxmXrO+6VE+QyTyo0aXulUm\n6R6q6CgFXRZ+D/cJXJRMFfbg3+weVkofGCxNkdyyl5UQ7Gs9UEwt9rY3/OCVyf77\nN+sCinF4gA7jisETRPhUQIAVMpJiChtoM0S40Uyvp30Ss3CePpq6Zs0LAoGBAN0s\nRRgiyY64flykditdFnnzhpmFsu76UGoFebP337PJNuLrjElPRvSIKes9sUEl7LQ5\ngj86N8YckqFbreeZVFlbVMlfiF8ihZZr9vs/i4XWxNZ+ve3kDfx+NZPbhrrjAm/I\nAC5r3MfMicedKysle5WtN2dTdACSlRUKMlLNdLXBAoGAXswc1996iZlIu6U0M14h\nDXssFvT3+6V/qsmMekpgWVv/IrqKtsFXDjivw/wCipYYiI90UvD+jXEz1lODfUxi\nQMEJDld8dwPGBmka3LBBDmX8OX9U8ap/jPYD0+LgVRfd9fNn37jH85KMKCyuqLZq\nd88x+3Eh0f/BehG4QAkkVmkCgYEAsNTqUUfO+usLLKrZSAZfKBXKsggxbdx4WqwS\nCA/QrX/dAETvntVhl7UePtTOWZYm0FmCh2bWnPDPU29ndVokFGVvPqnA8L87K28a\nwvqlZNa3/f7Lf93sSMnnsFW+PMGvewvfCT54AM7KAezpTM/lJWCegWA5yRNpWq8S\nWR2niIECgYA15TR6miyEumWnzx00TFM3d3SL/Kz1FJDRwPOjn9F//DUJgBtDQ+Al\njW911kzJof4A28UwAJ0XPo7yfLvlJhmiOSPa6/ckIlgA8KJEgixDLQ8y9B6jj1hB\nVfO8URR/qTdon17z/YGCOAS8sY46PJ/kbvGPbQZN8aYhaQHfbt+TIg==\n-----END RSA PRIVATE KEY-----\n"),
			"certificate":        []byte("-----BEGIN CERTIFICATE-----\nMIIDizCCAnOgAwIBAgIQXCsxM6Rzcq2URNQ8pH66JzANBgkqhkiG9w0BAQsFADAV\nMRMwEQYDVQQDEwprdWJlcm5ldGVzMB4XDTIxMTIyMzA2NDAyNloXDTIyMTIyMzA2\nNDAyNlowLTErMCkGA1UEAxMiYmFldHlsLXdlYmhvb2stc2VydmljZS5kZWZhdWx0\nLnN2YzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMxGak+4Awzd26Uc\n38o3NlqZvEg+PRjFDoWQt6pucvCtZ41GFYwE89Z2mmOTRyyfkWu2/yih3nC0a7oK\nq2ayw2sa11Ee7NFX1F630wPsfgP4hmKfgHlGz7NaLu5Id253mekpfb4AK5Rud6fV\n129ruYrTH1KHywqgdLgjXYUKD6ThObMaMnoTnai5gGq/wULUTts/3qCZMNwj9k6J\nGhpY47pQarMLUMmY9xO11p6MW1HzTECeP3AP1KUzvT0WySI2ocqtCH2SumIyinnS\nDruSSrqoJJ8Q40xpv4NN7jbny3I6L6rXDIWi0RaYCjShw7LiYHEmpCAXB1HywVk0\nve51XEsCAwEAAaOBvjCBuzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYB\nBQUHAwEwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBQltiuazoeynHeWIArVT/D+\nDxc5WjBlBgNVHREEXjBcghZiYWV0eWwtd2ViaG9vay1zZXJ2aWNlgh5iYWV0eWwt\nd2ViaG9vay1zZXJ2aWNlLmRlZmF1bHSCImJhZXR5bC13ZWJob29rLXNlcnZpY2Uu\nZGVmYXVsdC5zdmMwDQYJKoZIhvcNAQELBQADggEBAGnEL31pj+YvATksw/Oi+gAO\nykNnaXeggK0Kg1nofhW5yaj0G5OrTqrU0lbnRT10MH2+6ipoIHEO4IbI4gPJ4gbj\nYMobPu2Fx3tMWvI+SpK54BLOqefNU0BOWjpINUr+60iu8VZ48abuKgAcRbR6KbA2\nFZE7elgBGnrzeHu4yMtLvTB6UCs0Vg/F/FGUZ2ufu7lC9tUEsw4SvEHJzdFKFpM8\nwpfNANV9eD/WeX2Q4F9r/MfmWPWHKjIY9SG3usoUR9Jg8/fVEipgFTXCseJqDoTD\nOa7BzNuAgse9UFRk0sfFfrd3SxZEf4xYW90rkHSw+/G76/QKpXy3ZWAEMugfllk=\n-----END CERTIFICATE-----\n"),
			"signatureAlgorithm": []byte("SHA256-RSA"),
			"effectiveTime":      []byte("2021-12-23 06:40:26 +0000 UTC"),
			"expiredTime":        []byte("2022-12-23 06:40:26 +0000 UTC"),
			"serialNumber":       []byte("122513242306731234686214994664217295399"),
			"issuer":             []byte("kubernetes"),
			"fingerPrint":        []byte("A7:ED:24:3A:93:D0:66:43:E9:BB:E3:B4:A6:91:41:4B:61:10:BD:C3:DB:11:4D:7C:6B:A8:F9:1E:A8:F6:CE:92"),
		},
	}

	certUpdate := &specV1.Secret{
		Name:      "baetyl-tls-secret",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAzEZqT7gDDN3bpRzfyjc2Wpm8SD49GMUOhZC3qm5y8K1njUYV\njATz1naaY5NHLJ+Ra7b/KKHecLRrugqrZrLDaxrXUR7s0VfUXrfTA+x+A/iGYp+A\neUbPs1ou7kh3bneZ6Sl9vgArlG53p9XXb2u5itMfUofLCqB0uCNdhQoPpOE5sxoy\nehOdqLmAar/BQtRO2z/eoJkw3CP2TokaGljjulBqswtQyZj3E7XWnoxbUfNMQJ4/\ncA/UpTO9PRbJIjahyq0IfZK6YjKKedIOu5JKuqgknxDjTGm/g03uNufLcjovqtcM\nhaLRFpgKNKHDsuJgcSakIBcHUfLBWTS97nVcSwIDAQABAoIBAQCxnnSoNo9oa6yl\nyBKLCTEY3F5LztpfJFE5MBmTdxA4V/8XZx4rLs865svd01zq3yCaMXdyRDVfeHHH\nRam91ah4A9G/No2ihaUiauJvof38w+dN7/TM/xZwk/UEujwm8JYCmFI+mHVMTjV6\n2UTHHDssC+FLaMnSxqxWxGF34cSyjpouV6GqNf7x0BkieoSbxFSwW0kvkxblvzbe\nALjfqmIq2nZZ68VHI4gpdxvR/SOD61QW75AMpPjeo+s8ynifwt5tO8uuNBFA66bK\n9jDQnsOzEtyive7P/X/A0NX0GMPUG8tAf1L0SeWA8ZkpCtlSNJfUKMb03H7chPT8\nAXe/Dv2BAoGBAOxw+SQXGCIs/sMRfxb6AKihHwCsPbxmXrO+6VE+QyTyo0aXulUm\n6R6q6CgFXRZ+D/cJXJRMFfbg3+weVkofGCxNkdyyl5UQ7Gs9UEwt9rY3/OCVyf77\nN+sCinF4gA7jisETRPhUQIAVMpJiChtoM0S40Uyvp30Ss3CePpq6Zs0LAoGBAN0s\nRRgiyY64flykditdFnnzhpmFsu76UGoFebP337PJNuLrjElPRvSIKes9sUEl7LQ5\ngj86N8YckqFbreeZVFlbVMlfiF8ihZZr9vs/i4XWxNZ+ve3kDfx+NZPbhrrjAm/I\nAC5r3MfMicedKysle5WtN2dTdACSlRUKMlLNdLXBAoGAXswc1996iZlIu6U0M14h\nDXssFvT3+6V/qsmMekpgWVv/IrqKtsFXDjivw/wCipYYiI90UvD+jXEz1lODfUxi\nQMEJDld8dwPGBmka3LBBDmX8OX9U8ap/jPYD0+LgVRfd9fNn37jH85KMKCyuqLZq\nd88x+3Eh0f/BehG4QAkkVmkCgYEAsNTqUUfO+usLLKrZSAZfKBXKsggxbdx4WqwS\nCA/QrX/dAETvntVhl7UePtTOWZYm0FmCh2bWnPDPU29ndVokFGVvPqnA8L87K28a\nwvqlZNa3/f7Lf93sSMnnsFW+PMGvewvfCT54AM7KAezpTM/lJWCegWA5yRNpWq8S\nWR2niIECgYA15TR6miyEumWnzx00TFM3d3SL/Kz1FJDRwPOjn9F//DUJgBtDQ+Al\njW911kzJof4A28UwAJ0XPo7yfLvlJhmiOSPa6/ckIlgA8KJEgixDLQ8y9B6jj1hB\nVfO8URR/qTdon17z/YGCOAS8sY46PJ/kbvGPbQZN8aYhaQHfbt+TIg==\n-----END RSA PRIVATE KEY-----\n"),
			"certificate":        []byte("-----BEGIN CERTIFICATE-----\nMIIDizCCAnOgAwIBAgIQXCsxM6Rzcq2URNQ8pH66JzANBgkqhkiG9w0BAQsFADAV\nMRMwEQYDVQQDEwprdWJlcm5ldGVzMB4XDTIxMTIyMzA2NDAyNloXDTIyMTIyMzA2\nNDAyNlowLTErMCkGA1UEAxMiYmFldHlsLXdlYmhvb2stc2VydmljZS5kZWZhdWx0\nLnN2YzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMxGak+4Awzd26Uc\n38o3NlqZvEg+PRjFDoWQt6pucvCtZ41GFYwE89Z2mmOTRyyfkWu2/yih3nC0a7oK\nq2ayw2sa11Ee7NFX1F630wPsfgP4hmKfgHlGz7NaLu5Id253mekpfb4AK5Rud6fV\n129ruYrTH1KHywqgdLgjXYUKD6ThObMaMnoTnai5gGq/wULUTts/3qCZMNwj9k6J\nGhpY47pQarMLUMmY9xO11p6MW1HzTECeP3AP1KUzvT0WySI2ocqtCH2SumIyinnS\nDruSSrqoJJ8Q40xpv4NN7jbny3I6L6rXDIWi0RaYCjShw7LiYHEmpCAXB1HywVk0\nve51XEsCAwEAAaOBvjCBuzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYB\nBQUHAwEwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBQltiuazoeynHeWIArVT/D+\nDxc5WjBlBgNVHREEXjBcghZiYWV0eWwtd2ViaG9vay1zZXJ2aWNlgh5iYWV0eWwt\nd2ViaG9vay1zZXJ2aWNlLmRlZmF1bHSCImJhZXR5bC13ZWJob29rLXNlcnZpY2Uu\nZGVmYXVsdC5zdmMwDQYJKoZIhvcNAQELBQADggEBAGnEL31pj+YvATksw/Oi+gAO\nykNnaXeggK0Kg1nofhW5yaj0G5OrTqrU0lbnRT10MH2+6ipoIHEO4IbI4gPJ4gbj\nYMobPu2Fx3tMWvI+SpK54BLOqefNU0BOWjpINUr+60iu8VZ48abuKgAcRbR6KbA2\nFZE7elgBGnrzeHu4yMtLvTB6UCs0Vg/F/FGUZ2ufu7lC9tUEsw4SvEHJzdFKFpM8\nwpfNANV9eD/WeX2Q4F9r/MfmWPWHKjIY9SG3usoUR9Jg8/fVEipgFTXCseJqDoTD\nOa7BzNuAgse9UFRk0sfFfrd3SxZEf4xYW90rkHSw+/G76/QKpXy3ZWAEMugfllk=\n-----END CERTIFICATE-----\n"),
			"signatureAlgorithm": []byte("SHA256-RSA"),
			"effectiveTime":      []byte("2021-12-23 06:40:26 +0000 UTC"),
			"expiredTime":        []byte("2022-12-23 06:40:26 +0000 UTC"),
			"serialNumber":       []byte("122513242306731234686214994664217295399"),
			"issuer":             []byte("kubernetes"),
			"fingerPrint":        []byte("A7:ED:24:3A:93:D0:66:43:E9:BB:E3:B4:A6:91:41:4B:61:10:BD:C3:DB:11:4D:7C:6B:A8:F9:1E:A8:F6:CE:92"),
		},
	}
	sSecret.EXPECT().Get("default", "baetyl-tls-secret", "").Return(cert, nil).Times(1)
	sFacade.EXPECT().UpdateSecret("default", gomock.Any()).Return(certUpdate, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "cert.yaml")
	io.Copy(fw, strings.NewReader(updateCert))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_DeleteSecret(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Index = sIndex
	api.Facade = sFacade

	// good case: delete secret
	expectSecret := &specV1.Secret{
		Name:      "dcell",
		Namespace: "default",
		Labels: map[string]string{
			"secret": "dcell",
		},
		Annotations: map[string]string{
			"secret": "dcell",
		},
		Data: map[string][]byte{
			"password": []byte("1f2d1e2e67df"),
			"username": []byte("admin"),
		},
	}
	sSecret.EXPECT().Get("default", "dcell", "").Return(expectSecret, nil).Times(1)
	sIndex.EXPECT().ListAppIndexBySecret("default", "dcell").Return(nil, nil).Times(1)
	sFacade.EXPECT().DeleteSecret("default", "dcell").Return(nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "secret.yaml")
	io.Copy(fw, strings.NewReader(testSecret))
	w.Close()

	req, _ := http.NewRequest(http.MethodDelete, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: delete registry
	expectRegistry := &specV1.Secret{
		Name:      "myregistrykey",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretRegistry,
		},
		Data: map[string][]byte{
			"password": []byte("DOCKER_PASSWORD"),
			"username": []byte("DOCKER_USER"),
			"address":  []byte("DOCKER_REGISTRY_SERVER"),
		},
	}
	sSecret.EXPECT().Get("default", "myregistrykey", "").Return(expectRegistry, nil).Times(1)
	sIndex.EXPECT().ListAppIndexBySecret("default", "myregistrykey").Return(nil, nil).Times(1)
	sFacade.EXPECT().DeleteSecret("default", "myregistrykey").Return(nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "registry.yaml")
	io.Copy(fw, strings.NewReader(testRegistry))
	w.Close()

	req, _ = http.NewRequest(http.MethodDelete, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: delete cert
	expectCert := &specV1.Secret{
		Name:      "baetyl-tls-secret",
		Namespace: "default",
		Labels: map[string]string{
			specV1.SecretLabel: specV1.SecretCertificate,
		},
		Data: map[string][]byte{
			"key":                []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAzEZqT7gDDN3bpRzfyjc2Wpm8SD49GMUOhZC3qm5y8K1njUYV\njATz1naaY5NHLJ+Ra7b/KKHecLRrugqrZrLDaxrXUR7s0VfUXrfTA+x+A/iGYp+A\neUbPs1ou7kh3bneZ6Sl9vgArlG53p9XXb2u5itMfUofLCqB0uCNdhQoPpOE5sxoy\nehOdqLmAar/BQtRO2z/eoJkw3CP2TokaGljjulBqswtQyZj3E7XWnoxbUfNMQJ4/\ncA/UpTO9PRbJIjahyq0IfZK6YjKKedIOu5JKuqgknxDjTGm/g03uNufLcjovqtcM\nhaLRFpgKNKHDsuJgcSakIBcHUfLBWTS97nVcSwIDAQABAoIBAQCxnnSoNo9oa6yl\nyBKLCTEY3F5LztpfJFE5MBmTdxA4V/8XZx4rLs865svd01zq3yCaMXdyRDVfeHHH\nRam91ah4A9G/No2ihaUiauJvof38w+dN7/TM/xZwk/UEujwm8JYCmFI+mHVMTjV6\n2UTHHDssC+FLaMnSxqxWxGF34cSyjpouV6GqNf7x0BkieoSbxFSwW0kvkxblvzbe\nALjfqmIq2nZZ68VHI4gpdxvR/SOD61QW75AMpPjeo+s8ynifwt5tO8uuNBFA66bK\n9jDQnsOzEtyive7P/X/A0NX0GMPUG8tAf1L0SeWA8ZkpCtlSNJfUKMb03H7chPT8\nAXe/Dv2BAoGBAOxw+SQXGCIs/sMRfxb6AKihHwCsPbxmXrO+6VE+QyTyo0aXulUm\n6R6q6CgFXRZ+D/cJXJRMFfbg3+weVkofGCxNkdyyl5UQ7Gs9UEwt9rY3/OCVyf77\nN+sCinF4gA7jisETRPhUQIAVMpJiChtoM0S40Uyvp30Ss3CePpq6Zs0LAoGBAN0s\nRRgiyY64flykditdFnnzhpmFsu76UGoFebP337PJNuLrjElPRvSIKes9sUEl7LQ5\ngj86N8YckqFbreeZVFlbVMlfiF8ihZZr9vs/i4XWxNZ+ve3kDfx+NZPbhrrjAm/I\nAC5r3MfMicedKysle5WtN2dTdACSlRUKMlLNdLXBAoGAXswc1996iZlIu6U0M14h\nDXssFvT3+6V/qsmMekpgWVv/IrqKtsFXDjivw/wCipYYiI90UvD+jXEz1lODfUxi\nQMEJDld8dwPGBmka3LBBDmX8OX9U8ap/jPYD0+LgVRfd9fNn37jH85KMKCyuqLZq\nd88x+3Eh0f/BehG4QAkkVmkCgYEAsNTqUUfO+usLLKrZSAZfKBXKsggxbdx4WqwS\nCA/QrX/dAETvntVhl7UePtTOWZYm0FmCh2bWnPDPU29ndVokFGVvPqnA8L87K28a\nwvqlZNa3/f7Lf93sSMnnsFW+PMGvewvfCT54AM7KAezpTM/lJWCegWA5yRNpWq8S\nWR2niIECgYA15TR6miyEumWnzx00TFM3d3SL/Kz1FJDRwPOjn9F//DUJgBtDQ+Al\njW911kzJof4A28UwAJ0XPo7yfLvlJhmiOSPa6/ckIlgA8KJEgixDLQ8y9B6jj1hB\nVfO8URR/qTdon17z/YGCOAS8sY46PJ/kbvGPbQZN8aYhaQHfbt+TIg==\n-----END RSA PRIVATE KEY-----\n"),
			"certificate":        []byte("-----BEGIN CERTIFICATE-----\nMIIDizCCAnOgAwIBAgIQXCsxM6Rzcq2URNQ8pH66JzANBgkqhkiG9w0BAQsFADAV\nMRMwEQYDVQQDEwprdWJlcm5ldGVzMB4XDTIxMTIyMzA2NDAyNloXDTIyMTIyMzA2\nNDAyNlowLTErMCkGA1UEAxMiYmFldHlsLXdlYmhvb2stc2VydmljZS5kZWZhdWx0\nLnN2YzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMxGak+4Awzd26Uc\n38o3NlqZvEg+PRjFDoWQt6pucvCtZ41GFYwE89Z2mmOTRyyfkWu2/yih3nC0a7oK\nq2ayw2sa11Ee7NFX1F630wPsfgP4hmKfgHlGz7NaLu5Id253mekpfb4AK5Rud6fV\n129ruYrTH1KHywqgdLgjXYUKD6ThObMaMnoTnai5gGq/wULUTts/3qCZMNwj9k6J\nGhpY47pQarMLUMmY9xO11p6MW1HzTECeP3AP1KUzvT0WySI2ocqtCH2SumIyinnS\nDruSSrqoJJ8Q40xpv4NN7jbny3I6L6rXDIWi0RaYCjShw7LiYHEmpCAXB1HywVk0\nve51XEsCAwEAAaOBvjCBuzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYB\nBQUHAwEwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBQltiuazoeynHeWIArVT/D+\nDxc5WjBlBgNVHREEXjBcghZiYWV0eWwtd2ViaG9vay1zZXJ2aWNlgh5iYWV0eWwt\nd2ViaG9vay1zZXJ2aWNlLmRlZmF1bHSCImJhZXR5bC13ZWJob29rLXNlcnZpY2Uu\nZGVmYXVsdC5zdmMwDQYJKoZIhvcNAQELBQADggEBAGnEL31pj+YvATksw/Oi+gAO\nykNnaXeggK0Kg1nofhW5yaj0G5OrTqrU0lbnRT10MH2+6ipoIHEO4IbI4gPJ4gbj\nYMobPu2Fx3tMWvI+SpK54BLOqefNU0BOWjpINUr+60iu8VZ48abuKgAcRbR6KbA2\nFZE7elgBGnrzeHu4yMtLvTB6UCs0Vg/F/FGUZ2ufu7lC9tUEsw4SvEHJzdFKFpM8\nwpfNANV9eD/WeX2Q4F9r/MfmWPWHKjIY9SG3usoUR9Jg8/fVEipgFTXCseJqDoTD\nOa7BzNuAgse9UFRk0sfFfrd3SxZEf4xYW90rkHSw+/G76/QKpXy3ZWAEMugfllk=\n-----END CERTIFICATE-----\n"),
			"signatureAlgorithm": []byte("SHA256-RSA"),
			"effectiveTime":      []byte("2021-12-23 06:40:26 +0000 UTC"),
			"expiredTime":        []byte("2022-12-23 06:40:26 +0000 UTC"),
			"serialNumber":       []byte("122513242306731234686214994664217295399"),
			"issuer":             []byte("kubernetes"),
			"fingerPrint":        []byte("A7:ED:24:3A:93:D0:66:43:E9:BB:E3:B4:A6:91:41:4B:61:10:BD:C3:DB:11:4D:7C:6B:A8:F9:1E:A8:F6:CE:92"),
		},
	}
	sSecret.EXPECT().Get("default", "baetyl-tls-secret", "").Return(expectCert, nil).Times(1)
	sIndex.EXPECT().ListAppIndexBySecret("default", "baetyl-tls-secret").Return(nil, nil).Times(1)
	sFacade.EXPECT().DeleteSecret("default", "baetyl-tls-secret").Return(nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "cert.yaml")
	io.Copy(fw, strings.NewReader(testCert))
	w.Close()

	req, _ = http.NewRequest(http.MethodDelete, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_CreateConfig(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: create common kv config
	expectConfig := &specV1.Configuration{
		Name:      "common-cm",
		Namespace: "default",
		Labels: map[string]string{
			"abc": "abc",
		},
		Data: map[string]string{
			"example.property.1": "hello",
			"example.property.2": "world",
			"conf.yaml":          "property.1: value-1\nproperty.2: value-2\nproperty.3: value-3",
		},
	}
	sConfig.EXPECT().Get("default", "common-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", expectConfig).Return(expectConfig, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(commonCfg))
	w.Close()

	req, _ := http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create common object cfg bos current account
	objectCfg := &specV1.Configuration{
		Name:      "object-cm",
		Namespace: "default",
		Data: map[string]string{
			"_object_123.jpg": "{\"account\":\"current\",\"bucket\":\"bie-document\",\"object\":\"test-images/1.jpg\",\"source\":\"baidubos\",\"type\":\"object\",\"userID\":\"\"}",
		},
	}
	sConfig.EXPECT().Get("default", "object-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", objectCfg).Return(objectCfg, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(objectTypeCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create object cfg bos other account
	bosConfig := &specV1.Configuration{
		Name:      "bos-cm",
		Namespace: "default",
		Data: map[string]string{
			"_object_download1616485040271.zip": "{\"account\":\"other\",\"ak\":\"618466d309734d1b908590ec2ee46932\",\"bucket\":\"baetyl-test1616485040271\",\"endpoint\":\"http://bj-bos-sandbox.baidu-int.com\",\"object\":\"download1616485040271.zip\",\"sk\":\"9891e1e22a34463ea1d3d5c673e2bfd0\",\"source\":\"baidubos\",\"type\":\"object\",\"unpack\":\"zip\",\"userID\":\"\"}",
		},
	}
	sConfig.EXPECT().Get("default", "bos-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", bosConfig).Return(bosConfig, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(bosCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create object cfg http
	httpConfig := &specV1.Configuration{
		Name:      "http-cm",
		Namespace: "default",
		Data: map[string]string{
			"_object_tf_mnist.zip": "{\"object\":\"tf_mnist.zip\",\"source\":\"http\",\"type\":\"object\",\"unpack\":\"zip\",\"url\":\"https://doc.bce.baidu.com/bce-documentation/BIE/tf_mnist.zip\",\"userID\":\"\"}",
		},
	}
	sConfig.EXPECT().Get("default", "http-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", httpConfig).Return(httpConfig, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(httpCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create image cfg
	imageConfig := &specV1.Configuration{
		Name:      "image-cm",
		Namespace: "default",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-image",
		},
		Data: map[string]string{
			"address": "nginx:latest",
		},
	}
	sConfig.EXPECT().Get("default", "image-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", imageConfig).Return(imageConfig, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(imageCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create program cfg
	programConfig := &specV1.Configuration{
		Name:      "program-cm",
		Namespace: "default",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-program",
		},
		Data: map[string]string{
			"_object_linux-amd64": "{\"account\":\"current\",\"bucket\":\"bie-document\",\"object\":\"Easyedge-SDK/program1.zip\",\"platform\":\"linux-amd64\",\"source\":\"baidubos\",\"type\":\"object\",\"unpack\":\"zip\",\"userID\":\"\"}",
		},
	}
	sConfig.EXPECT().Get("default", "program-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", programConfig).Return(programConfig, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(programCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

	// good case: create function cfg type function
	functionConfig := &specV1.Configuration{
		Name:      "function-cm",
		Namespace: "default",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-function",
			"baetyl-function":    "",
		},
		Data: map[string]string{
			"_object_cfc-new-demo.zip": "{\"bucket\":\"baetyl-cloud-1cd2d7790b6f4347bbeb3ecee54eca6e\",\"function\":\"cfc-new-demo\",\"handler\":\"index.add2\",\"object\":\"9741a152f1282a514d72804927cc3f4c71d49aae86684195c2930c63c56f784b/cfc-new-demo.zip\",\"runtime\":\"python3\",\"source\":\"baidubos\",\"type\":\"function\",\"unpack\":\"zip\",\"userID\":\"\",\"version\":1}",
		},
	}
	sConfig.EXPECT().Get("default", "function-cm", "").Return(nil, nil).Times(1)
	sFacade.EXPECT().CreateConfig("default", functionConfig).Return(functionConfig, nil).Times(1)

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)
	fw, _ = w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(functionCfg))
	w.Close()

	req, _ = http.NewRequest(http.MethodPost, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re = httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateKVConfig(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: update common kv config
	expectConfig := &specV1.Configuration{
		Name:      "common-cm",
		Namespace: "default",
		Labels: map[string]string{
			"abc": "abc",
		},
		Data: map[string]string{
			"example.property.1": "hello",
			"example.property.2": "world",
			"conf.yaml":          "property.1: value-1\nproperty.2: value-2\nproperty.3: value-3",
		},
	}
	updateConfig := &specV1.Configuration{
		Name:      "common-cm",
		Namespace: "default",
		Labels: map[string]string{
			"abc": "abc",
		},
		Data: map[string]string{
			"example.property.1": "hello1",
			"example.property.2": "world1",
			"conf.yaml":          "property.1: value-11\nproperty.2: value-22\nproperty.3: value-33",
		},
	}
	sConfig.EXPECT().Get("default", "common-cm", "").Return(expectConfig, nil).Times(1)
	sFacade.EXPECT().UpdateConfig("default", gomock.Any()).Return(updateConfig, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(commonCfgUpdate))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)

}

func TestAPI_UpdateObjectConfig(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: update object config
	objectCfg := &specV1.Configuration{
		Name:      "object-cm",
		Namespace: "default",
		Data: map[string]string{
			"_object_123.jpg": "{\"account\":\"current\",\"bucket\":\"bie-document\",\"object\":\"test-images/1.jpg\",\"source\":\"baidubos\",\"type\":\"object\",\"userID\":\"\"}",
		},
	}
	objectCfgUpdate := &specV1.Configuration{
		Name:      "object-cm",
		Namespace: "default",
		Data: map[string]string{
			"_object_123.jpg": "{\"account\":\"current\",\"bucket\":\"bie-document\",\"object\":\"test-images/123.jpg\",\"source\":\"baidubos\",\"type\":\"object\",\"userID\":\"\"}",
		},
	}
	sConfig.EXPECT().Get("default", "object-cm", "").Return(objectCfg, nil).Times(1)
	sFacade.EXPECT().UpdateConfig("default", gomock.Any()).Return(objectCfgUpdate, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(objectTypeCfgUpdate))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_UpdateFunctionConfig(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Facade = sFacade

	// good case: update function cfg type function
	functionConfig := &specV1.Configuration{
		Name:      "function-cm",
		Namespace: "default",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-function",
			"baetyl-function":    "",
		},
		Data: map[string]string{
			"_object_cfc-new-demo.zip": "{\"bucket\":\"baetyl-cloud-1cd2d7790b6f4347bbeb3ecee54eca6e\",\"function\":\"cfc-new-demo\",\"handler\":\"index.add2\",\"object\":\"9741a152f1282a514d72804927cc3f4c71d49aae86684195c2930c63c56f784b/cfc-new-demo.zip\",\"runtime\":\"python3\",\"source\":\"baidubos\",\"type\":\"function\",\"unpack\":\"zip\",\"userID\":\"\",\"version\":1}",
		},
	}
	functionConfigUpdate := &specV1.Configuration{
		Name:      "function-cm",
		Namespace: "default",
		Labels: map[string]string{
			"baetyl-config-type": "baetyl-function",
			"baetyl-function":    "",
		},
		Data: map[string]string{
			"_object_cfc-new-demo.zip": "{\"bucket\":\"baetyl-cloud-1cd2d7790b6f4347bbeb3ecee54eca6e\",\"function\":\"cfc-new-demo\",\"handler\":\"index.sum\",\"object\":\"9741a152f1282a514d72804927cc3f4c71d49aae86684195c2930c63c56f784b/cfc-new-demo.zip\",\"runtime\":\"python3\",\"source\":\"baidubos\",\"type\":\"function\",\"unpack\":\"zip\",\"userID\":\"\",\"version\":1}",
		},
	}
	sConfig.EXPECT().Get("default", "function-cm", "").Return(functionConfig, nil).Times(1)
	sFacade.EXPECT().UpdateConfig("default", gomock.Any()).Return(functionConfigUpdate, nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(functionCfgUpdate))
	w.Close()

	req, _ := http.NewRequest(http.MethodPut, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}

func TestAPI_DeleteConfig(t *testing.T) {
	api, router, mockCtl := initYamlAPI(t)
	defer mockCtl.Finish()

	sApp := ms.NewMockApplicationService(mockCtl)
	sConfig := ms.NewMockConfigService(mockCtl)
	sSecret := ms.NewMockSecretService(mockCtl)
	sIndex := ms.NewMockIndexService(mockCtl)
	sFacade := mf.NewMockFacade(mockCtl)

	api.AppCombinedService = &service.AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}
	api.Index = sIndex
	api.Facade = sFacade

	// good case: delete kv config
	expectConfig := &specV1.Configuration{
		Name:      "common-cm",
		Namespace: "default",
		Labels: map[string]string{
			"abc": "abc",
		},
		Data: map[string]string{
			"example.property.1": "hello",
			"example.property.2": "world",
			"conf.yaml":          "property.1: value-1\nproperty.2: value-2\nproperty.3: value-3",
		},
	}
	sConfig.EXPECT().Get("default", "common-cm", "").Return(expectConfig, nil).Times(1)
	sIndex.EXPECT().ListAppIndexByConfig("default", "common-cm").Return(nil, nil).Times(1)
	sFacade.EXPECT().DeleteConfig("default", "common-cm").Return(nil).Times(1)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("file", "config.yaml")
	io.Copy(fw, strings.NewReader(commonCfg))
	w.Close()

	req, _ := http.NewRequest(http.MethodDelete, "/v1/yaml", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	re := httptest.NewRecorder()
	router.ServeHTTP(re, req)
	assert.Equal(t, http.StatusOK, re.Code)
}
