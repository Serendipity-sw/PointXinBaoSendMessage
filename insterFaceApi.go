package main

import (
	"net/url"
	"net/http"
	"github.com/smtc/glog"
	"io/ioutil"
	"crypto/tls"
	"bytes"
	"encoding/json"
	"net"
	"time"
)
/**
接口请求返回结构体
 */
type PointXinBaoApiData struct  {
	State string //接口请求状态 success表示成功 其他为失败
	Msg string  //接口请求获取的url地址
}



func init() {
	//deferinit.AddInit(zxGetMyUserJumpUrl,nil,999)
}

/**
获取用户短信内容
创建人:邵炜
创建时间:2016年5月27日11:03:51
输入参数:用户号码 价格 省份 渠道
 */
func getMyUserJumpUrl(mob, price, province, channel string) (*PointXinBaoApiData,error) {
	dataValue:=url.Values{}
	dataValue.Add("phone",mob)
	dataValue.Add("price",price)
	dataValue.Add("province",province)
	dataValue.Add("channel",channel)
	dataValueJson,err:=json.Marshal(dataValue)
	if err != nil {
		glog.Error("getMyUserJumpUrl data Marshal error! error: %s \n",err.Error())
		return nil,err
	}
	resp,err:=http.NewRequest("POST",addflow,bytes.NewReader(dataValueJson))
	if err != nil {
		glog.Error("getMyUserJumpUrl NewRequest error! error: %s sendUrl: %s dataValue: %s \n",err.Error(),addflow,string(dataValueJson))
		return nil,err
	}
	resp.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	httpClient.Transport=&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	respone, err := httpClient.Do(resp)
	if err != nil {
	glog.Error("getMyUserJumpUrl send http post error! err: %s sendUrl: %s dataValue: %s \n",err.Error(),addflow,string(dataValueJson))
		return nil,err
	}
	requestData,err:=ioutil.ReadAll(respone.Body)
	defer respone.Body.Close()
	if err != nil {
		glog.Error("getMyUserJumpUrl send http read request data error! error: %s sendUrl: %s dataValue: %s \n",err.Error(),addflow,string(dataValueJson))
		return nil,err
	}
	var model PointXinBaoApiData
	err=json.Unmarshal(requestData,&model)
	if err != nil {
		glog.Error("getMyUserJumpUrl request data unmarshal error! err: %s sendUrl: %s dataValue: %v \n",err.Error(),addflow,string(dataValueJson))
		return nil,err
	}
	return &model,nil
}
/**
短信发送接口
创建人:邵炜
创建时间:2016年5月27日10:51:28
输入参数:手机号码 短信发送的信息 渠道 备注
 */
func sendMessage(mob, msg, channel ,remark string)error {
	dataValues:=url.Values{}
	dataValues.Add("mob",mob)
	dataValues.Add("msg",msg)
	dataValues.Add("channel",channel)
	dataValues.Add("remark",remark)
	dataValues.Add("callPswd",smspassword)
	resp,err:=http.PostForm(sendMessageUrl,dataValues)
	if err != nil {
		glog.Error("sendMessage send message error! err: %s sendMessageUrl: %s dataValues: %v \n",err.Error(),sendMessageUrl,dataValues)
		return err
	}
	requestData,err:=ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		glog.Error("sendMessage read requestData error! err: %s sendMessageUrl: %s dataValues: %v \n",err.Error(),sendMessageUrl,dataValues)
		return err
	}
	glog.Info("sendMessage success! dataValues: %v messageApiData: %s \n",string(requestData))
	return nil
}
/**
追销用户短信内容
创建人:邵炜
创建时间:2016年5月27日11:03:51
输入参数:用户号码
 */
func zxGetMyUserJumpUrl(mob string) (*PointXinBaoApiData,error) {
	dataValues:=url.Values{}
	dataValues.Add("phone",mob)
	resp,err:=http.PostForm(doublekill,dataValues)
	if err != nil {
		glog.Error("zxGetMyUserJumpUrl send message error! err: %s sendMessageUrl: %s dataValues: %v \n",err.Error(),doublekill,dataValues)
		return nil,err
	}
	requestData,err:=ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		glog.Error("zxGetMyUserJumpUrl read requestData error! err: %s sendMessageUrl: %s dataValues: %v \n",err.Error(),doublekill,dataValues)
		return nil,err
	}
	var model PointXinBaoApiData
	err=json.Unmarshal(requestData,&model)
	if err != nil {
		glog.Error("getMyUserJumpUrl request data unmarshal error! err: %s sendUrl: %s dataValue: %v \n",err.Error(),doublekill,dataValues)
		return nil,err
	}
	glog.Info("getMyUserJumpUrl success! dataValues: %v messageApiData: %s \n",dataValues,string(requestData))
	return &model,nil
}