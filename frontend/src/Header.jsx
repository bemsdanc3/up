import { useState } from 'react'
import HomeIcon from './assets/HomeIcon.svg?react';
import SearchIcon from './assets/SearchIcon.svg?react';
import ProfileIcon from './assets/ProfileIcon.svg?react';
import UsersIcon from './assets/UsersIcon.svg?react';
import TrackListIcon from './assets/TrackListIcon.svg?react';

function Header({ profileFunc, homeFunc, loginData, usersFunc, tracksFunc }) {
  
  return (
    <>
      <header>
        <div id="headerLeft">
            <button onClick={()=>{homeFunc()}}>
              <HomeIcon />
              Главная
            </button>
            {loginData.role == 'admin' &&
            <>
            <button onClick={()=>{usersFunc()}}>
              <UsersIcon />
              Пользователи
            </button>
            <button onClick={()=>{tracksFunc()}}>
              <TrackListIcon />
              Треки
            </button>
            </>
            }
        </div>
        {/* <div id="headerCenter">
            <div id="searchDiv">
                <input type="text" />
                <button>
                  <SearchIcon />
                  Поиск
                </button>
            </div>
        </div> */}
        <div id="headerRight">
            <button onClick={()=>{profileFunc()}}>
              <ProfileIcon />
              Профиль
            </button>
        </div>
      </header>
    </>
  )
}

export default Header
