import { useState, useEffect } from 'react'

function Profile() {
    const [userData, setUserData] = useState({})
    const [userDataLoaded, setUserDataLoaded] = useState(false)

    const loadProfileInfo = async () => {
        try {
            const userProfileRes = await fetch(`http://localhost:8080/user/my/profile`,{
                method: 'GET',
                credentials: 'include',
                withCredentials: true,
            });
            if (userProfileRes.ok) {
                console.log("salamalekum")
                const resData = await userProfileRes.json()
                console.log(resData)
                setUserData(resData.User);
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
        if (!userDataLoaded) {
            loadProfileInfo();

        }
    }, [userDataLoaded])

    return (
        <>
            {userData && userDataLoaded &&
            <>
                <img src={userData.pfp} alt="img" />
                <h1>{userData.login}</h1>
                <h2>{userData.email}</h2>
            </>
            }
        </>
    )
}

export default Profile
