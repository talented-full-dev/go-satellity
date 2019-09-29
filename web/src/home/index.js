import style from './index.scss';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import React, {Component} from 'react';
import {Redirect} from 'react-router-dom';
import API from '../api/index.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
  }

  render() {
    return (
      <div>
        <h1 className={style.slogan}>
          {i18n.t('site.features')}
        </h1>
        <div className={style.features}>
          <div className={style.section}>
            <FontAwesomeIcon icon={['fa', 'chalkboard']} />
            <div className={style.desc}>
              {i18n.t('home.forum')}
            </div>
          </div>
        </div>
        <div>
        </div>
      </div>
    )
  }
}

export default Index;
