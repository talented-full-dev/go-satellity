import React, {useState, useEffect} from 'react';
import TimeAgo from 'react-timeago';
import showdown from 'showdown';
import PropTypes from 'prop-types';
import API from 'api/index.js';
import Avatar from 'users/avatar.js';
import New from './new.js';

import style from './index.module.scss';

const Index = (props) => {
  const i18n = window.i18n;
  const api = new API();
  const converter = new showdown.Converter();

  const {topicId, commentsCount} = props;

  const [comments, setComments] = useState([]);

  useEffect(() => {
    api.comment.index(topicId).then((resp) => {
      if (resp.error) {
        return;
      }
      const comments = resp.data.map((comment) => {
        comment.body = converter.makeHtml(comment.body);
        return comment;
      });
      setComments(comments);
    });
  }, [topicId]);

  const handleSubmit = (comment) => {
    comment.body = converter.makeHtml(comment.body);
    comments.push(comment);
    setComments(comments);
  };

  const commentsView = comments.map((comment) => {
    return (
      <li className={style.comment} key={comment.comment_id}>
        <div className={style.profile}>
          <Avatar user={comment.user} />
          <div className={style.detail}>
            {comment.user.nickname}
            <div className={style.time}>
              <TimeAgo date={comment.created_at} />
            </div>
          </div>
        </div>
        <article className='md' dangerouslySetInnerHTML={{__html: comment.body}}>
        </article>
      </li>
    );
  });

  const commentsContainer = (
    <div className={style.container}>
      <h3>{i18n.t('comment.count', {count: commentsCount})}</h3>
      <ul className={style.comments}>
        {commentsView}
      </ul>
    </div>
  );

  return (
    <div>
      {commentsCount > 0 && commentsContainer}
      <New
        topicId={topicId}
        handleSubmit={handleSubmit} />
    </div>
  );
};

Index.propTypes = {
  topicId: PropTypes.string,
  commentsCount: PropTypes.number,
};

export default Index;
