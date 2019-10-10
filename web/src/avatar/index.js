import style from './index.scss';
import React, {Component} from 'react';
import Avatar, {Piece} from 'avataaars';
import Colors from 'avataaars/dist/avatar/top/facialHair/Colors.js';

class Index extends Component {
  constructor(props) {
    super(props);

    this.state = {
      action: 'hair',
      avatar: 'Circle',
      top: 'ShortHairDreads01',
      accessory: 'Prescription02',
      hairColor: 'BrownDark',
      hatColor: 'PastelBlue',
      facialHair: 'BeardLight',
      facialHairColor: 'Red',
      clothe: 'GraphicShirt',
      clotheColor: 'PastelBlue',
      graphic: 'Resist',
      eye: 'Happy',
      eyebrow: 'Default',
      mouth: 'Smile',
      skin: 'Light',
      hair_shapes: ['NoHair','Eyepatch','Hat','Hijab','Turban','WinterHat1','WinterHat2','WinterHat3','WinterHat4','LongHairBigHair','LongHairBob','LongHairBun','LongHairCurly','LongHairCurvy','LongHairDreads','LongHairFrida','LongHairFro','LongHairFroBand','LongHairNotTooLong','LongHairShavedSides','LongHairMiaWallace','LongHairStraight','LongHairStraight2','LongHairStraightStrand','ShortHairDreads01','ShortHairDreads02','ShortHairFrizzle','ShortHairShaggyMullet','ShortHairShortCurly','ShortHairShortFlat','ShortHairShortRound','ShortHairShortWaved','ShortHairSides','ShortHairTheCaesar','ShortHairTheCaesarSidePart'],
      accessories: ['Blank','Kurt','Prescription01','Prescription02','Round','Sunglasses','Wayfarers'],
      hat_colors: ['Black','Blue01','Blue02','Blue03','Gray01','Gray02','Heather','PastelBlue','PastelGreen','PastelOrange','PastelRed','PastelYellow','Pink','Red','White'],
      hair_colors: ['Auburn','Black','Blonde','BlondeGolden','Brown','BrownDark','PastelPink','Platinum','Red','SilverGray'],
      facial_hairs: ['Blank','MoustacheFancy','MoustacheMagnum','BeardLight','BeardMedium', 'BeardMajestic'],
      facial_hair_colors: ['Auburn','Black','Blonde','BlondeGolden','Brown','BrownDark','Platinum','Red'],
      clothes: ['BlazerShirt','BlazerSweater','CollarSweater','GraphicShirt','Hoodie','Overall','ShirtCrewNeck','ShirtScoopNeck','ShirtVNeck'],
      clothe_colors: ['Black','Blue01','Blue02','Blue03','Gray01','Gray02','Heather','PastelBlue','PastelGreen','PastelOrange','PastelRed','PastelYellow','Pink','Red','White'],
      eyes: ['Close','Cry','Default','Dizzy','EyeRoll','Happy','Hearts','Side','Squint','Surprised','Wink','WinkWacky'],
      eyebrows: ['Angry','AngryNatural','Default','DefaultNatural','FlatNatural','RaisedExcited','RaisedExcitedNatural','SadConcerned','SadConcernedNatural','UnibrowNatural','UpDown','UpDownNatural'],
      mouths: ['Concerned','Default','Disbelief','Eating','Grimace','Sad','ScreamOpen','Serious','Smile','Tongue','Twinkle','Vomit'],
      skins: ['Tanned','Yellow','Pale','Light','Brown','DarkBrown','Black'],
      graphics: ['Blank','Skull','SkullOutline','Bat','Cumbia','Deer','Diamond','Hola','Selena','Pizza','Resist','Bear']
    };

    this.handleClick = this.handleClick.bind(this);
    this.handleActionClick = this.handleActionClick.bind(this);
  }

  handleClick(e, k, v) {
    e.preventDefault();
    this.setState({
      [k]: v
    });
  }

  handleActionClick(e, v) {
    e.preventDefault();
    this.setState({
      action: v
    });
  }

  render () {
    const state = this.state;

    const actions = ['hair','accessory','beard','clothe','eye','eyebrow','mouth', 'skin'].map((o) => {
      return (
        <span key={o} className={`${style.action} ${state.action === o ? style.current : ''}`} onClick={(e) => this.handleActionClick(e, o)}>{o}</span>
      )
    });

    const hairShapes = state.hair_shapes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'top', o)}>
          <Piece pieceType='top' pieceSize='100' topType={o} hairColor='Blank'/>
        </div>
      )
    });

    const accessories = state.accessories.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'accessory', o)}>
          <Piece pieceType='accessories' pieceSize='100' accessoriesType={o}/>
        </div>
      )
    });

    const facialHairs = state.facial_hairs.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'facialHair', o)}>
          <Piece pieceType='facialHair' pieceSize='100' facialHairType={o}/>
        </div>
      )
    });

    const clothes = state.clothes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'clothe', o)}>
          <Piece pieceType='clothe' pieceSize='100' clotheType={o}/>
        </div>
      )
    });

    const graphics = state.graphics.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'graphic', o)}>
          <Piece pieceType="graphics" pieceSize="100" graphicType={o} />
        </div>
      )
    });

    const eyes = state.eyes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'eye', o)}>
          <Piece pieceType='eyes' pieceSize='100' eyeType={o}/>
        </div>
      )
    });

    const eyebrows = state.eyebrows.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'eyebrow', o)}>
          <Piece pieceType='eyebrows' pieceSize='100' eyebrowType={o}/>
        </div>
      )
    });

    const mouths = state.mouths.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'mouth', o)}>
          <Piece pieceType='mouth' pieceSize='100' mouthType={o}/>
        </div>
      )
    });

    const skins = state.skins.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'skin', o)}>
          <Piece pieceType='skin' pieceSize='100' skinColor={o}/>
        </div>
      )
    });

    const hairColors = state.hair_colors.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'hairColor', o)}>
          <Piece pieceType='hairColor' pieceSize='48' hairColor={o}/>
        </div>
      )
    });

    const hatColors = state.hat_colors.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'hatColor', o)}>
          <Piece pieceType='hatColor' pieceSize='48' hatColor={o}/>
        </div>
      )
    });

    const facialHairColors = state.facial_hair_colors.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'facialHairColor', o)}>
          <Piece pieceType='facialHairColor' pieceSize='48' facialHairColor={o}/>
        </div>
      )
    });

    const clotheColors = state.clothe_colors.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'clotheColor', o)}>
          <Piece pieceType='clotheColor' pieceSize='48' clotheColor={o}/>
        </div>
      )
    });

    return (
      <div className={style.container}>
        <div className={style.canvas}>
          <div className={style.profile}>
            <div className={style.avatar}>
            <Avatar
              style={{width: '24rem', height: '24rem'}}
              avatarStyle={state.avatar}
              topType={state.top}
              accessoriesType={state.accessory}
              hairColor={state.hairColor}
              hatColor={state.hatColor}
              facialHairType={state.facialHair}
              facialHairColor={state.facialHairColor}
              clotheType={state.clothe}
              clotheColor={state.clotheColor}
              graphicType={state.graphic}
              eyeType={state.eye}
              eyebrowType={state.eyebrow}
              mouthType={state.mouth}
              skinColor={state.skin} />
            </div>
            {state.action === 'hair' && state.top !== 'NoHair' && state.top !== 'LongHairFrida' && state.top.includes('Hair') && <div> {hairColors} </div>}
              {state.action === 'hair' && state.top !== 'Eyepatch' && state.top !== 'Hat' && !state.top.includes('Hair') && <div> {hatColors} </div>}
          </div>
          <div>
            <div className={style.actions}>
              {actions}
            </div>
            {state.action === 'hair' && <div> {hairShapes} </div>}
            {state.action === 'accessory' && <div> {accessories} </div>}
            {state.action === 'beard' && <div> {facialHairs} </div>}
            {state.action === 'clothe' && <div> {clothes} </div>}
            {state.action === 'clothe' && state.clothe.includes('Graphic') && <div> {graphics} </div>}
            {state.action === 'eye' && <div> {eyes} </div>}
            {state.action === 'eyebrow' && <div> {eyebrows} </div>}
            {state.action === 'mouth' && <div> {mouths} </div>}
            {state.action === 'skin' && <div> {skins} </div>}
            {state.action === 'beard' && state.facialHair !== 'Blank' && <div> {facialHairColors} </div>}
            {state.action === 'clothe' && !state.clothe.includes('Blazer') && <div> {clotheColors} </div>}
          </div>
        </div>
      </div>
    )
  }
}

export default Index;
