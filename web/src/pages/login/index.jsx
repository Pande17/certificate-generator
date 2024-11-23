import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Login } from "../api middleware/index";
import backgroundImage from "../assets/Element.svg";

const LoginPage = () => {
  const [admin_name, setAdminName] = useState("");
  const [admin_password, setAdminPassword] = useState("");
  const [msg, setMsg] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const loginHandle = async (e) => {
    e.preventDefault();
    setMsg("");

    if (!admin_name || !admin_password) {
      setMsg("Username dan password harus diisi!");
      return;
    }

    try {
      const response = await Login.post(
        "/",
        { admin_name, admin_password },
        {
          validateStatus: (status) => status < 500,
        }
      );

   if (response.status === 200) {
     const token = response.data.token;
     localStorage.setItem("token", token);
     console.log("Token saved:", token); // Debugging log
     navigate("/dashboard");
   } else if (response.status === 401) {
     setMsg(response.data.message || "Username atau password salah");
   } else {
     setMsg("Terjadi kesalahan, coba lagi.");
   }
    } catch (err) {
      console.error("Error: ", err);
      setMsg("Terjadi kesalahan.");
    }
  };

  return (
    <div
      className="flex justify-center items-center min-h-screen bg-Background"
      style={{
        backgroundImage: `url(${backgroundImage})`,
        backgroundRepeat: "no-repeat",
        backgroundPosition: "center",
        backgroundSize: "cover",
      }}
    >
      <div className="max-h-fit justify-center align-middle bg-white p-6 text-center rounded-md border-black sm:w-2/6">
        <form onSubmit={loginHandle}>
          <p className="font-Poppins font-bold text-3xl text-Text mb-6 text-center">
            Login
          </p>
          {msg && <p className="text-red-600 mb-3">{msg}</p>}
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
              id="password"
            />
            <label className="flex gap-3 mt-3 items-center">
              <input
                type="checkbox"
                checked={showPassword}
                onChange={togglePasswordVisibility}
              />
              <span>Show password</span>
            </label>
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
  );
};

export default LoginPage;
