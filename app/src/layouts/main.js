import style from './main.module.scss';
import logo from '../assets/images/logo.svg';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import React, {Component} from 'react';
import {Route, Link, Routes} from 'react-router-dom';
import Config from '../components/config.js';
import API from '../api/index.js';
import Home from '../home/view.js';
import User from '../users/view.js';
import Topic from '../topics/view.js';
import Login from './login.js';

class MainLayout extends Component {
  constructor(props) {
    super(props);

    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    this.state = {p: ''};
  }

  render() {
    return (
      <div className={style.container}>
        <Header />
        <div className='wrapper'>
          <Routes>
            <Route index path='/' element={<Topic.Index />} />
            <Route exact path='/dashboard' element={<Home.Dashboard />} />
            <Route exact path='/categories/:id' element={<Topic.Index />} />
            <Route exact path='/user/edit' element={<User.Edit />} />
            <Route path='/users/:id' element={<User.Show />} />
            <Route exact path='/topics/new' element={<Topic.New />} />
            <Route path='/topics/:id/edit' element={<Topic.New />} />
            <Route path='/topics/:id' element={<Topic.Show />} />
          </Routes>
        </div>
      </div>
    );
  }
}

class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {logging: false};
    this.handleLoginClick = this.handleLoginClick.bind(this);
  }

  handleLoginClick(e) {
    const n = e.target.className;
    if (!(n.includes('close') || n.includes('modal') || n.includes('navi'))) {
      return;
    }
    this.setState({logging: !this.state.logging});
  }

  render() {
    const user = new API().user;
    let profile = <span className={style.navi} onClick={this.handleLoginClick}>Login</span>;
    if (user.loggedIn()) {
      profile = (
        <div className={style.navis}>
          <Link to='/topics/new' className={`${style.navi}`}> <FontAwesomeIcon icon={['fa', 'plus']} /> </Link>
          <Link to='/user/edit' className={`${style.navi} ${style.user}`}> {user.local().nickname} </Link>
        </div>
      );
    }

    return (
      <div>
        <header className={style.header}>
          <Link className={style.site} to='/'>
            <img className={style.logo} src={logo} alt={Config.Name} />
            <span className={style.name}>{Config.Name}</span>
          </Link>

          <div className={style.menus}>
            <Link className={`${style.menu} ${window.location.pathname === '/' ? style.current : ''}` } to='/'>
                Home
            </Link>
          </div>
          {profile}
        </header>
        {this.state.logging && <Login handleLoginClick={this.handleLoginClick} />}
      </div>
    );
  }
}

export default MainLayout;
