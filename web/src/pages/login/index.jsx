import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios"; 
import backgroundImage from "../assets/Element.svg"

const LoginPage = () => {
  const [admin_name, setAdminName] = useState(""); 
  const [admin_password, setAdminPassword] = useState(""); 
  const [msg, setMsg] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate(); 

const passwodVisibility = () => {
setShowPassword(!showPassword)
}

  const loginHandle = async (e) => {
    e.preventDefault();
    setMsg('')
    try {
      const response = await axios.post("http://127.0.0.1:3000/api/login", {
        admin_name,
        admin_password,
      }, {
        validateStatus: function (status) {
          return status<500
        }
      });

      if (response.status === 200) { 
        navigate("/dashboard");
        setMsg("Selamat datang!");
      } else if (response.status === 401){ 
        setMsg("Username atau password salah");
      } else {
        setMsg('terjadi kesalahan coba lagi')
      }
    } catch (err) {
         console.error("Error: ", err);
         setMsg("Terjadi kesalahan");
       
    }
  };

  return (
    <>
      <div className="flex justify-center items-center min-h-screen bg-Background"  style={{
        backgroundImage: `url(${backgroundImage})`,
        backgroundRepeat:"no-repeat",
        backgroundPosition: "center",
      }}>
        <div className="max-h-fit justify-center align-middle bg-white p-6 text-center rounded-md border-black sm:w-2/6">
          <form onSubmit={loginHandle}>
            <p className="font-Poppins font-bold text-3xl text-Text mb-6 text-center">
              Login
            </p>
            {msg && (<p className="text-red-600 mb-3">{msg}</p>)}
            <div className="mb-4 text-left">
              <p className="font-Poppins font-regular">Username</p>
              <input
                type="text"
                value={admin_name}
                placeholder="Username"
                onChange={(e) => setAdminName(e.target.value)}
                required
                className="border-2 p-1 rounded-md border-black w-full"
              />
            </div>
            <div className="mb-6 text-left">
              <p>Password</p>
              <input
                type={showPassword ? "text" : "password"}
                placeholder="Password"
                value={admin_password}
                onChange={(e) => setAdminPassword(e.target.value)}
                required
                className="border-2 p-1 rounded-md border-black w-full"
                id="passwod"
              />
              <div className="flex gap-3 mt-3">
                <input
                  type="checkbox"
                  checked={showPassword}
                  onChange={passwodVisibility}
                />
                <p className=" m-0">show password</p>
              </div>
            </div>
            <button
              type="submit"
              className="font-Poppins font-semibold text-white bg-Text pt-1 pb-1 pl-4 pr-4 rounded-md"
            >
              Login
            </button>
          </form>
        </div>
      </div>
    </>
  );
};

export default LoginPage;
