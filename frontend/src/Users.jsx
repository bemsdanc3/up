import { useEffect, useState } from 'react'

function Users({ userProfile }) {
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
  }, [])

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
                            <div className="userAdminBtns">
                                <button 
                                    onClick={()=>{console.log('btn working')}}
                                    className={user.role == 'admin' ? 'selected' : ''}
                                >
                                    Admin
                                </button>
                                <button 
                                    onClick={()=>{console.log('btn working')}}
                                    className={user.role == 'artist' ? 'selected' : ''}
                                >
                                    Artist
                                </button>
                                <button 
                                    onClick={()=>{console.log('btn working')}}
                                    className={user.role == 'user' ? 'selected' : ''}
                                >
                                    User
                                </button>
                            </div>
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
