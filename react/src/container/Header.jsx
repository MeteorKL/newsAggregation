import React from 'react'
import { Link } from 'react-router'
import { hashHistory } from 'react-router'
import { connect } from 'react-redux'
import { browserHistory } from 'react-router'
import crypto from 'crypto'

import './Header.scss'
import Dialog from '../components/Dialog.jsx'
import CheckBox from '../components/CheckBox.jsx'
import ajax from '../utils/ajax.js'

class Header extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      login: false,
      nickname: "",
      tags: [],
      userlayer: false,
      checkboxOption:[
        {
          text:'多选框-1',
          value: 'checkbox-1',
          checked: true
        },
        {
          text:'多选框-2',
          value: 'checkbox-2',
          checked: false
        },
        {
          text:'多选框-3',
          value: 'checkbox-3',
          checked: false
        }
      ]
    }
  }

  handleChangeCheckBox(index){
        let _checkboxOption = this.state.checkboxOption;

        if(_checkboxOption[index].checked){
            _checkboxOption[index].checked = false;
        }else{
            _checkboxOption[index].checked = true;
        }

        this.setState({checkboxOption: _checkboxOption})
    }

  componentDidMount() {
    ajax.Get("/api/checkLogin").then((r) => {
      if (r.status==0) {
        this.setState({login: true, nickname: r.data.nickname, tags: r.data.tags, userlayer: false});
      }
    }).catch((error) => {
      console.error(error);
    });
  }

  login(callback) {
    var nickname = document.getElementById("nickname").value;
    var password = document.getElementById("password").value;
    if (nickname == "") {
      callback({status: -1, message: "请输入昵称"});
      return
    }
    if (password == "") {
      callback({status: -1, message: "请输入密码"});
      return
    }
    password = crypto.createHash('md5').update(password).digest('hex');
    ajax.Get("/api/login?nickname="+nickname+"&password="+password).then((r) => {
      if (r.status==0) {
        console.log(r)
        this.setState({login: true, nickname: r.data.nickname, tags: r.data.tags, userlayer: false});
      } else {
        callback(r);
      }
    }).catch((error) => {
      console.error(error);
    });
  }

  register(callback) {
    var mail = document.getElementById("mail").value;
    var nickname = document.getElementById("nickname").value;
    var password = document.getElementById("password").value;
    var rePassword = document.getElementById("rePassword").value;

    var regex = /^([0-9A-Za-z\-_\.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g;
    if (!regex.test(mail)) {
      callback({status: -1, message: "邮箱的格式不合法"});
      return
    }
    if (nickname.length < 6) {
      callback({status: -1, message: "昵称长度小于6个字符"});
      return
    }
    if (password.length < 6) {
      callback({status: -1, message: "密码长度小于6个字符"});
      return
    }
    if (password != rePassword) {
      callback({status: -1, message: "两次输入的密码不一致"});
      return
    }
    password = crypto.createHash('md5').update(password).digest('hex');
    console.log(password);
    ajax.Get("/api/register?mail="+mail+"&nickname="+nickname+"&password="+password).then((r) => {
      callback(r);
    }).catch((error) => {
      console.error(error);
    });
  }
  logout() {
    ajax.Get("/api/logout").then((r) => {
      if (r.status==0) {
        this.setState({login: false, nickname: "", tags: []});
      }
    }).catch((error) => {
      console.error(error);
    });
  }
  updateTags(callback) {
    var param = ""
    var doms = document.getElementsByName("Tags");
    for (var i=0;i < doms.length;i++) {
      if (doms[i].checked){
        param += "&tags="+doms[i].value;
      }
    }
    ajax.Get("/api/updateTags?"+param.slice(1)).then((r) => {
      callback(r);
      this.setState({tags: r.data});
    }).catch((error) => {
      console.error(error);
    });
  }
  render() {
    return (
      <div className="Header">
        <div className="topbar">
        <div className="left clearfix">
          <a href="/download" >下载APP</a>
        </div>
        <div className="right">
          <ul className="item clearfix">
            {this.state.login?
            <li className='userbox'>
              <a id='userhead' className='userhead bold' onClick={()=>{
                this.setState({userlayer: !this.state.userlayer})
              }}>{this.state.nickname}</a>
              <div id='userlayer' className='userlayer' style={{"display": this.state.userlayer?"block":"none"}}>
                <ul>
                  <li>
                    <a id='triggerTags' className='layeritem'>我的关注</a>
                    <Dialog triggerID='triggerTags' title='我的关注' func={this.updateTags.bind(this)}>
                      {/*<CheckBox options={this.state.checkboxOption} onChange={this.handleChangeCheckBox.bind(this)}/>*/}
                      <label><input name="Tags" type="checkbox" value="国际" defaultChecked={this.state.tags.indexOf("国际")>-1?true:false} />国际 </label> 
                      <label><input name="Tags" type="checkbox" value="国内" defaultChecked={this.state.tags.indexOf("国内")>-1?true:false} />国内 </label> 
                      <label><input name="Tags" type="checkbox" value="军事" defaultChecked={this.state.tags.indexOf("军事")>-1?true:false} />军事 </label> 
                      <label><input name="Tags" type="checkbox" value="财经" defaultChecked={this.state.tags.indexOf("财经")>-1?true:false} />财经 </label> 
                      <label><input name="Tags" type="checkbox" value="科技" defaultChecked={this.state.tags.indexOf("科技")>-1?true:false} />科技 </label> 
                      <label><input name="Tags" type="checkbox" value="体育" defaultChecked={this.state.tags.indexOf("体育")>-1?true:false} />体育 </label> 
                      <label><input name="Tags" type="checkbox" value="娱乐" defaultChecked={this.state.tags.indexOf("娱乐")>-1?true:false} />娱乐 </label> 
                    </Dialog>
                  </li>
                  <li>
                    <a className='layeritem' onClick={()=>{this.logout()}}>退出</a>
                  </li>
                </ul>
              </div>
            </li>
            :
            <li>
              <a id='triggerRegister'>注册</a>
              <Dialog triggerID='triggerRegister' title='注册' func={this.register.bind(this)}>
                <input id='mail' className='input' placeholder='请输入邮箱'/>
                <input id='nickname' className='input' placeholder='昵称长度要求大于6个字符'/>
                <input id='password' className='input' type='password' placeholder='密码长度要求大于6个字符'/>
                <input id='rePassword' className='input' type='password' placeholder='再次输入密码'/>
              </Dialog>
            </li>
            }
            {this.state.login?<li/>:
            <li className="bold">
              <a id='triggerLogin'>登录</a>
              <Dialog triggerID='triggerLogin' title='登录' func={this.login.bind(this)}>
                <p className='label'>账号</p>
                <input id='nickname' className='input' placeholder='请输入昵称'/>
                <p className='label'>密码</p>
                <input id='password' className='input' type='password' placeholder='请输入密码'/>
              </Dialog>
            </li>
            }
            <li>
              <a id='triggerAbout'>关于</a>
              <Dialog triggerID='triggerAbout' title='关于' buttonDisable={true}>
                <p className='label'>TopNews</p>
                <p className='label'>作者：徐亮</p>
                <p className='label'>学号：3140102431</p>
              </Dialog>
            </li>
          </ul>
        </div>
        </div>
      </div>
    )
  }
}

export default Header