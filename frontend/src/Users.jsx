import { useEffect, useState } from 'react'

function Users({ userProfile, userData }) {
  const [users, setUsers] = useState([]);

  const getAllUsers = async () => {
    console.log('getting all tracks');
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      credentials: 'include',
      withCredentials: true,
    })
    console.log(response);
    const responseData = await response.json();
    console.log('users responseData');
    console.log(responseData);
    if (response.ok) {
      setUsers(responseData);
    } else {
      console.log('жопа');
    }
  }

  useEffect(()=>{
    getAllUsers();
    console.log(userData)
  }, [])

  const updUserRole = async (role, id) => {
    try {
      const updUserRoleRes = await fetch(`http://localhost:8080/users/update-role`,{
        method: 'PATCH',
        credentials: 'include',
        withCredentials: true,
        headers: {
          "Content-Type": "application/json", 
        },
        body: JSON.stringify({
          role: role,
          user_id: id
        }),
      });
      const responseData = await updUserRoleRes.json();
      console.log("responseData");
      console.log(responseData);
      if (updUserRoleRes.ok) {
        console.log("salamalekum")
        getAllUsers();
      } else {
        const errorData = await updUserRoleRes.json();
        console.log(errorData.error);
      }
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <> 
      <div id="usersPage">
        <h2>Пользователи: </h2>
        <div id="usersList">
            {users && users.length >=1 &&
                users.map((user)=>{
                    return (
                        <div className="user" key={user.id} onClick={()=>userProfile(user.id)}>
                            <div className="userLeftInfo">
                                <img src={user.pfp} alt="" />
                                <div className="userCardTextInfo">
                                    <span className='userLoginAndRole'>
                                        <span className='usrCardLogin'>{user.login}</span>
                                        <span className='usrCardRole'>{user.role}</span>
                                    </span>
                                    <span>{user.email}</span>
                                </div>
                            </div>
                            {userData.user_id != user.id &&
                            <div className="userAdminBtns">
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('admin', user.id)}}
                                    className={user.role == 'admin' ? 'selected' : ''}
                                >
                                    Admin
                                </button>
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('artist', user.id)}}
                                    className={user.role == 'artist' ? 'selected' : ''}
                                >
                                    Artist
                                </button>
                                <button 
                                    onClick={(e)=>{e.stopPropagation(); updUserRole('user', user.id)}}
                                    className={user.role == 'user' ? 'selected' : ''}
                                >
                                    User
                                </button>
                            </div>}
                        </div>
                    )
                })
            }
        </div>
      </div>
    </>
  )
}

export default Users
