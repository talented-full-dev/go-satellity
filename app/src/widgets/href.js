import style from './href.module.scss';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';

class Href extends Component {
  render() {
    let classes = this.props.class.split(' ').map((name) => {
      return style[[name]]
    });
    if (this.props.original) {
      return (
        <a href={this.props.action} className={classes}>{this.props.text}</a>
      );
    }
    return (
      <Link to={this.props.action} className={classes}>{this.props.text}</Link>
    );
  }
}

export default Href;
