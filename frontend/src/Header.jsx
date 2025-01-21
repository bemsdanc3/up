import { useState } from 'react'

function Header({ profileFunc, homeFunc }) {
  
  return (
    <>
      <header>
        <div id="headerLeft">
            <button onClick={()=>{homeFunc()}}>Home</button>
        </div>
        <div id="headerCenter">
            <div id="searchDiv">
                <input type="text" />
                <button>Search</button>
            </div>
        </div>
        <div id="headerRight">
            <button onClick={()=>{profileFunc()}}>Profile</button>
        </div>
      </header>
    </>
  )
}

export default Header
