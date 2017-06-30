import React from 'react'
import { Link } from 'react-router'
import { hashHistory } from 'react-router'
import { connect } from 'react-redux'
import { browserHistory } from 'react-router'

import './Body.scss'
import ajax from '../utils/ajax.js'

class Body extends React.Component {
  constructor(props) {
    super(props);
    this.tagIndex = {
      "推荐": 0,
      "国内": 1,
      "国际": 2,
      "军事": 3,
      "财经": 4,
      "科技": 5,
      "体育": 6,
    }
    this.tags = ["推荐", "国内", "国际", "军事", "财经", "科技", "体育"];
    var news = new Array(7);
    var time = Date.parse(new Date())/1000;
    for (var i=0;i<news.length;i++) {
      news[i] = {
        bottomTime: time,
        topTime: time,
        content: [],
      }
    }
    this.state = {
      tag: "推荐",
      news: news,
    }
  }

  loadNews(tag, top) {
    var param = "tag="+tag;
    // 当目前没有新闻的时候，加载历史新闻
    if (this.state.news[this.tagIndex[tag]].content.length==0) {
      top = false;
    }
    if (top) {
      var topTime = this.state.news[this.tagIndex[tag]].topTime;
      param += "&topTime="+topTime;
    } else {
      var bottomTime = this.state.news[this.tagIndex[tag]].bottomTime;
      param += "&bottomTime="+bottomTime;
    }
    ajax.Get("/api/feed?"+param, (r) => {
      console.log(r);
      if (r.status==0) {
        if (top) {
          console.log("topTime update");
          this.state.news.unshift.apply(this.state.news[this.tagIndex[tag]].content, r.data.content);
          this.state.news[this.tagIndex[tag]].topTime = r.data.newTime;
        } else {
          console.log("bottomTime update");
          this.state.news.push.apply(this.state.news[this.tagIndex[tag]].content, r.data.content);
          this.state.news[this.tagIndex[tag]].bottomTime = r.data.newTime;
        }
        console.log(tag, this.state.news[this.tagIndex[tag]].topTime, this.state.news[this.tagIndex[tag]].bottomTime);
      }
      this.setState({news: this.state.news, tag: tag});
    }, (error) => {
      console.error(error);
    });
  }

  componentDidMount() {
    this.loadNews(this.state.tag, false);
    document.addEventListener("scroll", () => {
      if (document.body.scrollTop + document.documentElement.clientHeight > document.body.scrollHeight - 800) {
        this.loadNews(this.state.tag, false);
      }
    })
  }

  formatMsgTime (timespan) {
    var dateTime = new Date(timespan);

    var year = dateTime.getFullYear();
    var month = dateTime.getMonth() + 1;
    var day = dateTime.getDate();
    var hour = dateTime.getHours();
    var minute = dateTime.getMinutes();
    var second = dateTime.getSeconds();
    var now = new Date();
    var now_new = Date.parse(now.toDateString());  //typescript转换写法

    var milliseconds = 0;
    var timeSpanStr;

    milliseconds = now_new - timespan;

    if (milliseconds <= 1000 * 60 * 1) {
      timeSpanStr = '刚刚';
    }
    else if (1000 * 60 * 1 < milliseconds && milliseconds <= 1000 * 60 * 60) {
      timeSpanStr = Math.round((milliseconds / (1000 * 60))) + '分钟前';
    }
    else if (1000 * 60 * 60 * 1 < milliseconds && milliseconds <= 1000 * 60 * 60 * 24) {
      timeSpanStr = Math.round(milliseconds / (1000 * 60 * 60)) + '小时前';
    }
    else if (1000 * 60 * 60 * 24 < milliseconds && milliseconds <= 1000 * 60 * 60 * 24 * 15) {
      timeSpanStr = Math.round(milliseconds / (1000 * 60 * 60 * 24)) + '天前';
    }
    else if (milliseconds > 1000 * 60 * 60 * 24 * 15 && year == now.getFullYear()) {
      timeSpanStr = month + '-' + day + ' ' + hour + ':' + minute;
    } else {
      timeSpanStr = year + '-' + month + '-' + day + ' ' + hour + ':' + minute;
    }
    return timeSpanStr;
  }

  render() {
    return (
      <div className="container">
        <div className="left channel">
          {this.tags.map((tag, index)=>{
            return (
              <li onClick={()=>{this.loadNews(tag, true);}} key={index}> 
                <a className={tag==this.state.tag?"wchannel-item active":"wchannel-item"}> 
                  <span>{tag}</span>
                </a> 
              </li>
            )
          })}
        </div>
    
    <div className="left content">
      <div className="bui-box slide"></div>
      
      <div className="feed-infinite-wrapper">
        <ul>

          {this.state.news[this.tagIndex[this.state.tag]].content.map((news, index)=>{
            var time = this.formatMsgTime(news.time*1000)
            return (
              
          <li key={index}>
            <div className="bui-box single-mode">
              <div className="bui-left single-mode-lbox">
                <a href={news.docurl} target="_blank" className="img-wrap">
                  <img className="lazy-load-img" src={news.imgurl}/>
                </a>
              </div> 
              <div className="single-mode-rbox">
                <div className="single-mode-rbox-inner">
                  <div className="title-box">
                    <a href={news.docurl} target="_blank" className="link">{news.title}</a>
                  </div>
                  <div className="bui-box footer-bar">
                    <div className="bui-left footer-bar-left">
                      <a onClick={()=>{this.loadNews(news.tag, true);}} target="_blank" className="footer-bar-action tag tag-style-other">{news.tag}</a> 
                      <a className="footer-bar-action source">{news.source}</a> 
                      <span className="footer-bar-action">{time}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </li>
            )
          })}
        </ul>
      </div>
    </div>
  </div>
    )
  }
}

export default Body