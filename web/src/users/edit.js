import style from './index.scss';
import 'react-image-crop/lib/ReactCrop.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import ReactCrop from 'react-image-crop';
import API from '../api/index.js';
import Button from '../widgets/button.js';

class Edit extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    const user = this.api.user.local();
    this.state = {
      avatar_url: user.avatar_url,
      nickname: user.nickname,
      biography: user.biography,
      submitting: false
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleClick = this.handleClick.bind(this);
    this.handleFileChange = this.handleFileChange.bind(this);
    this.onImageLoaded = this.onImageLoaded.bind(this);
    this.onCropComplete = this.onCropComplete.bind(this);
    this.onCropChange = this.onCropChange.bind(this);
    this.makeClientCrop = this.makeClientCrop.bind(this);
  }

  componentDidMount() {
    this.api.user.remote().then((resp) => {
      if (resp.error) {
        return
      }
      this.setState(resp.data);
    });
  }

  handleChange(e) {
    e.preventDefault();
    const {name, value} = e.target;
    this.setState({
      [name]: value
    });
  }

  handleClick(e) {
    this.refs.file.click();
  }

  handleFileChange(e) {
    if (e.target.files && e.target.files.length > 0) {
      const reader = new FileReader();
      reader.addEventListener('load', () => {
        this.setState({ avatar_url: reader.result });
      });
      reader.readAsDataURL(e.target.files[0]);
    }
  }

  onImageLoaded(image) {
    this.imageRef = image;
  }

  onCropComplete(crop) {
    this.makeClientCrop(crop);
  }

  onCropChange(crop, percentCrop) {
  };

 makeClientCrop(crop) {
    if (this.imageRef) {
      crop.width = this.imageRef.naturalWidth;
      crop.height = this.imageRef.naturalHeight;
      if (crop.width < 256 || crop.height < 256) {
        return;
      }
      if (crop.width > crop.height) {
        crop.width = crop.height;
      } else {
        crop.height = crop.width;
      }
      const canvas = document.createElement("canvas");
      let w = crop.width > 512 ? 512 : crop.width;
      canvas.width = crop.width;
      canvas.height = crop.height;
      const ctx = canvas.getContext("2d");
      ctx.drawImage(
        this.imageRef,
        0,
        0,
        crop.width,
        crop.height,
        0,
        0,
        w,
        w
      );
      this.setState({avatar_url: canvas.toDataURL()});
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    // TODO should use redirect
    const history = this.props.history;
    this.api.user.update(this.state).then((resp) => {
      this.setState({submitting: false});
      if (resp.error) {
        return
      }
      history.push('/');
    });
  }

  render() {
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

    let state = this.state;

    return (
      <div className='container'>
        <main className='column main'>
          <div className={style.profile}>
            <h2>{i18n.t('user.edit')}</h2>
            <form onSubmit={this.handleSubmit}>
              <div className={style.group}>
                <input type='file' ref='file' className={style.file} onChange={this.handleFileChange} />
                <ReactCrop
                  src={state.avatar_url}
                  crop={{aspect: 1}}
                  className={style.crop}
                  onImageLoaded={this.onImageLoaded}
                  onComplete={this.onCropComplete}
                  onChange={this.onCropChange}
                />
                <div className={style.image}>
                  {
                    state.avatar_url && (
                      <img src={state.avatar_url} onClick={this.handleClick} className={style.cover} />
                    )
                  }
                </div>
              </div>
              <div>
                <label name='nickname'>{i18n.t('user.nickname')}</label>
                <input type='text' name='nickname' value={state.nickname} autoComplete='off' onChange={this.handleChange} />
              </div>
              <div>
                <label name='biography'>{i18n.t('user.biography')}</label>
                <textarea type='text' name='biography' value={state.biography} onChange={this.handleChange} />
              </div>
              <div className='action'>
                <Button type='submit' class='submit' text={i18n.t('general.submit')} disabled={state.submitting} />
              </div>
            </form>
          </div>
        </main>
        <aside className='column aside'>
        </aside>
      </div>
    )
  }
}

export default Edit;
