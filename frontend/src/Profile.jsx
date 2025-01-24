import { useState, useEffect } from 'react'

function Profile({ userId, myData, handleTrackClick, playlistFunc, albumFunc }) {
    const [userData, setUserData] = useState({})
    const [tracks, setTracks] = useState([])
    const [albums, setAlbums] = useState([])
    const [playlists, setPlaylists] = useState([])
    const [userDataLoaded, setUserDataLoaded] = useState(false)

    const pfp = 'http://localhost:8080/uploads/covers-albums/5d3469bf-bcc9-4752-ac06-1cb888d92cc4.jpg'

    const loadProfileInfo = async () => {
        try {
            let userLink = `users/profile/${userId}`;
            const userProfileRes = await fetch(`http://localhost:8080/${userLink}`,{
                method: 'GET',
                credentials: 'include',
                withCredentials: true,
            });
            if (userProfileRes.ok) {
                console.log("salamalekum")
                const resData = await userProfileRes.json()
                console.log(resData)
                setUserData(resData.User || resData);
                setUserDataLoaded(true);
            } else {
                const errorData = await userProfileRes.json();
                console.log(errorData.error);
            }
        } catch (error) {
            console.log(error);
        }
    }

    const getAllTracks = async () => {
        console.log('getting all tracks');
        const response = await fetch(`http://localhost:8080/tracks/all/${userId}`, {
          method: 'GET',
          credentials: 'include',
          withCredentials: true,
        })
        console.log(response);
        const responseData = await response.json();
        console.log(responseData);
        if (response.ok) {
          setTracks(responseData);
        } else {
          console.log('жопа');
        }
      }

    const getAllAlbums = async () => {
        console.log('getting all tracks');
        const response = await fetch(`http://localhost:8080/album/all/${userId}`, {
          method: 'GET',
          credentials: 'include',
          withCredentials: true,
        })
        console.log(response);
        const responseData = await response.json();
        console.log(responseData);
        if (response.ok) {
          setAlbums(responseData);
        } else {
          console.log('жопа');
        }
      }

    const getAllPlaylists = async () => {
        console.log('getting all tracks');
        const response = await fetch(`http://localhost:8080/playlists/all/${userId}`, {
          method: 'GET',
          credentials: 'include',
          withCredentials: true,
        })
        console.log(response);
        const responseData = await response.json();
        console.log(responseData);
        if (response.ok) {
          setPlaylists(responseData);
        } else {
          console.log('жопа');
        }
      }
    
    useEffect(()=>{
        console.log(userId);
        loadProfileInfo();
        getAllAlbums();
        getAllPlaylists();
        getAllTracks();
    }, [userDataLoaded, userId])

    const formatDuration = (duration) => {
        const minutes = Math.floor(duration / 60);
        const seconds = duration % 60;
        return `${minutes}:${seconds.toString().padStart(2, '0')}`; // Форматирует секунды с ведущим нулём
      };

    return (
        <>
            <div id="profilePage">
                {userData && userDataLoaded &&
                <>
                    <div id="userInfo">
                        <img src={userData.pfp } alt="img" />
                        <div id="userTextInfo">
                            <h2>{userData.role}</h2>
                            <h1>{userData.login}</h1>
                            {(myData.role == 'admin' || myData.user_id == userData.id) &&
                                <h2>{userData.email}</h2>
                            }
                        </div>
                    </div>
                    {tracks && tracks.length >=1 &&
                    <>
                        <h2>Треки пользователя:</h2>
                        <div id="userTracks">
                                {tracks.map((track) => {
                                    return (
                                        <>
                                        <div className="track" onClick={()=>{handleTrackClick(track)}}>
                                            <div className="trackLeftInfo">
                                                <img src={track.cover} alt="" />
                                                <div className="trackTextInfo">
                                                <span>{track.title}</span>
                                                </div>
                                            </div>
                                            <div className="trackDurationInfo">
                                                <span>
                                                {formatDuration(track.duration)}
                                                </span>
                                            </div>
                                        </div>
                                        </>
                                    )
                                })}
                        </div>
                    </>
                    }
                    {albums && albums.length >= 1 &&
                    <>
                        <h2>Альбомы пользователя:</h2>
                        <div id="userAlbums">
                                {albums.map((album)=>{
                                    return (
                                        <div className="album" onClick={()=>albumFunc(album.id)}>
                                            <img src={album.cover} alt="" />
                                            <span>{album.title}</span>
                                        </div>
                                    )
                                })}
                        </div>
                    </>
                    }
                    {playlists && playlists.length >= 1 &&
                    <>
                        <h2>Плейлисты пользователя:</h2>
                        <div id="userPlaylists">
                                {playlists.map((playlist)=>{
                                    return (
                                        <div className="playlist" onClick={()=>playlistFunc(playlist.id)}>
                                            <img src={playlist.cover} alt="" />
                                            <span>{playlist.title}</span>
                                        </div>
                                    )
                                })}
                        </div>
                    </>
                    }
                </>
                }
            </div>
        </>
    )
}

export default Profile
