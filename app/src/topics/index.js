import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import Config from '../components/config.js';
import Base64 from '../components/base64.js';
import API from '../api/index.js';
import TopicItem from './item.js';
import Loading from '../widgets/loading.js';
import Widget from '../home/widget.js';

class Index extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.base64 = new Base64();
    this.pagination = 50;
    // TODO decode should categories in api;
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(this.base64.decode(d));
    }
    this.state = {
      id: props.match.params.id || 'latest',
      topics: [],
      categories: categories,
      category: {},
      category_id: 'latest',
      loading: true,
      offset: ''
    };

    this.loadTopics = this.loadTopics.bind(this);
  }

  componentDidMount() {
    this.api.category.index().then((resp) => {
      if (resp.error) {
        return
      }
      let category_id = 'latest';
      let current = {};
      let category = this.state.id;
      if (category !== 'latest') {
        for (let i = 0; i < resp.data.length; i++) {
          let c = resp.data[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            category_id = c.category_id;
            current = c;
            break;
          }
        }
      }
      this.setState({category: current, category_id: category_id, categories: resp.data}, () => {
        this.fetchTopics(category_id, true);
      });
    });
  }

  componentDidUpdate(prevProps) {
    if (this.props.match.params.id !== prevProps.match.params.id) {
      let category_id = 'latest';
      let current = {};
      let category = this.props.match.params.id || 'latest';
      if (category !== 'latest') {
        for (let i = 0; i < this.state.categories.length; i++) {
          let c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            category_id = c.category_id;
            current = c;
            break;
          }
        }
      }
      this.setState({category: current, category_id: category_id}, () => {
        this.fetchTopics(category_id, true);
      });
    }
  }

  loadTopics(e) {
    e.preventDefault();
    this.fetchTopics(this.state.category_id, false);
  }

  fetchTopics(category_id, replace) {
    this.setState({loading: replace, offset: ''});
    let request = category_id === 'latest' ? this.api.topic.index(this.state.offset) : this.api.category.topics(category_id, this.state.offset);

    request.then((resp) => {
      if (resp.error) {
        return
      }
      let data = resp.data;
      let offset = data.length > this.pagination ? data[data.length-1].updated_at : '' ;
      if (!replace) {
        data = this.state.topics.concat(data);
      }
      this.setState({category_id: category_id, loading: false, offset: offset, topics: data});
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;

    const loadingView = (
      <div className={style.loading}>
        <Loading />
      </div>
    )

    const topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    const categories = state.categories.map((category) => {
      return (
        <Link to={`/categories/${category.name}`} className={`${style.node} ${state.category_id === category.category_id ? style.current : ''}`} key={category.category_id} >
            {category.alias}
        </Link>
      )
    });

    let seoView;
    if (!!state.category.name) {
      seoView = (
        <Helmet>
          <title>{`${state.category.alias} - ${Config.Name}`}</title>
          <meta name='description' content={state.category.description} />
        </Helmet>
      )
    }

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='column main'>
          <div className={style.nodes}>
            <Link to='/'
              className={`${style.node} ${state.category_id === 'latest' ? style.current : ''}`}>{i18n.t('home.latest')}</Link>
            {categories}
          </div>
          {state.loading && loadingView}
          {!state.loading && <ul className={style.topics}> {topics} </ul>}
          {/* TODO i18n */}
          {state.offset !== '' && <div className={style.load}><span onClick={this.loadTopics}>Load More</span></div>}
        </main>
        <aside className='column aside'>
          <Widget />
        </aside>
      </div>
    );
  }
}

export default Index;
