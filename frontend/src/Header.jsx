import { useState } from 'react'
import HomeIcon from './assets/HomeIcon.svg?react';
import SearchIcon from './assets/SearchIcon.svg?react';
import ProfileIcon from './assets/ProfileIcon.svg?react';
import UsersIcon from './assets/UsersIcon.svg?react';
import TrackListIcon from './assets/TrackListIcon.svg?react';

function Header({ profileFunc, homeFunc }) {
  
  return (
    <>
      <header>
        <div id="headerLeft">
            <button onClick={()=>{homeFunc()}}>
              <HomeIcon />
              Home
            </button>
            <button onClick={()=>{console.log('Тут будет страница пользователей для админа')}}>
              <UsersIcon />
              Users
            </button>
            <button onClick={()=>{console.log('Тут будет страница треков для админа')}}>
              <TrackListIcon />
              Tracks
            </button>
        </div>
        <div id="headerCenter">
            <div id="searchDiv">
                <input type="text" />
                <button>
                  <SearchIcon />
                  Search
                </button>
            </div>
        </div>
        <div id="headerRight">
            <button onClick={()=>{profileFunc()}}>
              <ProfileIcon />
              Profile
            </button>
        </div>
      </header>
    </>
  )
}

export default Header
