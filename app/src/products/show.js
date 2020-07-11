import style from './show.module.scss';
import React, {Component} from 'react';
import showdown from 'showdown';
import showdownHighlight from 'showdown-highlight';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import SiteWidget from '../home/widget.js';
import Loading from '../components/loading.js';

export default class Show extends Component {
  constructor(props) {
    super(props);

    this.api = window.api;
    this.converter = new showdown.Converter({ extensions: ['header-anchors', showdownHighlight] });

    this.state = {
      product_id: props.match.params.id,
      loading: true,
      tags: [],
    };
  }

  componentDidMount() {
    this.api.product.show(this.state.product_id).then((resp) => {
      if (resp.error) {
        return
      }

      let data = resp.data;
      data.loading = false;
      data.html_body = this.converter.makeHtml(data.body);
      this.setState(data);
    });
  }

  render() {
    const state = this.state;
    const loadingView = (
      <div className={style.loading}>
        <Loading class='medium' />
      </div>
    );

    const productView = (
      <div className={style.product}>
        <div className={style.cover} style={{backgroundImage: `url(${state.cover_url})`}} />
        <div className={style.content}>
          <h1>
            {state.name}
          </h1>
          <div>
            {state.body !== '' && <article className={`md ${style.body}`} dangerouslySetInnerHTML={{__html: state.html_body}} />}
          </div>
          {
            state.source !== '' &&
            <div>
              Address: <a href={state.source} rel='nofollow noopener noreferrer' target='_blank'>{state.source}</a>
            </div>
          }
          <div className={style.tags}>
            <FontAwesomeIcon className={style.icon} icon={['fas', 'tags']} />
            {state.tags.join(', ')}
          </div>
        </div>
      </div>
    );

    return (
      <div className='container'>
        <main className='column main'>
          {state.loading && loadingView}
          {!state.loading && productView}
        </main>
        <aside className='column aside'>
          <SiteWidget />
        </aside>
      </div>
    )
  }
}
