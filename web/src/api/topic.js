class Topic {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }

  index(offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B')
    }
    return this.api.axios.get(`/topics?offset=${offset}`).then((resp) => {
      return resp.data;
    });
  }

  create(params) {
    const data = {title: params.title, body: params.body, category_id: params.category_id, draft: params.draft};
    return this.api.axios.post('/topics', data).then((resp) => {
      return resp.data;
    });
  }

  update(id, params) {
    const data = {title: params.title, body: params.body, category_id: params.category_id, draft: params.draft};
    return this.api.axios.post(`/topics/${id}`, data).then((resp) => {
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/topics/${id}`).then((resp) => {
      return resp.data;
    });
  }

  action(action, id) {
    if (action !== 'like' || action !== 'unlike' || action !== 'bookmark' || action !== 'unsave') {
      return this.api.axios.post(`/topics/${id}/${action}`, {}).then((resp) => {
        return resp.data;
      });
    }
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/topics').then((resp) => {
      return resp.data;
    })
  }
}

export default Topic;
