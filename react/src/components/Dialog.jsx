import React from 'react'
import { Link } from 'react-router'
import { browserHistory } from 'react-router'

import './Dialog.scss'

export default class Dialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hidden: true,
      alerm: false,
      hideForm: false,
      alermContent: "",
    }
  }

  componentDidMount() {
    var triggerButton = document.getElementById(this.props.triggerID);
    triggerButton.addEventListener('click', ()=>{
      this.setState({hidden: false});
    });
  }

  render() {
    return (
      <div>
      {this.state.hidden?<div/>:
      <div>
        {this.state.hideForm?<div/>:
        <div className='MaskLayer' onClick={()=>{this.setState({hidden: !this.state.hidden});}}></div>
        }
        {!this.state.alerm?<div/>:
        <div className='Alerm'>{this.state.alermContent}</div>
        }
        {this.state.hideForm?<div/>:
        <div className='Form'>
          <div className='header'>
            <span className='title'>{this.props.title}</span>
            <a className="headerbtn" onClick={()=>{this.setState({hidden: !this.state.hidden});}}>x</a>
          </div>
          <div className='body'>
            {this.props.children}
          </div>
          {this.props.buttonDisable?<div/>:
          <div className='footer'>
            <button onClick={()=>{this.setState({hidden: !this.state.hidden});}}>取消</button>
            <button className='ok' onClick={()=>{this.props.func(
              (data) => {
                this.setState({alerm: true, alermContent: data.message});
                setTimeout(() => {
                  this.setState({alerm: false});
                }, 3000);
                if (data.status==0) {
                  this.setState({hideForm: true});
                  setTimeout(() => {
                    this.setState({hideForm: false});
                    this.setState({hidden: true});
                  }, 3000);
                }
              }
              )}}>确认</button>
          </div>
          }
        </div>
        }
      </div>
      }
      </div>
    )
  }
}
