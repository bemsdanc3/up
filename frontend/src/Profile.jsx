import { useState, useEffect } from 'react'

function Profile({ userId }) {
    const [userData, setUserData] = useState({})
    const [userDataLoaded, setUserDataLoaded] = useState(false)

    const pfp = 'http://localhost:8080/uploads/covers-albums/5d3469bf-bcc9-4752-ac06-1cb888d92cc4.jpg'

    const loadProfileInfo = async () => {
        try {
            let userLink;
            if (userId == 'self') {
                userLink = 'user/my/profile'
            } else {
                userLink = `users/profile/${userId}`
            }
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
    
    useEffect(()=>{
        console.log(userId);
        loadProfileInfo();
    }, [userDataLoaded, userId])

    return (
        <>
            <div id="profilePage">
                {userData && userDataLoaded &&
                <>
                    <div id="userInfo">
                        <img src={pfp} alt="img" />
                        <div id="userTextInfo">
                            <h1>{userData.login}</h1>
                            {userData.email &&<h2>{userData.email}</h2>}
                        </div>
                    </div>
                    <h2>User tracks:</h2>
                    <div id="userTracks">
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                        <div className="track" >
                        <div className="trackLeftInfo">
                            <img src="" alt="" />
                            <div className="trackTextInfo">
                            <span>title</span>
                            </div>
                        </div>
                        <div className="trackDurationInfo">
                            <span>
                            5:55
                            </span>
                        </div>
                        </div>
                    </div>
                    <h2>User albums:</h2>
                    <div id="userAlbums">
                        <div className="album">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="album">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="album">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="album">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                    </div>
                    <h2>User playlists:</h2>
                    <div id="userPlaylists">
                        <div className="playlist">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="playlist">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="playlist">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                        <div className="playlist">
                            <img src="" alt="" />
                            <span>Title</span>
                        </div>
                    </div>
                </>
                }
            </div>
        </>
    )
}

export default Profile
